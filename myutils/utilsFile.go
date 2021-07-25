package myutils

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func GetFileMd5(file string) string {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Open", err)
		return ""
	}
	defer f.Close()
	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		fmt.Println("Copy", err)
		return ""
	}

	return fmt.Sprintf("%x", md5hash.Sum(nil))
}

//获取单个文件的大小
func GetFileSize(file string) int64 {
	if fileInfo, err := os.Stat(file); err == nil {
		return fileInfo.Size()
	}
	return 0
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return false, err
}

//获取文件真实名称 大小写不铭感
func GetFileRealName(path string) string {
	path, _ = filepath.Abs(path)
	_, err := os.Stat(path)
	if err == nil {
		return path
	}
	if os.IsNotExist(err) {
		dir := filepath.Dir(path)
		for {
			if ok, _ := PathExists(dir); ok {
				break
			}
			dir = filepath.Dir(dir)
		}
		files := GetFileList(dir)
		if file, ok := files[strings.ToLower(path)]; ok {
			return file
		}
	}
	return ""
}

func GetFileList(path string) map[string]string {
	list := make(map[string]string)
	_ = filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		fileName := path
		if f.IsDir() {
			//files := GetFileList(fileName)
			//for k, v := range files {
			//	list[k] = v
			//}
			return nil
		}
		list[strings.ToLower(fileName)] = fileName
		return nil
	})
	return list
}

func GetFileContentAsStringLines(filePath string) ([]string, error) {
	result := []string{}
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return result, err
	}
	s := string(b)
	for _, lineStr := range strings.Split(s, "\n") {
		lineStr = strings.TrimSpace(lineStr)
		if lineStr == "" {
			continue
		}
		result = append(result, lineStr)
	}
	return result, nil
}

func WriteStringsToFile(filePath string, text []string) error {
	log.Println("WriteStringsToFile", filePath)
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	if _, err := f.WriteString(strings.Join(text, "\n")); err != nil {
		return err
	}
	return nil
}

func WriteToFile(name string, data []byte) error {
	os.MkdirAll(path.Dir(name), os.ModePerm)
	f, err := os.Create(name)
	if err != nil {
		log.Println(err.Error())
	}

	_, err = f.Write(data)
	if err != nil {
		log.Println(err.Error())
	}
	f.Close()
	return err
}

func GetAllFile(pathname string, depth, limit int) ([]string, error) {
	if rd, err := ioutil.ReadDir(pathname); err == nil {
		ret := make([]string, 0)
		for _, fi := range rd {
			if fi.IsDir() {
				if depth < limit {
					if files, err := GetAllFile(pathname+fi.Name()+"\\", depth+1, limit); err == nil {
						ret = append(ret, files...)
					} else {
						return ret, err
					}
				}
			} else {
				ret = append(ret, fmt.Sprintf("%s/%s", pathname, fi.Name()))
			}
		}
		return ret, nil
	} else {
		return nil, nil
	}
}

package myutils

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"strings"

	"io/ioutil"
)

func GbkToUtf8(s []byte) []byte {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil
	}
	return d
}

func Utf8ToGbk(s []byte) []byte {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil
	}
	return d
}

func ByteString(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}

// fromHexChar converts a hex character into its value and a success flag.
func fromHexChar(c byte) (byte, bool) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}

	return 0, false
}

func SliceToHexStringExt(data []byte) string {
	lss := make([]string, 0)
	lss = append(lss, fmt.Sprintf("\n===16进制数据===>"))
	for i := 0; i < len(data); i += 0x10 {
		buf := data[i:]
		if len(buf) > 0x10 {
			buf = buf[:0x10]
		}
		hexStr := fmt.Sprintf("% x", buf)
		gbkBuf := GbkToUtf8(buf)
		for i := 0; i < len(gbkBuf); i++ {
			if gbkBuf[i] <= 0x20 {
				gbkBuf[i] = '.'
			}
		}
		gbkStr := fmt.Sprintf("| %s", gbkBuf)
		lss = append(lss, fmt.Sprintf("\t==%.04x=>%-48s  %-16s", i, hexStr, gbkStr))
	}
	return strings.Join(lss, "\n")
}

//数据包扩展格式 转16进制数据格式
func HexExtToHex(hexext string) string {
	hexs := make([]string, 0)
	lss := strings.Split(hexext, "\n")
	for _, v := range lss {
		hexLine := GetCenterText(v, "=>", "|")
		if hexLine != "" {
			hexs = append(hexs, strings.Trim(hexLine, " "))
		}
	}
	return strings.Join(hexs, " ")
}

//取文本中间
func GetCenterText(s, start, stop string) string {
	index1 := 0
	index2 := len(s)
	if start != "" {
		index1 = strings.Index(s, start) + len(start)
	}
	if stop != "" {
		index2 = strings.Index(s, stop)
	}
	if index2 > index1 && index1 >= 0 {
		return s[index1:index2]
	}
	return ""
}

package myutils

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
	"strings"
)

func GetRunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	name := f.Name()
	short := f.Name()
	for i := len(name) - 1; i > 0; i-- {
		if name[i] == '.' {
			short = name[i+1:]
			break
		}
	}
	return short
}

func GetRunFileLine() string {
	var file string
	var line int
	var ok bool
	_, file, line, ok = runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short
	return fmt.Sprintf("%s:%d", file, line)
}

func GetObjName(obj interface{}) string {
	return reflect.ValueOf(obj).Type().Name()
}

func GetObjType(obj interface{}) string {
	return reflect.ValueOf(obj).Type().String()
}

func GetBetweenStr(str, start, end string) string {
	n := 0
	if start == "" {
		n = 0
	} else {
		n = strings.Index(str, start)
		if n == -1 {
			return ""
		}
	}
	str = string([]byte(str)[n+len(start):])
	m := 0
	if end == "" {
		m = len(str)
	} else {

		m = strings.Index(str, end)
		if m == -1 {
			return ""
		}
	}
	str = string([]byte(str)[:m])
	return str
}

func ParsePortStr(txt string) (n int, err error) {
	txts := strings.Split(txt, ";")
	portList := &SubnetPortList{}
	for _, v := range txts {
		portItem, err := NewSubnetPort(v)
		if err != nil {
			return 0, err
		} else {
			if err := portList.CheckRepettion(portItem); err != nil {
				return 0, err
			} else {
				portList.Add(portItem)
			}
		}
	}
	return portList.GetPortCount(), nil
}

func AbsInt(x int32) int32 {
	if x >= 0 {
		return x
	}
	return -x
}

func MaxInt(a, b int32) int32 {
	if a >= b {
		return a
	}
	return b
}

func MaxUint32(a, b uint32) uint32 {
	if a < b {
		return b
	}
	return a
}

func MaxInt32(a, b int32) int32 {
	if a < b {
		return b
	}
	return a
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MinUint32(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func MinInt32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

//开方
func Sqrt(x uint32) uint32 {
	f := float64(x)
	ff := math.Sqrt(f)
	return uint32(ff)
}

func GetDistanceU32Magic(x, y, x2, y2 uint32) uint32 {
	d1 := uint32(MaxInt(AbsInt(int32(x-x2)), AbsInt(int32(y-y2))))
	d2 := uint32(MinInt32(AbsInt(int32(x-x2)), AbsInt(int32(y-y2))))
	return d1 + d2/2
}

func GetDistanceU32(x, y, x2, y2 uint32) uint32 {
	return uint32(MaxInt(AbsInt(int32(x-x2)), AbsInt(int32(y-y2))))
}

func GetDistance(x, y, x2, y2 int32) int32 {
	return MaxInt(AbsInt(x-x2), AbsInt(y-y2))
}

func GetMinDistanceU32(x, y, x2, y2 uint32) uint32 {
	return uint32(MinInt32(AbsInt(int32(x-x2)), AbsInt(int32(y-y2))))
}

func ListStrIndex(ls []string, s string) int {
	for k, v := range ls {
		//fmt.Println(v)
		if v == s {
			return k
		}
	}
	return -1
}

func Listu32Index(ls []uint32, s uint32) int {
	for k, v := range ls {
		//fmt.Println(v)
		if v == s {
			return k
		}
	}
	return -1
}

func ListIntIndex(ls []int, s int) int {
	for k, v := range ls {
		//fmt.Println(v)
		if v == s {
			return k
		}
	}
	return -1
}

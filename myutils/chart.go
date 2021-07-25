package myutils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//duration 时间长度 单位分钟
func GetChartForNodeState(durationMinute int, text ...string) []int {
	str := strings.Join(text, ";")
	strings.Count(str, ";")
	ret := make([]int, durationMinute)
	nowTime := time.Now().In(time.FixedZone("", 0))
	currentIndex := nowTime.Hour()*60 + nowTime.Minute()
	strLen := len(str)
	pi, pn := 0, 0
	index := -1
	for i := 0; i < strLen; i++ {
		switch str[i] {
		case ':':
			buf := str[pi:i]
			if n, err := strconv.ParseInt(buf, 16, 16); err == nil {
				index = int(n)
			}
			pn = i + 1
		case ';':
			if index >= 0 {
				if n, err := strconv.ParseInt(str[pn:i], 16, 16); err == nil {
					newIndex := (index + durationMinute - currentIndex - 1) % durationMinute
					if newIndex >= durationMinute {
						fmt.Println(newIndex)
					}
					ret[newIndex] = int(n)
					index = -1
					pi = i + 1
				}
			} else {
				pi = i + 1
				pn = pi
			}
		default:

		}
	}
	return ret
}

//duration 时间长度 单位分钟
func GetChartForText(text string, durationMinute int) []int {
	ret := make([]int, durationMinute)
	strLen := len(text)
	pi, pn := 0, 0
	index := -1
	for i := 0; i < strLen; i++ {
		switch text[i] {
		case ':':
			buf := text[pi:i]
			if n, err := strconv.ParseInt(buf, 36, 32); err == nil {
				index = int(n)
			}
			pn = i + 1
		case ';':
			if index >= 0 {
				if int(index) < durationMinute {
					if n, err := strconv.ParseInt(text[pn:i], 36, 32); err == nil {
						ret[index] = int(n)
						index = -1
						pi = i + 1
					}
				}
			} else {
				pi = i + 1
				pn = pi
			}
		default:

		}
	}
	return ret
}

//duration 时间长度 单位分钟
func GetChartLabels(t *time.Time, durationMinute int) []string {
	ret := make([]string, durationMinute)
	var nowTime time.Time
	if t == nil {
		nowTime = time.Now().In(time.FixedZone("", 8*3600))
		nowTime = nowTime.Add(-time.Duration(durationMinute-1) * time.Minute)
	} else {
		nowTime = *t
	}
	for i := 0; i < durationMinute; i++ {
		ret[i] = nowTime.Format("01-02 15:04")
		nowTime = nowTime.Add(time.Minute)
	}
	return ret
}

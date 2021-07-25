package myutils

import "strings"

func StringInList(s string, list []string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func StringInListEx(s string, list []string) bool {
	for _, name := range list {
		if s == name {
			return true
		} else {
			isPick := true
			pnames := strings.Split(name, "*")
			for j := 0; j < len(pnames); j++ {
				cpName := strings.Trim(pnames[j], " ")
				if cpName != "" {
					if index := strings.Index(s, cpName); index >= 0 {
						s = s[index+len(cpName):]
					} else {
						isPick = false
						break
					}
				}
			}
			if isPick {
				return true
			}
		}
	}
	return false
}

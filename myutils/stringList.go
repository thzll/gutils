package myutils

import (
	"strings"
	"sync"
)

type StringList struct {
	list map[string]int
	l    sync.Mutex
}

func (s *StringList) Add(name string) {
	s.l.Lock()
	defer s.l.Unlock()
	if s.list == nil {
		s.list = make(map[string]int)
	}
	s.list[name] = 1
}

func (s *StringList) Clear() {
	s.l.Lock()
	defer s.l.Unlock()
	s.list = make(map[string]int)
}

func (s *StringList) Count() int {
	s.l.Lock()
	defer s.l.Unlock()
	return len(s.list)
}

func (s *StringList) IsEffactiveName(name string) int {
	isPick := true
	for k, v := range s.list {
		pNmae := k
		if name == pNmae {
			return v
		} else {
			isPick = true
			pnames := strings.Split(pNmae, "*")
			for j := 0; j < len(pnames); j++ {
				cpName := strings.Trim(pnames[j], " ")
				if cpName != "" {
					if index := strings.Index(name, cpName); index >= 0 {
						name = name[index+len(cpName):]
					} else {
						isPick = false
						break
					}
				}
			}
			if isPick {
				return v
			}
		}
	}
	return 0
}

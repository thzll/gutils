package myutils

import (
	"fmt"
	"strconv"
	"strings"
)

type SubnetPort struct {
	PortBegin int
	PortEnd   int
	ToPort    int
	PortCount int
	Txt       string
}

type SubnetPortList struct {
	List []*SubnetPort
}

func NewSubnetPort(txt string) (*SubnetPort, error) {
	portExTxt := ""
	toPort := 0
	portExs := strings.Split(txt, "->")
	if len(portExs) == 2 {
		portExTxt = portExs[0]
		if v, err := strconv.Atoi(portExs[1]); err != nil {
			return nil, fmt.Errorf("无效的整数", portExs[1])
		} else {
			toPort = v
		}
	} else if len(portExs) == 1 {
		portExTxt = portExs[0]
	} else {
		return nil, fmt.Errorf("格式错误[%s] 例如(200-300;400,500->600)", txt)
	}
	ports := strings.Split(portExTxt, "-")
	if len(ports) == 1 {
		port, err := strconv.Atoi(ports[0])
		if err != nil {
			return nil, fmt.Errorf("%s 不是有效个整数格式", txt)
		} else {
			return &SubnetPort{
				PortBegin: port,
				PortEnd:   port,
				ToPort:    toPort,
				PortCount: 1,
				Txt:       txt,
			}, nil
		}
	} else if len(ports) == 2 {
		portBegin, err := strconv.Atoi(ports[0])
		if err == nil {
			portEnd, err := strconv.Atoi(ports[1])
			if err == nil {
				if portEnd >= portBegin {
					return &SubnetPort{
						PortBegin: portBegin,
						PortEnd:   portEnd,
						ToPort:    toPort,
						PortCount: portEnd - portBegin + 1,
						Txt:       txt,
					}, nil
				} else {
					return nil, fmt.Errorf("Err %s", txt)
				}
			} else {
				return nil, fmt.Errorf("%s 不是有效个整数格式", ports[1])
			}
		} else {
			return nil, fmt.Errorf("%s 不是有效个整数格式", ports[0])
		}
	} else {
		return nil, fmt.Errorf("格式错误[%s] 例如(200-300;400)", txt)
	}
}

func (self *SubnetPortList) Add(item *SubnetPort) {
	self.List = append(self.List, item)
}

func (self *SubnetPortList) CheckRepettion(item *SubnetPort) error {
	for _, v := range self.List {
		if item.PortBegin >= v.PortBegin && item.PortBegin <= v.PortEnd {
			return fmt.Errorf("Repetition %s, have %s", v.Txt, item.Txt)
		}
		if item.PortEnd >= v.PortBegin && item.PortEnd <= v.PortEnd {
			return fmt.Errorf("Repetition %s, have %s", v.Txt, item.Txt)
		}

	}
	return nil
}

func (self *SubnetPortList) GetPortCount() int {
	ret := 0
	for _, v := range self.List {
		ret += v.PortCount
	}
	return ret
}

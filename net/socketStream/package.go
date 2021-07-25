package socketStream

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type TPackage struct {
	IsSum    bool   //校验码要求
	Tag      uint16 //标签验证 0xffff
	Op       uint16 //指令码
	Seq      uint32 //4字节附加数据
	BodySize uint16
	Body     []byte
}

func NewTPackage() *TPackage {
	return &TPackage{
		Tag: 0xFFFF,
	}
}

func (s *TPackage) ParseBuf(buf []byte) error {
	len := len(buf)
	if len < 10 {
		return fmt.Errorf("Err bufLen(%d) < 10", len)
	} else {
		buff := bytes.NewBuffer(buf)
		binary.Read(buff, binary.LittleEndian, &s.Tag)
		binary.Read(buff, binary.LittleEndian, &s.Op)
		binary.Read(buff, binary.LittleEndian, &s.Seq)
		binary.Read(buff, binary.LittleEndian, &s.BodySize)
		if int(s.BodySize) <= len-10 {
			s.Body = make([]byte, s.BodySize)
			if err := binary.Read(buff, binary.LittleEndian, &s.Body); err == nil {
				return nil
			} else {
				return fmt.Errorf("Err %s\n", err.Error())
			}
		}
	}
	return fmt.Errorf("xx")
}

func (s *TPackage) tobuf() []byte {
	buff := make([]byte, 10+int(s.BodySize), getSliceMaxLen(10+int(s.BodySize)))
	binary.LittleEndian.PutUint16(buff[0:2], s.Tag)
	binary.LittleEndian.PutUint16(buff[2:4], s.Op)
	binary.LittleEndian.PutUint32(buff[4:8], s.Seq)
	binary.LittleEndian.PutUint16(buff[8:10], s.BodySize)
	if s.BodySize > 0 {
		copy(buff[10:], s.Body)
	}
	return buff
}

func (s *TPackage) ToBuff() []byte {
	buff := s.tobuf()
	if s.IsSum {
		sum := s.GetSum(buff[2:])
		binary.LittleEndian.PutUint16(buff[0:2], uint16(sum&0xffff))
	}
	return buff
}

//获取数据校验码
func (s *TPackage) GetSum(b []byte) uint32 {
	if b == nil {
		b = s.tobuf()[2:]
	}
	var sum uint32 = 0
	len1 := len(b) / 2
	len2 := len(b) % 2
	for i := 0; i < len1; i++ {
		sum += uint32(binary.LittleEndian.Uint16(b[i*2 : i*2+2]))
	}
	if len2 > 0 {
		sum += uint32(b[len1*2])
	}
	return sum
}

func NewTPackageFromBuf(buf []byte) (*TPackage, error) {
	len := len(buf)
	if len < 10 {
		return nil, fmt.Errorf("Err bufLen(%d) < 10", len)
	} else {
		buff := bytes.NewBuffer(buf)
		pkg := new(TPackage)
		binary.Read(buff, binary.LittleEndian, &pkg.Tag)
		binary.Read(buff, binary.LittleEndian, &pkg.Op)
		binary.Read(buff, binary.LittleEndian, &pkg.Seq)
		binary.Read(buff, binary.LittleEndian, &pkg.BodySize)
		if int(pkg.BodySize) <= len-10 {
			pkg.Body = make([]byte, pkg.BodySize)
			if err := binary.Read(buff, binary.LittleEndian, &pkg.Body); err == nil {
				return pkg, nil
			} else {
				return nil, fmt.Errorf("Err %s\n", err.Error())
			}
		}
	}
	return nil, fmt.Errorf("xx")
}

func getSliceMaxLen(l int) int {
	ret := 1
	for ret < l {
		ret *= 2
	}
	return ret
}

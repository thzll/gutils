package socketStream

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

const (
	READ_TIMEOUT  = 0
	WRITE_TIMEOUT = 1
)

func Readn(conn net.Conn, buf []byte) error {
	length := len(buf)
	hot := buf
	for {
		if length <= 0 {
			return nil
		}
		n, err := conn.Read(hot)
		if err != nil {
			return err
		}
		length -= n
		hot = hot[n:]
	}
	return nil
}
func Writen(conn net.Conn, buf []byte) error {
	length := len(buf)
	hot := buf
	for {
		if length <= 0 {
			return nil
		}
		n, err := conn.Write(hot)
		if err != nil {
			return err
		}
		length -= n
		hot = hot[n:]
	}
	return nil
}

func CheckBodySize(bodySize int) error {
	if bodySize > PACKAGE_MAX_BODY_SIZE {
		return fmt.Errorf("bodySize is too big[%d] maxBodySize [%d]", bodySize, PACKAGE_MAX_BODY_SIZE)
	} else {
		return nil
	}
}

/***********************************
t：超时时间
dir：0 设置写超时，1，设置读超时
***********************************/
func SetSocketTimeout(conn net.Conn, t int64, dir int) (err error) {
	if dir == WRITE_TIMEOUT {
		err = conn.SetWriteDeadline(time.Now().Add(time.Duration(t) * time.Millisecond))
	} else {
		err = conn.SetReadDeadline(time.Now().Add(time.Duration(t) * time.Millisecond))
	}
	return err
}

func ReadPackage(conn net.Conn, timeout int) (*TPackage, error) {
	hdr := new(TPackage)
	rbuff := make([]byte, 10)
	rbuff[0] = 0xff
	SetSocketTimeout(conn, int64(timeout), READ_TIMEOUT)
	err := Readn(conn, rbuff)
	if err != nil {
		return nil, err
	}
	buff := bytes.NewBuffer(rbuff)
	binary.Read(buff, binary.LittleEndian, &hdr.Tag)
	binary.Read(buff, binary.LittleEndian, &hdr.Op)
	binary.Read(buff, binary.LittleEndian, &hdr.Seq)
	binary.Read(buff, binary.LittleEndian, &hdr.BodySize)
	//fmt.Println("Err rbuff",rbuff)
	if hdr.Tag != PACKAGE_TAG {
		fmt.Println("Err rbuff", rbuff)
		return nil, fmt.Errorf("Tag[%04x] is Not [%x] %s", hdr.Tag, PACKAGE_TAG, conn.RemoteAddr())
	}
	if err := CheckBodySize(int(hdr.BodySize)); err != nil {
		return nil, err
	}

	if hdr.BodySize == 0 {
		hdr.Body = nil
		return hdr, nil
	}
	SetSocketTimeout(conn, int64(timeout), READ_TIMEOUT)
	body := make([]byte, hdr.BodySize)
	err = Readn(conn, body)
	hdr.Body = body
	return hdr, err
}

func WritePackage(conn net.Conn, hdr *TPackage) (int, error) {
	hdr.Tag = 0xFFFF
	buff := bytes.NewBuffer(nil)
	binary.Write(buff, binary.LittleEndian, hdr.Tag)
	binary.Write(buff, binary.LittleEndian, hdr.Op)
	binary.Write(buff, binary.LittleEndian, hdr.Seq)
	binary.Write(buff, binary.LittleEndian, hdr.BodySize)
	if hdr.BodySize > 0 {
		binary.Write(buff, binary.BigEndian, hdr.Body)
	}
	return conn.Write(buff.Bytes())
}

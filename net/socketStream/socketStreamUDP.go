package socketStream

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net"
	"time"
)

type TSocketUDP struct {
	DestIP  string
	localIP string
	port    string
	conn    net.Conn
	recvBuf []byte
}

func (s *TSocketUDP) RemoteAddrInfo() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *TSocketUDP) LocalAddrInfo() net.Addr {
	return s.conn.LocalAddr()
}

func (s *TSocketUDP) SetSocketReadDeadline(timeout int) error {
	return SetSocketTimeout(s.conn, int64(timeout), READ_TIMEOUT)
}
func (s *TSocketUDP) SetSocketWriteDeadline(timeout int) error {
	return SetSocketTimeout(s.conn, int64(timeout), WRITE_TIMEOUT)
}
func (s *TSocketUDP) EncryptionSocketChannel() error {
	tls_con, err := createTlsConn(s.conn)
	if err != nil {
		s.conn.Close()
		return fmt.Errorf("Error:failed to encryp socket")
	}
	s.conn = tls_con
	return nil
}
func (s *TSocketUDP) Close() {
	if s != nil && s.conn != nil {
		s.conn.Close()
	}
}

func (s *TSocketUDP) Read(buf []byte, timeout int) (n int, err error) {
	var tDelay int
	tDelay = timeout
	if tDelay == 0 {
		tDelay = 60
	}
	s.SetSocketReadDeadline(tDelay)
	return s.conn.Read(buf)
}

func (s *TSocketUDP) Write(buf []byte, timeout int) (n int, err error) {
	var tDelay int
	tDelay = timeout
	if tDelay == 0 {
		tDelay = 60
	}
	s.SetSocketWriteDeadline(tDelay)
	return s.conn.Write(buf)
}

func (s *TSocketUDP) WritePackage(pkg *TPackage, timeout int) (int, error) {
	var tDelay int
	tDelay = timeout
	if tDelay == 0 {
		tDelay = 60
	}
	s.SetSocketWriteDeadline(tDelay)
	n, err := WritePackage(s.conn, pkg)
	return n, err
}
func (s *TSocketUDP) ReadPackage(timeout int) (*TPackage, error) {
	var tDelay int
	tDelay = timeout
	if tDelay == 0 {
		tDelay = 1000
	}
	hdr := new(TPackage)
	rbuff := make([]byte, 1500)
	rbuff[0] = 0xff
	SetSocketTimeout(s.conn, int64(timeout), READ_TIMEOUT)
	//n, raddr, err := net.UDPConn(s.conn).ReadFromUDP(rbuff[0:])
	n, err := s.conn.Read(rbuff)
	if err != nil {
		return nil, err
	}
	buff := bytes.NewBuffer(rbuff[:n])
	if n < 10 {
		return nil, fmt.Errorf("readLen:%d < 10", n)
	}
	binary.Read(buff, binary.LittleEndian, &hdr.Tag)
	binary.Read(buff, binary.LittleEndian, &hdr.Op)
	binary.Read(buff, binary.LittleEndian, &hdr.Seq)
	binary.Read(buff, binary.LittleEndian, &hdr.BodySize)
	if hdr.Tag != PACKAGE_TAG {
		fmt.Println("Err rbuff", rbuff)
		return nil, fmt.Errorf("Tag[%04x] is Not [%x] %s", hdr.Tag, PACKAGE_TAG, s.conn.RemoteAddr())
	}
	if err := CheckBodySize(int(hdr.BodySize)); err != nil {
		return nil, err
	}
	if hdr.BodySize == 0 {
		hdr.Body = nil
		return hdr, nil
	} else {
		if n >= 10+int(hdr.BodySize) {
			hdr.Body = rbuff[10 : 10+hdr.BodySize]
			return hdr, nil
		} else {
			return nil, fmt.Errorf("readLen:%d < 10+bodylen:%d", n, hdr.BodySize)
		}

	}
}

func (s *TSocketUDP) GetConn() net.Conn {
	return s.conn
}

func CreateSocketUDP(destIp, destPort, localIP, localPort string) (*TSocketUDP, error) {
	destaddr := fmt.Sprintf("%s:%s", destIp, destPort)
	if conn, err := net.Dial("udp", destaddr); err == nil {
		conn.SetDeadline(time.Now().Add(20 * time.Second))
		masterChannelSocket := new(TSocketUDP)
		masterChannelSocket.conn = conn
		masterChannelSocket.DestIP = destIp
		masterChannelSocket.localIP = localIP
		masterChannelSocket.port = destPort
		return masterChannelSocket, nil
	} else {
		return nil, err
	}
	//return CreateSocketUDP_(fmt.Sprintf("%s:%s", destIp, destPort), localIP, localPort)
}

func CreateSocketUDP_(destaddr, localIP, localPort string) (*TSocketUDP, error) {
	masterChannelSocket := new(TSocketUDP)
	toAddr := destaddr
	destAddr, err := net.ResolveUDPAddr("udp", toAddr)
	if err != nil || destAddr == nil {
		return nil, err
	}

	fromAddr := localIP + ":" + localPort
	loaclAddr, err := net.ResolveUDPAddr("udp", fromAddr)
	if err != nil || loaclAddr == nil {
		return nil, err
	}
	con, err := net.DialUDP("udp", loaclAddr, destAddr)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Error:failed to connect dest IP,%s", err.Error()))
		return nil, err
	}

	con.SetDeadline(time.Now().Add(20 * time.Second))
	destIp, destPort, _ := net.SplitHostPort(destaddr)
	masterChannelSocket.conn = con
	masterChannelSocket.DestIP = destIp
	masterChannelSocket.localIP = localIP
	masterChannelSocket.port = destPort
	return masterChannelSocket, nil

}

func CreateSocketUDPex(destIp, destPort string, timeout int) (*TSocketUDP, error) {
	toAddr := destIp + ":" + destPort
	masterChannelSocket := new(TSocketUDP)
	con, err := net.DialTimeout("udp", toAddr, time.Millisecond*time.Duration(timeout))
	if err != nil {
		log.Println("Error:failed to connect dest IP,", err.Error())
		return nil, err
	}

	con.SetDeadline(time.Now().Add(20 * time.Second))
	masterChannelSocket.conn = con
	return masterChannelSocket, nil

}

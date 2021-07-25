package socketStream

import (
	"fmt"
	"go.uber.org/zap"
	"log"
	"net"
	"time"
)

type TSocket struct {
	DestIP    string
	localIP   string
	port      string
	conn      net.Conn
	recvBuf   []byte
	IsSkipSum bool
}

func (this *TSocket) RemoteAddrInfo() net.Addr {
	return this.conn.RemoteAddr()
}

func (this *TSocket) LocalAddrInfo() net.Addr {
	return this.conn.LocalAddr()
}

func (this *TSocket) SetSocketReadDeadline(timeout int) error {
	return SetSocketTimeout(this.conn, int64(timeout), READ_TIMEOUT)
}
func (this *TSocket) SetSocketWriteDeadline(timeout int) error {
	return SetSocketTimeout(this.conn, int64(timeout), WRITE_TIMEOUT)
}
func (this *TSocket) EncryptionSocketChannel() error {
	tls_con, err := createTlsConn(this.conn)
	if err != nil {
		this.conn.Close()
		return fmt.Errorf("Error:failed to encryp socket")
	}
	this.conn = tls_con
	return nil
}
func (this *TSocket) Close() {
	if this != nil && this.conn != nil {
		this.conn.Close()
	}
}

func (this *TSocket) Read(buf []byte, timeout int) (n int, err error) {
	var tDelay int
	tDelay = timeout
	if tDelay == 0 {
		tDelay = 60
	}
	this.SetSocketReadDeadline(tDelay)
	return this.conn.Read(buf)
}

func (this *TSocket) Write(buf []byte, timeout int) (n int, err error) {
	var tDelay int
	tDelay = timeout
	if tDelay == 0 {
		tDelay = 60
	}
	this.SetSocketWriteDeadline(tDelay)
	return this.conn.Write(buf)
}

func (this *TSocket) WritePackage(pkg *TPackage, timeout int) (int, error) {
	var tDelay int
	tDelay = timeout
	if tDelay == 0 {
		tDelay = 60
	}
	this.SetSocketWriteDeadline(tDelay)
	if !this.IsSkipSum {
		pkg.IsSum = true
	}
	n, err := WritePackage(this.conn, pkg)
	return n, err
}
func (this *TSocket) ReadPackage(timeout int) (*TPackage, error) {
	var tDelay int
	tDelay = timeout
	if tDelay == 0 {
		tDelay = 1000
	}
	return ReadPackage(this.conn, tDelay)
}

func (this *TSocket) GetConn() net.Conn {
	return this.conn
}

func InitSocket(conn net.Conn) *TSocket {
	c := new(TSocket)
	c.conn = conn
	return c
}
func CreateSocket(destIp, destPort, localIP, localPort string) (*TSocket, error) {
	return CreateSocket_(fmt.Sprintf("%s:%s", destIp, destPort), localIP, localPort)
}

func CreateSocket_(destaddr, localIP, localPort string) (*TSocket, error) {
	masterChannelSocket := new(TSocket)

	toAddr := destaddr
	destAddr, err := net.ResolveTCPAddr("tcp", toAddr)
	if err != nil || destAddr == nil {
		return nil, err
	}

	fromAddr := localIP + ":" + localPort
	loaclAddr, err := net.ResolveTCPAddr("tcp", fromAddr)
	if err != nil || loaclAddr == nil {
		return nil, err
	}
	con, err := net.DialTCP("tcp", loaclAddr, destAddr)
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

func CreateSocketex(toAddr string, timeout int) (*TSocket, error) {
	masterChannelSocket := new(TSocket)
	con, err := net.DialTimeout("tcp", toAddr, time.Millisecond*time.Duration(timeout))
	if err != nil {
		log.Println("Error:failed to connect dest IP,", err.Error())
		return nil, err
	}

	con.SetDeadline(time.Now().Add(20 * time.Second))
	masterChannelSocket.conn = con
	return masterChannelSocket, nil

}

func CreateSocketTimeout(destIp, destPort string, timeout int) (*TSocket, error) {
	masterChannelSocket := new(TSocket)

	toAddr := destIp + ":" + destPort

	con, err := net.DialTimeout("tcp", toAddr, time.Duration(timeout)*time.Second)
	if err != nil {
		log.Println("Error:failed to connect dest IP,", err.Error())
		return nil, err
	}

	con.SetDeadline(time.Now().Add(20 * time.Second))
	masterChannelSocket.conn = con
	masterChannelSocket.DestIP = destIp
	masterChannelSocket.port = destPort
	return masterChannelSocket, nil

}

package socketStream

import (
	"fmt"
	"github.com/thzll/gutils/myutils"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type InSocket interface {
	OnGoReader([]byte)
	OnTimer()
	OnExitReader()
	OnDisConnected()
}

type SocketHandler struct {
	socket          *TSocket //用于连接代理
	isRunding       bool     //DataHaldler 防止次函数只被调用一次
	writerTimerTime int      //定时器
	socketTimeOut   int      //网络接口数据超时时间
	lock_           sync.Mutex
}

func (self *SocketHandler) DataHaldler(obj InSocket) {
	self.lock_.Lock()
	if self.isRunding {
		self.lock_.Unlock()
		return
	} else {
		self.isRunding = true
	}
	self.lock_.Unlock()
	self.sGoReader(obj)
	self.closeConnAndNil()
}

func (self *SocketHandler) sGoReader(obj InSocket) {
	socket := self.GetSocket()
	//var t *time.Timer = nil
	lastOnTimer := int64(0)
	if self.writerTimerTime > 0 {
		//t = time.NewTimer(time.Millisecond * time.Duration(self.writerTimerTime))
		lastOnTimer = int64(time.Now().UnixNano() + int64(time.Millisecond*time.Duration(self.writerTimerTime)))
	}
	buf := make([]byte, 8096)
	for {
		if lastOnTimer > 0 {
			if cTime := time.Now().UnixNano(); cTime > lastOnTimer {
				lastOnTimer = cTime + int64(time.Millisecond*time.Duration(self.writerTimerTime))
				obj.OnTimer()
			}
		}
		//if t != nil {
		//	select {
		//	case <-t.C:
		//		t.Reset(time.Millisecond * time.Duration(self.writerTimerTime))
		//		obj.OnTimer()
		//	default:
		//	}
		//}
		timeout := myutils.GetMinInt(self.socketTimeOut, self.writerTimerTime)
		if n, err := socket.Read(buf, timeout); err == nil {
			obj.OnGoReader(buf[:n])
		} else {
			if ne, ok := err.(net.Error); ok {
				if ne.Timeout() {
					continue
				} else {
					log.Println("Error: sGoReader  ne.Error()", ne.Error())
					if strings.Index(ne.Error(), "use of closed network connection") >= 0 {
						obj.OnDisConnected()
					}
					break
				}
			} else if err.Error() == "EOF" ||
				strings.Index(err.Error(), "use of closed network connection") >= 0 {
				obj.OnDisConnected()
				break
			} else {
				log.Println("Error: sGoReader  err", err.Error())
				break
			}
		}
	}
	self.CloseConn()
	obj.OnExitReader()
}

func (self *SocketHandler) SetTimer(time int) {
	self.writerTimerTime = time
}

func (self *SocketHandler) InitData() {
	self.socketTimeOut = 10

}

func (self *SocketHandler) InitSocket(conn net.Conn) error {
	self.closeConnAndNil()
	self.lock_.Lock()
	defer self.lock_.Unlock()
	self.socket = InitSocket(conn)
	self.InitData()
	return nil
}

func (self *SocketHandler) Connetct(destIp string, destPort int, localIP, localPort string) error {
	self.closeConnAndNil()
	self.lock_.Lock()
	defer self.lock_.Unlock()
	if socket, err := CreateSocket(destIp, fmt.Sprintf("%d", destPort), localIP, localIP); err == nil {
		self.socket = socket
		self.InitData()
		return nil
	} else {
		log.Println("Err SocketHandler.Connetct", err)
		return err
	}
}

func (self *SocketHandler) ConnetctTimeOut(destIp string, destPort int, timeout int) error {
	self.closeConnAndNil()
	self.lock_.Lock()
	defer self.lock_.Unlock()
	if socket, err := CreateSocketex(fmt.Sprintf("%s:%d", destIp, destPort), timeout); err == nil {
		self.socket = socket
		self.InitData()
		return nil
	} else {
		log.Println("Err SocketHandler.Connetct", err)
		return err
	}
}

func (self *SocketHandler) CloseConn() {
	self.lock_.Lock()
	defer self.lock_.Unlock()
	if self.socket != nil {
		self.socket.Close()
	}
}
func (self *SocketHandler) closeConnAndNil() {
	self.lock_.Lock()
	defer self.lock_.Unlock()
	if self.socket != nil {
		self.socket.Close()
		self.socket = nil
	}
}

func (self *SocketHandler) Send(b []byte) error {
	self.lock_.Lock()
	defer self.lock_.Unlock()
	if self.socket != nil {
		if _, err := self.socket.Write(b, self.socketTimeOut); err != nil {
			return err
			if err == io.EOF {
				log.Println("Err SocketHandler.send err:", err)
				return err
			}
			if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			}
		} else {
			//log.Println("write: ", n)
			return nil
		}
	} else {
		return fmt.Errorf("Err SocketHandler.socket is nil")
	}
	return nil
}

func (self *SocketHandler) GetSocket() *TSocket {
	self.lock_.Lock()
	defer self.lock_.Unlock()
	return self.socket
}

func (self *SocketHandler) EncryptSocket() error {
	return self.socket.EncryptionSocketChannel()
}

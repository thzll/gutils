package socketStream

import (
	"fmt"
	"go.uber.org/zap"
	"github.com/thzll/gutils/myutils"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type InSocketPkg interface {
	OnGoReader(pkg *TPackage)
	OnTimer()
	OnExitReader()
	OnError(err error)
}

type SocketHandlerPkg struct {
	socket          *TSocket //用于连接代理
	isRunding       bool     //DataHaldler 防止次函数只被调用一次
	writerTimerTime int      //定时器单位毫秒
	socketTimeOut   int      //网络接口数据超时时间
	lock_           sync.Mutex
}

func (self *SocketHandlerPkg) DataHaldlerPkg(obj InSocketPkg) {
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
	self.isRunding = false
}

func (self *SocketHandlerPkg) sGoReader(obj InSocketPkg) {
	socket := self.GetSocket()
	//var t *time.Timer = nil
	lastOnTimer := int64(0)
	if self.writerTimerTime > 0 {
		//t = time.NewTimer(time.Second * time.Duration(self.writerTimerTime))
		lastOnTimer = int64(time.Now().UnixNano() + int64(time.Millisecond*time.Duration(self.writerTimerTime)))
	}
	for {
		//if t != nil {
		//	select {
		//	case <-t.C:
		//		t.Reset(time.Second * time.Duration(self.writerTimerTime))
		//		obj.OnTimer()
		//	default:
		//	}
		//}
		if lastOnTimer > 0 {
			if cTime := time.Now().UnixNano(); cTime > lastOnTimer {
				lastOnTimer = cTime + int64(time.Millisecond*time.Duration(self.writerTimerTime))
				obj.OnTimer()
			}
		}
		timeout := myutils.GetMinInt(self.socketTimeOut, self.writerTimerTime)
		if pkg, err := socket.ReadPackage(timeout); err == nil {
			obj.OnGoReader(pkg)
		} else {
			if pkg != nil {
				obj.OnError(fmt.Errorf("SocketHandlerPkg.sGoReader pkg:%+v err:%s", pkg, err.Error()))
			}
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				//log.Println("Error:SocketHandlerPkg.sGoReader err1 ", err)
				continue
			} else {
				obj.OnError(fmt.Errorf("SocketHandlerPkg.sGoReader err:%s", err.Error()))
				break
			}
		}
	}
	self.CloseConn()
	obj.OnExitReader() //结束
}

func (self *SocketHandlerPkg) SetTimer(time int) {
	self.writerTimerTime = time
}

func (self *SocketHandlerPkg) InitData() {
	self.socketTimeOut = 200

}

func (self *SocketHandlerPkg) InitSocket(conn net.Conn) error {
	self.closeConnAndNil()
	self.lock_.Lock()
	defer self.lock_.Unlock()
	self.socket = InitSocket(conn)
	self.InitData()
	return nil
}
func (self *SocketHandlerPkg) Connetct(destIp string, destPort int, localIP, localPort string) error {
	self.closeConnAndNil()
	self.lock_.Lock()
	defer self.lock_.Unlock()
	if socket, err := CreateSocket(destIp, fmt.Sprintf("%d", destPort), localIP, localIP); err == nil {
		self.socket = socket
		self.InitData()
		return nil
	} else {
		zap.L().Error("Err SocketHandler.Connetct", zap.String("err", err.Error()))
		return err
	}
}

func (self *SocketHandlerPkg) Connetct_(destAddr string, localIP, localPort string) error {
	self.closeConnAndNil()
	self.lock_.Lock()
	defer self.lock_.Unlock()
	if socket, err := CreateSocket_(destAddr, localIP, localIP); err == nil {
		self.socket = socket
		self.InitData()
		return nil
	} else {
		zap.L().Error("Err SocketHandler.Connetct", zap.String("err", err.Error()))
		return err
	}
}

// timeOut millsecond
func (self *SocketHandlerPkg) ConnetctEx(destAddr string, timeout int) error {
	self.closeConnAndNil()
	self.lock_.Lock()
	defer self.lock_.Unlock()
	if socket, err := CreateSocketex(destAddr, timeout); err == nil {
		self.socket = socket
		self.InitData()
		return nil
	} else {
		zap.L().Error("Err SocketHandler.Connetct", zap.String("err", err.Error()))
		return err
	}
}

func (self *SocketHandlerPkg) CloseConn() {
	self.lock_.Lock()
	defer self.lock_.Unlock()
	if self.socket != nil {
		self.socket.Close()
	}
}
func (self *SocketHandlerPkg) closeConnAndNil() {
	self.lock_.Lock()
	defer self.lock_.Unlock()
	if self.socket != nil {
		self.socket.Close()
		self.socket = nil
	}
}

func (self *SocketHandlerPkg) SendPackage(pkg *TPackage) error {
	self.lock_.Lock()
	defer self.lock_.Unlock()
	ctime := time.Now().UnixNano()
	if self.socket != nil {
		if _, err := self.socket.WritePackage(pkg, self.socketTimeOut); err != nil {
			if err == io.EOF {
				log.Println("Err SocketHandlerPkg.SendPackage err:", err)
				return err
			}
			if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
				if time.Now().UnixNano()-ctime >= 2000*int64(time.Millisecond) {
					return err
				}
			}
		} else {
			//log.Println("write: ", n)
			return nil
		}
	} else {
		return fmt.Errorf("Err SocketHandlerPkg.socket is nil")
	}
	return nil
}

func (self *SocketHandlerPkg) GetSocket() *TSocket {
	self.lock_.Lock()
	defer self.lock_.Unlock()
	return self.socket
}

func (self *SocketHandlerPkg) EncryptSocket() error {
	return self.socket.EncryptionSocketChannel()
}

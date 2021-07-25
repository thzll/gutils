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

type SocketUDPHandlerPkg struct {
	socket          *TSocketUDP //用于连接代理
	isRunding       bool        //DataHaldler 防止次函数只被调用一次
	writerTimerTime int         //定时器单位毫秒
	socketTimeOut   int         //网络接口数据超时时间
	lock_           sync.Mutex
}

func (s *SocketUDPHandlerPkg) DataHaldlerPkg(obj InSocketPkg) {
	s.lock_.Lock()
	if s.isRunding {
		s.lock_.Unlock()
		return
	} else {
		s.isRunding = true
	}
	s.lock_.Unlock()
	s.sGoReader(obj)
	s.closeConnAndNil()
	s.isRunding = false
}

func (s *SocketUDPHandlerPkg) sGoReader(obj InSocketPkg) {
	socket := s.GetSocket()
	//var t *time.Timer = nil
	lastOnTimer := int64(0)
	if s.writerTimerTime > 0 {
		//t = time.NewTimer(time.Second * time.Duration(s.writerTimerTime))
		lastOnTimer = int64(time.Now().UnixNano() + int64(time.Millisecond*time.Duration(s.writerTimerTime)))
	}
	for {
		//if t != nil {
		//	select {
		//	case <-t.C:
		//		t.Reset(time.Second * time.Duration(s.writerTimerTime))
		//		obj.OnTimer()
		//	default:
		//	}
		//}
		if lastOnTimer > 0 {
			if cTime := time.Now().UnixNano(); cTime > lastOnTimer {
				lastOnTimer = cTime + int64(time.Millisecond*time.Duration(s.writerTimerTime))
				obj.OnTimer()
			}
		}
		timeout := myutils.GetMinInt(s.socketTimeOut, s.writerTimerTime)
		if timeout == 0 {
			timeout = s.socketTimeOut
		}
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
	s.CloseConn()
	obj.OnExitReader() //结束
}

func (s *SocketUDPHandlerPkg) SetTimer(time int) {
	s.writerTimerTime = time
}

func (s *SocketUDPHandlerPkg) InitData() {
	s.socketTimeOut = 200

}

func (s *SocketUDPHandlerPkg) Connetct(destIp string, destPort int, localIP, localPort string) error {
	s.closeConnAndNil()
	s.lock_.Lock()
	defer s.lock_.Unlock()
	if socket, err := CreateSocketUDP(destIp, fmt.Sprintf("%d", destPort), localIP, localIP); err == nil {
		s.socket = socket
		s.InitData()
		return nil
	} else {
		zap.L().Error("Err SocketUDPHandler.Connetct", zap.String("err", err.Error()))
		return err
	}
}

func (s *SocketUDPHandlerPkg) Connetct_(destAddr string, localIP, localPort string) error {
	s.closeConnAndNil()
	s.lock_.Lock()
	defer s.lock_.Unlock()
	if socket, err := CreateSocketUDP_(destAddr, localIP, localIP); err == nil {
		s.socket = socket
		s.InitData()
		return nil
	} else {
		zap.L().Error("Err SocketHandler.Connetct", zap.String("err", err.Error()))
		return err
	}
}

func (s *SocketUDPHandlerPkg) CloseConn() {
	s.lock_.Lock()
	defer s.lock_.Unlock()
	log.Println("close=======================")
	if s.socket != nil {
		s.socket.Close()
	}
}
func (s *SocketUDPHandlerPkg) closeConnAndNil() {
	s.lock_.Lock()
	defer s.lock_.Unlock()
	if s.socket != nil {
		s.socket.Close()
		s.socket = nil
	}
}

func (s *SocketUDPHandlerPkg) SendPackage(pkg *TPackage) error {
	s.lock_.Lock()
	defer s.lock_.Unlock()
	ctime := time.Now().UnixNano()
	if s.socket != nil {
		if _, err := s.socket.WritePackage(pkg, s.socketTimeOut); err != nil {
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

func (S *SocketUDPHandlerPkg) GetSocket() *TSocketUDP {
	S.lock_.Lock()
	defer S.lock_.Unlock()
	return S.socket
}

func (S *SocketUDPHandlerPkg) EncryptSocket() error {
	return S.socket.EncryptionSocketChannel()
}

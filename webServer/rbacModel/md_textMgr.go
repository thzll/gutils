package rbacModel

import (
	"go.uber.org/zap"
	"log"
	"sync"
	"time"
)

type TextMgr struct {
	list           map[int64]*Text
	lock           sync.Mutex
	IsRun          bool
	unixupdatetime int64
	dbupdatetime   int64 //数据库更新数据时间 单位秒
	active         Active
}

var (
	gTextMgr *TextMgr = &TextMgr{list: make(map[int64]*Text)}
)

func (self *TextMgr) Start() {
	self.lock.Lock()
	defer self.lock.Unlock()
	if !self.IsRun {
		go self.Run()
		self.IsRun = true
	}
}

func (self *TextMgr) Run() {
	self.active.setActive()
	for {
		if self.active.checkActive() {
			self.loadDataFromDb(self.dbupdatetime)
		}
		time.Sleep(time.Second * 5)
	}
}

func (self *TextMgr) loadDataFromDb(utime int64) {
	list := GetTextListByUpdateTime(utime)
	rUpdatetime := int64(0)
	self.lock.Lock()
	defer self.lock.Unlock()
	if len(list) > 0 {
		for _, vn := range list {
			updatetime := vn.UpdatedAt.Unix()
			if updatetime > rUpdatetime {
				rUpdatetime = updatetime
			}
			self.list[vn.Id] = vn
			zap.L().Info("Info: TextMgr.loadDataFromDb Update ",
				zap.String("name", vn.Name),
				zap.String("updateAt", vn.UpdatedAt.String()),
			)
		}
		self.dbupdatetime = rUpdatetime
	}
	self.unixupdatetime = time.Now().UnixNano()
}

func (self *TextMgr) CheckUpdateData() {
	if len(self.list) == 0 {
		self.loadDataFromDb(0)
	}
	self.active.setActive()
	self.Start()
}

func (self *TextMgr) GetText(nodeName, name string) *Text {
	log.Println("GetText", name, nodeName)
	self.CheckUpdateData()
	for _, v := range self.list {
		if v.Name == name && v.NodeName == nodeName {
			return v
		}
	}
	return nil
}

//get node list
func GetTextByName(nodeName, name string) string {
	text := gTextMgr.GetText(nodeName, name)
	if text != nil {
		return text.Text
	}
	return ""
}

package rbacModel

import (
	"log"
	"sort"
	"sync"
	"time"
)

type NodeMgr struct {
	nodes          map[int64]*Node
	lock           sync.Mutex
	IsRun          bool
	unixupdatetime int64
}

var (
	gNodeMgr *NodeMgr
)

type NodeSlice []*Node
type NodeSliceById []*Node

func (s NodeSlice) Len() int           { return len(s) }
func (s NodeSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s NodeSlice) Less(i, j int) bool { return s[i].Sort < s[j].Sort }

func (s NodeSliceById) Len() int           { return len(s) }
func (s NodeSliceById) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s NodeSliceById) Less(i, j int) bool { return s[i].Id < s[j].Id }

func NewNodeMgr() *NodeMgr {
	return &NodeMgr{nodes: make(map[int64]*Node, 0)}
}

func (self *NodeMgr) Run() {
	self.IsRun = true
	for {
		updatetime, _ := self.loadDataFromDb(self.unixupdatetime)
		if updatetime > self.unixupdatetime {
			self.unixupdatetime = updatetime
		}
		time.Sleep(5 * time.Second)
		break
	}
	self.IsRun = false
}

func (self *NodeMgr) loadDataFromDb(unixupdatetime int64) (int64, bool) {
	var bRs bool = false
	var rUpdatetime int64 = 0
	list := GetNodeListByUpdateTime(unixupdatetime)
	self.lock.Lock()
	defer self.lock.Unlock()

	for _, vn := range list {
		updatetime := vn.UpdatedAt.Unix()
		if updatetime > rUpdatetime {
			rUpdatetime = updatetime
		}
		self.nodes[vn.Id] = vn
		log.Printf("Info: NodeMgr.loadDataFromDb Update groupName:%s Id:%d\n", vn.Name, vn.Id)
		bRs = true
	}
	return rUpdatetime, bRs
}

func (self *NodeMgr) GetNodeByName(name string) *Node {
	self.lock.Lock()
	defer self.lock.Unlock()

	for _, vn := range self.nodes {
		if vn.Name == name {
			return vn
		}
	}
	return nil
}

func InitCheck() {
	if gNodeMgr == nil {
		gNodeMgr = NewNodeMgr()
		gNodeMgr.loadDataFromDb(0)
		go gNodeMgr.Run()
	}
}
func GetWebNodelist() (list []*Node) {
	log.Println("GetGrouplist")
	InitCheck()
	for _, v := range gNodeMgr.nodes {
		elem := v
		list = append(list, elem)
		log.Println("WebNode", elem)
	}
	sort.Sort(NodeSliceById(list))
	return list
}

func GetWebNodeById(id int64) *Node {
	log.Println("GetWebNodeById ", id)
	InitCheck()
	return gNodeMgr.nodes[id]
}

func GetWebNodeByName(name string) *Node {
	log.Println("GetWebNodeByName ", name)
	InitCheck()
	return gNodeMgr.GetNodeByName(name)
}

func GetNodeTree(pid int64, level int64) (list []*Node) {
	InitCheck()
	for _, v := range gNodeMgr.nodes {
		elem := v
		if elem.Pid == pid && elem.Level == level && elem.Status == 2 {
			list = append(list, elem)
			//log.Println("GetNodeTree Group", elem)
		}
	}
	sort.Sort(NodeSlice(list))
	return list
}

func GetNodeTreeByUserId(pid int64, level int64, userid int64) []*Node {
	if userid == 1 { //admin超级管理员
		return GetNodeTree(pid, level)
	}
	InitCheck()
	list := make([]*Node, 0)
	for _, v := range gNodeMgr.nodes {
		elem := v
		if elem.Pid == pid && elem.Level == level && elem.Status == 2 {
			list = append(list, elem)
			//log.Println("GetNodeTreeByUserId Group", elem)
		}
	}
	sort.Sort(NodeSlice(list))
	return list
}

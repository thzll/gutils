package rbacModel

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type NodeAccessRuleMgr struct {
	elems          map[int64]*NodeAccessRule
	lock           sync.Mutex
	IsRun          bool
	unixupdatetime int64
}

var (
	gNodeAccessRuleMgr *NodeAccessRuleMgr
)

func NewNodeAccessRuleMgr() *NodeAccessRuleMgr {
	return &NodeAccessRuleMgr{elems: make(map[int64]*NodeAccessRule, 0)}
}

func InitCheckNodeAccessRuleMgr() {
	if gNodeAccessRuleMgr == nil {
		gNodeAccessRuleMgr = NewNodeAccessRuleMgr()
		gNodeAccessRuleMgr.loadDataFromDb(0)
		go gNodeAccessRuleMgr.Run()
	}
}

func (self *NodeAccessRuleMgr) Run() {
	self.IsRun = true
	for {
		updatetime, _ := self.loadDataFromDb(self.unixupdatetime)
		if updatetime > self.unixupdatetime {
			self.unixupdatetime = updatetime
		}
		time.Sleep(5 * time.Second)
	}
}

func (self *NodeAccessRuleMgr) loadDataFromDb(unixupdatetime int64) (int64, bool) {
	var bRs bool = false
	var rUpdatetime int64 = 0
	list := GetNodeAccessRuleListByUpdateTime(unixupdatetime)
	self.lock.Lock()
	defer self.lock.Unlock()

	for _, vn := range list {
		updatetime := vn.UpdatedAt.Unix()
		if updatetime > rUpdatetime {
			rUpdatetime = updatetime
		}
		self.elems[vn.RoleId*10000+vn.NodeId] = vn
		log.Printf("Info: NodeAccessRuleMgr.loadDataFromDb Update groupName:%d Id:%d\n", vn.NodeId, vn.RoleId)
		bRs = true
	}
	return rUpdatetime, bRs
}

func GetNodeAccessRuleByRoleId(roleId int64) []*NodeAccessRule {
	fmt.Println("GetNodeAccessRuleByRoleId", roleId)
	InitCheckNodeAccessRuleMgr()
	list := make([]*NodeAccessRule, 0)
	log.Printf("GetNodeAccessRuleByRoleId %d \n", roleId)
	for _, v := range gNodeAccessRuleMgr.elems {
		elem := v
		if elem.RoleId == roleId {
			fmt.Println("yes", elem)
			list = append(list, elem)
		} else {
			fmt.Println("not", elem)
		}
	}
	return list
}

func GetNodeAccessRuleByRoleIdNodeId(roleId int64, nodeId int64) *NodeAccessRule {
	fmt.Println("GetNodeAccessRuleByRoleId", roleId)
	InitCheckNodeAccessRuleMgr()
	log.Printf("GetNodeAccessRuleByRoleId %d \n", roleId)
	for _, v := range gNodeAccessRuleMgr.elems {
		elem := v
		if elem.RoleId == roleId && elem.NodeId == nodeId {
			return elem
		} else {

		}
	}
	return nil
}

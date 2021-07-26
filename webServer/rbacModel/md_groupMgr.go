package rbacModel

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

type GroupMgr struct {
	groups         map[int64]*Group
	lock           sync.Mutex
	IsRun          bool
	unixupdatetime int64
}

var (
	gGroupMgr *GroupMgr
)

func initGroupMgr() {
	if gGroupMgr == nil {

		gGroupMgr = NewGroupMgr()
		gGroupMgr.loadDataFromDb(0)
		go gGroupMgr.Run()
	}
}

func NewGroupMgr() *GroupMgr {
	log.Printf("Info: NewGroupMgr")
	return &GroupMgr{groups: make(map[int64]*Group, 0)}
}

func (self *GroupMgr) Run() {
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

func (self *GroupMgr) loadDataFromDb(unixupdatetime int64) (int64, bool) {
	var bRs bool = false
	var rUpdatetime int64 = 0
	list := GetGroupListByUpdateTime(unixupdatetime)
	self.lock.Lock()
	defer self.lock.Unlock()
	for _, vn := range list {
		updatetime := vn.UpdatedAt.Unix()
		if updatetime > rUpdatetime {
			rUpdatetime = updatetime
		}
		self.groups[vn.Id] = vn
		log.Printf("Info: GroupMgr.loadDataFromDb Update groupName:%s Id:%d\n", vn.Name, vn.Id)
		bRs = true
	}
	return rUpdatetime, bRs
}

func (self *GroupMgr) GetAccessNodesByGroupId(gid int64) (list []*Node) {
	group, _ := self.groups[gid]
	if group != nil {
		log.Println("group!=nil", gid)
		role := GetRoleByRoleId(group.Role.Id)
		if role != nil {
			accRules := GetNodeAccessRuleByRoleId(role.Id)
			if len(accRules) > 0 {
				nodeids := make([]string, 0, len(accRules))
				for _, v := range accRules {
					nodeids = append(nodeids, fmt.Sprintf("%d", v.NodeId))
				}
				ids := strings.Join(nodeids, ",")
				where := fmt.Sprintf("id in(%s)", ids)
				list = GetNodeListByWhere(where)
			}
			return list
		} else {
		}
	} else {
	}
	return list
}
func (self *GroupMgr) GetAccessRuleByGroupId(gid int64) (list []*NodeAccessRule) {
	fmt.Println("GetAccessRuleByGroupId", gid)
	group, _ := self.groups[gid]
	if group != nil {

		list = GetNodeAccessRuleByRoleId(group.RoleId)
		return list
	} else {
		log.Println("Err group==nil", gid)
	}
	return list
}

func GetAccessRuleByGroupIdNodeName(userGroupId int64, nodeName string) *NodeAccessRule {
	initGroupMgr()
	node := GetWebNodeByName(nodeName)
	if node != nil {
		list := gGroupMgr.GetAccessRuleByGroupId(userGroupId)
		for _, v := range list {
			if v.NodeId == node.Id {
				return v
			}
		}
	} else {
		log.Println("Err GetWebNodeByName ==nil NodeName:", nodeName)
	}
	return nil
}

//get group list
func GetGrouplist() (list []*Group) {
	log.Println("GetGrouplist")
	initGroupMgr()
	log.Println(gGroupMgr.groups)
	for _, v := range gGroupMgr.groups {
		elem := v
		list = append(list, elem)
		log.Println("Group", elem)
	}
	return list
}

func AccessList(userGroupId int64) (list []*NodeAccessRule) {
	initGroupMgr()
	list = gGroupMgr.GetAccessRuleByGroupId(userGroupId)
	for _, v := range list {
		node := GetWebNodeById(v.NodeId)
		if node != nil {
			v.NodeName = node.Name
		}
	}
	return list

}

func CheckReqIsAccess(userGroupId int64, nodeName string, action string) bool {
	if userGroupId == 1 { //超级管理员组
		return true
	}
	accessRule := GetAccessRuleByGroupIdNodeName(userGroupId, nodeName)
	if accessRule != nil {
		if accessRule.AccessRule == "*" || accessRule.AccessRule == "" {
			return true
		} else {
			rules := strings.Split(accessRule.AccessRule, ";")
			for _, v := range rules {
				if v == action {
					return false
				}
			}
			return true
		}
	} else {
		log.Println("Err GetAccessRuleByGroupIdNodeName ==nil NodeName:", nodeName, "userGroupId:", userGroupId)
	}
	return true
}

func CheckReqIsRead(userGroupId int64, nodeName string, field string) bool {
	accessRule := GetAccessRuleByGroupIdNodeName(userGroupId, nodeName)
	if accessRule != nil {
		if accessRule.ReadRule == "" {
			return true
		} else {
			rules := strings.Split(accessRule.ReadRule, ";")
			for _, v := range rules {
				if v == field {
					return false
				}
			}
			return true
		}
	} else {
		log.Println("Err GetAccessRuleByGroupIdNodeName ==nil NodeName:", nodeName, "userGroupId:", userGroupId)
	}
	return true
}

func CheckReqIsWrite(userGroupId int64, nodeName string, field string) bool {
	accessRule := GetAccessRuleByGroupIdNodeName(userGroupId, nodeName)
	if accessRule != nil {
		if accessRule.WriteRule == "" {
			return true
		} else {
			rules := strings.Split(accessRule.WriteRule, ";")
			for _, v := range rules {
				if v == field {
					return false
				}
			}
			return true
		}
	} else {
		log.Println("Err GetAccessRuleByGroupIdNodeName ==nil NodeName:", nodeName, "userGroupId:", userGroupId)
	}
	return true
}

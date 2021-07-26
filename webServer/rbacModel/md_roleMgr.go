package rbacModel

import (
	"go.uber.org/zap"
	"sync"
	"time"
)

type RoleMgr struct {
	roles          map[int64]*Role
	lock           sync.Mutex
	IsRun          bool
	unixupdatetime int64
}

var (
	gRoleMgr *RoleMgr
)

func initRoleMgr() {
	if gRoleMgr == nil {
		gRoleMgr = NewRoleMgr()
		gRoleMgr.loadDataFromDb(0)
		go gRoleMgr.Run()
	}
}

func NewRoleMgr() *RoleMgr {
	return &RoleMgr{roles: make(map[int64]*Role, 0)}
}

func (self *RoleMgr) Run() {
	self.IsRun = true
	for {
		updatetime, _ := self.loadDataFromDb(self.unixupdatetime)
		if updatetime > self.unixupdatetime {
			self.unixupdatetime = updatetime
		}
		time.Sleep(5 * time.Second)
	}
}

func (self *RoleMgr) loadDataFromDb(unixupdatetime int64) (int64, bool) {
	var bRs bool = false
	var rUpdatetime int64 = 0
	list := GetRoleListByUpdateTime(unixupdatetime)
	self.lock.Lock()
	defer self.lock.Unlock()

	for _, vn := range list {
		updatetime := vn.UpdatedAt.Unix()
		if updatetime > rUpdatetime {
			rUpdatetime = updatetime
		}
		//vn.Related()
		self.roles[vn.Id] = vn
		zap.L().Info("Info: RoleMgr.loadDataFromDb Update ",
			zap.String("name", vn.Name),
			zap.Int64("id", vn.Id),
		)
		bRs = true
	}
	return rUpdatetime, bRs
}

//get group list
func GetRolelist() (list []*Role) {
	initRoleMgr()
	for _, v := range gRoleMgr.roles {
		elem := v
		list = append(list, elem)
	}
	return list
}

func GetRoleByRoleId(Id int64) *Role {
	initRoleMgr()
	role := gRoleMgr.roles[Id]
	return role
}

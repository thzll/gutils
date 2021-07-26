package rbacModel

import "time"

const TABLE_HEADER = "rbac_"

type Active struct {
	LiveTime   int64 //生命时间 超过此值代表不活动 单位是秒
	activeTime int64 //活动时间单位是秒
}

func (self *Active) checkActive() bool {
	ctime := time.Now().Unix()
	if self.LiveTime == 0 {
		self.LiveTime = 60
	}
	return ctime-self.activeTime < self.LiveTime
}

func (self *Active) setActive() {
	self.activeTime = time.Now().Unix()
}

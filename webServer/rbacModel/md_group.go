package rbacModel

import (
	"fmt"
	"github.com/thzll/gutils/webServer/dbOrm"
	"go.uber.org/zap"
)

//分组表
type Group struct {
	Id              int64
	Name            string `gorm:"size(100)"`
	Title           string `gorm:"size(100)"`
	Status          int64  `gorm:"default(2)"`
	Sort            int64  `gorm:"default(1)"`
	Role            *Role
	RoleId          int64
	dbOrm.TimeModel //时间模型

}

func (n *Group) TableName() string {
	return TABLE_HEADER + "groups"
}

//===============
//add
func AddGroup(u *Group, args map[string]interface{}) (int64, error) {
	return dbOrm.Add(u)
}

func DelGroupById(Id int64) (int64, error) {
	where := fmt.Sprintf("id=%d", Id)
	return dbOrm.DelByWhere(&Group{}, where)
}

//del list
func DelGroupByIds(Ids []int64) (int64, error) {
	where := ""
	for _, v := range Ids {
		if where == "" {
			where = fmt.Sprintf("id=%d", v)
		} else {
			where = fmt.Sprintf("%s or id=%d", where, v)
		}
	}
	if where != "" {
		return dbOrm.DelByWhere(&Group{}, where)
	}
	return 0, fmt.Errorf("Err where == ''")
}

func GetGroupById(id int64) (item *Group) {
	if _, err := dbOrm.GetByWhere(item, fmt.Sprintf("id=%d", id)); err != nil {
		zap.L().Error("GetGroupById", zap.String("err", err.Error()))
	}
	return item
}

//get list
func GetGroupList() (list []*Group) {
	if _, err := dbOrm.GetList(&Group{}, &list); err != nil {
		zap.L().Error("GetGrouplist", zap.String("err", err.Error()))
	}
	return list
}

func GetGroupListByWhere(where string) (list []*Group) {
	if _, err := dbOrm.GetListByWhere(&Group{}, &list, where); err != nil {
		zap.L().Error("GetGrouplist",
			zap.String("where", where),
			zap.String("err", err.Error()),
		)
	}
	return list
}

//获取实例 通过更新时间
func GetGroupListByUpdateTime(updateTime int64) (list []*Group) {
	if _, err := dbOrm.GetListByUpdateTime(&Group{}, &list, updateTime); err != nil {
		zap.L().Error("GetGrouplist",
			zap.Int64("updateTime", updateTime),
			zap.String("err", err.Error()),
		)
	}
	return list
}

//get  Page
func GetGroupListWithPage(pageSize int, offset int, sort string, sortOrder string, where string) (page dbOrm.Page) {
	var vList []Group
	return dbOrm.GetListWithPage(&Group{}, &vList, pageSize, offset, sort, sortOrder, where)
}

func UpdateGroupById(id int64, args map[string]interface{}) (int64, error) {
	where := fmt.Sprintf("id=%d", id)
	return dbOrm.UpdateByWhere(&Group{}, "", "", where, args)
}

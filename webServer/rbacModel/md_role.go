package rbacModel

import (
	"fmt"
	"github.com/thzll/gutils/webServer/dbOrm"
	"go.uber.org/zap"
)

//角色表
type Role struct {
	Id              int64
	Title           string `gorm:"size(100)" form:"Title"  valid:"Required"`
	Name            string `gorm:"size(100)" form:"Name"  valid:"Required"`
	Remark          string `gorm:"null;size(200)" form:"Remark" valid:"MaxSize(200)"`
	Status          int64  `gorm:"default(2)" form:"Status" valid:"Range(1,2)"`
	dbOrm.TimeModel        //时间模型
}

func (s *Role) TableName() string {
	return TABLE_HEADER + "roles"
}

//===============
//add
func AddRole(u *Role, args map[string]interface{}) (int64, error) {
	return dbOrm.Add(u)
}

func DelRoleById(Id int64) (int64, error) {
	where := fmt.Sprintf("id=%d", Id)
	return dbOrm.DelByWhere(&Role{}, where)
}

//del list
func DelRoleByIds(Ids []int64) (int64, error) {
	where := ""
	for _, v := range Ids {
		if where == "" {
			where = fmt.Sprintf("id=%d", v)
		} else {
			where = fmt.Sprintf("%s or id=%d", where, v)
		}
	}
	if where != "" {
		return dbOrm.DelByWhere(&Role{}, where)
	}
	return 0, fmt.Errorf("Err where == ''")
}

func GetRoleById(id int64) (item *Role) {
	if _, err := dbOrm.GetByWhere(item, fmt.Sprintf("id=%d", id)); err != nil {
		zap.L().Error("GetRoleById", zap.String("err", err.Error()))
	}
	return item
}

//get list
func GetRoleList() (list []*Role) {
	if _, err := dbOrm.GetList(&Role{}, &list); err != nil {
		zap.L().Error("GetRolelist", zap.String("err", err.Error()))
	}
	return list
}

func GetRoleListByWhere(where string) (list []*Role) {
	if _, err := dbOrm.GetListByWhere(&Role{}, &list, where); err != nil {
		zap.L().Error("GetRolelist",
			zap.String("where", where),
			zap.String("err", err.Error()),
		)
	}
	return list
}

//获取实例 通过更新时间
func GetRoleListByUpdateTime(updateTime int64) (list []*Role) {
	if _, err := dbOrm.GetListByUpdateTime(&Role{}, &list, updateTime); err != nil {
		zap.L().Error("GetRolelist",
			zap.Int64("updateTime", updateTime),
			zap.String("err", err.Error()),
		)
	}
	return list
}

//get  Page
func GetRoleListWithPage(pageSize int, offset int, sort string, sortOrder string, where string) (page dbOrm.Page) {
	var vList []Role
	return dbOrm.GetListWithPage(&Role{}, &vList, pageSize, offset, sort, sortOrder, where)
}

func UpdateRole(role *Role) (int64, error) {
	return dbOrm.UpdateByWhere(role, "", "", "", nil)
}

func UpdateRoleById(id int64, args map[string]interface{}) (int64, error) {
	where := fmt.Sprintf("id=%d", id)
	return dbOrm.UpdateByWhere(&Role{}, "", "", where, args)
}

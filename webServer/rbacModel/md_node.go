package rbacModel

import (
	"fmt"
	"github.com/thzll/gutils/webServer/dbOrm"
	"go.uber.org/zap"
)

//节点表
type Node struct {
	Id              int64
	Name            string  `gorm:"size(100)"`
	Title           string  `gorm:"size(100)"`
	Level           int64   `gorm:"default(1)"`
	Sort            int64   `gorm:"default(1)"`
	Pid             int64   `form:"Pid"`
	Class           string  `gorm:"size(100)"`
	Remark          string  `gorm:"size(200)"`
	Status          int64   `gorm:"default(2)"`
	GroupId         int64   `gorm:"default(1)`
	Role            []*Role `gorm:"rel(m2m)"`
	nn              []byte
	dbOrm.TimeModel //时间模型
	//其他属性
	RoleId     int64  `gorm:"-"`
	IsAccess   int64  `gorm:"-"` //是否可访问
	AccessRule string `gorm:"-"` //可访问规则
	ReadRule   string `gorm:"-"` //可读数据包字段
	WriteRule  string `gorm:"-"` //可读数据包字段

}

func (n *Node) TableName() string {
	return TABLE_HEADER + "nodes"
}

//===============
//add
func AddNode(u *Node, args map[string]interface{}) (int64, error) {
	return dbOrm.Add(u)
}

func DelNodeById(Id int64) (int64, error) {
	where := fmt.Sprintf("id=%d", Id)
	return dbOrm.DelByWhere(&Node{}, where)
}

//del list
func DelNodeByIds(Ids []int64) (int64, error) {
	where := ""
	for _, v := range Ids {
		if where == "" {
			where = fmt.Sprintf("id=%d", v)
		} else {
			where = fmt.Sprintf("%s or id=%d", where, v)
		}
	}
	if where != "" {
		return dbOrm.DelByWhere(&Node{}, where)
	}
	return 0, fmt.Errorf("Err where == ''")
}

func GetNodeById(id int64) (item *Node) {
	if _, err := dbOrm.GetByWhere(item, fmt.Sprintf("id=%d", id)); err != nil {
		zap.L().Error("GetNodeById", zap.String("err", err.Error()))
	}
	return item
}

//get list
func GetNodeList() (list []*Node) {
	if _, err := dbOrm.GetList(&Node{}, &list); err != nil {
		zap.L().Error("GetNodelist", zap.String("err", err.Error()))
	}
	return list
}

func GetNodeListByWhere(where string) (list []*Node) {
	if _, err := dbOrm.GetListByWhere(&Node{}, &list, where); err != nil {
		zap.L().Error("GetNodelist",
			zap.String("where", where),
			zap.String("err", err.Error()),
		)
	}
	return list
}

//获取实例 通过更新时间
func GetNodeListByUpdateTime(updateTime int64) (list []*Node) {
	if _, err := dbOrm.GetListByUpdateTime(&Node{}, &list, updateTime); err != nil {
		zap.L().Error("GetNodelist",
			zap.Int64("updateTime", updateTime),
			zap.String("err", err.Error()),
		)
	}
	return list
}

//get  Page
func GetNodeListWithPage(pageSize int, offset int, sort string, sortOrder string, where string) (page dbOrm.Page) {
	var vList []Node
	return dbOrm.GetListWithPage(&Node{}, &vList, pageSize, offset, sort, sortOrder, where)
}

func UpdateNodeById(id int64, args map[string]interface{}) (int64, error) {
	where := fmt.Sprintf("id=%d", id)
	return dbOrm.UpdateByWhere(&Node{}, "", "", where, args)
}

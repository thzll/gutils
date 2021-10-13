package rbacModel

import (
	"fmt"
	"github.com/thzll/gutils/webServer/dbOrm"
	"go.uber.org/zap"
	//"github.com/astaxie/beego/validation"
)

//节点表
type NodeAccessRule struct {
	NodeId          int64  `gorm:""`          //节点ID
	RoleId          int64  `gorm:""`          //规则Id
	IsAccess        int64  `gorm:"" `         //是否可访问
	AccessRule      string `gorm:"size(255)"` //可访问规则
	ReadRule        string `gorm:"size(255)"` //可读数据包字段
	WriteRule       string `gorm:"size(255)"` //可读数据包字段
	dbOrm.TimeModel        //时间模型

	NodeName string `gorm:"-;"` //节点ID
}

func (n *NodeAccessRule) TableName() string {
	return TABLE_HEADER + "node_access_rules"
}

//===============
//add
func AddNodeAccessRule(u *NodeAccessRule, args map[string]interface{}) (int64, error) {
	return dbOrm.Add(u)
}

func InserOrUpdateNodeAccessRule(args map[string]interface{}) (int64, error) {
	obj := new(NodeAccessRule)
	dbOrm.SetObjValue(obj, args)
	return dbOrm.InsertOrUpdate(obj, args)
}

func DelNodeAccessRuleById(Id int64) (int64, error) {
	where := fmt.Sprintf("id=%d", Id)
	return dbOrm.DelByWhere(&NodeAccessRule{}, where)
}

//del list
func DelNodeAccessRuleByIds(Ids []int64) (int64, error) {
	where := ""
	for _, v := range Ids {
		if where == "" {
			where = fmt.Sprintf("id=%d", v)
		} else {
			where = fmt.Sprintf("%s or id=%d", where, v)
		}
	}
	if where != "" {
		return dbOrm.DelByWhere(&NodeAccessRule{}, where)
	}
	return 0, fmt.Errorf("Err where == ''")
}

func GetNodeAccessRuleById(id int64) (item *NodeAccessRule) {
	if _, err := dbOrm.GetByWhere(item, fmt.Sprintf("id=%d", id)); err != nil {
		zap.L().Error("GetNodeAccessRuleById", zap.String("err", err.Error()))
	}
	return item
}

//get list
func GetNodeAccessRuleList() (list []*NodeAccessRule) {
	if _, err := dbOrm.GetList(&NodeAccessRule{}, &list); err != nil {
		zap.L().Error("GetNodeAccessRulelist", zap.String("err", err.Error()))
	}
	return list
}

func GetNodeAccessRuleListByWhere(where string) (list []*NodeAccessRule) {
	if _, err := dbOrm.GetListByWhere(&NodeAccessRule{}, &list, where); err != nil {
		zap.L().Error("GetNodeAccessRulelist",
			zap.String("where", where),
			zap.String("err", err.Error()),
		)
	}
	return list
}

//获取实例 通过更新时间
func GetNodeAccessRuleListByUpdateTime(updateTime int64) (list []*NodeAccessRule) {
	if _, err := dbOrm.GetListByUpdateTime(&NodeAccessRule{}, &list, updateTime); err != nil {
		zap.L().Error("GetNodeAccessRulelist",
			zap.Int64("updateTime", updateTime),
			zap.String("err", err.Error()),
		)
	}
	return list
}

//get  Page
func GetNodeAccessRuleListWithPage(pageSize int, offset int, sort string, sortOrder string, where string) (page dbOrm.Page) {
	var vList []NodeAccessRule
	return dbOrm.GetListWithPage(&NodeAccessRule{}, &vList, pageSize, offset, sort, sortOrder, where)
}

func UpdateNodeAccessRuleById(id int64, args map[string]interface{}) (int64, error) {
	where := fmt.Sprintf("id=%d", id)
	return dbOrm.UpdateByWhere(&NodeAccessRule{}, "", "", where, args)
}

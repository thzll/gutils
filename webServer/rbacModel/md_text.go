package rbacModel

import (
	"fmt"
	"go.uber.org/zap"
	"goplot/common/dbOrm"
)

//网站界面文本分组 按节点路径命名
type Text struct {
	Id              int64
	Name            string `gorm:"size(32)"`
	Title           string `gorm:"size(32)"`
	NodeName        string `gorm:"size(32)"`
	Text            string `gorm:"size(32)"`
	dbOrm.TimeModel        //时间模型
}

func (n *Text) TableName() string {
	return TABLE_HEADER + "texts"
}

//===============
//add
func AddText(u *Text, args map[string]interface{}) (int64, error) {
	return dbOrm.Add(u)
}

func DelTextById(Id int64) (int64, error) {
	where := fmt.Sprintf("id=%d", Id)
	return dbOrm.DelByWhere(&Text{}, where)
}

//del list
func DelTextByIds(Ids []int64) (int64, error) {
	where := ""
	for _, v := range Ids {
		if where == "" {
			where = fmt.Sprintf("id=%d", v)
		} else {
			where = fmt.Sprintf("%s or id=%d", where, v)
		}
	}
	if where != "" {
		return dbOrm.DelByWhere(&Text{}, where)
	}
	return 0, fmt.Errorf("Err where == ''")
}

func GetTextById(id int64) (item *Text) {
	if _, err := dbOrm.GetByWhere(item, fmt.Sprintf("id=%d", id)); err != nil {
		zap.L().Error("GetTextById", zap.String("err", err.Error()))
	}
	return item
}

//get list
func GetTextList() (list []*Text) {
	if _, err := dbOrm.GetList(&Text{}, &list); err != nil {
		zap.L().Error("GetTextlist", zap.String("err", err.Error()))
	}
	return list
}

func GetTextListByWhere(where string) (list []*Text) {
	if _, err := dbOrm.GetListByWhere(&Text{}, &list, where); err != nil {
		zap.L().Error("GetTextlist",
			zap.String("where", where),
			zap.String("err", err.Error()),
		)
	}
	return list
}

//获取实例 通过更新时间
func GetTextListByUpdateTime(updateTime int64) (list []*Text) {
	if _, err := dbOrm.GetListByUpdateTime(&Text{}, &list, updateTime); err != nil {
		zap.L().Error("GetTextlist",
			zap.Int64("updateTime", updateTime),
			zap.String("err", err.Error()),
		)
	}
	return list
}

//get  Page
func GetTextListWithPage(pageSize int, offset int, sort string, sortOrder string, where string) (page dbOrm.Page) {
	var vList []Text
	return dbOrm.GetListWithPage(&Text{}, &vList, pageSize, offset, sort, sortOrder, where)
}

func UpdateTextById(id int64, args map[string]interface{}) (int64, error) {
	where := fmt.Sprintf("id=%d", id)
	return dbOrm.UpdateByWhere(&Text{}, "", "", where, args)
}

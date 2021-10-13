package rbacModel

import (
	"fmt"
	"github.com/thzll/gutils/webServer/dbOrm"
	"go.uber.org/zap"
)

//用户表
type Order struct {
	Id              int64
	OrderId         int64  `gorm:"default(0)"` //订单编号 时间加编号
	PayOrderId      string `gorm:"size(32)"`   //订单编号 网络支付编号
	AccountId       string `gorm:"size(32)"`   //订单关联的账号
	Mode            int64  `gorm:"default(0)"` //0：充值  1：消费
	Val             int64  `gorm:"default(0)"` //价值
	Text            string `gorm:"size(1024)"` //描述
	Descs           string `gorm:"size(1024)"` //描述
	Status          int64  `gorm:"default(0)"` //0：订单未完成待支付 1：订单已经完成
	dbOrm.TimeModel        //时间模型
}

func (n *Order) TableName() string {
	return TABLE_HEADER + "orders"
}

//===============
//add
func AddOrder(u *Order, args map[string]interface{}) (int64, error) {
	return dbOrm.Add(u)
}

func DelOrderById(Id int64) (int64, error) {
	where := fmt.Sprintf("id=%d", Id)
	return dbOrm.DelByWhere(&Order{}, where)
}

//del list
func DelOrderByIds(Ids []int64) (int64, error) {
	where := ""
	for _, v := range Ids {
		if where == "" {
			where = fmt.Sprintf("id=%d", v)
		} else {
			where = fmt.Sprintf("%s or id=%d", where, v)
		}
	}
	if where != "" {
		return dbOrm.DelByWhere(&Order{}, where)
	}
	return 0, fmt.Errorf("Err where == ''")
}

func GetOrderById(id int64) (item *Order) {
	item = &Order{}
	if _, err := dbOrm.GetByWhere(item, fmt.Sprintf("id=%d", id)); err != nil {
		zap.L().Error("GetOrderById", zap.String("err", err.Error()))
		return nil
	}
	return item
}

func GetOrderByOrderId(id int64) (item *Order) {
	item = &Order{}
	if _, err := dbOrm.GetByWhere(item, fmt.Sprintf("order_id=%d", id)); err != nil {
		zap.L().Error("GetOrderByOrderId", zap.String("err", err.Error()))
		return nil
	}
	return item
}

//get list
func GetOrderList() (list []*Order) {
	if _, err := dbOrm.GetList(&Order{}, &list); err != nil {
		zap.L().Error("GetOrderlist", zap.String("err", err.Error()))
	}
	return list
}

func GetOrderListByWhere(where string) (list []*Order) {
	if _, err := dbOrm.GetListByWhere(&Order{}, &list, where); err != nil {
		zap.L().Error("GetOrderlist",
			zap.String("where", where),
			zap.String("err", err.Error()),
		)
	}
	return list
}

//获取实例 通过更新时间
func GetOrderListByUpdateTime(updateTime int64) (list []*Order) {
	if _, err := dbOrm.GetListByUpdateTime(&Order{}, &list, updateTime); err != nil {
		zap.L().Error("GetOrderlist",
			zap.Int64("updateTime", updateTime),
			zap.String("err", err.Error()),
		)
	}
	return list
}

//get  Page
func GetOrderListWithPage(pageSize int, offset int, sort string, sortOrder string, where string) (page dbOrm.Page) {
	var vList []Order
	return dbOrm.GetListWithPage(&Order{}, &vList, pageSize, offset, sort, sortOrder, where)
}

func SetOrderFinish(Id int64) (int64, error) {
	fmt.Println("SetOrderFinish", Id)
	order := GetOrderByOrderId(Id)
	if order == nil {
		return 0, fmt.Errorf("can't find order by id %d", Id)
	} else {
		switch order.Mode {
		case 0: //充值
			{
				err := AddMoney(order.AccountId, order.Val)
				if err != nil {
					fmt.Println("Err ", err.Error())
					return 0, err
				}
			}
		case 1: //消费
			{
				err := AddMoney(order.AccountId, 0-order.Val)
				if err != nil {
					fmt.Println("Err ", err.Error())
					return 0, err
				}
			}
		default:
			{
				return 0, fmt.Errorf("Err mode:%d faild", order.Mode)
			}
		}
	}
	data := make(map[string]interface{})
	data["Status"] = 1
	num, err := UpdateOrderById(Id, data)
	return num, err
}

func UpdateOrderById(id int64, args map[string]interface{}) (int64, error) {
	return dbOrm.UpdateByWhere(&Order{Id: id}, "", "", "", args)
}

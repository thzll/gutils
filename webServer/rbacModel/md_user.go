package rbacModel

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"go.uber.org/zap"
	"goplot/common/dbOrm"
	"goplot/common/myutils"
	"log"
	"time"
)

//用户表
type User struct {
	Id              int64
	Username        string    `gorm:"unique;size(32)" json:"username"`
	Password        string    `gorm:"size(32)" json:"password"`
	Repassword      string    `gorm:"-" json:"repassword"`
	Nickname        string    `gorm:"unique;size(32)"`
	Email           string    `gorm:"size(32)" json:"email"`
	Qq              string    `gorm:"size(32)" json:"qq"`
	Remark          string    `gorm:"null;size(200)" json:"remark"`                //评论
	Verification    string    `gorm:"null;size(200)" json:"verification"`          //验证码 用于找回密码
	ValidTime       time.Time `gorm:"DEFAULT:CURRENT_TIMESTAMP" json:"valid_time"` //有效时间
	Money           int64     `gorm:"default(0)" json:"money"`                     //账号余额
	Discount        int64     `gorm:"default(100)" json:"discount"`                //账号折扣 100代表不打折 90代表9折
	Consume         int64     `gorm:"default(0)" json:"consume"`                   //消费
	Recharge        int64     `gorm:"default(0)" json:"recharge"`                  //充值
	Ip              string    `gorm:"null;size(200)" json:"ip"`
	Status          int64     `gorm:"default(2)" json:"status"`
	Telephone       string    `gorm:"size(32)" json:"telephone"` //电话号码
	Descs           string    `gorm:"size(32)" json:"descs"`     //描述字段
	CardId          string    `gorm:"size(32)" json:"card_id"`   //身份证号码
	CardName        string    `gorm:"size(32)" json:"card_name"` //身份证人名
	Ui              int64     `gorm:"default(0)"`                //Ui权限
	LastLoginAt     time.Time `gorm:"DEFAULT:CURRENT_TIMESTAMP"`
	Group           *Group
	GroupId         int64
	dbOrm.TimeModel //时间模型

	InstanceCount        int64 `gorm:"-"` //实例个数
	ActiveInstanceCount  int64 `gorm:"-"` //活动的实例个数
	InvalidInstanceCount int64 `gorm:"-"` //失效的实例个数
}

var g_UserCache map[string]*User //用户缓存

func (n *User) TableName() string {
	return TABLE_HEADER + "users"
}

func AddMoney(username string, money int64) error {
	fmt.Println("AddMoney username:", username, "money:", money)
	u := GetUserByName(username)
	if u.Id > 0 {
		u.Money += money
		if u.Money >= 0 {
			params := make(map[string]interface{})
			params["Money"] = u.Money
			_, err := UpdateUserByName(username, params)
			return err
		} else {
			return fmt.Errorf("Err 余额不足")
		}
	} else {
		return fmt.Errorf("Err 账号不存在")
	}
}

//===============
//add
func AddUser(u *User, args map[string]interface{}) (int64, error) {
	return dbOrm.Add(u)
}

func DelUserById(Id int64) (int64, error) {
	where := fmt.Sprintf("id=%d", Id)
	return dbOrm.DelByWhere(&User{}, where)
}

//del list
func DelUserByIds(Ids []int64) (int64, error) {
	where := ""
	for _, v := range Ids {
		if where == "" {
			where = fmt.Sprintf("id=%d", v)
		} else {
			where = fmt.Sprintf("%s or id=%d", where, v)
		}
	}
	if where != "" {
		return dbOrm.DelByWhere(&User{}, where)
	}
	return 0, fmt.Errorf("Err where == ''")
}

func DelUserByName(name string) (int64, error) {
	where := fmt.Sprintf("username='%s'", name)
	return dbOrm.DelByWhere(&User{}, where)
}

func GetUserById(id int64) (item *User) {
	if _, err := dbOrm.GetByWhere(item, fmt.Sprintf("id=%d", id)); err != nil {
		zap.L().Error("GetUserById", zap.String("err", err.Error()))
	}
	return item
}

func GetUserByName(account string) *User {
	user := &User{}
	if _, err := dbOrm.GetByWhere(user, fmt.Sprintf("username='%s'", account)); err != nil {
		zap.L().Error("GetUserById", zap.String("err", err.Error()))
		return nil
	}
	//if _, err := dbOrm.Related(user, user.Group); err != nil {
	//	zap.L().Error("User Related", zap.String("err", err.Error()))
	//}
	return user
}

//get list
func GetUserList() (list []*User) {
	if _, err := dbOrm.GetList(&User{}, &list); err != nil {
		zap.L().Error("GetUserlist", zap.String("err", err.Error()))
	}
	return list
}

func GetUserListByWhere(where string) (list []*User) {
	if _, err := dbOrm.GetListByWhere(&User{}, &list, where); err != nil {
		zap.L().Error("GetUserlist",
			zap.String("where", where),
			zap.String("err", err.Error()),
		)
	}
	return list
}

//获取实例 通过更新时间
func GetUserListByUpdateTime(updateTime int64) (list []*User) {
	if _, err := dbOrm.GetListByUpdateTime(&User{}, &list, updateTime); err != nil {
		zap.L().Error("GetUserlist",
			zap.Int64("updateTime", updateTime),
			zap.String("err", err.Error()),
		)
	}
	return list
}

//get  Page
func GetUserListWithPage(pageSize int, offset int, sort string, sortOrder string, where string) (page dbOrm.Page) {
	var vList []User
	log.Println(vList)
	return dbOrm.GetListWithPage(&User{}, &vList, pageSize, offset, sort, sortOrder, where)
}

func ResetPassword(username string, password string, verification string) error {
	u := GetUserByName(username)
	if u.Id > 0 {
		if time.Now().Unix() > u.ValidTime.Unix() {
			fmt.Println(time.Now().Unix() - u.ValidTime.Unix())
			return fmt.Errorf("重置密码连接失效")
		}
		if u.Verification != verification {
			return fmt.Errorf("重置密码校验码验证失败")
		}
		passwordMd5 := myutils.Strtomd5(password)
		if u.Password == passwordMd5 {
			return fmt.Errorf("新密码并不能与源密码相同")
		}
		user := make(map[string]interface{})
		user["Password"] = passwordMd5
		user["ValidTime"] = time.Now()
		_, err := UpdateUserByName(username, user)
		return err
	} else {
		return fmt.Errorf("账号错误")
	}
}

func UpdateUserById(id int64, args map[string]interface{}) (int64, error) {
	where := fmt.Sprintf("id=%d", id)
	return dbOrm.UpdateByWhere(&User{}, "", "", where, args)
}

func UpdateUserByName(account string, args map[string]interface{}) (int64, error) {
	where := fmt.Sprintf("username='%s'", account)
	return dbOrm.UpdateByWhere(&User{}, "", "", where, args)
}

func UpdateLoginStatus(username string) (int64, error) {
	user := make(orm.Params)
	user["lastlogintime"] = time.Now()
	return UpdateUserByName(username, user)
}

//func MailResetPassword(username string) error {
//	user := make(map[string]interface{})
//	verification := fmt.Sprintf("%d", rand.Uint64())
//	validTime := time.Now().Add(time.Second * 600) //10分组内有效
//	user["Verification"] = verification
//	user["ValidTime"] = validTime
//	if _, err := UpdateUserByName(username, user); err != nil {
//		return err
//	}
//	u := GetUserByName(username)
//	err := mail.SendMailUsingTLS("admin@chuanliukeji.com",
//		"fgpjlybxwsklbhbd",
//		"smtp.qq.com:465",
//		u.Email,
//		"找回密码",
//		fmt.Sprintf("点击这里修改密码：http://console.cldun.com/public/forgot_password?a=reset_password&username=%s&verification=%s", username, verification),
//		"html",
//	)
//	fmt.Println("sendMail ", err)
//	return err
//}

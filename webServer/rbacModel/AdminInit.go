package rbacModel

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/zap"
	"goplot/common/dbOrm"
	"goplot/common/myutils"
	"time"
	// import _ "github.com/jinzhu/gorm/dialects/postgres"
	// import _ "github.com/jinzhu/gorm/dialects/sqlite"
	// import _ "github.com/jinzhu/gorm/dialects/mssql"
)

func getDb() *gorm.DB {
	return dbOrm.GetDb("")
}

func Syncdb(_force bool, skipTab ...string) {
	db := getDb()
	if _force {
		//db.DropTable(&User{})
		db.DropTable(&Group{})
		db.DropTable(&Node{})
		db.DropTable(&Role{})
		db.DropTable(&Order{})
		db.DropTable(&Text{})
		db.DropTable(&NodeAccessRule{})
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Group{})
	db.AutoMigrate(&Node{})
	db.AutoMigrate(&Role{})
	db.AutoMigrate(&Order{})
	db.AutoMigrate(&Text{})
	db.AutoMigrate(&NodeAccessRule{})
	if _force {
		insertNodes(db)
		insertRole(db)
		inserNodeAccessRule(db)
		insertGroup(db)
		insertUser(db)
		time.Sleep(time.Second * 10)
		fmt.Println("database init is complete.\nPlease restart the application")
	}
}

//网站数据库连接

func inserNodeAccessRule(db *gorm.DB) {
	return
	//管理员
	for i := 1; i < 38; i++ {
		if i == 1 || (i >= 4 && i <= 8) {
			continue
		}
		AccessRule := &NodeAccessRule{NodeId: int64(i),
			RoleId:     2,
			IsAccess:   1,
			AccessRule: "",
			ReadRule:   "",
			WriteRule:  "",
		}
		if _, err := AddNodeAccessRule(AccessRule, nil); err != nil {
			zap.L().Error("AddRole", zap.String("err", err.Error()))
		}
	}

	//普通组长
	for i := 1; i < 38; i++ {
		if i == 3 || (i >= 31 && i <= 37) {
			AccessRule := &NodeAccessRule{NodeId: int64(i),
				RoleId:     3,
				IsAccess:   1,
				AccessRule: "",
				ReadRule:   "",
				WriteRule:  "",
			}
			if _, err := AddNodeAccessRule(AccessRule, nil); err != nil {
				zap.L().Error("AddRole", zap.String("err", err.Error()))
			}
		}
	}

	//普通用户
	for i := 1; i < 38; i++ {
		if i == 3 || (i >= 31 && i <= 37) {
			AccessRule := &NodeAccessRule{NodeId: int64(i),
				RoleId:     4,
				IsAccess:   1,
				AccessRule: "",
				ReadRule:   "",
				WriteRule:  "",
			}
			if _, err := AddNodeAccessRule(AccessRule, nil); err != nil {
				zap.L().Error("AddRole", zap.String("err", err.Error()))
			}
		}
	}
}

func insertGroup(*gorm.DB) {
	fmt.Println("insert group ...")
	groups := [4]Group{
		{Name: "Admin", Title: "超级管理员", Status: 2},
		{Name: "Administrators", Title: "管理员", Status: 2},
		{Name: "Leader", Title: "普通组长", Status: 2},
		{Name: "User", Title: "普通用户", Status: 2},
	}
	for k, v := range groups {
		g := new(Group)
		g.Name = v.Name
		g.Status = v.Status
		g.Title = v.Title
		g.Sort = 1
		g.RoleId = int64(k + 1)
		if _, err := AddGroup(g, nil); err != nil {
			zap.L().Error("AddGroup", zap.String("err", err.Error()))
		}
	}
	fmt.Println("insert group end", len(groups))
}

func insertUser(*gorm.DB) {
	fmt.Println("insert user ...")
	u := new(User)
	u.Username = "admin"
	u.Nickname = "adminManster"
	u.Password = myutils.Pwdhash("Lin643119")
	u.Email = "cldun@gmail.com"
	u.Remark = "I m admin"
	u.Status = 2
	u.GroupId = 1
	if _, err := AddUser(u, nil); err != nil {
		zap.L().Error("AddUser", zap.String("err", err.Error()))
	}
	fmt.Println("insert user end", 1)
}

func insertRole(*gorm.DB) {
	fmt.Println("insert role ...")
	roles := [4]Role{
		{Name: "Admin", Title: "超级管理员规则", Remark: "超级管理员", Status: 2},
		{Name: "Administrators", Title: "管理员规则", Remark: "管理员", Status: 2},
		{Name: "Leader", Title: "普通组长规则", Remark: "普通组长", Status: 2},
		{Name: "User", Title: "普通用户规则", Remark: "用户", Status: 2},
	}
	for _, v := range roles {
		r := new(Role)
		r.Name = v.Name
		r.Remark = v.Remark
		r.Status = v.Status
		r.Title = v.Title
		if _, err := AddRole(r, nil); err != nil {
			zap.L().Error("AddRole", zap.String("err", err.Error()))
		}
	}
	fmt.Println("insert role end", len(roles))
}

func insertNodes(*gorm.DB) {
	fmt.Println("insert node ...")
	//nodes := make([20]Node)
	nodes := [12]Node{
		{Id: 1, Name: "/rbac/", Title: "权限管理", Remark: "", Level: 1, Pid: 0, Class: "icon-key", Sort: 3, Status: 2, GroupId: 1},
		{Id: 2, Name: "/admin/", Title: "川流盾管理后台", Remark: "", Level: 1, Pid: 0, Class: "icon-cogs", Sort: 1, Status: 2, GroupId: 1},
		{Id: 3, Name: "/user/", Title: "客户管理后台", Remark: "", Level: 1, Pid: 0, Class: "icon-user", Sort: 2, Status: 2, GroupId: 1},

		{Id: 11, Name: "/rbac/node/", Title: "节点列表", Remark: "", Level: 2, Pid: 1, Class: "icon-sitemap", Sort: 1, Status: 2, GroupId: 1},
		{Id: 12, Name: "/rbac/user/", Title: "用户列表", Remark: "", Level: 2, Pid: 1, Class: "icon-user-md", Sort: 2, Status: 2, GroupId: 1},
		{Id: 13, Name: "/rbac/group/", Title: "用户分组", Remark: "", Level: 2, Pid: 1, Class: "icon-group", Sort: 3, Status: 2, GroupId: 1},
		{Id: 14, Name: "/rbac/role/", Title: "权限规则", Remark: "", Level: 2, Pid: 1, Class: "icon-lock", Sort: 4, Status: 2, GroupId: 1},
		{Id: 15, Name: "/rbac/node_text/", Title: "文本分组", Remark: "", Level: 2, Pid: 1, Class: "icon-lock", Sort: 5, Status: 2, GroupId: 1},

		{Id: 21, Name: "/customer/", Title: "客户管理", Remark: "", Level: 2, Pid: 2, Class: "icon-group", Sort: 0, Status: 0, GroupId: 1},
		{Id: 22, Name: "/customer/", Title: "账号", Remark: "", Level: 2, Pid: 2, Class: "icon-sitemap", Sort: 10, Status: 0, GroupId: 1}, //9
		{Id: 23, Name: "/commodity/", Title: "产品管理", Remark: "", Level: 2, Pid: 2, Class: "icon-shield", Sort: 20, Status: 0, GroupId: 1},
		{Id: 24, Name: "/general/", Title: "通用属性", Remark: "", Level: 2, Pid: 2, Class: "icon-certificate", Sort: 30, Status: 0, GroupId: 1},

		//{Id: 211, Name: "/customer/tests", Title: "测试", Remark: "", Level: 2, Pid: 2, Class: "icon-reorder", Sort: 1, Status: 2, GroupId: 1},
		//{Id: 212, Name: "/customer/game_accounts", Title: "游戏账号", Remark: "", Level: 2, Pid: 2, Class: "icon-cogs", Sort: 2, Status: 2, GroupId: 1},
		//{Id: 213, Name: "/customer/lua_cmds", Title: "Lua命令", Remark: "", Level: 2, Pid: 2, Class: "icon-cogs", Sort: 3, Status: 2, GroupId: 1},
		//{Id: 214, Name: "/customer/net_pack", Title: "数据包历史", Remark: "", Level: 2, Pid: 2, Class: "icon-cogs", Sort: 104, Status: 2, GroupId: 1},
		//{Id: 215, Name: "/customer/net_protocol", Title: "数据包协议", Remark: "", Level: 2, Pid: 2, Class: "icon-cogs", Sort: 105, Status: 2, GroupId: 1},
		//{Id: 216, Name: "/customer/npc_infos", Title: "Npc信息", Remark: "", Level: 2, Pid: 2, Class: "icon-cogs", Sort: 6, Status: 2, GroupId: 1},
		//{Id: 217, Name: "/customer/item_rules", Title: "物品规则", Remark: "", Level: 2, Pid: 2, Class: "icon-cogs", Sort: 6, Status: 2, GroupId: 1},
		//{Id: 218, Name: "/customer/role_rules", Title: "角色规则", Remark: "", Level: 2, Pid: 2, Class: "icon-cogs", Sort: 6, Status: 2, GroupId: 1},
		//{Id: 219, Name: "/customer/guaji_rules", Title: "挂机点", Remark: "", Level: 2, Pid: 2, Class: "icon-cogs", Sort: 6, Status: 2, GroupId: 1},
		//{Id: 220, Name: "/customer/server_names", Title: "服务器名", Remark: "", Level: 2, Pid: 2, Class: "icon-cogs", Sort: 6, Status: 2, GroupId: 1},
		//{Id: 221, Name: "/customer/magic_names", Title: "技能名", Remark: "", Level: 2, Pid: 2, Class: "icon-cogs", Sort: 6, Status: 2, GroupId: 1},
	}
	for _, v := range nodes {
		n := new(Node)
		n.Name = v.Name
		n.Title = v.Title
		n.Remark = v.Remark
		n.Level = v.Level
		n.Sort = v.Sort
		n.Pid = v.Pid
		n.Class = v.Class
		n.Status = v.Status
		n.GroupId = v.GroupId
		if _, err := AddNode(n, nil); err != nil {
			zap.L().Error("AddNode", zap.String("err", err.Error()))
		}

	}
	fmt.Println("insert node end count", len(nodes))
}

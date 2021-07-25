package dbOrm

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/thzll/gutils/webServer/settings"
	"log"
	// import _ "github.com/jinzhu/gorm/dialects/postgres"
	// import _ "github.com/jinzhu/gorm/dialects/sqlite"
	// import _ "github.com/jinzhu/gorm/dialects/mssql"
)

var gdbmap map[string]*gorm.DB

func init() {
	gdbmap = make(map[string]*gorm.DB)
}

func Init(nameSpace string, config *settings.MySQLConfig) (err error) {
	if nameSpace == "" {
		nameSpace = "default"
	}
	db_type := config.DbType
	db_host := config.Host
	db_port := config.Port
	db_user := config.User
	db_pass := config.Password
	db_name := config.DbName
	db_path := config.DbPath
	db_sslmode := config.DbSslmode
	log.Println(config)
	var dsn string
	var db *gorm.DB
	switch db_type {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local", //这里不指定时区 不然服务器和本地测试环境可能不一样
			//dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&loc=Local&parseTime=true",
			db_user, db_pass, db_host, db_port, db_name)
		db, err = gorm.Open("mysql", dsn)
		break
	case "postgres":
		dsn = fmt.Sprintf("dbname=%s host=%s  user=%s  password=%s  port=%d  sslmode=%s",
			db_name, db_host, db_user, db_pass, db_port, db_sslmode)
		db, err = gorm.Open("postgres", dsn)
	case "sqlite3":
		if db_path == "" {
			db_path = "./"
		}
		dsn = fmt.Sprintf("%s%s.db", db_path, db_name)
		db, err = gorm.Open("sqlite3", dsn)
	default:
		return fmt.Errorf("Database driver is not allowed:", db_type)
	}
	if err != nil {
		return err
	} else {
		db.DB().SetMaxIdleConns(config.MaxIdleConns)
		db.DB().SetMaxOpenConns(config.MaxOpenConns)
		db.Callback().Create().Remove("gorm:update_time_stamp")
		db.Callback().Update().Remove("gorm:update_time_stamp")
		gdbmap[nameSpace] = db
		return nil
	}
}

func DbClose() {
	for _, v := range gdbmap {
		v.Close()
	}
}

func GetDb(nameSpace string) *gorm.DB {
	return getDb(nameSpace)
}

func getDb(nameSpace string) *gorm.DB {
	if nameSpace == "" {
		nameSpace = "default"
	}
	return gdbmap[nameSpace]
}

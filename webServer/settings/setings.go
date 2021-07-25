package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"time"
)

var Conf = new(AppConfig)

type AppConfig struct {
	Name           string `mapstructure:"name"`
	Mode           string `mapstructure:"mode"`
	Version        string `mapstructure:"version"`
	Port           int    `mapstructure:"port"`
	*LogConfig     `mapstructure:"log"`
	*MySQLConfig   `mapstructure:"mysql"`
	MySQLConfigLog *MySQLConfig `mapstructure:"mysql_log"`
	*BeeGoConfig   `mapstructure:"beego"`
	*HttpServer    `mapstructure:"server"`
	*App           `mapstructure:"app"`
	*GameProxy     `mapstructure:"game_proxy"`
	*GameServer    `mapstructure:"game_server"`
}

type App struct {
	RunMode     string `mapstructure:"run_mode"`
	IdentityKey string `mapstructure:"identity_key"`
}

type HttpServer struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"db_name"`
	DbType       string `mapstructure:"db_type"`
	DbSslmode    string `mapstructure:"db_sslmode"` //postgres 数据库用
	DbPath       string `mapstructure:"db_path"`    //sqlite3 数据库用
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type BeeGoConfig struct {
	TemplateType string `mapstructure:"template_type"`
	ViewsPath    string `mapstructure:"views_path"`
	//admin用户名 此用户登录不用认证
	AdminUser string `mapstructure:"admin_user"`
	//域名白名单 白名单IP可以直接通过IP访问网站内容否则提示不能访问
	DomanWhiteList string `mapstructure:"doman_white_list"`
	//默认需要认证模块
	NotAuthPackage string `mapstructure:"not_auth_package"`
	//默认认证类型 0 不认证 1 登录认证 2 实时认证
	UserAuthType int `mapstructure:"user_auth_type"`
	//默认登录网关
	AuthGateway string `mapstructure:"auth_gateway"`
	//用户主页
	UserIndex string `mapstructure:"user_index"`
}

type GameProxy struct {
	Port       int    `mapstructure:"port"`
	MappedHost string `mapstructure:"mapped_host"` /*映射IP*/
	MappedPort int    `mapstructure:"mapped_port"` /*映射端口*/
}

type GameServer struct {
	IsSavePackage bool `mapstructure:"is_save_package"`
}

func Init() (err error) {
	viper.SetConfigName("config")  //指定配置文件名称 （不需要后缀)
	viper.SetConfigType("yaml")    //指定配置文件类型
	viper.AddConfigPath("../conf") // 指定查找配置文件的路径 （这里使用相对路径）
	viper.AddConfigPath("./conf")  // 指定查找配置文件的路径 （这里使用相对路径）
	err = viper.ReadInConfig()
	if err != nil {
		//读取配置信息失败
		fmt.Printf("viper.ReadInConfig() faild, err:%v\n", err)
		return
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		fmt.Printf("viper.Unmarshal faild, err:%v\n", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了")
		if err := viper.Unmarshal(&Conf); err != nil {
			fmt.Printf("viper.Unmarshal faild, err:%v\n", err)
		}
	})
	return err
}

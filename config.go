package main

import "github.com/BurntSushi/toml"

// Conf Conf
var Conf = &Config{}

// Config Config
type Config struct {
	DataDir    string
	GeekUsers  []GeekUser
	SaveStatic bool     // 保存静态文件
	SaveJSON   bool     // 保存json数据
	Force      bool     // 强制更新html
	Emails     []string // 发送passwd的email地址列表
	CronEntry  string   // 定时配置
	Admin      string   // adminuserpass
	HTTP       HTTP
	SMTP       SMTP
}

// GeekUser GeekUser
type GeekUser struct {
	User string
	Pass string
}

// HTTP HTTP
type HTTP struct {
	Listen    string
	BasicAuth []string
}

// SMTP SMTP
type SMTP struct {
	Host string
	Port int
	User string
	Pass string
}

// InitConfig InitConfig
func InitConfig(fpath string) {
	var err error
	if _, err = toml.DecodeFile(fpath, Conf); err != nil {
		panic(err)
	}

	if len(Conf.DataDir) == 0 {
		Conf.DataDir = "./data"
	}
}

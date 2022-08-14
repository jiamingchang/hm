package setting

import (
	"github.com/go-ini/ini"
	"log"
)

// app.ini配置文件读取初始化模块

var Conf *ini.File

type App struct {
	JwtSecret 	string
}

type Server struct {
	HttpPort       string
	RunMode        string
	PageSize       int
}

var AppSettings = &App{}
var ServerSettings = &Server{}

func init() {
	var err error
	Conf, err = ini.Load("app.ini")
	if err != nil {
		log.Fatalf("setting.go, 读取项目配置文件 'app.ini'失败: %v", err)
	}

	mapTo("app", AppSettings)
	mapTo("serve", ServerSettings)

}

func mapTo(section string, v interface{}) {
	err := Conf.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}

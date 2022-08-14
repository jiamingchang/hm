package models

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"hm/setting"
	"log"
)
// 全局数据库对象
var db *gorm.DB

func init() {
	var (
		err                          error
		dbName, user, password, host string
		// 数据库名称，数据库用户名，数据库用户密码，数据库host地址
	)

	sec, err := setting.Conf.GetSection("database")
	if err != nil {
		log.Fatal(2, "从项目配置文件'app.ini'中读取数据库相关配置信息失败: %v", err)
	}

	dbName = sec.Key("Name").String()
	user = sec.Key("User").String()
	password = sec.Key("Password").String()
	host = sec.Key("Host").String()

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName)
	fmt.Println(dsn)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{PrepareStmt: true})
	if err != nil {
		log.Fatal(2, "连接数据库失败: %v", err)
	}
	defer fmt.Println("数据库初始化成功！")

	AutoMigrate(&User{})
	AutoMigrate(&BasicInfo{})
	AutoMigrate(&MedInfoList{})
	AutoMigrate(&Prescription{})
}

func AutoMigrate(model interface{}){
	err := db.AutoMigrate(model)
	if err != nil {
		log.Fatalln(err)
	}
}
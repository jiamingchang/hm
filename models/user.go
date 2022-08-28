package models

import (
	"gorm.io/gorm"
)

// User 用户表
type User struct {
	gorm.Model
	Name string		 `json:"name" form:"name" validate:"required"`
	Password string  `json:"password" form:"password" validate:"required,max=20,min=6"`
	TrueName string  `json:"true_name" form:"true_name" validate:"required"`
	Sex int          `json:"sex" form:"sex" validate:"required"`
	Age int          `json:"age" form:"age" validate:"required"`
	Phone string	 `json:"phone" form:"phone" validate:"required,len=11"`
	Authority string `gorm:"default:user"`
}

// Login 登录账号
func Login(name, password string) (user User, err error){
	err = db.Where("name = ? AND password = ?", name, password).First(&user).Error
	return
}

// AddUser 注册账号
func AddUser(user User) (err error) {
	err = db.First(&User{}, "name", user.Name).Error
	if err != nil{
		db.Create(&user)
	}
	return
}

func GetUser(id interface{})(user User, err error){
	err = db.First(&user, id).Error
	return
}

// GetUserList 获取用户列表
func GetUserList(name string)(users []User){
	result := db.Model(&User{})
	if name != ""{
		result = result.Where("true_name", name)
	}
	result.Order("id desc").Where("authority != ?", "admin").Find(&users)
	return
}

// DeleteUser 注销账号
func (u *User)DeleteUser() error{
	return db.Delete(&User{}, u.ID).Error
}

func init(){
	user := User{
		Name: "admin",
		Password: "admin",
		TrueName: "admin",
		Sex: 1,
		Age: 18,
		Phone: "14750183889",
		Authority: "admin",
	}
	db.FirstOrCreate(&user)
}


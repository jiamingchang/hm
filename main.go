package main

import (
	"github.com/gin-gonic/gin"
	"hm/api"
	"hm/middleware"
	"log"
)

func main(){
	router := gin.Default()

	// 用户登录
	router.POST("/login", api.Login)
	// 用户注册
	router.POST("/addUser", api.AddUser)
	// ----user----
	v1 := router.Group("/user/", middleware.AuthRequired())
	{
		// 用户注销
		v1.DELETE("deleteUser", api.DeleteUser)
		// 获取用户信息
		v1.GET("getUser", api.GetUser)
		// 获取用户处方
		v1.GET("getprescriptionList", api.GetprescriptionListv1)
		// 更换用户处方状态
		v1.POST("editprescription", api.Editprescriptionv1)
	}
	// ----admin----
	v2 := router.Group("/admin/", middleware.AuthRequired(), middleware.Authority("admin"))
	{
		// 获取用户列表
		v2.GET("getUserList", api.GetUserList)
		// 获取处方列表
		v2.GET("getprescriptionList", api.GetprescriptionListv2)
		// 增加处方
		v2.POST("addprescription", api.Addprescription)
		// 更换用户处方状态
		v2.POST("editprescription", api.Editprescriptionv2)
		// 删除处方
		v2.DELETE("deleteprescription", api.Deleteprescription)
	}

	err := router.Run(":8080")
	if err != nil {
		log.Fatalln(err)
	}

}

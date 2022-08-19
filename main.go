package main

import (
	"github.com/gin-gonic/gin"
	"hm/api"
	"hm/middleware"
	"log"
)

func main(){
	router := gin.Default()

	router.GET("/chat",api.Chat)
	router.POST("/login", api.Login)
	router.POST("/addUser", api.AddUser)
	v1 := router.Group("/user/", middleware.AuthRequired())
	{
		v1.DELETE("deleteUser", api.DeleteUser)
		v1.GET("getUser", api.GetUser)
		v1.GET("getprescriptionList", api.GetprescriptionListv1)
		v1.POST("editprescription", api.Editprescriptionv1)
	}
	v2 := router.Group("/admin/", middleware.AuthRequired(), middleware.Authority("admin"))
	{
		v2.GET("getUserList", api.GetUserList)
		v2.GET("getprescriptionList", api.GetprescriptionListv2)
		v2.POST("addprescription", api.Addprescription)
		v2.POST("editprescription", api.Editprescriptionv2)
		v2.DELETE("deleteprescription", api.Deleteprescription)
	}

	err := router.Run(":8089")
	if err != nil {
		log.Fatalln(err)
	}

}

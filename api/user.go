package api

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"hm/middleware"
	"hm/models"
	"log"
	"net/http"
)

// Login 登录账号
func Login(c *gin.Context){
	var json struct{
		Name string `json:"name" form:"name" validate:"required"`
		Password string `json:"password" form:"password" validate:"required"`
	}
	if !BindAndValid(c, &json){
		return
	}

	user, err := models.Login(json.Name, json.Password)
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"message": "用户不存在或密码错误",
			"error": err.Error(),
		})
		return
	}
	token, err := middleware.GenerateToken(cast.ToString(user.ID), user.Authority)
	if err != nil {
		log.Println(err)
	}
	c.JSON(http.StatusOK,gin.H{
		"isSuccess": true,
		"token": token,
	})
}

// AddUser 注册账号
func AddUser(c *gin.Context){
	var json models.User
	if !BindAndValid(c, &json){
		return
	}
	err := models.AddUser(json)
	if err == nil {
		c.JSON(http.StatusOK,gin.H{
			"message": "用户名已存在",
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"isSuccess": true,
	})
}

// DeleteUser 注销账号
func DeleteUser(c *gin.Context){
	user := Current(c)
	err := user.DeleteUser()
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"message": "用户不存在",
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"isSuccess": true,
	})
}

// GetUser 获取用户信息
func GetUser(c *gin.Context){
	user := Current(c)
	c.JSON(http.StatusOK,gin.H{
		"user": user,
	})
}

// GetUserList 获取用户列表
func GetUserList(c *gin.Context){
	users := models.GetUserList(c.Query("true_name"))
	c.JSON(http.StatusOK,gin.H{
		"users": users,
	})
}

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"hm/middleware"
	"hm/models"
	"hm/pkg"
	"log"
	"net/http"
)

func Login(c *gin.Context){
	name := c.PostForm("name")
	password := c.PostForm("password")
	user, err := models.Login(name, password)
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

func DeleteUser(c *gin.Context){
	user := pkg.Current(c)
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

func GetUser(c *gin.Context){
	user := pkg.Current(c)
	c.JSON(http.StatusOK,gin.H{
		"user": user,
	})
}

func GetUserList(c *gin.Context){
	users := models.GetUserList(c.Query("true_name"))
	c.JSON(http.StatusOK,gin.H{
		"users": users,
	})
}

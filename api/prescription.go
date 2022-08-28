package api

import (
	"github.com/gin-gonic/gin"
	"hm/models"
	"log"
	"net/http"
)

// Addprescription 增加处方
func Addprescription(c *gin.Context) {
	var json struct{
		Prescription models.Prescription  `json:"prescription" validate:"required"`
		BasicInfo    models.BasicInfo     `json:"basic_info" validate:"required"`
		MedInfoLists []models.MedInfoList `json:"med_info_lists" validate:"required"`
	}
	if !BindAndValid(c, &json){
		return
	}

	json.Prescription.BasicInfo = json.BasicInfo
	err := json.Prescription.Addprescription(json.MedInfoLists)
	if err != nil {
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"isSuccess": true,
	})
}

// Deleteprescription 删除处方
func Deleteprescription(c *gin.Context){
	err := models.Deleteprescription(c.Query("id"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"isSuccess": true,
	})
}

// Editprescriptionv1 编辑处方（user）
func Editprescriptionv1(c *gin.Context){
	var json struct{
		Id uint `json:"id" form:"id" validate:"required"`
		State string `json:"state" form:"state" validate:"required"`
	}
	if !BindAndValid(c, &json){
		return
	}

	err := models.Editprescriptionv1(Current(c).ID, json.Id, json.State)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"isSuccess": true,
	})
}

// Editprescriptionv2 编辑处方（admin）
func Editprescriptionv2(c *gin.Context){
	var json struct{
		Id uint `json:"id" form:"id" validate:"required"`
		State string `json:"state" form:"state" validate:"required"`
	}
	if !BindAndValid(c, &json){
		return
	}

	err := models.Editprescriptionv2(json.Id, json.State)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"isSuccess": true,
	})
}

// GetprescriptionListv1 获取处方（user）
func GetprescriptionListv1(c *gin.Context){
	prescriptionList := models.GetprescriptionListv1(Current(c).ID, c.Query("state"))
	c.JSON(http.StatusOK, gin.H{
		"prescriptionList": prescriptionList,
	})
}

// GetprescriptionListv2 获取处方（admin）
func GetprescriptionListv2(c *gin.Context){
	prescriptionList := models.GetprescriptionListv2(c.Query("state"))
	c.JSON(http.StatusOK, gin.H{
		"prescriptionList": prescriptionList,
	})
}

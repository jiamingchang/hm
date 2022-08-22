package api

import (
	"github.com/gin-gonic/gin"
	"hm/models"
	"hm/pkg"
	"log"
	"net/http"
)

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

func Editprescriptionv1(c *gin.Context){
	var json struct{
		Id uint `json:"id" form:"id" validate:"required"`
		State string `json:"state" form:"state" validate:"required"`
	}
	if !BindAndValid(c, &json){
		return
	}

	err := models.Editprescriptionv1(pkg.Current(c).ID, json.Id, json.State)
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

func Editprescriptionv2(c *gin.Context){
	var json struct{
		Id uint `json:"id" form:"id" validate:"required"`
		State string `json:"state" form:"state" validate:"required"`
	}
	if !BindAndValid(c, &json){
		return
	}

	err := models.Editprescriptionv1(pkg.Current(c).ID, json.Id, json.State)
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

func GetprescriptionListv1(c *gin.Context){
	prescriptionList := models.GetprescriptionListv1(pkg.Current(c).ID, c.Query("state"))
	c.JSON(http.StatusOK, gin.H{
		"prescriptionList": prescriptionList,
	})
}
func GetprescriptionListv2(c *gin.Context){
	prescriptionList := models.GetprescriptionListv2(c.Query("state"))
	c.JSON(http.StatusOK, gin.H{
		"prescriptionList": prescriptionList,
	})
}

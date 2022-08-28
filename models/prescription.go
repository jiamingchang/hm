package models

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Prescription 药单
type Prescription struct {
	gorm.Model
	PrescriptionID string	`json:"prescription_id" validate:"required"`
	State string		    `json:"state" validate:"required"`

	BasicInfoId uint
	BasicInfo BasicInfo		`json:"basic_info" validate:"-"`
	Mifs []MedInfoList  	`gorm:"foreignKey:PrescriptionId"`
	UserId uint				`json:"user_id" validate:"required"`
	User User				`json:"user" validate:"-"`
}

// BasicInfo 药单基础信息
type BasicInfo struct {
	gorm.Model
	RecordNumber string			`json:"record_number" validate:"required"`
	PrescriptionNumber string 	`json:"prescription_number" validate:"required"`
	OpenDate string 			`json:"open_date" validate:"required"`
	BedNumber int				`json:"bed_number" validate:"required"`
	Category string				`json:"category" validate:"required"`
	ClinicalDiagnosis string	`json:"clinical_diagnosis" validate:"required"`
	AuditDoctor string			`json:"audit_doctor" validate:"required"`
	DeploymentDoctor string		`json:"deployment_doctor" validate:"required"`
	CheckDoctor string			`json:"check_doctor" validate:"required"`
	Doctor string				`json:"doctor" validate:"required"`
}

// MedInfoList 药品信息
type MedInfoList struct {
	gorm.Model
	MedID string	`json:"med_id" validate:"required"`
	MedName string	`json:"med_name" validate:"required"`
	Dosage string	`json:"dosage" validate:"required"`
	UseType string	`json:"use_type" validate:"required"`
	TotalNum string	`json:"total_num" validate:"required"`
	PrescriptionId uint
}

// Addprescription 增加处方
func (p *Prescription)Addprescription(Ids []MedInfoList) (err error){
	user, _:= GetUser(p.UserId)
	p.User = user
	db.Create(&p)
	err = db.Model(&p).Association("Mifs").Append(&Ids)
	return
}

// Deleteprescription 删除处方
func Deleteprescription(id interface{}) error{
	var prescription Prescription
	result := db.First(&prescription, id)
	if result.RowsAffected == 0{
		return errors.New("处方不存在")
	}
	db.Select(clause.Associations).Delete(&BasicInfo{}, prescription.BasicInfoId)
	db.Select(clause.Associations).Delete(&MedInfoList{}, "prescription_id", prescription.ID)
	db.Select(clause.Associations).Delete(&prescription)
	return nil
}

// Editprescriptionv1 编辑处方（user）
func Editprescriptionv1(userid, id, state interface{})(err error){
	var p Prescription
	err = db.First(&p, id).Error
	if err!=nil {
		return errors.New("更改状态失败")
	}
	if p.UserId != userid {
		return errors.New("不是你的处方")
	}
	db.Model(&p).Update("state", state)
	return nil
}

// Editprescriptionv2 编辑处方（admin）
func Editprescriptionv2(id, state interface{})(err error){
	var p Prescription
	err = db.First(&p, id).Error
	if err!=nil {
		return errors.New("更改状态失败")
	}
	db.Model(&p).Update("state", state)
	return
}

// GetprescriptionListv1 获取处方（user）
func GetprescriptionListv1(userid, state interface{}) (prescriptionList []Prescription) {
	result := db.Order("id desc")
	if state !=""{
		result = result.Where("state", state)
	}
	result.Where("user_id", userid).Preload(clause.Associations).Find(&prescriptionList)
	return
}

// GetprescriptionListv2 获取处方（admin）
func GetprescriptionListv2(state string) (prescriptionList []Prescription) {
	result := db.Order("id desc")
	if state !=""{
		result = result.Where("state", state)
	}
	result.Preload(clause.Associations).Find(&prescriptionList)
	return
}

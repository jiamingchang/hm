package models

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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

type MedInfoList struct {
	gorm.Model
	MedID string	`json:"med_id" validate:"required"`
	MedName string	`json:"med_name" validate:"required"`
	Dosage string	`json:"dosage" validate:"required"`
	UseType string	`json:"use_type" validate:"required"`
	TotalNum string	`json:"total_num" validate:"required"`
	PrescriptionId uint
}

func (p *Prescription)Addprescription(Ids []MedInfoList) (err error){
	user, _:= GetUser(p.UserId)
	p.User = user
	db.Create(&p)
	err = db.Model(&p).Association("Mifs").Append(&Ids)
	return
}

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

func Editprescriptionv2(id, state interface{})(err error){
	var p Prescription
	err = db.First(&p, id).Error
	if err!=nil {
		return errors.New("更改状态失败")
	}
	db.Model(&p).Update("state", state)
	return
}

func (b *BasicInfo)AddbasicInfo(){
	db.Create(&b)
}

func (m *MedInfoList)AddmedInfo(){
	db.Create(&m)
}

func GetprescriptionListv1(userid, state interface{}) (prescriptionList []Prescription) {
	db.Order("id desc").Where("state", state).Where("user_id", userid).Preload(clause.Associations).Find(&prescriptionList)
	return
}

func GetprescriptionListv2(state string) (prescriptionList []Prescription) {
	db.Order("id desc").Where("state", state).Preload(clause.Associations).Find(&prescriptionList)
	return
}

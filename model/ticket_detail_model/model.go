package ticket_detail_model

import (
	"github/j92z/go_kaadda/model"
	"github.com/jinzhu/gorm"
)

type TicketDetail struct {
	gorm.Model
	CompanyID      string `gorm:"index:company_id" json:"CompanyID"`
	CompanyName    string `json:"CompanyName"`
	TicketType     int    `json:"JobType"`
	TicketTypeName string `json:"JobTypeName"`
	Submitter      string `json:"Submitter"`
	SubmitterName  string `json:"SubmitterName"`
	Introduce      string `json:"Introduce"`
	Content        string `json:"Content"`
	Status         int    `json:"Status"`
	Contact        string `json:"Contact"`
}

func (m *TicketDetail) TableName() string {
	return "ticket_detail"
}

func Setup() {
	if !model.DB.HasTable(&TicketDetail{}) {
		model.DB.AutoMigrate(&TicketDetail{})
	}
}

func (m *TicketDetail) Add() error {
	return model.DB.Create(m).Error
}

func (m *TicketDetail) Update() error {
	return model.DB.Save(m).Error
}

func (m *TicketDetail) Remove() error {
	return model.DB.Delete(m).Error
}

func FindById(id int) *TicketDetail {
	var detail TicketDetail
	model.DB.Where("id = ?", id).First(&detail)
	return &detail
}

func GetList(page int, size int, status int) *[]TicketDetail {
	var list []TicketDetail
	tx := model.DB.Offset(page * size).Limit(size)
	if status >= 0 {
		tx.Where("status = ?", status)
	}
	tx.Find(&list)
	return &list
}

func GetListCount(status int) int {
	var count int
	tx := model.DB.Model(TicketDetail{})
	if status >= 0 {
		tx.Where("status = ?", status)
	}
	tx.Count(&count)
	return count
}

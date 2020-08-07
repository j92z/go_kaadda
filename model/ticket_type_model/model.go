package ticket_type_model

import (
	"github/j92z/go_kaadda/model"
	"github.com/jinzhu/gorm"
)

type TicketType struct {
	gorm.Model
	Name string `json:"Name"`
}

func (m *TicketType) TableName() string {
	return "ticket_type"
}

func Setup() {
	if !model.DB.HasTable(&TicketType{}) {
		model.DB.AutoMigrate(&TicketType{})
		(&TicketType{Name: "个人设备硬件故障"}).AddOne()
		(&TicketType{Name: "操作系统故障"}).AddOne()
		(&TicketType{Name: "系统配置支持"}).AddOne()
	}
}

func GetAll() []*TicketType {
	var list []*TicketType
	model.DB.Find(&list)
	return list
}

func (m *TicketType) AddOne() {
	model.DB.Create(m)
}

func FindById(id int) *TicketType {
	var ticketType TicketType
	model.DB.Where("id = ?", id).First(&ticketType)
	return &ticketType
}

package banner_model

import (
	"github/j92z/go_kaadda/model"
	"github.com/jinzhu/gorm"
)

type Banner struct {
	gorm.Model
	Fid   int    `json:"Fid"`
	FPath string `json:"FPath"`
	Aid   int    `json:"Aid"`
	APath string `json:"APath"`
}

func (m *Banner) TableName() string {
	return "banner"
}

func Setup() {
	if !model.DB.HasTable(&Banner{}) {
		model.DB.AutoMigrate(&Banner{})
	}
}

func (m *Banner) Add() error {
	return model.DB.Create(m).Error
}

func (m *Banner) Update() error {
	return model.DB.Save(m).Error
}

func (m *Banner) Remove() error {
	return model.DB.Delete(m).Error
}

func GetList(page int, size int) []*Banner {
	var list []*Banner
	model.DB.Offset(page * size).Limit(size).Find(&list)
	return list
}
func GetListOrder(page int, size int, order string) []*Banner {
	var list []*Banner
	model.DB.Offset(page * size).Limit(size).Order(order).Find(&list)
	return list
}

func GetListCount() int {
	var count int
	model.DB.Model(Banner{}).Count(&count)
	return count
}

func FindById(id int) *Banner {
	var info Banner
	model.DB.Where("id = ?", id).First(&info)
	return &info
}

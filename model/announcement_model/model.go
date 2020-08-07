package announcement_model

import (
	"github/j92z/go_kaadda/model"
	"github.com/jinzhu/gorm"
)

type Announcement struct {
	gorm.Model
	Title       string `json:"Title"`
	Content     string `json:"Content"`
	Creator     string `json:"Creator"`
	CreatorName string `json:"CreatorName"`
	FIds        string `json:"FIds"`
	ViewCount   int    `json:"ViewCount"`
	StarCount   int    `json:"StarCount"`
	Status      int    `json:"Status"`
}

func (m *Announcement) TableName() string {
	return "announcement"
}

func Setup() {
	if !model.DB.HasTable(&Announcement{}) {
		model.DB.AutoMigrate(&Announcement{})
	}
}

func (m *Announcement) Add() error {
	return model.DB.Create(m).Error
}

func (m *Announcement) Update() error {
	return model.DB.Save(m).Error
}

func (m *Announcement) Remove() error {
	return model.DB.Delete(m).Error
}

func GetList(page int, size int, status int) []*Announcement {
	var list []*Announcement
	tx := model.DB.Offset(page * size).Limit(size)
	if status >= 0 {
		tx = tx.Where("status = ?", status)
	}
	tx.Find(&list)
	return list
}

func GetListCount(status int) int {
	var count int
	tx := model.DB.Model(Announcement{})
	if status >= 0 {
		tx = tx.Where("status = ?", status)
	}
	tx.Count(&count)
	return count
}

func FindById(id int) *Announcement {
	var info Announcement
	model.DB.Where("id = ?", id).First(&info)
	return &info
}

func ExistById(id int) bool {
	var count int
	model.DB.Model(Announcement{}).Where("id = ?", id).Count(&count)
	return count > 0
}

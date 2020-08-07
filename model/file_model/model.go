package file_model

import (
	"github/j92z/go_kaadda/model"
	"github.com/jinzhu/gorm"
)

type File struct {
	gorm.Model
	Name     string `json:"Name"`
	Path     string `json:"Path"`
	MimeType string `json:"MimeType"`
}

func (m *File) TableName() string {
	return "file"
}

func Setup() {
	if !model.DB.HasTable(&File{}) {
		model.DB.AutoMigrate(&File{})
	}
}

func (m *File) Add() error {
	return model.DB.Create(m).Error
}

func (m *File) Remove() error {
	return model.DB.Delete(m).Error
}

func FindById(id int) *File {
	var info File
	model.DB.Where("id = ?", id).First(&info)
	return &info
}

func FindByIds(ids string) []*File {
	var list []*File
	model.DB.Where("id IN (?)", ids).First(&list)
	return list
}

func ExistById(id int) bool {
	var count int
	model.DB.Model(File{}).Where("id = ?", id).Count(&count)
	return count > 0
}

func CountByIds(ids []string) int {
	var count int
	model.DB.Model(File{}).Where("id IN (?)", ids).Count(&count)
	return count
}

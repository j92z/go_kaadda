package star_relation_model

import (
	"github/j92z/go_kaadda/model"
)

type StarRelation struct {
	ID  uint   `gorm:"primary_key"`
	Uid string `json:"Uid"`
	Aid int    `json:"Aid"`
}

func Setup() {
	if !model.DB.HasTable(&StarRelation{}) {
		model.DB.AutoMigrate(&StarRelation{})
	}
}

func (m *StarRelation) TableName() string {
	return "star_relation"
}

func (m *StarRelation) Add() error {
	return model.DB.Create(m).Error
}

func (m *StarRelation) Remove() error {
	return model.DB.Delete(m).Error
}

func FindUidByAid(aid int) *[]string {
	var list []string
	model.DB.Model(StarRelation{}).Where("aid = ?", aid).Pluck("uid", &list)
	return &list
}

func FindUidAndAid(uid string, aid int) *StarRelation {
	var info StarRelation
	model.DB.Model(StarRelation{}).Where("aid = ? AND uid = ?", aid, uid).First(&info)
	return &info
}

func ExistByUidAndAid(uid string, aid int) bool {
	var count int
	model.DB.Model(StarRelation{}).Where("aid = ? AND uid = ?", aid, uid).Count(&count)
	return count > 0
}

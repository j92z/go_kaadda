package announcement_service

import (
	"github/j92z/go_kaadda/model"
	"github/j92z/go_kaadda/model/announcement_model"
	"github/j92z/go_kaadda/model/star_relation_model"
	"errors"
)

func StarOrNot(uid string, aid int) error {
	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	announcementModel := announcement_model.FindById(aid)
	if announcementModel.ID <= 0 {
		return errors.New("找不到当前公告")
	}
	starRelationModel := star_relation_model.FindUidAndAid(uid, aid)
	if starRelationModel.ID > 0 {
		if err := tx.Delete(starRelationModel).Error; err != nil {
			return err
		}
		announcementModel.StarCount = announcementModel.StarCount - 1
	} else {
		newStarRelation := &star_relation_model.StarRelation{
			Uid: uid,
			Aid: aid,
		}
		if err := tx.Create(newStarRelation).Error; err != nil {
			return err
		}
		announcementModel.StarCount = announcementModel.StarCount + 1
	}
	if err := tx.Save(announcementModel).Error; err != nil {
		return err
	}
	return tx.Commit().Error
}

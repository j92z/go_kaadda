package model_init

import (
	"github/j92z/go_kaadda/model/announcement_model"
	"github/j92z/go_kaadda/model/banner_model"
	"github/j92z/go_kaadda/model/file_model"
	"github/j92z/go_kaadda/model/star_relation_model"
	"github/j92z/go_kaadda/model/ticket_detail_model"
	"github/j92z/go_kaadda/model/ticket_type_model"
)

func TableInit() {
	ticket_type_model.Setup()
	ticket_detail_model.Setup()
	file_model.Setup()
	announcement_model.Setup()
	banner_model.Setup()
	star_relation_model.Setup()
}

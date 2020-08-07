package main

import (
	"github/j92z/go_kaadda/model"
	"github/j92z/go_kaadda/model/model_init"
	"github/j92z/go_kaadda/setting"
)

func init() {
	setting.Setup()
	model.Setup()
}

func main() {
	model_init.TableInit()

}

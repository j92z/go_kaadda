package main

import (
	"github/j92z/go_kaadda/httpx"
	"github/j92z/go_kaadda/model"
	"github/j92z/go_kaadda/setting"
)

func init() {
	setting.Setup()
	model.Setup()
}

func main() {
	httpx.Run()

}

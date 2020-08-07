package model

import (
	"github/j92z/go_kaadda/setting"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

func Setup() {
	var err error
	DB, err = gorm.Open("sqlite3", setting.EnvSetting.SQLite.Path)
	if err != nil {
		panic("failed to connect database")
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string { //default table prefix
		return setting.EnvSetting.SQLite.TablePrefix + defaultTableName
	}
	DB.LogMode(setting.EnvSetting.SQLite.LogMode) //调试开发模式

}

func CloseDB() {
	defer DB.Close()
}

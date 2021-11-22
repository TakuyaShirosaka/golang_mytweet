package migration

import (
	"errors"
	"fmt"
	"mytweet/db"
	"mytweet/models/tweet"
	"mytweet/models/user"
)

// DbInit  DBの初期化
func DbInit() error {
	db := db.GormConnect()

	// コネクション解放
	defer db.Close()

	//構造体に基づいてテーブルを作成
	if err := db.AutoMigrate(&tweet.Tweet{}).Error; err != nil {
		fmt.Println("Unable autoMigrateDB Tweet - " + err.Error())
		return errors.New("Unable autoMigrateDB Tweet - " + err.Error())
	}

	if err := db.AutoMigrate(&user.User{}).Error; err != nil {
		fmt.Println("Unable autoMigrateDB User - " + err.Error())
		return errors.New("Unable autoMigrateDB User - " + err.Error())
	}
	return nil
}

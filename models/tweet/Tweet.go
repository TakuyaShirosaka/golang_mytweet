package tweet

import (
	"github.com/jinzhu/gorm"
	"mytweet/db"
)

type Tweet struct {
	gorm.Model
	Content string `form:"content" json:"content" binding:"required"`
}

// Insert データインサート処理
func Insert(content string) {
	db := db.GormConnect()

	defer db.Close()
	// Insert処理
	db.Create(&Tweet{Content: content})
}

// Update DB更新
func Update(id int, tweetText string) {
	db := db.GormConnect()
	var tweet Tweet
	db.First(&tweet, id)
	tweet.Content = tweetText
	db.Save(&tweet)
	db.Close()
}

// GetAll 全件取得
func GetAll() []Tweet {
	db := db.GormConnect()
	defer db.Close()
	var tweets []Tweet
	// FindでDB名を指定して取得した後、orderで登録順に並び替え
	db.Order("created_at desc").Find(&tweets)
	db.Close()

	return tweets
}

// GetOne DB一つ取得
func GetOne(id int) Tweet {
	db := db.GormConnect()
	var tweet Tweet
	db.First(&tweet, id)
	db.Close()
	return tweet
}

// Delete DB削除
func Delete(id int) {
	db := db.GormConnect()
	var tweet Tweet
	db.First(&tweet, id)
	db.Delete(&tweet)
	db.Close()
}

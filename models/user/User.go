package user

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"mytweet/db"
	"mytweet/middlewares/crypto"
)

type User struct {
	gorm.Model
	Username string `form:"username" json:"username" binding:"required" gorm:"unique;not null"`
	Password string `form:"password" json:"password" binding:"required"`
}

// CreateUser ユーザー登録処理
func CreateUser(username string, password string) []error {
	passwordEncrypt, _ := crypto.PasswordEncrypt(password)
	connect := db.GormConnect()
	defer connect.Close()
	// Insert処理
	if err := connect.Create(&User{Username: username, Password: passwordEncrypt}).GetErrors(); len(err) != 0 {
		fmt.Println("----CreateUser ユーザー登録処理----")
		fmt.Println(err)
		return err
	}
	return nil
}

// GetUser ユーザーを一件取得
func GetUser(username string) User {
	connect := db.GormConnect()
	var user User
	connect.First(&user, "username = ?", username)
	defer connect.Close()
	return user
}

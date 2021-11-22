package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" //直接的な記述が無いが、ソース内でインポートしたいものに対しては"_"を頭につける決まり
	"github.com/jinzhu/gorm"
	"mytweet/config"
)

func GormConnect() *gorm.DB {

	DBMS := "mysql"
	USER := config.GetConfig().Get("db.USER").(string)
	PASS := config.GetConfig().Get("db.PASS").(string)
	DBNAME := config.GetConfig().Get("db.DBNAME").(string)
	connectTemplate := "%s:%s@%s/%s%s"

	// MySQLだと文字コードの問題で"?parseTime=true"を末尾につける必要がある
	CONNECT := fmt.Sprintf(connectTemplate, USER, PASS, "tcp(127.0.0.1:3306)", DBNAME, "?parseTime=true")
	db, err := gorm.Open(DBMS, CONNECT)

	if err != nil {
		panic(err.Error())
	}
	return db

}


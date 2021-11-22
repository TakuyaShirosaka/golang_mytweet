package main

import (
	"flag" // コマンドライン引数
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"mytweet/config"
	"mytweet/migration"
	"mytweet/server"
	"os"
)

func main() {

	// 環境変数読込 mytweet.exe -e testな感じでコントロール
	environment := flag.String("e", "development", "")

	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}

	flag.Parse()
	config.Init(*environment)

	// ★環境変数の使用方法
	user := config.GetConfig().Get("user")          // 親要素を指定
	userName := config.GetConfig().Get("user.name") // 子要素を直接指定、ピリオド区切り
	fmt.Println(user)
	fmt.Println(userName)

	// DBの初期化
	migration.DbInit()

	router := server.Router()

	router.Run()

}

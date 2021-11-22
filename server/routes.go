package server

import (
	"github.com/gin-gonic/gin"
	"mytweet/controllers/api_tweet_controller"
	"mytweet/controllers/tweet_controller"
	"mytweet/middlewares/auth"
	"mytweet/middlewares/cors"
)

func Router() *gin.Engine {

	router := gin.Default()
	router.LoadHTMLGlob("views/*.html")

	// ここからCorsの設定
	router.Use(cors.GetCorsConfig())

	ctrl := tweet_controller.Controller{}
	apiCtrl := api_tweet_controller.Controller{}

	/** モノシリックなアプリケーション？ **/
	// 一覧 特定パスにmiddleware適用
	router.GET("/", AuthMiddleware.AuthMiddleware(), ctrl.Index)

	// グループの作成＋全体にmiddleware適用
	auth := router.Group("/").Use(AuthMiddleware.AuthMiddleware())
	{
		//登録
		auth.POST("/new", ctrl.New)

		//投稿詳細
		auth.GET("/detail/:id", ctrl.Detail)

		//更新
		auth.POST("/update/:id", ctrl.Update)

		//削除確認
		auth.GET("/delete_check/:id", ctrl.DeleteCheck)

		//削除
		auth.POST("/delete/:id", ctrl.Delete)
	}

	// ユーザー登録画面
	router.GET("/signup", func(c *gin.Context) {
		c.HTML(200, "signup.html", gin.H{})
	})

	// ユーザー登録
	router.POST("/signup", ctrl.CreateUser)

	// ユーザーログイン画面
	router.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{})
	})

	// ユーザーログイン
	router.POST("/login", ctrl.Login)
	/*******************************************************/


	/** REST API **/
	api := router.Group("/api")
	{
		api.POST("/signup", apiCtrl.CreateUser)
		api.POST("/login", apiCtrl.Login)
	}

	authApi := router.Group("/api").Use(AuthMiddleware.AuthMiddleware())
	{
		authApi.POST("/", apiCtrl.Index)

		//登録
		authApi.POST("/new", apiCtrl.New)

		//投稿詳細
		authApi.POST("/detail/:id", apiCtrl.Detail)

		//更新
		authApi.POST("/update/:id", apiCtrl.Update)

		//削除
		authApi.POST("/delete/:id", apiCtrl.Delete)
	}

	return router
}

package api_tweet_controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	awsCommon "mytweet/middlewares/aws"
	"mytweet/middlewares/crypto"
	"mytweet/models/tweet"
	"mytweet/models/user"
	"net/http"
	"strconv"
	"time"
)

type Controller struct{}

func (pc Controller) CreateUser(c *gin.Context) {
	log.Println("CreateUser")
	var body user.User
	// バリデーション処理
	if err := c.ShouldBindJSON(&body); err != nil {
		c.HTML(http.StatusInternalServerError, "signup.html", gin.H{"err": err})
		return
	} else {
		username := body.Username
		password := body.Password
		// 登録ユーザーが重複していた場合にはじく処理
		if err := user.CreateUser(username, password); err != nil {
			c.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": err})
		}
		c.JSON(http.StatusOK, gin.H{"status": "you are signup"})
	}
}

func (pc Controller) Login(c *gin.Context) {
	log.Println("Login")
	var body user.User
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Abort"})
		panic("ログインできませんでした")
	}

	// リクエストから取得したユーザー名、パスワード
	username := body.Username
	password := body.Password

	// DBから取得したユーザーパスワード(Hash)
	dbPassword := user.GetUser(username).Password

	// ユーザーパスワードの比較
	if err := crypto.CompareHashAndPassword(dbPassword, password); err != nil {
		log.Println("ログインできませんでした")
		c.JSON(http.StatusBadRequest, gin.H{"status": "Abort"})
		return
	} else {
		log.Println("ログインできました")

		cfg, err := awsCommon.GetAwsCredential()
		if err != nil {
			log.Fatalf("unable to load SDK config, %v", err)
		}

		// Using the Config value, create the DynamoDB client
		svc := dynamodb.NewFromConfig(cfg)
		uuid, err := uuid.NewRandom()
		fmt.Println("new uuid:" + uuid.String())

		// 現在時刻+1時間後
		now := time.Now().Add(1 * time.Hour).Unix()
		n := strconv.FormatInt(now, 10)
		fmt.Println("new expires:" + n)

		output, err := svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String("sessions"),
			Item: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: uuid.String()},
				"data": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
					"username": &types.AttributeValueMemberS{Value: username},
					"login":    &types.AttributeValueMemberBOOL{Value: true},
				}},
				"expires": &types.AttributeValueMemberN{Value: n},
			},
		})

		if err != nil {
			log.Fatalf("failed to PutItem, %v", err)
			return
		}

		jsonE, _ := json.Marshal(output)
		fmt.Println(string(jsonE))
		c.SetSameSite(http.SameSiteNoneMode)
		c.SetCookie("uuid", uuid.String(), 60000, "/", c.Request.URL.Hostname(), true, true)
		c.SetCookie("username", username, 60000, "/", c.Request.URL.Hostname(), true, true)
		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	}
}

func (pc Controller) Index(c *gin.Context) {
	tweets := tweet.GetAll()
	sessions, _ := c.Get("sessions")
	c.JSON(http.StatusOK, gin.H{"tweets": tweets, "sessions": sessions})
}

func (pc Controller) New(c *gin.Context) {
	var body tweet.Tweet
	// ここがバリデーション部分
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err})
	} else {
		content := body.Content
		tweet.Insert(content)
		c.JSON(http.StatusOK, gin.H{"status": "New"})
	}
}

func (pc Controller) Detail(c *gin.Context) {
	n := c.Param("id")
	id, err := strconv.Atoi(n)
	if err != nil {
		panic(err)
	}
	tweet := tweet.GetOne(id)
	c.JSON(http.StatusOK, gin.H{"tweet": tweet})
}

func (pc Controller) Update(c *gin.Context) {
	var body tweet.Tweet
	n := c.Param("id")
	id, err := strconv.Atoi(n)
	if err != nil {
		panic("ERROR")
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err})
	} else {
		postTweet := body.Content
		tweet.Update(id, postTweet)
		c.JSON(http.StatusOK, gin.H{"status": "update"})
	}
}

func (pc Controller) Delete(c *gin.Context) {
	n := c.Param("id")
	id, err := strconv.Atoi(n)
	if err != nil {
		panic("ERROR")
	}
	tweet.Delete(id)
	c.JSON(http.StatusOK, gin.H{"status": "delete"})
}

package tweet_controller

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

func (pc Controller) Index(c *gin.Context) {
	tweets := tweet.GetAll()
	sessions, _ := c.Get("sessions")
	c.HTML(200, "index.html", gin.H{"tweets": tweets, "sessions": sessions})
}

func (pc Controller) New(c *gin.Context) {
	var form tweet.Tweet
	// ここがバリデーション部分
	if err := c.Bind(&form); err != nil {
		tweets := tweet.GetAll()
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"tweets": tweets, "err": err})
		c.Abort()
	} else {
		content := c.PostForm("content")
		tweet.Insert(content)
		c.Redirect(302, "/")
	}
}

func (pc Controller) Detail(c *gin.Context) {
	n := c.Param("id")
	id, err := strconv.Atoi(n)
	if err != nil {
		panic(err)
	}
	tweet := tweet.GetOne(id)
	c.HTML(200, "detail.html", gin.H{"tweet": tweet})
}

func (pc Controller) Update(c *gin.Context) {
	n := c.Param("id")
	id, err := strconv.Atoi(n)
	if err != nil {
		panic("ERROR")
	}
	postTweet := c.PostForm("tweet")
	tweet.Update(id, postTweet)
	c.Redirect(302, "/")
}

func (pc Controller) DeleteCheck(c *gin.Context) {
	n := c.Param("id")
	id, err := strconv.Atoi(n)
	if err != nil {
		panic("ERROR")
	}
	tweet := tweet.GetOne(id)
	c.HTML(200, "delete.html", gin.H{"tweet": tweet})
}

func (pc Controller) Delete(c *gin.Context) {
	n := c.Param("id")
	id, err := strconv.Atoi(n)
	if err != nil {
		panic("ERROR")
	}
	tweet.Delete(id)
	c.Redirect(302, "/")
}

func (pc Controller) CreateUser(c *gin.Context) {
	var form user.User
	// バリデーション処理
	if err := c.Bind(&form); err != nil {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": err})
		c.Abort()
	} else {
		username := c.PostForm("username")
		password := c.PostForm("password")
		// 登録ユーザーが重複していた場合にはじく処理
		if err := user.CreateUser(username, password); err != nil {
			c.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": err})
		}
		c.Redirect(302, "/")
	}
}

func (pc Controller) Login(c *gin.Context) {
	// DBから取得したユーザーパスワード(Hash)
	dbPassword := user.GetUser(c.PostForm("username")).Password
	log.Println(dbPassword)

	// フォームから取得したユーザーパスワード
	formPassword := c.PostForm("password")

	// ユーザーパスワードの比較
	if err := crypto.CompareHashAndPassword(dbPassword, formPassword); err != nil {
		log.Println("ログインできませんでした")
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"err": err})
		c.Abort()
	} else {
		log.Println("ログインできました")

		cfg, err := awsCommon.GetAwsCredential()
		if err != nil {
			log.Fatalf("unable to load SDK config, %v", err)
		}

		// Using the Config value, create the DynamoDB client
		svc := dynamodb.NewFromConfig(cfg)

		// Build the request with its input parameters
		resp, err := svc.ListTables(context.TODO(), &dynamodb.ListTablesInput{
			Limit: aws.Int32(5),
		})
		if err != nil {
			log.Fatalf("failed to list tables, %v", err)
		}

		fmt.Println("Tables:")
		for _, tableName := range resp.TableNames {
			fmt.Println(tableName)
		}
		/** AWS TEST */

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
					"username": &types.AttributeValueMemberS{Value: c.PostForm("username")},
					"login":  &types.AttributeValueMemberBOOL{Value: true},
				}},
				"expires": &types.AttributeValueMemberN{Value: n},
			},
		})

		if err != nil {
			log.Fatalf("failed to PutItem, %v", err)
		}

		jsonE, _ := json.Marshal(output)
		fmt.Println(string(jsonE))
		//c.SetSameSite(http.SameSiteNoneMode)
		c.SetCookie("uuid", uuid.String(), 60000 , "/", c.Request.URL.Hostname(), false, true)
		c.Redirect(302, "/")
	}
}

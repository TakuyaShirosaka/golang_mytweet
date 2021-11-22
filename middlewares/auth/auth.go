package AuthMiddleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"mytweet/middlewares/aws/dynamodb"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("-----AuthMiddleware-----")
		uuid, error := c.Cookie("uuid")
		if error != nil {
			fmt.Println("this user is not login")
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.AbortWithStatus(http.StatusTemporaryRedirect)
			return
		}

		fmt.Println("uuid is:" + uuid)
		output, _ := dynamodb.GetItemById("sessions", uuid)

		fmt.Println("--- output -> item -> data -> Value -> login || username -> Value ---")
		var result map[string]interface{}
		json.Unmarshal(output, &result)

		// output -> item -> data -> Value
		fish := result["data"].(map[string]interface{})["Value"].(map[string]interface{})
		for key, value := range fish {
			// login || username -> Value , 恐らくこんな方法は実用的ではないとは思うが
			objectValue := value.(map[string]interface{})["Value"]
			fmt.Println(key, " : ", objectValue)
		}
		c.Set("sessions", string(output))
		c.Next()
	}
}

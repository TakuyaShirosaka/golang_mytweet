package dynamodb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	awsCommon "mytweet/middlewares/aws"
)

func GetItemById(tableName string, id string) ([]byte, error) {

	cfg, err := awsCommon.GetAwsCredential()
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Using the Config value, create the DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)
	output, err := svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		fmt.Printf("this seesion is disabled, %v\n", err)
		return nil, err
	}

	//
	debug, _ := json.Marshal(output)
	fmt.Println(string(debug))
	//

	jsonE, _ := json.Marshal(output.Item)

	return jsonE, nil
}

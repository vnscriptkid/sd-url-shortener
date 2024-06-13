package handlers

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
	"github.com/vnscriptkid/sd-url-shortener/core/models"
)

func MakeCreateShortURLHandler(db *dynamodb.DynamoDB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var mapping models.URLMapping
		if err := c.ShouldBindJSON(&mapping); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		mapping.CreatedAt = time.Now().Format(time.RFC3339)
		mapping.UsageCount = 0
		mapping.IsActive = true

		av, err := dynamodbattribute.MarshalMap(mapping)
		if err != nil {
			c.JSON(500, gin.H{"error": "Could not marshal mapping"})
			return
		}

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(table),
		}

		if _, err = db.PutItem(input); err != nil {
			c.JSON(500, gin.H{"error": "Could not insert item"})
			return
		}

		// Optionally, add to Redis cache
		// rdb.Set(ctx, mapping.ShortCode, mapping.OriginalURL, 0)

		c.JSON(200, mapping)
	}
}

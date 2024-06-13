package handlers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/vnscriptkid/sd-url-shortener/core/models"
)

func MakeGetOriginalURLHandler(db *dynamodb.DynamoDB, rdb *redis.Client) func(c *gin.Context) {
	return func(ctx *gin.Context) {
		shortCode := ctx.Param("shortCode")

		// Check Redis cache first
		url, err := rdb.Get(ctx, shortCode).Result()
		if err == redis.Nil {
			// Not found in cache, check DynamoDB
			result, err := db.GetItem(&dynamodb.GetItemInput{
				TableName: aws.String(table),
				Key: map[string]*dynamodb.AttributeValue{
					"shortCode": {
						S: aws.String(shortCode),
					},
				},
			})
			if err != nil || result.Item == nil {
				ctx.JSON(404, gin.H{"error": "URL not found"})
				return
			}

			var mapping models.URLMapping
			if err := dynamodbattribute.UnmarshalMap(result.Item, &mapping); err != nil {
				ctx.JSON(500, gin.H{"error": "Could not unmarshal result"})
				return
			}

			// Update Redis cache
			rdb.Set(ctx, shortCode, mapping.OriginalURL, 0)

			ctx.JSON(200, gin.H{"originalUrl": mapping.OriginalURL})
		} else if err != nil {
			ctx.JSON(500, gin.H{"error": "Redis error"})
		} else {
			ctx.JSON(200, gin.H{"originalUrl": url})
		}
	}
}

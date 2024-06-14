package handlers

import (
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
	"github.com/vnscriptkid/sd-url-shortener/core/models"
)

func MakeGetURLsByUserHandler(db *dynamodb.DynamoDB) func(c *gin.Context) {
	return func(ctx *gin.Context) {
		userId, err := strconv.Atoi(ctx.Param("userId"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId"})
			return
		}

		input := &dynamodb.QueryInput{
			TableName:              aws.String(table),
			IndexName:              aws.String("userId-createdAt-index"), // Specify the GSI
			KeyConditionExpression: aws.String("userId = :userId"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":userId": {
					N: aws.String(strconv.Itoa(userId)),
				},
			},
			ScanIndexForward: aws.Bool(false), // Sort results in ascending order
		}

		result, err := db.Query(input)
		if err != nil || len(result.Items) == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "No URLs found for user"})
			return
		}

		var mappings []models.URLMapping
		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &mappings)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not unmarshal results"})
			return
		}

		ctx.JSON(http.StatusOK, mappings)
	}

}

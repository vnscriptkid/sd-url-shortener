package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/vnscriptkid/sd-url-shortener/core/handlers"
)

var (
	db  *dynamodb.DynamoDB
	rdb *redis.Client
)

func init() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize DynamoDB
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}))
	db = dynamodb.New(sess)

	// Initialize Redis
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}

func main() {
	r := gin.Default()

	r.POST("/create", handlers.MakeCreateShortURLHandler(db))
	r.GET("/get/:shortCode", handlers.MakeGetOriginalURLHandler(db, rdb))
	r.GET("/geturls/:userId", handlers.MakeGetURLsByUserHandler(db))
	r.Run()
}

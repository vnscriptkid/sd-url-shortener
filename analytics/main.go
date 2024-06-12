package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/olivere/elastic/v7"
)

type VisitEvent struct {
	ShortURL  string `json:"short_url"`
	VisitedAt string `json:"visited_at"`
	IPAddress string `json:"ip_address"`
	Referrer  string `json:"referrer"`
	UserAgent string `json:"user_agent"`
}

func main() {
	// Initialize Kafka consumer
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to start Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Initialize Elasticsearch client
	esClient, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v", err)
	}

	// Consume messages from Kafka
	partitionConsumer, err := consumer.ConsumePartition("url_visits", 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Failed to start Kafka partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var visitEvent VisitEvent
			if err := json.Unmarshal(msg.Value, &visitEvent); err != nil {
				log.Printf("Failed to unmarshal visit event: %v", err)
				continue
			}

			_, err := esClient.Index().
				Index("visits").
				BodyJson(visitEvent).
				Do(context.Background())
			if err != nil {
				log.Printf("Failed to index visit event to Elasticsearch: %v", err)
			}

		case <-signals:
			return
		}
	}
}

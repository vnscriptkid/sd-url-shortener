package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/olivere/elastic/v7"
)

type VisitEvent struct {
	ShortURL    string   `json:"short_url"`
	OriginalURL string   `json:"original_url"`
	VisitedAt   string   `json:"visited_at"`
	IPAddress   string   `json:"ip_address"`
	Referrer    string   `json:"referrer"`
	UserAgent   string   `json:"user_agent"`
	GeoLocation GeoPoint `json:"geo_location"`
	Country     string   `json:"country"`
	Region      string   `json:"region"`
	City        string   `json:"city"`
	Browser     string   `json:"browser"`
	OS          string   `json:"os"`
	DeviceType  string   `json:"device_type"`
}

type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

var (
	urls        = []string{"abc123", "def456", "ghi789", "jkl012", "mno345"}
	originals   = []string{"http://example.com/1", "http://example.com/2", "http://example.com/3"}
	referrers   = []string{"google.com", "bing.com", "yahoo.com", "duckduckgo.com", "ask.com"}
	userAgents  = []string{"Mozilla/5.0", "Chrome/90.0", "Safari/537.36", "Edge/91.0", "Opera/76.0"}
	countries   = []string{"USA", "Canada", "UK", "Germany", "France"}
	regions     = []string{"California", "Ontario", "England", "Bavaria", "Île-de-France"}
	cities      = []string{"Los Angeles", "Toronto", "London", "Munich", "Paris"}
	browsers    = []string{"Chrome", "Firefox", "Safari", "Edge", "Opera"}
	oses        = []string{"Windows", "macOS", "Linux", "iOS", "Android"}
	deviceTypes = []string{"Desktop", "Mobile", "Tablet"}
)

func main() {
	// Initialize Elasticsearch client
	esClient, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v", err)
	}

	// Create index if it doesn't exist
	exists, err := esClient.IndexExists("visits").Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to check if index exists: %v", err)
	}
	if !exists {
		createIndex(esClient)
	}

	// Generate and save visit data
	for i := 0; i < 1000; i++ {
		visitEvent := generateVisitData()
		_, err := esClient.Index().
			Index("visits").
			BodyJson(visitEvent).
			Do(context.Background())
		if err != nil {
			log.Printf("Failed to index visit event to Elasticsearch: %v", err)
		} else {
			fmt.Printf("Indexed visit event: %+v\n\n", visitEvent)
		}
	}
}

func createIndex(client *elastic.Client) {
	mapping := `{
		"mappings": {
			"properties": {
				"short_url": { "type": "keyword" },
				"original_url": { "type": "text" },
				"visited_at": { "type": "date" },
				"ip_address": { "type": "ip" },
				"referrer": { "type": "text" },
				"user_agent": { "type": "text" },
				"geo_location": { "type": "geo_point" },
				"country": { "type": "keyword" },
				"region": { "type": "keyword" },
				"city": { "type": "keyword" },
				"browser": { "type": "keyword" },
				"os": { "type": "keyword" },
				"device_type": { "type": "keyword" }
			}
		}
	}`
	_, err := client.CreateIndex("visits").BodyString(mapping).Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to create index: %v", err)
	}
}

func generateVisitData() VisitEvent {
	shortURL := urls[rand.Intn(len(urls))]
	originalURL := originals[rand.Intn(len(originals))]
	visitedAt := time.Now().Add(-time.Duration(rand.Intn(30*24)) * time.Hour).Format(time.RFC3339)
	ipAddress := fmt.Sprintf("192.168.%d.%d", rand.Intn(256), rand.Intn(256))
	referrer := referrers[rand.Intn(len(referrers))]
	userAgent := userAgents[rand.Intn(len(userAgents))]
	geoLocation := GeoPoint{
		Lat: -90 + rand.Float64()*180,
		Lon: -180 + rand.Float64()*360,
	}
	country := countries[rand.Intn(len(countries))]
	region := regions[rand.Intn(len(regions))]
	city := cities[rand.Intn(len(cities))]
	browser := browsers[rand.Intn(len(browsers))]
	os := oses[rand.Intn(len(oses))]
	deviceType := deviceTypes[rand.Intn(len(deviceTypes))]

	return VisitEvent{
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		VisitedAt:   visitedAt,
		IPAddress:   ipAddress,
		Referrer:    referrer,
		UserAgent:   userAgent,
		GeoLocation: geoLocation,
		Country:     country,
		Region:      region,
		City:        city,
		Browser:     browser,
		OS:          os,
		DeviceType:  deviceType,
	}
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	_ "github.com/lib/pq"
)

// Database connection details
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "postgres"
)

type TicketServer struct {
	db *sql.DB
}

func NewTicketServer(db *sql.DB) *TicketServer {
	return &TicketServer{db: db}
}

func (ts *TicketServer) FetchNextUniqueRandomIntWithLock() (int, error) {
	tx, err := ts.db.Begin()
	if err != nil {
		return 0, err
	}

	var id, currentNumber, startRange, endRange int
	// PostgreSQL will wait for the lock to be released by the transaction that currently holds it
	// By default, PostgreSQL will wait indefinitely for the lock to be released
	err = tx.QueryRow(`
		WITH selected_range AS (
			SELECT id, start_range, end_range, current_number
			FROM int_ranges
			ORDER BY RANDOM()
			LIMIT 1
			FOR UPDATE
		)
		UPDATE int_ranges
		SET current_number = selected_range.current_number + 1
		FROM selected_range
		WHERE int_ranges.id = selected_range.id
		RETURNING int_ranges.id, int_ranges.current_number, int_ranges.start_range, int_ranges.end_range
	`).Scan(&id, &currentNumber, &startRange, &endRange)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Ensure the current number does not exceed the range
	if currentNumber > endRange {
		tx.Rollback()
		return 0, fmt.Errorf("range exceeded for range id %d", id)
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return currentNumber, nil
}

func (ts *TicketServer) FetchNextUniqueRandomIntWithSkipLock() (int, error) {
	const maxRetries = 10
	var id, currentNumber, startRange, endRange int

	for retries := 0; retries < maxRetries; retries++ {
		tx, err := ts.db.Begin()
		if err != nil {
			return 0, err
		}

		// `SKIP LOCKED` allows you to fetch rows that are not currently locked by other transactions,
		// which can improve concurrency and throughput in high-contention environments.
		err = tx.QueryRow(`
			WITH selected_range AS (
				SELECT id, start_range, end_range, current_number
				FROM int_ranges
				ORDER BY RANDOM()
				LIMIT 1
				FOR UPDATE SKIP LOCKED
			)
			UPDATE int_ranges
			SET current_number = selected_range.current_number + 1
			FROM selected_range
			WHERE int_ranges.id = selected_range.id
			RETURNING int_ranges.id, int_ranges.current_number, int_ranges.start_range, int_ranges.end_range
		`).Scan(&id, &currentNumber, &startRange, &endRange)

		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				// If no rows are found, retry after a short delay
				time.Sleep(10 * time.Millisecond)
				continue
			}
			return 0, err
		}

		// Ensure the current number does not exceed the range
		if currentNumber > endRange {
			tx.Rollback()
			return 0, fmt.Errorf("range exceeded for range id %d", id)
		}

		if err = tx.Commit(); err != nil {
			return 0, err
		}

		return currentNumber, nil
	}

	return 0, fmt.Errorf("max retries exceeded, could not fetch a unique random integer")
}

func main() {
	// Connect to the PostgreSQL database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)

	// FIX: pq: sorry, too many clients already
	db.SetMaxOpenConns(20) // Adjust based on server capacity
	db.SetMaxIdleConns(5)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("> connected to postgres")

	ticketServer := NewTicketServer(db)

	wg := sync.WaitGroup{}

	// Channel to limit the number of concurrent workers
	numWorkers := 20
	workerChan := make(chan struct{}, numWorkers)

	// Count number of successes and failures using atomic counters
	var successCount, failureCount atomic.Uint32

	startTime := time.Now()
	for i := 0; i < 500; i++ {
		wg.Add(1)

		// This effectively "occupies" a slot in the channel.
		// If the channel is full (i.e., 10 goroutines are already running), this operation will block until a slot is available.
		workerChan <- struct{}{}

		go func() {
			defer wg.Done()

			// ensures that when the goroutine completes, it removes an empty struct from workerChan, freeing up a slot for another goroutine to run.
			defer func() { <-workerChan }()

			uniqueID, err := ticketServer.FetchNextUniqueRandomIntWithLock()
			if err != nil {
				fmt.Println("Error fetching unique random integer:", err)
				failureCount.Add(1)
				return
			}
			fmt.Printf("Fetched unique random integer: %d\n", uniqueID)
			successCount.Add(1)
		}()
	}

	wg.Wait()
	fmt.Printf("All goroutines finished\n")
	fmt.Printf("Successes: %d, Failures: %d\n", successCount.Load(), failureCount.Load())
	fmt.Printf("Time taken: %v\n", time.Since(startTime))
}

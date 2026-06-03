package main

import (
	"log"
	"time"

	"github.com/felipedenardo/go-worker-pool-lab/internal/orders"
)

const workerCount = 2

func main() {
	start := time.Now()

	batch := []orders.Order{
		{ID: 1, Price: 100},
		{ID: 2, Price: 250},
		{ID: 3, Price: 75},
	}

	results := orders.Process(batch, workerCount)

	for _, r := range results {
		log.Printf("FINAL: %+v", r)
	}
	log.Printf("elapsed: %s", time.Since(start))
}

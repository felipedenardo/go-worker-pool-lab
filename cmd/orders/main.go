package main

import (
	"context"
	"log"
	"time"

	"github.com/felipedenardo/go-worker-pool-lab/internal/orders"
)

const (
	workerCount = 2
	jobBuffer   = 1
)

func main() {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	batch := []orders.Order{
		{ID: 1, Price: 100},
		{ID: 2, Price: 250},
		{ID: 3, Price: 75},
		{ID: 4, Price: 0},
	}

	results := orders.Process(ctx, batch, workerCount, jobBuffer)

	for _, r := range results {
		if r.Err != nil {
			log.Printf("pedido %d: ERRO: %v", r.OrderID, r.Err)
			continue
		}
		log.Printf("FINAL: %+v", r)
	}
	log.Printf("elapsed: %s", time.Since(start))
}

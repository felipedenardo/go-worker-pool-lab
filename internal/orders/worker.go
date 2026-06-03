package orders

import (
	"context"
	"log"
	"sync"
)

// jobs <-chan Job: o worker só RECEBE (lê) daqui
// results chan<- Result: o worker só ENVIA (escreve) aqui
func worker(ctx context.Context, id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	// `range` num canal vazio BLOQUEIA e espera; encerra quando jobs é fechado.
	for job := range jobs {
		log.Printf("worker %d: pedido %d", id, job.Order.ID)
		results <- job.run(ctx)
	}
}

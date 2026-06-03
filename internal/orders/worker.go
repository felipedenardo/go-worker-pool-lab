package orders

import (
	"log"
	"sync"
)

// jobs <-chan Job: o worker só RECEBE (lê) daqui
// results chan<- jobResult: o worker só ENVIA (escreve) aqui
func worker(id int, jobs <-chan Job, results chan<- jobResult, wg *sync.WaitGroup) {
	for job := range jobs {
		log.Printf("worker %d: %s order %d", id, job.Type, job.Order.ID)
		results <- job.run()
		wg.Done()
	}
}

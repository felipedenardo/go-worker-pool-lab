package orders

import (
	"context"
	"sync"
)

// Process roda um worker pool sobre orders e devolve os Result na ordem de entrada.
//
// Lab: jobBuffer é exposto só como exemplo de como o buffer dos channels poderia
// ser usado. Como orders já está em memória, um channel unbuffered resolveria —
// o buffer aqui serve pra demonstrar o desacoplamento entre producer e workers.
func Process(ctx context.Context, orders []Order, workerCount, jobBuffer int) []Result {
	// jobs = esteira de ENTRADA;
	//results = esteira de SAÍDA.
	jobs := make(chan Job, jobBuffer)
	results := make(chan Result, jobBuffer)

	// A concorrência é limitada a workerCount.
	var wg sync.WaitGroup
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker(ctx, i, jobs, results, &wg)
	}

	// Em goroutine pra rodar em paralelo com o coletor de results (abaixo) e não
	// travar: enfileirar jobs e drenar results no mesmo fluxo daria deadlock.
	go producer(ctx, orders, jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	// Fan-in: os resultados chegam fora de ordem, então coletamos por OrderID e reordenamos pela ordem de entrada.
	byID := make(map[uint]Result, len(orders))
	for r := range results {
		byID[r.OrderID] = r
	}

	out := make([]Result, 0, len(orders))
	for _, o := range orders {
		if r, ok := byID[o.ID]; ok {
			out = append(out, r)
		} else {
			out = append(out, Result{OrderID: o.ID, Err: ctx.Err()})
		}
	}
	return out
}

// producer cria os Job a partir dos orders e os enfileira no channel jobs.
// jobs é chan<- Job (só-envio): explicita que o producer apenas ESCREVE.
//
// Use quando a geração dos jobs for lenta/IO-bound (ex.: buscar no banco) e você
// quer overlap com o processamento dos workers. Neste lab é só exemplo: orders
// já está em memória, então o select com ctx.Done() abaixo nunca trava de fato.
func producer(ctx context.Context, orders []Order, jobs chan<- Job) {
	defer close(jobs)
	for _, o := range orders {
		// select: torna o envio cancelável — se o ctx morrer, desiste em vez de travar.
		select {
		case jobs <- Job{Order: o}:
		case <-ctx.Done():
			return
		}
	}
}

package orders

import (
	"context"
	"sync"
)

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

	// (goroutine à parte pra não travar o consumo de results — se não: deadlock)
	// select: torna o envio CANCELÁVEL: tenta enfileirar, mas se o ctx for
	// cancelado primeiro, desiste e sai em vez de travar esperando enviar.
	go func() {
		defer close(jobs)
		for _, o := range orders {
			select {
			case jobs <- Job{Order: o}:
			case <-ctx.Done():
				return
			}
		}
	}()

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

package orders

import "sync"

func Process(orders []Order, workerCount int) []Result {
	totalJobs := len(orders) * len(jobTypes)

	// jobs  = esteira de ENTRADA (trabalho a fazer).
	// results = esteira de SAÍDA (trabalho já feito).
	// O buffer = totalJobs faz os envios nunca bloquearem.
	jobs := make(chan Job, totalJobs)
	results := make(chan jobResult, totalJobs)

	// workerCount goroutines em paralelo.
	// "2 JOBS por vez" — quem estiver livre pega o próximo job do canal.
	// Iniciamos antes de enfileirar: eles só dormem no `range jobs` até chegar trabalho.
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		go worker(i, jobs, results, &wg)
	}

	// Fan-out: cada Order vira um Job por JobType, todos enfileirados no canal.
	// close(jobs) sinaliza "acabou o trabalho" e encerra o `range jobs` dos workers.
	wg.Add(totalJobs)
	for _, o := range orders {
		for _, t := range jobTypes {
			jobs <- Job{Order: o, Type: t}
		}
	}
	close(jobs)

	// Goroutine: espera todos os jobs e fecha results (encerra o range abaixo) sem bloquear o fan-in.
	go func() {
		wg.Wait()
		close(results)
	}()

	// Fan-in: os resultados de cada job chegam fora de ordem, então reagrupamos por OrderID.
	merged := make(map[uint]*Result, len(orders))
	for _, o := range orders {
		merged[o.ID] = &Result{OrderID: o.ID}
	}

	// OrderID = qual pedido; Type = qual campo preencher.
	// Assim validate e price do mesmo pedido caem no mesmo Result.
	for r := range results {
		agg := merged[r.OrderID]
		switch r.Type {
		case ValidateJob:
			agg.Valid = r.Valid
		case PriceJob:
			agg.Price = r.Price
		}
	}

	out := make([]Result, 0, len(orders))
	for _, o := range orders {
		out = append(out, *merged[o.ID])
	}
	return out
}

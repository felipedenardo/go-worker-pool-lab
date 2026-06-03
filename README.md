# go-worker-pool-lab

Laboratório de estudo de **concorrência em Go**. Implementa um **worker pool**
que processa um lote de pedidos (`Order`) em paralelo, explorando na prática:

- **goroutines** — workers rodando concorrentemente;
- **canais com buffer** — uma esteira de entrada (`jobs`) e uma de saída (`results`);
- **fan-out / fan-in** — distribuir o trabalho entre vários workers e depois reagrupar os resultados;
- **`sync.WaitGroup`** — saber quando todo o trabalho terminou para fechar os canais.

## Padrões usados

| Padrão | Onde | Para que serve |
|---|---|---|
| **Worker pool** | `pool.go` (loop que cria `workerCount` goroutines) | Limitar quantas tarefas rodam ao mesmo tempo, sem estourar recursos (memória, conexões, rate limit). |
| **Fan-out** | `pool.go` (loop que enfileira os `Job` no canal) | Espalhar o trabalho de uma fonte para vários workers em paralelo. |
| **Fan-in** | `pool.go` (`merged` + `range results`) | Juntar os resultados parciais (que chegam fora de ordem) num resultado por pedido. |
| **Canais direcionais** | `worker.go` (`<-chan` / `chan<-`) | Comunicar entre goroutines sem lock; a direção trava o uso errado em tempo de compilação. |
| **Canal com buffer** | `pool.go` (`make(chan ..., totalJobs)`) | Desacoplar produtor e consumidor; os envios não bloqueiam. |
| **`sync.WaitGroup`** | `pool.go` (`Add`/`Done`/`Wait`) | Esperar todos os jobs terminarem antes de fechar `results`. |

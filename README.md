# go-worker-pool-lab

Laboratório de estudo de **concorrência em Go**. Implementa um **worker pool**
que processa um lote de pedidos (`Order`) em paralelo, explorando na prática:

- **goroutines** — workers rodando concorrentemente;
- **canais (com e sem buffer)** — uma esteira de entrada (`jobs`) e uma de saída (`results`);
- **fan-out / fan-in** — distribuir o trabalho entre vários workers e depois reagrupar os resultados;
- **`sync.WaitGroup`** — saber quando todo o trabalho terminou para fechar os canais;
- **`context`** — cancelamento e timeout;
- **backpressure** — segurar o produtor quando os workers estão ocupados.

## Padrões usados

| Padrão | Onde | Para que serve |
|---|---|---|
| **Worker pool** | `pool.go` (loop que cria `workerCount` goroutines) | Limitar quantas tarefas rodam ao mesmo tempo, sem estourar recursos (memória, conexões, rate limit). |
| **Fan-out** | `pool.go` (produtor enfileira os `Job` no canal) | Espalhar o trabalho de uma fonte para vários workers em paralelo. |
| **Fan-in** | `pool.go` (`byID` + `range results`) | Coletar os resultados (que chegam fora de ordem) e reordená-los pela ordem de entrada. |
| **Canais direcionais** | `worker.go` (`<-chan` / `chan<-`) | Comunicar entre goroutines sem lock; a direção trava o uso errado em tempo de compilação. |
| **Backpressure** | `pool.go` (`jobBuffer` pequeno + produtor em goroutine) | Com buffer pequeno o envio bloqueia quando os workers estão saturados — controla o fluxo. |
| **`context`** | `pool.go` / `job.go` (`select` com `ctx.Done()`) | Cancelar/dar timeout: o produtor para de enfileirar e os jobs em andamento abortam. |
| **Tratamento de erro** | `job.go` / `Result.Err` | Cada pedido carrega seu próprio erro; uma falha não derruba o lote inteiro. |
| **`sync.WaitGroup`** | `pool.go` (`Add`/`Done`/`Wait`) | Esperar todos os workers terminarem antes de fechar `results`. |

> Concorrência aqui ajuda porque o trabalho é **I/O-bound** (o `time.Sleep` simula
> uma chamada de rede): os workers sobrepõem os tempos de espera. Concorrência melhora
> principalmente **throughput**.

## Rodar

```bash
go run ./cmd/orders          # executa o lab
go test -race ./...          # testa com o detector de race
```

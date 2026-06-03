# go-worker-pool-lab

Laboratório de estudo de **concorrência em Go**. Implementa um **worker pool**
que processa um lote de pedidos (`Order`) em paralelo, explorando na prática:

- **goroutines** — workers rodando concorrentemente;
- **canais com buffer** — uma esteira de entrada (`jobs`) e uma de saída (`results`);
- **fan-out / fan-in** — distribuir o trabalho entre vários workers e depois reagrupar os resultados;
- **`sync.WaitGroup`** — saber quando todo o trabalho terminou para fechar os canais.

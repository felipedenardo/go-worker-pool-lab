package orders

import (
	"context"
	"fmt"
	"time"
)

const processingDelay = 100 * time.Millisecond

const priceDiscount = 5

type Job struct {
	Order Order
}

// run processa um pedido de ponta a ponta (valida e precifica). Respeita o
// cancelamento via ctx e devolve um Result que pode carregar um erro.
func (j Job) run(ctx context.Context) Result {
	res := Result{OrderID: j.Order.ID}

	// Simula trabalho de I/O, mas aborta na hora se o ctx for cancelado.
	select {
	case <-time.After(processingDelay):
	case <-ctx.Done():
		res.Err = ctx.Err()
		return res
	}

	if j.Order.Price <= 0 {
		res.Err = fmt.Errorf("pedido %d: preço inválido (%.2f)", j.Order.ID, j.Order.Price)
		return res
	}
	res.Valid = true
	res.Price = j.Order.Price - priceDiscount
	return res
}

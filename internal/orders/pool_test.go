package orders

import (
	"context"
	"testing"
)

func TestProcessAggregatesEveryOrder(t *testing.T) {
	batch := []Order{
		{ID: 1, Price: 100},
		{ID: 2, Price: 250},
		{ID: 3, Price: 75},
	}

	got := Process(context.Background(), batch, 2, 1)

	if len(got) != len(batch) {
		t.Fatalf("retornou %d resultados, esperava %d", len(got), len(batch))
	}

	for i, o := range batch {
		r := got[i]
		if r.OrderID != o.ID {
			t.Errorf("resultado %d: OrderID = %d, esperava %d (a ordem deve ser preservada)", i, r.OrderID, o.ID)
		}
		if r.Err != nil {
			t.Errorf("pedido %d: erro inesperado: %v", o.ID, r.Err)
		}
		if !r.Valid {
			t.Errorf("pedido %d: Valid = false, esperava true", o.ID)
		}
		if want := o.Price - priceDiscount; r.Price != want {
			t.Errorf("pedido %d: Price = %v, esperava %v", o.ID, r.Price, want)
		}
	}
}

func TestProcessInvalidOrder(t *testing.T) {
	got := Process(context.Background(), []Order{{ID: 1, Price: 0}}, 2, 1)

	if len(got) != 1 {
		t.Fatalf("retornou %d resultados, esperava 1", len(got))
	}
	if got[0].Err == nil {
		t.Errorf("esperava erro para preço inválido, veio nil")
	}
	if got[0].Valid {
		t.Errorf("pedido inválido não deveria ter Valid = true")
	}
}

func TestProcessRespectsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancela antes de processar

	got := Process(ctx, []Order{{ID: 1, Price: 100}, {ID: 2, Price: 250}}, 2, 1)

	for _, r := range got {
		if r.Err == nil {
			t.Errorf("pedido %d: esperava erro de cancelamento, veio nil", r.OrderID)
		}
	}
}

func TestProcessEmptyBatch(t *testing.T) {
	if got := Process(context.Background(), nil, 2, 1); len(got) != 0 {
		t.Fatalf("retornou %d resultados, esperava 0", len(got))
	}
}

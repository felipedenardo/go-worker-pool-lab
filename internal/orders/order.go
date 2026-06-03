package orders

type Order struct {
	ID    uint
	Price float64
}

type Result struct {
	OrderID uint
	Valid   bool
	Price   float64
	Err     error
}

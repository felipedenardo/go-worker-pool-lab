package orders

import "time"

const processingDelay = 100 * time.Millisecond

const priceDiscount = 5

type JobType string

const (
	ValidateJob JobType = "validate"
	PriceJob    JobType = "price"
)

var jobTypes = []JobType{ValidateJob, PriceJob}

type Job struct {
	Order Order
	Type  JobType
}

type jobResult struct {
	OrderID uint
	Type    JobType
	Valid   bool
	Price   float64
}

// run é onde mora a regra de negócio de cada job: o único lugar que sabe COMO
// validar e COMO precificar. O worker não conhece essa lógica, só chama run().
func (j Job) run() jobResult {
	time.Sleep(processingDelay)

	res := jobResult{OrderID: j.Order.ID, Type: j.Type}
	switch j.Type {
	case ValidateJob:
		res.Valid = true
	case PriceJob:
		res.Price = j.Order.Price - priceDiscount
	}
	return res
}

package order

type Order struct {
	ID          int64
	UserID      int64
	Status      string
	TotalAmount float64
}
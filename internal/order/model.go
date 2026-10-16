package order

type Order struct {
	ID          int64
	UserID      int64
	Status      string
	TotalAmount float64
}

type OrderItem struct {
	ID        int64
	OrderID   int64
	ProductID int64
	Quantity  int
	UnitPrice float64
}
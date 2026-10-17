package order

type Order struct {
    ID          int64   `json:"id"`
    UserID      int64   `json:"user_id"`
    Status      string  `json:"status"`
    TotalAmount float64 `json:"total_amount"`
}

type OrderItem struct {
	ID        int64
	OrderID   int64
	ProductID int64
	Quantity  int
	UnitPrice float64
}
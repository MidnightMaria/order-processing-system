package httpapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"order-processing-system/internal/order"
)

type OrderService interface {
	CreateOrder(
		ctx context.Context,
		userID int64,
		productID int64,
		quantity int,
	) (*order.Order, error)
}

type OrderHandler struct {
	orderService OrderService
}

func NewOrderHandler(orderService OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

type createOrderRequest struct {
	UserID    int64 `json:"user_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Endpoint ini hanya menerima POST.
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Batasi ukuran request body agar payload tidak terlalu besar.
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MiB

	var req createOrderRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "invalid JSON request", http.StatusBadRequest)
		return
	}

	// Tolak body yang berisi lebih dari satu JSON object.
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		http.Error(w, "request body must contain a single JSON object", http.StatusBadRequest)
		return
	}

	// Validasi field wajib sebelum memanggil service.
	if req.UserID <= 0 || req.ProductID <= 0 || req.Quantity <= 0 {
		http.Error(w, "user_id, product_id, and quantity must be positive", http.StatusBadRequest)
		return
	}

	createdOrder, err := h.orderService.CreateOrder(
		r.Context(),
		req.UserID,
		req.ProductID,
		req.Quantity,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdOrder); err != nil {
		// Response sudah mulai dikirim, jadi jangan mencoba menulis
		// response error kedua di sini.
		return
	}
}

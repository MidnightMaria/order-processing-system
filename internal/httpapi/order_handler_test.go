package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"order-processing-system/internal/order"
)

type mockOrderService struct {
	createdOrder *order.Order
	err          error
}

func (m *mockOrderService) CreateOrder(
	ctx context.Context,
	userID int64,
	productID int64,
	quantity int,
) (*order.Order, error) {
	return m.createdOrder, m.err
}

func TestCreateOrderRejectsInvalidRequests(t *testing.T) {
	handler := NewOrderHandler(nil)

	tests := []struct {
		name           string
		method         string
		body           string
		expectedStatus int
	}{
		{
			name:           "reject GET request",
			method:         http.MethodGet,
			body:           "",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "reject malformed JSON",
			method:         http.MethodPost,
			body:           `{"user_id":`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "reject zero quantity",
			method:         http.MethodPost,
			body:           `{"user_id":1,"product_id":2,"quantity":0}`,
			expectedStatus: http.StatusBadRequest,
		},

		{
			name:           "reject unknown JSON field",
			method:         http.MethodPost,
			body:           `{"user_id":1,"product_id":2,"quantity":1,"unknown":true}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "reject multiple JSON objects",
			method:         http.MethodPost,
			body:           `{"user_id":1,"product_id":2,"quantity":1} {"user_id":1,"product_id":2,"quantity":1}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "reject negative user ID",
			method:         http.MethodPost,
			body:           `{"user_id":-1,"product_id":2,"quantity":1}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				tt.method,
				"/orders",
				strings.NewReader(tt.body),
			)
			rec := httptest.NewRecorder()

			handler.CreateOrder(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf(
					"expected status %d, got %d; body: %s",
					tt.expectedStatus,
					rec.Code,
					rec.Body.String(),
				)
			}
		})
	}
}

func TestCreateOrderReturnsJSON(t *testing.T) {
	mockService := &mockOrderService{
		createdOrder: &order.Order{
			ID:          123,
			UserID:      1,
			Status:      "pending",
			TotalAmount: 25000,
		},
	}

	handler := NewOrderHandler(mockService)

	body := `{"user_id":1,"product_id":2,"quantity":1}`
	req := httptest.NewRequest(
		http.MethodPost,
		"/orders",
		strings.NewReader(body),
	)
	rec := httptest.NewRecorder()

	handler.CreateOrder(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d; body: %s",
			http.StatusCreated,
			rec.Code,
			rec.Body.String(),
		)
	}

	var got order.Order
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	if got.ID != 123 {
		t.Errorf("expected ID 123, got %d", got.ID)
	}

	if got.UserID != 1 {
		t.Errorf("expected UserID 1, got %d", got.UserID)
	}

	if got.Status != "pending" {
		t.Errorf("expected status pending, got %q", got.Status)
	}

	if got.TotalAmount != 25000 {
		t.Errorf("expected total amount 25000, got %f", got.TotalAmount)
	}

	if gotContentType := rec.Header().Get("Content-Type"); gotContentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", gotContentType)
	}
}

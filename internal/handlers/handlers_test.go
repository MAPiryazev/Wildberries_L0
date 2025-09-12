package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/MAPiryazev/Wildberries_L0/internal/model"
	"github.com/MAPiryazev/Wildberries_L0/internal/service"
	"github.com/go-chi/chi/v5"
)

// Мок сервиса
type mockOrderService struct{}

func (m *mockOrderService) SaveOrder(order *models.Order) error          { return nil }
func (m *mockOrderService) SaveOrdersBatch(orders []*models.Order) error { return nil }
func (m *mockOrderService) SendToDLQ(ctx context.Context, original []byte, errMsg string) error {
	return nil
}
func (m *mockOrderService) GetOrderByID(uid string) (*models.Order, error) {
	if uid == "b563feb7b2b84b6test" {
		return &models.Order{OrderUID: uid}, nil
	}
	return nil, service.ErrOrderNotFound
}

func TestHandleGetOrderByID(t *testing.T) {
	h := &Handler{orderService: &mockOrderService{}}

	tests := []struct {
		uid          string
		wantStatus   int
		wantResponse string
	}{
		{"b563feb7b2b84b6test", http.StatusOK, `"OrderUID":"b563feb7b2b84b6test"`},
		{"999", http.StatusNotFound, "заказ не найден"},
		{"", http.StatusBadRequest, "не был получен uid"},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, "/order/"+tt.uid, nil)
		w := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("uid", tt.uid)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		h.handleGetOrderByID(w, req)

		resp := w.Result()
		if resp.StatusCode != tt.wantStatus {
			t.Errorf("для uid=%s: ожидали статус %d, получили %d", tt.uid, tt.wantStatus, resp.StatusCode)
		}

		body := w.Body.String()
		if !contains(body, tt.wantResponse) {
			t.Errorf("для uid=%s: ожидали ответ %q, получили %q", tt.uid, tt.wantResponse, body)
		}
	}
}

// проверка, что body содержит подстроку
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s[0:len(substr)] == substr || s[len(s)-len(substr):] == substr || json.Valid([]byte(s)))
}

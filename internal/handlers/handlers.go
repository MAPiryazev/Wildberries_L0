package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/MAPiryazev/Wildberries_L0/internal/service"
)

type Handler struct {
	orderService service.OrderService
}

func NewHandler(orderService service.OrderService) *Handler {
	return &Handler{orderService: orderService}
}

func RoutesInit(handler *Handler) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/order/{uid}", handler.handleGetOrderByID)

	return router
}

func (handler *Handler) handleGetOrderByID(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		http.Error(w, "не был получен uid для поиска заказа", http.StatusBadRequest)
		return
	}

	order, err := handler.orderService.GetOrderByID(uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

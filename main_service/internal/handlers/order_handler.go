package handlers

import (
	"encoding/json"
	"log"
	"main_service/internal/grpcclient"
	"main_service/internal/models"
	"main_service/internal/storage"
	"net/http"
	"strconv"
)

type OrderHandler struct {
	storage *storage.PostgresStorage
}

func NewOrderHandler(storage *storage.PostgresStorage) *OrderHandler {
	return &OrderHandler{storage: storage}
}

func (h *OrderHandler) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	eventIDStr := r.PathValue("event_id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Некорректный event_id", http.StatusBadRequest)
		return
	}

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}

	// Вызов gRPC-клиента
	resp, err := grpcclient.CallCreateOrder(eventID, order.Zone, order.Row, order.Seat, order.Email)
	if err != nil {
		http.Error(w, "Ошибка при создании заказа", http.StatusInternalServerError)
		return
	}

	log.Printf("Заказ успешно создан: %s", resp.Status)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": resp.Status})
}

func (h *OrderHandler) GetAvailableSeatsHandler(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.PathValue("event_id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Некорректный event_id", http.StatusBadRequest)
		return
	}

	resp, err := grpcclient.CallGetAvailableSeats(eventID)
	if err != nil {
		http.Error(w, "Ошибка получения мест", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

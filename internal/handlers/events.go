package handlers

import (
	"encoding/json"
	"net/http"
	"ticket-booking/internal/models"
	"ticket-booking/internal/storage"
	"time"
)

type EventHandler struct {
	storage *storage.PostgresStorage
}

func NewEventHandler(storage *storage.PostgresStorage) *EventHandler {
	return &EventHandler{storage: storage}
}

func (h *EventHandler) AddEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var event models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}

	h.storage.AddEvent(event)
	w.WriteHeader(http.StatusCreated)
}

func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	queryDate := r.URL.Query().Get("date")

	var events []models.Event
	if queryDate != "" {
		date, err := time.Parse("2006-01-02", queryDate)
		if err != nil {
			http.Error(w, "Некорректный формат даты", http.StatusBadRequest)
			return
		}
		events, _ = h.storage.GetEventsByDate(date)
	} else {
		events, _ = h.storage.GetAllEvents()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

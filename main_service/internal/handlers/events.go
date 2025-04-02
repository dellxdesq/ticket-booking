package handlers

import (
	"encoding/json"
	"log"
	"main_service/internal/models"
	"main_service/internal/storage"
	"net/http"
	"time"
)

type EventHandler struct {
	storage *storage.PostgresStorage
}

func NewEventHandler(storage *storage.PostgresStorage) *EventHandler {
	return &EventHandler{storage: storage}
}

func (h *EventHandler) AddEvent(w http.ResponseWriter, r *http.Request) {
	var rawEvent struct {
		Type     string `json:"type"`
		Title    string `json:"title"`
		DateTime string `json:"dateTime"` // Получаем дату и время в одном поле
		Tickets  int    `json:"tickets"`
	}

	if err := json.NewDecoder(r.Body).Decode(&rawEvent); err != nil {
		log.Println("Ошибка декодирования JSON:", err)
		http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}

	parsedDateTime, err := time.Parse("2006-01-02T15:04", rawEvent.DateTime) // Парсим дату и время
	if err != nil {
		log.Println("Ошибка парсинга даты и времени:", err)
		http.Error(w, "Некорректный формат даты и времени", http.StatusBadRequest)
		return
	}

	event := models.Event{
		Type:     rawEvent.Type,
		Title:    rawEvent.Title,
		DateTime: parsedDateTime, // Сохраняем в БД как time.Time
		Tickets:  rawEvent.Tickets,
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

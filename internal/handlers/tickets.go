package handlers

import (
	"encoding/json"
	"net/http"
	"ticket-booking/internal/models"
	"ticket-booking/internal/storage"
)

type TicketHandler struct {
	storage *storage.MemoryStorage
}

func NewTicketHandler(storage *storage.MemoryStorage) *TicketHandler {
	return &TicketHandler{storage: storage}
}

func (h *TicketHandler) AddTicketTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var template models.TicketTemplate
	if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
		http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}

	event, found := h.storage.GetEventByID(template.EventID)
	if !found {
		http.Error(w, "Мероприятие не найдено", http.StatusNotFound)
		return
	}

	template.Title = event.Title
	template.EventDate = event.Date

	h.storage.AddTicketTemplate(template)
	w.WriteHeader(http.StatusCreated)
}

func (h *TicketHandler) GetTickets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	tickets := h.storage.GetAllTickets()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tickets)
}

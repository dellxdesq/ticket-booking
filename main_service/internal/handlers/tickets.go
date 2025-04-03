package handlers

import (
	"encoding/json"
	"log"
	"main_service/internal/models"
	"main_service/internal/storage"
	"net/http"
)

type TicketHandler struct {
	storage *storage.PostgresStorage
}

func NewTicketHandler(storage *storage.PostgresStorage) *TicketHandler {
	return &TicketHandler{storage: storage}
}

func (h *TicketHandler) AddTicketTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var template models.TicketTemplate
	if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
		log.Println("Ошибка декодирования JSON:", err)
		http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
		return
	}

	// Получаем мероприятие
	event, found := h.storage.GetEventByID(template.EventID)
	if !found {
		http.Error(w, "Мероприятие не найдено", http.StatusNotFound)
		return
	}

	// Передаем дату и время мероприятия в шаблон билета
	template.Title = event.Title
	template.EventDate = event.DateTime // Теперь сохраняем полную дату и время

	if err := h.storage.AddTicketTemplate(template); err != nil {
		http.Error(w, "Ошибка добавления шаблона билета", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *TicketHandler) GetTickets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	tickets, err := h.storage.GetAllTickets()
	if err != nil {
		http.Error(w, "Ошибка получения билетов", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tickets)
}

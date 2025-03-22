package main

import (
	"log"
	"net/http"
	"ticket-booking/internal/handlers"
	"ticket-booking/internal/storage"
)

func main() {
	store := storage.NewMemoryStorage()
	eventHandler := handlers.NewEventHandler(store)
	ticketHandler := handlers.NewTicketHandler(store)

	http.HandleFunc("/events", eventHandler.GetEvents)
	http.HandleFunc("/events/add", eventHandler.AddEvent)
	http.HandleFunc("/tickets/add", ticketHandler.AddTicketTemplate)
	http.HandleFunc("/tickets", ticketHandler.GetTickets)

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

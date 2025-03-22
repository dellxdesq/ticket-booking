package main

import (
	"log"
	"net/http"
	"ticket-booking/internal/handlers"
	"ticket-booking/internal/storage"
)

func main() {
	// Замените строку подключения на вашу
	dataSourceName := "postgres://postgres:@localhost:5432/afishadb?sslmode=disable"
	store, err := storage.NewPostgresStorage(dataSourceName)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	if err := store.InitDB(); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}

	eventHandler := handlers.NewEventHandler(store)
	ticketHandler := handlers.NewTicketHandler(store)

	http.HandleFunc("/events", eventHandler.GetEvents)
	http.HandleFunc("/events/add", eventHandler.AddEvent)
	http.HandleFunc("/tickets/add", ticketHandler.AddTicketTemplate)
	http.HandleFunc("/tickets", ticketHandler.GetTickets)

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

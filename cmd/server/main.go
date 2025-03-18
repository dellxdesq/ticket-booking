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

	http.HandleFunc("/events", eventHandler.GetEvents)
	http.HandleFunc("/events/add", eventHandler.AddEvent)

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

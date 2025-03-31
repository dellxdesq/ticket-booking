package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	_ "os"
	"ticket-booking/internal/handlers"
	"ticket-booking/internal/storage"

	_ "github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)
	store, err := storage.NewPostgresStorage(dataSourceName)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	if err := store.InitDB(); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}

	orderHandler := handlers.NewOrderHandler(store)
	eventHandler := handlers.NewEventHandler(store)
	ticketHandler := handlers.NewTicketHandler(store)

	http.HandleFunc("/events", eventHandler.GetEvents)
	http.HandleFunc("/events/add", eventHandler.AddEvent)
	http.HandleFunc("/events/{event_id}/seats", orderHandler.GetAvailableSeatsHandler)   //TODO реализовать методы вычисления свобоных мест вместо заглушки
	http.HandleFunc("/events/{event_id}/tickets/order", orderHandler.CreateOrderHandler) //TODO реализовать реальный заказ билета вместо заглушки
	http.HandleFunc("/tickets/add", ticketHandler.AddTicketTemplate)
	http.HandleFunc("/tickets", ticketHandler.GetTickets)

	//ошибка badRequest
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		http.ServeFile(w, r, "static/error404.png")
		w.WriteHeader(http.StatusNotFound)
	})
	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

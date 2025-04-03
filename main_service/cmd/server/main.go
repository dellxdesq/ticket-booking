package main

import (
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	"log"
	"main_service/internal/handlers"
	middleware "main_service/internal/service"
	"main_service/internal/storage"
	"net/http"
	"os"
	_ "os"
)

// a
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла:", err)
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

	http.HandleFunc("/auth/register", handlers.RegisterHandler)
	http.HandleFunc("/auth/login", handlers.LoginHandler)

	http.HandleFunc("/events", eventHandler.GetEvents)
	http.HandleFunc("/events/{event_id}/seats", orderHandler.GetAvailableSeatsHandler)
	http.HandleFunc("/events/add", middleware.AuthMiddleware(eventHandler.AddEvent))
	http.HandleFunc("/events/{event_id}/tickets/order", middleware.AuthMiddleware(orderHandler.CreateOrderHandler))
	http.HandleFunc("/tickets/add", middleware.AuthMiddleware(ticketHandler.AddTicketTemplate))
	http.HandleFunc("/tickets", middleware.AuthMiddleware(ticketHandler.GetTickets))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		http.ServeFile(w, r, "static/error404.png")
		w.WriteHeader(http.StatusNotFound)
	})
	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

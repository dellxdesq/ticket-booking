package main

import (
	"log"
	"main_service/internal/config"
	"main_service/internal/handlers"
	"main_service/internal/router"
	"main_service/internal/storage"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	store, err := storage.InitStorage(cfg)
	if err != nil {
		log.Fatalf("Ошибка инициализации хранилища: %v", err)
	}

	orderHandler := handlers.NewOrderHandler(store)
	eventHandler := handlers.NewEventHandler(store)
	ticketHandler := handlers.NewTicketHandler(store)
	authHandler := handlers.NewAuthHandler(store)

	router.SetupRoutes(orderHandler, eventHandler, ticketHandler, authHandler)

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

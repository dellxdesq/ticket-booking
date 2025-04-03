package router

import (
	"main_service/internal/handlers"
	service "main_service/internal/service"
	"net/http"
)

func SetupRoutes(orderHandler *handlers.OrderHandler, eventHandler *handlers.EventHandler, ticketHandler *handlers.TicketHandler, authHandler *handlers.AuthHandler) {
	http.HandleFunc("/auth/register", authHandler.RegisterHandler)
	http.HandleFunc("/auth/login", authHandler.LoginHandler)

	http.HandleFunc("/events", eventHandler.GetEvents)
	http.HandleFunc("/events/{event_id}/seats", orderHandler.GetAvailableSeatsHandler)
	http.HandleFunc("/events/add", service.AuthMiddleware(eventHandler.AddEvent))
	http.HandleFunc("/events/{event_id}/tickets/order", service.AuthMiddleware(orderHandler.CreateOrderHandler))
	http.HandleFunc("/tickets/add", service.AuthMiddleware(ticketHandler.AddTicketTemplate))
	http.HandleFunc("/tickets", service.AuthMiddleware(ticketHandler.GetTickets))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		http.ServeFile(w, r, "static/error404.png")
		w.WriteHeader(http.StatusNotFound)
	})
}

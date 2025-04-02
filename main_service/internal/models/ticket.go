package models

import "time"

type TicketTemplate struct {
	ID           int       `json:"id"`
	EventID      int       `json:"event_id"`
	Title        string    `json:"title"`
	Price        float64   `json:"price"`
	Row          *int      `json:"row,omitempty"`
	Seat         *int      `json:"seat,omitempty"`
	Zone         *string   `json:"zone,omitempty"`
	EventDate    time.Time `json:"event_date"` // Дата + время мероприятия
	Location     string    `json:"location"`
	TicketNumber int       `json:"ticket_number"`
}

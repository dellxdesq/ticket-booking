package models

import "time"

type Event struct {
	ID       int       `json:"id"`
	Type     string    `json:"type"`
	Title    string    `json:"title"`
	DateTime time.Time `json:"dateTime"` // Дата и время теперь хранятся в одном поле
	Tickets  int       `json:"tickets"`
}

package models

import "time"

type Event struct {
	ID      int       `json:"id"`
	Type    string    `json:"type"`
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
	Tickets int       `json:"tickets"`
}

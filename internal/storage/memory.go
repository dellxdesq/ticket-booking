package storage

import (
	"sync"
	"ticket-booking/internal/models"
	"time"
)

type MemoryStorage struct {
	events []models.Event
	mu     sync.Mutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (s *MemoryStorage) AddEvent(event models.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	event.ID = len(s.events) + 1
	s.events = append(s.events, event)
}

func (s *MemoryStorage) GetEventsByDate(date time.Time) []models.Event {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []models.Event
	for _, e := range s.events {
		if e.Date.Format("2006-01-02") == date.Format("2006-01-02") {
			result = append(result, e)
		}
	}
	return result
}

func (s *MemoryStorage) GetAllEvents() []models.Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.events
}

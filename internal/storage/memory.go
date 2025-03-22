package storage

import (
	"database/sql"
	"fmt"
	"sync"
	"ticket-booking/internal/models"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)
//ticketTemplates []models.TicketTemplate
func (s *PostgresStorage) InitDB() error {
	query := `
	CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		type VARCHAR(255),
		title VARCHAR(255),
		date TIMESTAMP,
		tickets INT
	);`

	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы: %w", err)
	}
	return nil
}

type PostgresStorage struct {
	db *sql.DB
	mu sync.Mutex
}

func NewPostgresStorage(dataSourceName string) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) AddEvent(event models.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event.ID = 0
	query := `INSERT INTO events (type, title, date, tickets) VALUES ($1, $2, $3, $4) RETURNING id`
	err := s.db.QueryRow(query, event.Type, event.Title, event.Date, event.Tickets).Scan(&event.ID)
	return err
}

func (s *PostgresStorage) GetEventsByDate(date time.Time) ([]models.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []models.Event
	query := `SELECT id, type, title, date, tickets FROM events WHERE date::date = $1`
	rows, err := s.db.Query(query, date.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e models.Event
		if err := rows.Scan(&e.ID, &e.Type, &e.Title, &e.Date, &e.Tickets); err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}

func (s *PostgresStorage) GetAllEvents() ([]models.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []models.Event
	query := `SELECT id, type, title, date, tickets FROM events`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e models.Event
		if err := rows.Scan(&e.ID, &e.Type, &e.Title, &e.Date, &e.Tickets); err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}
//тут пересечение
func (s *MemoryStorage) AddTicketTemplate(template models.TicketTemplate) {
	s.mu.Lock()
	defer s.mu.Unlock()
	template.ID = len(s.ticketTemplates) + 1
	s.ticketTemplates = append(s.ticketTemplates, template)
}

func (s *MemoryStorage) GetTemplatesByEvent(eventID int) []models.TicketTemplate {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []models.TicketTemplate
	for _, t := range s.ticketTemplates {
		if t.EventID == eventID {
			result = append(result, t)
		}
	}
	return result
}

func (s *MemoryStorage) GetEventByID(eventID int) (models.Event, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, e := range s.events {
		if e.ID == eventID {
			return e, true
		}
	}
	return models.Event{}, false
}

func (s *MemoryStorage) GetAllTickets() []models.TicketTemplate {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ticketTemplates
}

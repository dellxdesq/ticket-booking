package storage

import (
	"database/sql"
	"fmt"
	"ticket-booking/internal/models"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(dataSourceName string) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) InitDB() error {
	eventQuery := `
	CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		type VARCHAR(255),
		title VARCHAR(255),
		date TIMESTAMP,
		tickets INT
	);`

	_, err := s.db.Exec(eventQuery)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы событий: %w", err)
	}

	ticketQuery := `
	CREATE TABLE IF NOT EXISTS tickets (
		id SERIAL PRIMARY KEY,
		event_id INT REFERENCES events(id),
		title VARCHAR(255),
		event_date TIMESTAMP,
		price DECIMAL,
		row INT,
		seat INT,
		zone VARCHAR(255),
		location VARCHAR(255),
		ticket_number INT
	);`

	_, err = s.db.Exec(ticketQuery)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы билетов: %w", err)
	}

	return nil
}

func (s *PostgresStorage) AddEvent(event models.Event) error {
	event.ID = 0
	query := `INSERT INTO events (type, title, date, tickets) VALUES ($1, $2, $3, $4) RETURNING id`
	err := s.db.QueryRow(query, event.Type, event.Title, event.Date, event.Tickets).Scan(&event.ID)
	return err
}

func (s *PostgresStorage) GetEventByID(eventID int) (models.Event, bool) {

	var e models.Event
	query := `SELECT id, type, title, date, tickets FROM events WHERE id = $1`
	err := s.db.QueryRow(query, eventID).Scan(&e.ID, &e.Type, &e.Title, &e.Date, &e.Tickets)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Event{}, false
		}
		return models.Event{}, false
	}
	return e, true
}

func (s *PostgresStorage) GetEventsByDate(date time.Time) ([]models.Event, error) {

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

func (s *PostgresStorage) AddTicketTemplate(template models.TicketTemplate) error {

	query := `INSERT INTO tickets (event_id, title, event_date, price, row, seat, zone, location, ticket_number) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	_, err := s.db.Exec(query, template.EventID, template.Title, template.EventDate, template.Price, template.Row, template.Seat, template.Zone, template.Location, template.TicketNumber)
	return err
}

func (s *PostgresStorage) GetAllTickets() ([]models.TicketTemplate, error) {

	var result []models.TicketTemplate
	query := `SELECT id, event_id, title, event_date, price, row, seat, zone, location, ticket_number FROM tickets`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t models.TicketTemplate
		if err := rows.Scan(&t.ID, &t.EventID, &t.Title, &t.EventDate, &t.Price, &t.Row, &t.Seat, &t.Zone, &t.Location, &t.TicketNumber); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

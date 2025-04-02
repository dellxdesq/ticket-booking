package storage

import (
	"database/sql"
	"log"
	"sort"
	"strconv"
	"strings"

	_ "github.com/jackc/pgx/v4/stdlib"
	pb "order_service/proto/grpc/order"
)

type Storage struct {
	DB *sql.DB
}

func (s *Storage) GetAvailableSeats(eventID int64) ([]*pb.Zone, error) {
	// получаю кол-во мест, рядов и какие зоны
	query := `
		SELECT zone, row, seat
		FROM tickets
		WHERE event_id = $1;
	`
	var zonesStr, rowsStr, seatsStr string
	err := s.DB.QueryRow(query, eventID).Scan(&zonesStr, &rowsStr, &seatsStr)
	if err != nil {
		return nil, err
	}

	// Разбиваем строки
	zones := strings.Split(zonesStr, "") // Теперь зоны "ABC" = ["A", "B", "C"], пока хз будет по буквам но наверно надо разделитель сделать
	rowsCount, _ := strconv.Atoi(rowsStr)
	seatsCount, _ := strconv.Atoi(seatsStr)

	//чек занятых мест
	occupiedSeats := make(map[string]map[int][]int) // типо так map[zone]map[row][]seat
	queryOccupied := `
		SELECT zone, row, seat FROM order_tickets WHERE event_id = $1;
	`
	rows, err := s.DB.Query(queryOccupied, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var zone string
		var row, seat int
		if err := rows.Scan(&zone, &row, &seat); err != nil {
			return nil, err
		}
		if occupiedSeats[zone] == nil {
			occupiedSeats[zone] = make(map[int][]int)
		}
		occupiedSeats[zone][row] = append(occupiedSeats[zone][row], seat)
	}

	//запись свободных мест
	zonesMap := make(map[string]map[int][]int)
	for _, zone := range zones {
		zonesMap[zone] = make(map[int][]int)
		for row := 1; row <= rowsCount; row++ {
			for seat := 1; seat <= seatsCount; seat++ {
				if contains(occupiedSeats[zone][row], seat) {
					continue
				}
				zonesMap[zone][row] = append(zonesMap[zone][row], seat)
			}
		}
	}

	//конверт формата под grpc
	var result []*pb.Zone
	for zone, rows := range zonesMap {
		var rowList []*pb.Row
		rowNumbers := make([]int, 0, len(rows))

		for rowNumber := range rows {
			rowNumbers = append(rowNumbers, rowNumber)
		}

		sort.Ints(rowNumbers) //рябы по порядку

		for _, rowNumber := range rowNumbers {
			rowList = append(rowList, &pb.Row{
				Number: int64(rowNumber),
				Seats:  toInt64Slice(rows[rowNumber]),
			})
		}

		result = append(result, &pb.Zone{
			Name: zone,
			Rows: rowList,
		})
	}

	//сорт зон по имени
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result, nil
}

func contains(arr []int, item int) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

func toInt64Slice(arr []int) []int64 {
	res := make([]int64, len(arr))
	for i, v := range arr {
		res[i] = int64(v)
	}
	return res
}

func (s *Storage) GetZoneRowSeat(eventID int64) (string, int64, int64, error) {
	query := `
		SELECT zone, row, seat 
		FROM tickets 
		WHERE event_id = $1
		LIMIT 1;
	`

	var zone string
	var row, seat int64

	err := s.DB.QueryRow(query, eventID).Scan(&zone, &row, &seat)
	if err != nil {
		return "", 0, 0, err
	}

	return zone, row, seat, nil
}

// запись в заказы
func (s *Storage) CreateOrder(eventID int64, zone string, row int64, seat int64, email string) error {
	query := `INSERT INTO order_tickets (event_id, zone, row, seat, user_email) VALUES ($1, $2, $3, $4, $5)`
	_, err := s.DB.Exec(query, eventID, zone, row, seat, email)
	if err != nil {
		log.Printf("Ошибка при создании заказа: %v", err)
		return err
	}
	return nil
}

// чекаю на NULL в полях tickets
func (s *Storage) CheckEventStructure(eventID int64) (existsZone bool, existsRow bool, existsSeat bool, err error) {
	query := `
		SELECT 
			(MAX(CASE WHEN zone IS NOT NULL THEN 1 ELSE 0 END) = 1),
			(MAX(CASE WHEN row IS NOT NULL THEN 1 ELSE 0 END) = 1),
			(MAX(CASE WHEN seat IS NOT NULL THEN 1 ELSE 0 END) = 1)
		FROM tickets 
		WHERE event_id = $1;
	`
	err = s.DB.QueryRow(query, eventID).Scan(&existsZone, &existsRow, &existsSeat)
	return
}

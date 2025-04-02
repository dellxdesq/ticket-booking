package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"order_service/internal/storage"
	nt "order_service/proto/grpc/notifications"
	pb "order_service/proto/grpc/order"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

var host string = "notifications:50052"

type orderServer struct {
	pb.UnimplementedOrderServiceServer
	store *storage.Storage
}

func (s *orderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	eventID := req.EventId
	zone := req.Zone
	row := req.Row
	seat := req.Seat
	email := req.Email

	zones, rows, seats, err := s.store.GetZoneRowSeat(eventID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения значений tickets: %w", err)
	}

	// Проверка входных данных
	if utf8.RuneCountInString(zone) != 1 {
		return nil, fmt.Errorf("неверное значение zone: %s", zone)
	}

	if !strings.Contains(zones, zone) {
		return nil, fmt.Errorf("зона %s отсутствует в доступных зонах: %s", zone, zones)
	}

	if (row <= 0 || row > rows) && row < 0 {
		return nil, fmt.Errorf("номер ряда %d некорректен, должно быть от 1 до %d", row, rows)
	}

	if (seat <= 0 || seat > seats) && seat < 0 {
		return nil, fmt.Errorf("номер места %d некорректен, должно быть от 1 до %d", seat, seats)
	}

	existsZone, existsRow, existsSeat, err := s.store.CheckEventStructure(eventID)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки структуры события: %w", err)
	}

	//сколько надо ввести параметров
	expectedFields := 0
	if existsZone {
		expectedFields++
	}
	if existsRow {
		expectedFields++
	}
	if existsSeat {
		expectedFields++
	}

	//ввёдённые параметры в json
	providedFields := 0
	if zone != "" {
		providedFields++
	}
	if row != 0 {
		providedFields++
	}
	if seat != 0 {
		providedFields++
	}

	if expectedFields != providedFields {
		return nil, fmt.Errorf("неверное количество параметров: ожидается %d, получено %d", expectedFields, providedFields)
	}

	err = s.store.CreateOrder(eventID, zone, row, seat, email)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать заказ: %w", err)
	}

	log.Printf("Заказ создан: eventID=%d, zone=%s, row=%d, seat=%d, email=%s", eventID, zone, row, seat, email)

	go sendNotification(email, eventID, zone, row, seat)

	eventTime, err := s.store.GetEventTime(eventID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить время события: %w", err)
	}

	sendToQueue(email, eventID, eventTime)

	return &pb.CreateOrderResponse{
		Status: fmt.Sprintf("Заказ для события %d, зона %s, ряд %d, место %d успешно создан.", eventID, zone, row, seat),
	}, nil
}

func sendToQueue(email string, eventID int64, eventTime time.Time) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Ошибка подключения к RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Ошибка открытия канала: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"notification_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Ошибка создания очереди: %v", err)
	}

	message := map[string]interface{}{
		"email":      email,
		"event_id":   eventID,
		"event_time": eventTime.Format(time.RFC3339),
	}
	body, _ := json.Marshal(message)

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Fatalf("Ошибка отправки сообщения: %v", err)
	}

	log.Printf("Сообщение отправлено: %s", body)
}

func sendNotification(email string, eventID int64, zone string, row, seat int64) {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Printf("Ошибка подключения к Notification Service: %v", err)
		return
	}
	defer conn.Close()

	client := nt.NewNotificationServiceClient(conn)

	subject := "Подтверждение бронирования"
	body := fmt.Sprintf("Здравствуйте! Спасибо за покупку билетов. Ваш билет на мероприятие %d, зона %s, ряд %d, место %d подтвержден!", eventID, zone, row, seat)

	log.Printf("Отправка email: to=%s, subject=%s", email, subject)

	_, err = client.SendEmail(context.Background(), &nt.EmailRequest{
		Email:   email,
		Subject: subject,
		Body:    body,
	})
	if err != nil {
		log.Printf("Ошибка отправки email: %v", err)
	}
}

func (s *orderServer) GetAvailableSeats(ctx context.Context, req *pb.GetAvailableSeatsRequest) (*pb.GetAvailableSeatsResponse, error) {
	eventID := req.EventId

	seatsData, err := s.store.GetAvailableSeats(eventID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения мест: %w", err)
	}

	return &pb.GetAvailableSeatsResponse{
		EventId: eventID,
		Zones:   seatsData,
	}, nil
}

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
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	store := &storage.Storage{DB: db}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Ошибка создания слушателя: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, &orderServer{store: store})
	reflection.Register(grpcServer)
	log.Println("Order Service запущен на порту :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка работы gRPC сервера: %v", err)
	}
}

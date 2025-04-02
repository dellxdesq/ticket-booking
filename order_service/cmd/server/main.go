package main

import (
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"ticket-booking/order_service/internal/storage"
	pb "ticket-booking/proto/grpc/order"

	"google.golang.org/grpc"
)

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

	return &pb.CreateOrderResponse{
		Status: fmt.Sprintf("Заказ для события %d, зона %s, ряд %d, место %d успешно создан.", eventID, zone, row, seat),
	}, nil
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
	connStr := "postgres://postgres:@localhost:5432/afishadb?sslmode=disable"
	db, err := sql.Open("pgx", connStr)
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

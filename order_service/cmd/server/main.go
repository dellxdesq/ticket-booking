package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	nt "ticket-booking/proto/grpc/notifications"
	pb "ticket-booking/proto/grpc/order"

	"google.golang.org/grpc"
)

type orderServer struct {
	pb.UnimplementedOrderServiceServer
}

func (s *orderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	eventID := req.EventId
	zone := req.Zone
	row := req.Row
	seat := req.Seat
	email := req.Email

	// Логика бронирования (имитация)
	log.Printf("Заказ: eventID=%d, zone=%s, row=%d, seat=%d, email=%s", eventID, zone, row, seat, email)

	go sendNotification(email, eventID, zone, row, seat)

	return &pb.CreateOrderResponse{
		Status: fmt.Sprintf("Заказ для события %d, зона %s, ряд %d, место %d успешно создан.", eventID, zone, row, seat),
	}, nil
}

func sendNotification(email string, eventID int64, zone string, row, seat int64) {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Printf("Ошибка подключения к Notification Service: %v", err)
		return
	}
	defer conn.Close()

	client := nt.NewNotificationServiceClient(conn)

	subject := "Подтверждение бронирования"
	body := fmt.Sprintf("Ваш билет на мероприятие %d, зона %s, ряд %d, место %d подтвержден!", eventID, zone, row, seat)

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

	// Эмуляция данных, вместо этого должен быть запрос в БД.
	seatsData := []*pb.Zone{
		{
			Name: "A",
			Rows: []*pb.Row{
				{Number: 1, Seats: []int64{1, 2, 3, 4, 5, 6, 7, 10, 14, 15}},
				{Number: 4, Seats: []int64{12, 13, 14, 15}},
			},
		},
		{
			Name: "B",
			Rows: []*pb.Row{
				{Number: 1, Seats: []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}},
			},
		},
	}

	return &pb.GetAvailableSeatsResponse{
		EventId: eventID,
		Zones:   seatsData,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Ошибка создания слушателя: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, &orderServer{})
	reflection.Register(grpcServer)
	log.Println("Order Service запущен на порту :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка работы gRPC сервера: %v", err)
	}
}

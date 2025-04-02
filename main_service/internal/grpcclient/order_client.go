package grpcclient

import (
	"context"
	"log"
	"time"

	pb "main_service/proto/grpc/order"

	"google.golang.org/grpc"
)

var host string = "order_service:50051"

func CallCreateOrder(eventID int64, zone string, row, seat int64, email string) (*pb.CreateOrderResponse, error) {
	conn, err := grpc.Dial(host, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.CreateOrder(ctx, &pb.CreateOrderRequest{
		EventId: eventID,
		Zone:    zone,
		Row:     row,
		Seat:    seat,
		Email:   email,
	})
	if err != nil {
		log.Printf("Ошибка вызова CreateOrder: %v", err)
		return nil, err
	}

	return resp, nil
}

func CallGetAvailableSeats(eventID int64) (*pb.GetAvailableSeatsResponse, error) {
	conn, err := grpc.Dial(host, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetAvailableSeats(ctx, &pb.GetAvailableSeatsRequest{EventId: eventID})
	if err != nil {
		log.Printf("Ошибка вызова GetAvailableSeats: %v", err)
		return nil, err
	}

	return resp, nil
}

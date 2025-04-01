package main

import (
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	pb "ticket-booking/proto/grpc/notifications"

	"google.golang.org/grpc"
)

type notificationsServer struct {
	//pb.UnimplementedNotificationsServiceServer
}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Ошибка создания слушателя: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, &notificationsServer{})
	reflection.Register(grpcServer)
	log.Println("Order Service запущен на порту :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка работы gRPC сервера: %v", err)
	}
}

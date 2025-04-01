package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "ticket-booking/proto/grpc/notifications"
)

type notificationServer struct {
	pb.UnimplementedNotificationServiceServer
}

func (s *notificationServer) SendEmail(ctx context.Context, req *pb.EmailRequest) (*pb.EmailResponse, error) {
	err := sendEmail(req.Email, req.Subject, req.Body)
	if err != nil {
		return &pb.EmailResponse{Status: "Ошибка отправки"}, err
	}
	return &pb.EmailResponse{Status: "Отправлено"}, nil
}

func sendEmail(to, subject, body string) error {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Println("Ошибка загрузки .env файла:", err)
	}

	smtpHost := os.Getenv("SMTP_HOST")
	fmt.Println(smtpHost)
	smtpPort := os.Getenv("SMTP_PORT")
	fmt.Println(smtpPort)
	smtpUser := os.Getenv("SMTP_USER")

	smtpPassword := os.Getenv("SMTP_PASSWORD")

	from := smtpUser
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		log.Printf("Ошибка отправки email: %v", err)
		return err
	}

	log.Printf("Email отправлен на %s", to)
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Ошибка создания слушателя: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, &notificationServer{})
	reflection.Register(grpcServer)

	log.Println("Notification Service запущен на порту :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка работы gRPC сервера: %v", err)
	}
}

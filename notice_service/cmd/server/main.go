package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net"
	"net/smtp"

	"google.golang.org/grpc"
	pb "notice_service/proto/grpc/notice"
	"os"
)

type notificationServer struct {
	pb.UnimplementedNotificationServiceServer
}

func (s *notificationServer) SendEmail(ctx context.Context, req *pb.EmailRequest) (*pb.EmailResponse, error) {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	to := []string{req.Email}

	encodedSubject := "=?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte(req.Subject)) + "?="
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		smtpUser, req.Email, encodedSubject, req.Body)

	msgBytes := []byte(msg)

	serverAddr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	err := smtp.SendMail(serverAddr, auth, smtpUser, to, msgBytes)
	if err != nil {
		log.Printf("Ошибка отправки email: %v", err)
		return &pb.EmailResponse{Status: "Ошибка отправки"}, err
	}

	log.Printf("Email успешно отправлен на: %s", req.Email)
	return &pb.EmailResponse{Status: "Отправлено"}, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла:", err)
	}

	listener, err := net.Listen("tcp", ":50054") // Новый порт
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, &notificationServer{})

	log.Println("Notice Service запущен на порту :50054")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка работы gRPC сервера: %v", err)
	}
}

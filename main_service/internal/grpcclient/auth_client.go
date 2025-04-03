package grpcclient

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "main_service/proto/grpc/auth"
)

var authHost = "auth_service:50053"

func RegisterUser(email, password string) (string, error) {
	conn, err := grpc.Dial(authHost, grpc.WithInsecure(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Println("Ошибка подключения к auth_service:", err)
		return "", err
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)
	resp, err := client.Register(context.Background(), &pb.RegisterRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Println("Ошибка регистрации:", err)
		return "", err
	}

	return resp.Status, nil
}

func LoginUser(email, password string) (string, error) {
	conn, err := grpc.Dial(authHost, grpc.WithInsecure(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Println("Ошибка подключения к auth_service:", err)
		return "", err
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)
	resp, err := client.Login(context.Background(), &pb.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Println("Ошибка входа:", err)
		return "", err
	}

	return resp.Token, nil
}

func ValidateToken(ctx context.Context, token string) (bool, error) {
	conn, err := grpc.Dial(authHost, grpc.WithInsecure(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Println("Ошибка подключения к auth_service:", err)
		return false, err
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)
	resp, err := client.ValidateToken(ctx, &pb.ValidateTokenRequest{Token: token})
	if err != nil {
		log.Println("Ошибка валидации токена:", err)
		return false, err
	}

	return resp.Valid, nil
}

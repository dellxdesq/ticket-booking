package server

import (
	"auth_service/internal/storage"
	"context"
	"fmt"
	"log"
	"net"

	pb "auth_service/proto/grpc/auth"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"time"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	db        *storage.Storage
	jwtSecret string
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthService(db *storage.Storage, secret string) *AuthService {
	return &AuthService{db: db, jwtSecret: secret}
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error generating password hash: %v", err)
	}

	_, err = s.db.DB.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", req.Email, string(hash))
	if err != nil {
		return nil, fmt.Errorf("error inserting user into DB: %v", err)
	}

	return &pb.AuthResponse{Status: "User registered successfully"}, nil
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var hash string
	err := s.db.DB.QueryRow("SELECT password FROM users WHERE email = $1", req.Email).Scan(&hash)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials: %v", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Email: req.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(24 * time.Hour)},
		},
	})
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("error signing token: %v", err)
	}

	_, err = s.db.DB.Exec("UPDATE users SET token = $1 WHERE email = $2", tokenString, req.Email)
	if err != nil {
		return nil, fmt.Errorf("error saving token to DB: %v", err)
	}

	return &pb.LoginResponse{Token: tokenString}, nil
}

func RunServer(db *storage.Storage, jwtSecret string) {
	authService := NewAuthService(db, jwtSecret)

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, authService)
	reflection.Register(s)

	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Auth service running on port 50053")
	if err := s.Serve(listener); err != nil {
		log.Fatal("failed to serve: ", err)
	}
}

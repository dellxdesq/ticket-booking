package main

import (
	"auth_service/internal/storage"
	"context"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
	"time"

	pb "auth_service/proto/grpc/auth"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	db        *sql.DB
	jwtSecret string
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthService(db *sql.DB, secret string) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: secret,
	}
}

//func generateSecret() string {
//	bytes := make([]byte, 32)
//	rand.Read(bytes)
//	return hex.EncodeToString(bytes)
//}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error generating password hash: %v", err)
	}

	_, err = s.db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", req.Email, string(hash))
	if err != nil {
		return nil, fmt.Errorf("error inserting user into DB: %v", err)
	}

	return &pb.AuthResponse{Status: "User registered successfully"}, nil
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var hash string
	err := s.db.QueryRow("SELECT password FROM users WHERE email = $1", req.Email).Scan(&hash)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials: %v", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Email: req.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(24 * time.Hour),
			},
		},
	})
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("error signing token: %v", err)
	}

	// Сохраняем токен в БД
	_, err = s.db.Exec("UPDATE users SET token = $1 WHERE email = $2", tokenString, req.Email)
	if err != nil {
		return nil, fmt.Errorf("error saving token to DB: %v", err)
	}

	return &pb.LoginResponse{Token: tokenString}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	token, err := jwt.ParseWithClaims(req.Token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	// Проверяем, соответствует ли токен тому, что в БД
	var storedToken string
	err = s.db.QueryRow("SELECT token FROM users WHERE email = $1", claims.Email).Scan(&storedToken)
	if err != nil || storedToken != req.Token {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	return &pb.ValidateTokenResponse{Valid: true}, nil
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

	store, err := storage.NewStorage(dataSourceName)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET не задан в .env")
	}
	authService := NewAuthService(store.DB, jwtSecret)

	s := grpc.NewServer()
	//authService := NewAuthService(store.DB, generateSecret())
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

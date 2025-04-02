FROM golang:1.24 as builder

WORKDIR /app
COPY . .

RUN go mod tidy

EXPOSE 50052

CMD ["go", "run", "./cmd/server/main.go"]
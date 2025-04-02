FROM golang:1.24 as builder

WORKDIR /app
COPY . .

RUN go mod tidy

EXPOSE 8080

CMD ["go", "run", "./cmd/server/main.go"]
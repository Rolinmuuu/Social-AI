FROM golang:1.23.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o socialai main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/socialai .

EXPOSE 8080

CMD ["./socialai"]
# syntax=docker/dockerfile:1
FROM golang:buster

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o ./bin/main ./main.go
CMD ["./bin/main"]
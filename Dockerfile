FROM golang:1.23 as builder

WORKDIR /app
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main ./cmd/main.go

FROM alpine:3.18.4

COPY --from=builder /app/main .

RUN chmod +x /main

EXPOSE 8080

FROM golang:1.26.1 AS builder

WORKDIR /app

COPY . .

ENV CGO_ENABLED=0
RUN go build -o /app/main ./cmd/arci/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]

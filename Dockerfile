FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /go-calc-server ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /go-calc-server .
COPY --from=builder /app/gen/openapi/calculator.swagger.json ./gen/openapi/calculator.swagger.json

EXPOSE 50051 8080

CMD ["./go-calc-server"]
package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/goCalcProj/gen/pb"
	"github.com/goCalcProj/internal/calculator"
	"github.com/rs/cors"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// gRPC-сервис
type server struct {
	pb.UnimplementedCalculatorServiceServer
}

// Обработка запроса на выполнение инструкций
func (s *server) ProcessInstructions(ctx context.Context, req *pb.ProcessInstructionsRequest) (*pb.ProcessInstructionsResponse, error) {
	cb := calculator.NewBuilder()
	return cb.ProcessInstructions(ctx, req)
}

func main() {
	// Запуск gRPC
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Ошибка прослушивания порта: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	go func() {
		log.Println("gRPC-сервер слушает на :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Ошибка запуска gRPC: %v", err)
		}
	}()

	// Контекст и HTTP-прокси через grpc-gateway
	ctx := context.Background()
	gwMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err = pb.RegisterCalculatorServiceHandlerFromEndpoint(ctx, gwMux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("Ошибка регистрации gateway: %v", err)
	}

	// Обработчик для swagger.json
	mux := http.NewServeMux()
	mux.Handle("/", gwMux)
	mux.HandleFunc("/swagger.json", serveSwagger)

	// Применение CORS ко всему HTTP-серверу
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Разрешаем всё
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Accept"},
		AllowCredentials: true,
	}).Handler(mux)

	log.Println("HTTP-сервер слушает на :8080")
	if err := http.ListenAndServe(":8080", corsHandler); err != nil {
		log.Fatalf("Ошибка запуска HTTP: %v", err)
	}
}

// serveSwagger обслуживает файл спецификации Swagger
func serveSwagger(w http.ResponseWriter, r *http.Request) {
	log.Printf("Swagger-запрос от %s (%s %s)", r.RemoteAddr, r.Method, r.URL.Path)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, r, "gen/openapi/calculator.swagger.json")
}

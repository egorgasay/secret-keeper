package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"secret-keeper/internal/server/config"
	grpchandler "secret-keeper/internal/server/handler/grpc"
	"secret-keeper/internal/server/storage"
	"secret-keeper/internal/server/usecase"
	"secret-keeper/pkg/api/server"
	"syscall"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to initialize config: %v", err)
	}

	store, err := storage.New(cfg.DBConfig)
	if err != nil {
		log.Fatalf("failed to initialize storage: %v", err)
	}
	defer store.Close()

	logic, err := usecase.New(store)
	if err != nil {
		log.Fatalf("failed to initialize logic: %v", err)
	}

	go func() {
		log.Println("Server is running on grpc://" + cfg.Host)
		grpcServer := grpc.NewServer()
		ghandler := grpchandler.New(logic)

		lis, err := net.Listen("tcp", cfg.Host)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		server.RegisterSecretKeeperServer(grpcServer, ghandler)

		err = grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("grpcServer Serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutdown Server ...")
}

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/saddiestlxrd/TaskFlow/proto/taskflow.pb"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
)

func main() {
	dbLink := os.Getenv("dbLink")
	if dbLink == "" {
		log.Fatalf("dbLink environment variable not set")
	}

	dbConn, err := pgx.Connect(context.Background(), dbLink)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbConn.Close(context.Background())

	err = dbConn.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	fmt.Println("Successfully connected to database")

	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 5000: %v", err)
	}
}

package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "TaskFlow/proto"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type TaskFlowServer struct {
	pb.UnimplementedTaskServiceServer
	db *pgx.Conn
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbLink := os.Getenv("dbLink")
	if dbLink == "" {
		log.Fatalf("dbLink environment variable not set")
	}

	dbConn, err := pgx.Connect(context.Background(), dbLink)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer func() {
		if err := dbConn.Close(context.Background()); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

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
	pb.RegisterTaskServiceServer(grpcServer, &TaskFlowServer{db: dbConn})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 5000: %v", err)
	}

	fmt.Println("Server started on port :5000") //should be in app so when i run main.go it run app.go and do it's logic
}

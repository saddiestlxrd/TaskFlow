package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"TaskFlow/models"
	pb "TaskFlow/proto"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
)

type TaskFlowServer struct {
	pb.UnimplementedTaskServiceServer
	db *pgx.Conn
}

func (s *TaskFlowServer) createUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := &models.User{
		Username:  req.Name,
		Password:  req.Password, //TODO: add hashing
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := s.db.Exec(context.Background(), "INSERT INTO users (username, password, email, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		user.Username, user.Password, user.Email, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{User: &pb.User{
		Id:       fmt.Sprintf("%d", user.Id), //Think that, this is bad decision
		Name:     user.Username,
		Email:    user.Email,
		Password: user.Password,
	}}, nil
}

func main() {

	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatalf("Error loading .env file")
	//}

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
	pb.RegisterTaskServiceServer(grpcServer, &TaskFlowServer{db: dbConn})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 5000: %v", err)
	}
}

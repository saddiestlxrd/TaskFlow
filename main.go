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
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)

type TaskFlowServer struct {
	pb.UnimplementedTaskServiceServer
	db *pgx.Conn
}

func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	return string(hashedPassword)
}

func (s *TaskFlowServer) createUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := &models.User{
		Username:  req.Name,
		Password:  HashPassword(req.Password), //may be use just var user models.User and add hash into models
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := s.db.Exec(ctx, "INSERT INTO users (username, password, email, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		user.Username, user.Password, user.Email, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{User: &pb.User{
		Id:       user.Id,
		Name:     user.Username,
		Email:    user.Email,
		Password: user.Password,
	}}, nil
}

func (s *TaskFlowServer) loginUser(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user models.User
	err := s.db.QueryRow(ctx, "SELECT * FROM users WHERE email = $1", req.Email).Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{User: &pb.User{
		Id:       user.Id,
		Name:     user.Username,
		Email:    user.Email,
		Password: user.Password,
	}}, nil
}

func (s *TaskFlowServer) createTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	var task models.Task
	err := s.db.QueryRow(ctx, "INSERT INTO tasks (title, description, status, created_at, updated_at, user_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		task.Title, task.Description, task.Status, task.CreatedAt, task.UpdatedAt, task.UserId).Scan(&task.Id)
	if err != nil {
		return nil, err
	}

	return &pb.CreateTaskResponse{Task: &pb.Task{
		Id:          task.Id,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
		Owner:       req.Owner,
	}}, nil
}

func (s *TaskFlowServer) deleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	_, err := s.db.Exec(ctx, "DELETE FROM tasks WHERE id = $1 AND user_id = $2", req.Id, req.Owner)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteTaskResponse{}, nil
}

func (s *TaskFlowServer) updateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	_, err := s.db.Exec(ctx, "UPDATE tasks SET title = $2, description = $3, status = $4, updated_at = $5 WHERE id = $1 AND user_id = $6",
		req.Id, req.Title, req.Description, req.Status, time.Now(), req.Owner)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateTaskResponse{}, nil
}

func (s *TaskFlowServer) getTasks(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	rows, err := s.db.Query(ctx, "SELECT * FROM tasks WHERE user_id = $1", req.Owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*pb.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt, &task.UserId)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, &pb.Task{
			Id:          task.Id,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			CreatedAt:   task.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
			Owner:       req.Owner,
		})
	}

	return &pb.GetTaskResponse{Tasks: tasks}, nil //used ethernet to change proto little bit add repeated to getTaskResponse
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
	defer func() {
		if err := dbConn.Close(context.Background()); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}() //stole it from ethernet think it is good to handle error

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

	fmt.Println("Server started on port :5000")
}

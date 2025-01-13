package app

import (
	pb "TaskFlow/proto"
	"github.com/jackc/pgx/v5"
)

type TaskFlowServer struct {
	pb *pb.UnimplementedTaskServiceServer
	db *pgx.Conn
}

func NewApp(pb *pb.UnimplementedTaskServiceServer, db *pgx.Conn) *TaskFlowServer {
	return &TaskFlowServer{
		pb: pb,
		db: db,
	}
} //something is wrong idk how to do this need more information

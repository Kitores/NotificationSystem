package main

import (
	pb "NotificationSystem/lib/api/user_v1"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"log"
	"net"
)

const port = ":50051"

type UserServer struct {
	conn *pgx.Conn
	pb.UnimplementedUserV1Server
}

func NewUserServer() *UserServer {
	return &UserServer{}
}

func (server *UserServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Fail to listen: %v", lis.Addr())
	}
	s := grpc.NewServer()
	pb.RegisterUserV1Server(s, server)
	log.Printf("Starting listening at %v", lis.Addr())
	return s.Serve(lis)
}

func (server *UserServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	created_user := &pb.User{FirstName: in.GetFirstName(), LastName: in.LastName, PhoneNumber: in.GetPhoneNumber()}
	tx, err := server.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("conn.Begin failed: %v", err)
	}

	_, err = tx.Exec(context.Background(), "insert into users(phone_number, first_name, last_name) values($1, $2, $3)", created_user.PhoneNumber, created_user.FirstName, created_user.LastName)
	if err != nil {
		log.Fatalf("tx.Exec failed: %v", err)
	}
	tx.Commit(context.Background())
	return created_user, err
}

func (server *UserServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	var user_list *pb.UserList = &pb.UserList{}
	rows, err := server.conn.Query(context.Background(), "select * from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := pb.User{}
		err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.PhoneNumber)
		if err != nil {
			return nil, err
		}
		user_list.Users = append(user_list.Users, &user)
	}
	return user_list, nil
}

func main() {
	//dataBaseUrl := "postgres://postgres:pass@localhost:5432/postgres"
	connStr := fmt.Sprintf("host=localhost port=5432 user=postgres password=pass dbname=subusers sslmode=disable")
	var userServer *UserServer = NewUserServer()
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to get connection: %v", err)
	}
	userServer.conn = conn
	if err := userServer.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

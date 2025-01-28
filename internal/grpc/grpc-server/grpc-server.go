package grpc_server

import (
	"NotificationSystem/lib/logger/sl"
	"NotificationSystem/pkg/user_v1"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

const (
	port = "0.0.0.0:50051"
)

type UserServer struct {
	conn *pgx.Conn
	user_v1.UnimplementedUserV1Server
}

func NewUserServer() *UserServer {
	return &UserServer{}
}

func (server *UserServer) Run(log *slog.Logger) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		slog.Error(fmt.Sprintf("Fail to listen: %s", port), sl.Err(err))
	}
	//opts := []grpc_slog.Option{
	//	grpc_slog.WithDurationField(func(duration time.Duration) slog.Field {
	//		return slog.F("grpc.time_ns", duration.Nanoseconds())
	//	}),
	//}
	//sl := *slog
	//srv := grpc.NewServer(
	//	grpc.UnaryInterceptor(
	//		grpc_middleware.ChainUnaryServer(
	//			grpc_slog.UnaryServerInterceptor(sl),
	//		),
	//	),
	//	grpc.StreamInterceptor(
	//		grpc_middleware.ChainStreamServer(
	//			grpc_slog.StreamServerInterceptor(grpcLogger(true), opts...),
	//		),
	//	),
	//)

	s := grpc.NewServer()
	user_v1.RegisterUserV1Server(s, server)
	log.Info("Starting listening at ", lis.Addr())
	return s.Serve(lis)
}

func (server *UserServer) CreateNewUser(ctx context.Context, in *user_v1.NewUser) (*user_v1.User, error) {
	fmt.Println("CreateNewUserDone!\n")
	createdUser := &user_v1.User{FirstName: in.FirstName, LastName: in.LastName, PhoneNumber: in.PhoneNumber, TelegramId: in.TelegramId, Mail: in.Mail}
	//tx, err := server.conn.Begin(context.Background())
	//if err != nil {
	//	log.Fatalf("conn.Begin failed: %v", err)
	//}
	//
	//_, err = tx.Exec(context.Background(), "insert into users(firstname, lastname, phonenumber, telegramid, mail) values($1, $2, $3, $4, $5)", createdUser.FirstName, createdUser.LastName, createdUser.PhoneNumber, createdUser.TelegramId, createdUser.Mail)
	//if err != nil {
	//	log.Printf("tx.Exec failed: %v", err)
	//} else {
	//	log.Printf("Created new user: %v", createdUser)
	//}
	//err = tx.Commit(context.Background())
	//if err != nil {
	//	log.Fatalf("tx.Commit failed: %v", err)
	//}
	//return createdUser, err
	return createdUser, nil
}

func (server *UserServer) SendNotification(ctx context.Context, in *user_v1.Notification) (*user_v1.UserList, error) {
	var userList = &user_v1.UserList{}
	fmt.Println("sendNoteDone!\n")
	//rows, err := server.conn.Query(context.Background(), "select * from users")
	//if err != nil {
	//
	//	return nil, err
	//}
	//defer rows.Close()
	//for rows.Next() {
	//	user := user_v1.User{}
	//	err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.PhoneNumber, &user.TelegramId, &user.Mail)
	//	if err != nil {
	//		return nil, err
	//	}
	//	tgbot.SendNotificationToBot(in.NotificationText, user.GetTelegramId())
	//
	//	//var subs = []string{"mihail_yermolayev@mail.ru", "mihailyermolayev@gmail.com"}
	//	var arr []string
	//	arr = append(arr, user.Mail)
	//	email_sender.SendMailFunc("mixa-erm2005@mail.ru", arr, "Normik", in.NotificationText, "1YakxkLAy7tf02C3cRHQ")
	//	userList.Users = append(userList.Users, &user)
	//}
	////log.Printf("User List: %v", user_list)
	return userList, nil
}

func RunGRPCServe(log *slog.Logger) {
	connStr := fmt.Sprintf("host=localhost port=5432 user=postgres password= dbname=postgres sslmode=disable")

	var userServer = NewUserServer()
	ExampleInterceptorLogger()
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Error("Unable to get connection: ", sl.Err(err))
		//os.Exit(1)
	}
	userServer.conn = conn
	if err = userServer.Run(log); err != nil {
		log.Error("failed to serve: ", sl.Err(err))
	}
}

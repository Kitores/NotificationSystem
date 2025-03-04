package grpc_server

import (
	"context"
	"fmt"
	"github.com/Kitores/NotificationSystem/notification-service/internal/config"
	"github.com/Kitores/NotificationSystem/notification-service/internal/tgbot"
	email_sender "github.com/Kitores/NotificationSystem/notification-service/lib/api/email-sender"
	"github.com/Kitores/NotificationSystem/notification-service/lib/logger/sl"
	"github.com/Kitores/NotificationSystem/notification-service/pkg/user_v1"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"log"
	"log/slog"
	"net"
)

type UserServer struct {
	MailPass         string
	MailHost         string
	MailPort         string
	MailFrom         string
	TelegramBotToken string
	conn             *pgx.Conn
	user_v1.UnimplementedUserV1Server
}

func NewUserServer() *UserServer {
	return &UserServer{}
}

func (server *UserServer) Run(cfg *config.Config, log *slog.Logger) error {
	address := cfg.GRPC.Host + ":" + cfg.GRPC.Port
	lis, err := net.Listen("tcp", address)
	if err != nil {
		slog.Error(fmt.Sprintf("Fail to listen: %s", address), sl.Err(err))
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
	server.MailPort = cfg.SMTP.Port
	server.MailFrom = cfg.SMTP.Mail_from
	server.MailPass = cfg.SMTP.MailPassword
	server.MailHost = cfg.SMTP.Host
	server.TelegramBotToken = cfg.BotToken
	s := grpc.NewServer()
	user_v1.RegisterUserV1Server(s, server)
	log.Info("Starting listening at ", lis.Addr())
	return s.Serve(lis)
}

func (server *UserServer) CreateNewUser(ctx context.Context, in *user_v1.NewUser) (*user_v1.User, error) {
	createdUser := &user_v1.User{FirstName: in.FirstName, LastName: in.LastName, PhoneNumber: in.PhoneNumber, TelegramId: in.TelegramId, Mail: in.Mail}
	tx, err := server.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("conn.Begin failed: %v", err)
	}

	_, err = tx.Exec(context.Background(), "insert into users(first_name, last_name, phone_number, telegram_id, email) values($1, $2, $3, $4, $5)", createdUser.FirstName, createdUser.LastName, createdUser.PhoneNumber, createdUser.TelegramId, createdUser.Mail)
	if err != nil {
		log.Printf("tx.Exec failed: %v", err)
	} else {
		log.Printf("Created new user: %v", createdUser)
	}
	err = tx.Commit(context.Background())
	if err != nil {
		log.Fatalf("tx.Commit failed: %v", err)
	}
	fmt.Println("CreateNewUserDone!\n")
	return createdUser, err
}

func (server *UserServer) SendNotification(ctx context.Context, in *user_v1.Notification) (*user_v1.UserList, error) {
	var userList = &user_v1.UserList{}
	rows, err := server.conn.Query(context.Background(), "select * from users")
	if err != nil {

		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := user_v1.User{}
		err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.PhoneNumber, &user.TelegramId, &user.Mail)
		if err != nil {
			return nil, err
		}

		err = tgbot.SendNotificationToBot(in.NotificationText, user.GetTelegramId(), server.TelegramBotToken)
		if err != nil {
			return nil, err
		}

		var arr []string
		arr = append(arr, user.Mail)

		err = email_sender.SendMailFunc(arr, "Notification", in.NotificationText, server.MailFrom, server.MailPass, server.MailHost, server.MailPort)
		if err != nil {
			return nil, err
		}
		userList.Users = append(userList.Users, &user)
	}
	//log.Printf("User List: %v", user_list)
	fmt.Println("sendNoteDone!\n")
	return userList, nil
}

func RunGRPCServe(log *slog.Logger, cfg *config.Config) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.UserDb, cfg.Storage.Password, cfg.Storage.Dbname, cfg.Storage.SSLmode)
	var userServer = NewUserServer()
	ExampleInterceptorLogger()
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Error("Unable to get connection: ", sl.Err(err))
		//os.Exit(1)
	}
	userServer.conn = conn
	if err = userServer.Run(cfg, log); err != nil {
		log.Error("failed to serve: ", sl.Err(err))
	}

}

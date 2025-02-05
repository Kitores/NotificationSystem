package main

import (
	"NotificationSystem/notification-service/pkg/user_v1"
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

// package main
type NotificationRequest struct {
	Message string `json:"message"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

const addres = "localhost:50051"

func handleNotificationRequest(msg *NotificationRequest, client user_v1.UserV1Client, ctx context.Context) {
	fmt.Println("penis")
	userList, err := client.SendNotification(ctx, &user_v1.Notification{NotificationText: msg.Message})
	if err != nil {
		log.Fatalf("Unable to save new user: %v", err)
	}
	fmt.Println(userList)
}

func consumeMessages(conn *amqp.Connection, ch *amqp.Channel, queueName string, client user_v1.UserV1Client, ctx context.Context) {
	fmt.Println("jopa")
	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			var request NotificationRequest
			fmt.Println(string(d.Body))
			err := json.Unmarshal(d.Body, &request)
			if err != nil {
				log.Printf("Failed to unmarshal message: %s", err)
				continue
			}

			handleNotificationRequest(&request, client, ctx)
		}
	}()
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	conn_grpc, err := grpc.NewClient(addres, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn_grpc.Close()
	client := user_v1.NewUserV1Client(conn_grpc)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}
	consumeMessages(conn, ch, q.Name, client, ctx)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

//TODO: запустить докер контейнер RabbitMQ, gRPC-сервер и проверить цепочку отправки сообщений

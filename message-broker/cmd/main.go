package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kitores/NotificationSystem/message-broker/pkg/user_v1"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"time"
)

type NotificationRequest struct {
	Message string `json:"message"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func handleNotificationRequest(msg *NotificationRequest, client user_v1.UserV1Client, ctx context.Context) {
	fmt.Println("Request handle")
	userList, err := client.SendNotification(ctx, &user_v1.Notification{NotificationText: msg.Message})
	if err != nil {
		log.Fatalf("Unable to save new user: %v", err)
	}
	fmt.Println(userList)
}

func consumeMessages(conn *amqp.Connection, ch *amqp.Channel, queueName string, client user_v1.UserV1Client, ctx context.Context) {
	fmt.Println("message Consumed...")
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
	url := os.Getenv("QUEUE_URL")
	address := os.Getenv("QUEUE_ADDRESS")
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	connGrpc, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer connGrpc.Close()
	client := user_v1.NewUserV1Client(connGrpc)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"notification_queue", // name
		false,                // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
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
//TODO: Дописать докерфайл(или проверить) и организовать работу с конфигом для очереди(хранение юзернейма и пароля)

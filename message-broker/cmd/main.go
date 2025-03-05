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
		log.Fatalf("Unable send notification: %v", err)
	}
	fmt.Println(userList)
}

func consumeMessages(body []byte, client user_v1.UserV1Client, ctx context.Context) {
	var request NotificationRequest
	fmt.Println(string(body))
	err := json.Unmarshal(body, &request)
	if err != nil {
		log.Printf("Failed to unmarshal message: %s", err)
	}
	handleNotificationRequest(&request, client, ctx)

	fmt.Println("message Consumed...")
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
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

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			consumeMessages(d.Body, client, ctx)
			//handleNotificationRequest(d.Body, client, ctx)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

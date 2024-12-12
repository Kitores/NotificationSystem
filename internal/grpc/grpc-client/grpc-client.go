package main

import (
	"NotificationSystem/internal/grpc/user_v1"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const (
	addres = "localhost:50051"
)

func main() {
	conn, err := grpc.NewClient(addres, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := user_v1.NewUserV1Client(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = SaveNewUser("12345", "", "", c, ctx)
	if err != nil {
		log.Fatalf("Unable to save new user: %v", err)
	}
}

func SaveNewUser(phoneNumber, firstName, lastName string, c user_v1.UserV1Client, ctx context.Context) error {

	r, err := c.CreateNewUser(ctx, &user_v1.NewUser{
		FirstName:   firstName,
		LastName:    lastName,
		PhoneNumber: phoneNumber,
	})
	if err != nil {
		log.Fatalf("could not create user: %v", err)
	}
	log.Printf("created user %v", r)
	return err
}

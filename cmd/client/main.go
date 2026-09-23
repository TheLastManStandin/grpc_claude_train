package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	courierv1 "github.com/moondoggy/courier/proto/courier/v1"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("new client:", err)
	}
	defer conn.Close()

	client := courierv1.NewCourierServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	created, err := client.CreateOrder(ctx, &courierv1.CreateOrderRequest{
		SellerAddress:   "Улица Пушкина Дом Колотушкина",
		ReceiverAddress: "Baker street 221B",
		WeightGrams:     1000,
	})
	if err != nil {
		log.Fatal("create order:", err)
	}
	log.Printf("created order: id - %v ; date - %s",
		created.GetId(), created.GetCreatedAt().AsTime())

	got, err := client.GetOrder(ctx, &courierv1.GetOrderRequest{
		Id: created.GetId(),
	})
	if err != nil {
		log.Fatal("get order:", err)
	}
	log.Printf("got order %q / %q / %d",
		got.GetReceiverAddress(), got.GetSellerAddress(), got.GetWeightGrams())

	if _, err := client.GetOrder(ctx, &courierv1.GetOrderRequest{Id: 9999}); err != nil {
		st, ok := status.FromError(err)

		if !ok {
			// Не gRPC-ошибка: например, контекст отменили локально,
			// до того как запрос вообще ушёл в сеть.
			log.Fatalf("не gRPC-ошибка: %v", err)
		}

		log.Printf("ожидаемо: code=%s message=%q", st.Code(), st.Message())

		switch st.Code() {
		case codes.NotFound:
			log.Println("→ такого заказа нет, ретраить бессмысленно")
		case codes.Unavailable:
			log.Println("→ сервер лежит, вот это можно ретраить")
		default:
			log.Println("→ непонятно, пробрасываем наверх")
		}
	}
}

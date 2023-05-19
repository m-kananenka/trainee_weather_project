package service

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"telegram_service/gRPC/pb"
)

type TgService struct{}

func (s TgService) GetWeather() string {

	conn, err := grpc.Dial("localhost:8083", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Создание клиентского экземпляра
	client := pb.NewGetWeatherClient(conn)

	// Вызов удаленной процедуры
	req := &pb.Request{
		User:    "World",
		Weather: "gege",
	}

	res, err := client.Get(context.Background(), req)
	if err != nil {
		log.Printf("Failed to call MyMethod: %v", err)
		return "something wrong"
	}

	return res.String()
}

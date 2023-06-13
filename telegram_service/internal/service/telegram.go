package service

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	pb2 "telegram_service/cmd/weather/pb"
)

type TgService struct{}

func (t *TgService) GetWeather(city string) (string, error) {

	conn, err := grpc.Dial("localhost:8083", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return "", err
	}
	defer conn.Close()

	weatherClient := pb2.NewGetWeatherClient(conn)

	req := &pb2.Request{
		City: city,
	}

	res, err := weatherClient.Get(context.Background(), req)
	if err != nil {
		log.Printf("Failed to call GetWeather: %v", err)
		return "", err
	}

	return res.String(), nil
}

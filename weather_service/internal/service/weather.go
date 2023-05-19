package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"weather_service/gRPC/proto/pb"
	"weather_service/internal/config"
)

type GRPCServer struct {
	cfg    *config.Config
	logger *logrus.Logger
}

func NewGRPCServer(cfg *config.Config, logger *logrus.Logger) *GRPCServer {
	return &GRPCServer{
		cfg:    cfg,
		logger: logger,
	}
}

func (g *GRPCServer) Get(context.Context, *pb.Request) (*pb.Response, error) {
	temps := g.GetTemperatureByCity()
	return &pb.Response{Response: temps}, nil
}

type respBody struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func (g *GRPCServer) GetTemp(city string) string {

	resp, err := http.Get(fmt.Sprintf(g.cfg.URL, g.cfg.APIKey, city))
	if err != nil {
		g.logger.Printf("request to openweathermap failed: %s\n", err.Error())
	}

	var data respBody
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		g.logger.Printf("failed to decode response body: %s\n", err.Error())
	}

	return fmt.Sprintf("City: %s, Temp: %.1f", city, kelvinToCelsius(data.Main.Temp))

}

func kelvinToCelsius(temp float64) float64 {
	const kelvinConstant = 273
	return temp - kelvinConstant
}

func (g *GRPCServer) GetTemperatureByCity() []string {

	cities := []string{"Minsk", "Gomel", "Mogilev", "Brest", "Grodno", "Vitebsk"}
	var result = make([]string, 0)
	for i := range cities {
		temp := g.GetTemp(cities[i])
		result = append(result, temp)
	}

	return result
}

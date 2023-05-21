package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"weather_service/api/pb"
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

func (g *GRPCServer) Get(ctx context.Context, req *pb.Request) (*pb.Response, error) {

	temps := g.GetWeather(req.GetCity())
	return &pb.Response{Response: temps}, nil
}

type respBody struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func (g *GRPCServer) GetWeather(city string) string {

	resp, err := http.Get(fmt.Sprintf(g.cfg.URL, g.cfg.APIKey, city))
	if err != nil {
		g.logger.Printf("request to openweathermap failed: %s\n", err.Error())
		return "incorrect name of city"
	}

	var data respBody
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		g.logger.Printf("failed to decode response body: %s\n", err.Error())
		return "incorrect name of city"
	}

	if data.Main.Temp == 0 {
		return ""
	}

	return fmt.Sprintf("City: %s, Temp: %.1f", city, kelvinToCelsius(data.Main.Temp))

}

func kelvinToCelsius(temp float64) float64 {
	const kelvinConstant = 273
	return temp - kelvinConstant
}

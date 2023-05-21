package server

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"net"
	"weather_service/api/pb"
	"weather_service/internal/config"
	"weather_service/internal/service"
)

type WeatherServer struct {
	logger  *logrus.Logger
	cfg     *config.Config
	client  *grpc.Server
	service *service.GRPCServer
}

func NewWeatherServer(logger *logrus.Logger, cfg *config.Config, service *service.GRPCServer) *WeatherServer {
	return &WeatherServer{
		logger:  logger,
		cfg:     cfg,
		client:  grpc.NewServer(),
		service: service,
	}
}

func (w *WeatherServer) Register() {

	w.client.RegisterService(&pb.GetWeather_ServiceDesc, w.service)
}

func (w *WeatherServer) Start() {

	l, err := net.Listen("tcp", w.cfg.Port)
	if err != nil {
		log.Fatal(err)
	}

	if err := w.client.Serve(l); err != nil {
		log.Fatal(err)
	}
}

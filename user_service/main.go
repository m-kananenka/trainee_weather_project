package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"user_service/internal/config"
	"user_service/internal/user/repository"
	"user_service/internal/user/server"
	"user_service/internal/user/service"
)

func main() {

	logger := logrus.New()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal(err)
	}

	cfg := &config.Config{}
	err = cfg.Process()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info(cfg.DB.Driver)

	db, err := sqlx.Connect(cfg.DB.Driver, fmt.Sprintf("user=%s dbname=%s  sslmode=%s password= %s", cfg.DB.User,
		cfg.DB.Name, cfg.DB.SSLMode, cfg.DB.Password))

	if err != nil {
		logger.Fatal(err)
	}

	userRepo := repository.NewUserRepo(db, cfg.DB)

	err = userRepo.RunMigrations()
	if err != nil {
		logger.Warning(err)
	}

	controller := service.NewController(userRepo, cfg)

	srv := server.NewServer(cfg.Port, echo.New(), logger, controller, cfg)

	srv.Register()

	srv.RegisterRoutes()

	go srv.StartGRPC()

	srv.StartRouter()

}

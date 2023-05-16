package server

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"trainee_weather_project/weather_service/internal/config"
)

type WeatherServer struct {
	logger *logrus.Logger
	cfg    *config.Config
	client http.Client
}

func NewWeatherServer(logger *logrus.Logger, cfg *config.Config) *WeatherServer {
	return &WeatherServer{
		logger: logger,
		cfg:    cfg,
		client: http.Client{},
	}
}

type respBody struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func (w WeatherServer) GetTemp(city string) {

	resp, err := w.client.Get(fmt.Sprintf(w.cfg.URL, w.cfg.APIKey, city))
	if err != nil {
		w.logger.Printf("request to openweathermap failed: %s\n", err.Error())
	}

	var data respBody
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		w.logger.Printf("failed to decode response body: %s\n", err.Error())
	}

	fmt.Printf("City: %s, Temp: %.1f\n", city, kelvinToCelsius(data.Main.Temp))

}

func kelvinToCelsius(temp float64) float64 {
	const kelvinConstant = 273
	return temp - kelvinConstant
}

func (w WeatherServer) PrintTemps() []string {
	cities := []string{"Minsk", "Gomel", "Mogilev", "Brest", "Grodno", "Vitebsk"}
	for i := range cities {

		w.GetTemp(cities[i])
	}
	return cities
}

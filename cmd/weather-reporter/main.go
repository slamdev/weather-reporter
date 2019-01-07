package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	h "net/http"
	"os"
	"os/signal"
	"syscall"
	"weather-reporter/internal"
	"weather-reporter/internal/http"
	"weather-reporter/internal/weather"
	"weather-reporter/internal/weather/providers"
)

var httpServer http.HttpServer

func init() {
	config := internal.NewConfig()
	log.WithField("config", config).Debug("application config")

	yahooWeatherProvider := providers.NewYahooWeatherProvider(h.Client{
		Timeout:   config.HttpClientTimeout,
		Transport: http.InstrumentHttpTransport("yahoo", h.DefaultTransport),
	})

	openWeatherMapWeatherProvider := providers.NewOpenWeatherMapWeatherProvider(h.Client{
		Timeout:   config.HttpClientTimeout,
		Transport: http.InstrumentHttpTransport("openWeatherMap", h.DefaultTransport),
	}, config.OpenWeatherMapAppID)

	cache := weather.NewWeatherCache(config.CacheExpiration)

	weatherProcessor := weather.NewWeatherService(cache, yahooWeatherProvider, openWeatherMapWeatherProvider)
	handler := func(weather string) (interface{}, error) {
		return weatherProcessor.GetCurrentWeather(weather)
	}

	httpServer = http.NewHttpServer(config.HttpPort, http.CreateWeatherHttpRouter(handler))
}

func main() {
	go func() {
		log.Info("starting http server")
		if err := httpServer.Start(); err != nil {
			log.WithField("error", fmt.Sprintf("%+v", err)).Fatal("http server failed to start")
		}
	}()

	waitForShutdown(func() {
		if err := httpServer.Stop(); err != nil {
			log.WithField("error", fmt.Sprintf("%+v", err)).Fatal("http server failed to stop")
		}
		log.Info("http server stopped")
	})
}

func waitForShutdown(shutdownHook func()) {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	sig := <-gracefulStop
	log.WithField("signal", sig).Info("shutdown signal received")
	shutdownHook()
}

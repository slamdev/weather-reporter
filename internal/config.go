package internal

import (
	"github.com/namsral/flag"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type Config struct {
	HttpPort            int
	HttpClientTimeout   time.Duration
	OpenWeatherMapAppID string
	CacheExpiration     time.Duration
}

func NewConfig() Config {
	var config Config

	flag.IntVar(&config.HttpPort, "http_port", 8080, "The port for the http server to listen on")

	flag.DurationVar(&config.HttpClientTimeout, "http_client_timeout", time.Second*2, "The timeout for http client requests")

	flag.StringVar(&config.OpenWeatherMapAppID, "open_weather_map_app_id", "_REPLACE_",
		"The App ID for the OpenWeatherMap provider")

	flag.DurationVar(&config.CacheExpiration, "cache_expiration", time.Second*60, "The weather cache expiration time")

	var logFormat string
	flag.StringVar(&logFormat, "log_format", "text", "The format of the logs. Either text, or json")

	var debug bool
	flag.BoolVar(&debug, "debug", true, "Enable debug logging")

	flag.Parse()

	var formatter log.Formatter = &log.TextFormatter{}
	if logFormat == strings.ToLower("json") {
		formatter = &log.JSONFormatter{}
	}
	log.SetFormatter(formatter)

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	return config
}

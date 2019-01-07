package providers

import (
	"encoding/json"
	"fmt"
	"github.com/oliveagle/jsonpath"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"weather-reporter/internal/weather"
)

const yahooUrl = "https://query.yahooapis.com/v1/public/yql"

func NewYahooWeatherProvider(client http.Client) weather.Provider {
	return &yahooWeatherProvider{
		client: client,
	}
}

type yahooWeatherProvider struct {
	client http.Client
}

func (p *yahooWeatherProvider) Get(city string) (weather.Weather, error) {
	query := `select item.condition, wind from weather.forecast where woeid in (select woeid from geo.places(1) where text="%v")`
	params := url.Values{}
	params.Set("format", "json")
	params.Set("q", fmt.Sprintf(query, strings.ToLower(city)))
	urlString := yahooUrl + "?" + params.Encode()
	log.WithField("url", urlString).
		WithField("provider", "yahoo").
		Debug("sending http request")
	r, err := p.client.Get(urlString)
	if err != nil {
		return weather.Weather{}, errors.Wrapf(err, "yahoo: failed to get %v weather", city)
	}
	if r.StatusCode != 200 {
		return weather.Weather{}, errors.Errorf("yahoo: request failed with message: %v", r.Status)
	}
	defer r.Body.Close()
	return p.toWeather(r.Body)
}

func (p *yahooWeatherProvider) toWeather(data io.Reader) (weather.Weather, error) {
	var jsonData interface{}
	err := json.NewDecoder(data).Decode(&jsonData)
	if err != nil {
		return weather.Weather{}, errors.Wrap(err, "yahoo: failed to unmarshal json response")
	}
	windSpeedStr, err := jsonpath.JsonPathLookup(jsonData, "$.query.results.channel.wind.speed")
	if err != nil {
		return weather.Weather{}, errors.Wrapf(err, "yahoo: failed to extract wind speed from %v", jsonData)
	}
	windSpeed, err := strconv.Atoi(windSpeedStr.(string))
	if err != nil {
		return weather.Weather{}, errors.Wrapf(err, "yahoo: failed to convert %v to int", windSpeedStr)
	}
	temperatureDegreesStr, err := jsonpath.JsonPathLookup(jsonData, "$.query.results.channel.item.condition.temp")
	if err != nil {
		return weather.Weather{}, errors.Wrapf(err, "yahoo: failed to extract temperature degrees from %v", jsonData)
	}
	temperatureDegrees, err := strconv.Atoi(temperatureDegreesStr.(string))
	if err != nil {
		return weather.Weather{}, errors.Wrapf(err, "yahoo: failed to convert %v to int", temperatureDegreesStr)
	}
	w := weather.Weather{
		WindSpeed:          windSpeed,
		TemperatureDegrees: p.toCelsius(temperatureDegrees),
	}
	log.WithField("weather", w).
		WithField("provider", "yahoo").
		Debug("got weather data")
	return w, nil
}

func (p *yahooWeatherProvider) toCelsius(fahrenheit int) int {
	celsius := (float64(fahrenheit) - 32) * 5 / 9
	return int(math.Round(celsius))
}

package providers

import (
	"encoding/json"
	"github.com/oliveagle/jsonpath"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"weather-reporter/internal/weather"
)

const openWeatherMapUrl = "http://api.openweathermap.org/data/2.5/weather"

func NewOpenWeatherMapWeatherProvider(client http.Client, appID string) weather.Provider {
	return &openWeatherMapWeatherProvider{
		client: client,
		appID:  appID,
	}
}

type openWeatherMapWeatherProvider struct {
	client http.Client
	appID  string
}

func (p *openWeatherMapWeatherProvider) Get(city string) (weather.Weather, error) {
	params := url.Values{}
	params.Set("appid", p.appID)
	params.Set("units", "metric")
	params.Set("q", strings.ToLower(city))
	urlString := openWeatherMapUrl + "?" + params.Encode()
	log.WithField("url", urlString).
		WithField("provider", "openWeatherMap").
		Debug("sending http request")
	r, err := p.client.Get(urlString)
	if err != nil {
		return weather.Weather{}, errors.Wrapf(err, "openWeatherMap: failed to get %v weather", city)
	}
	if r.StatusCode != 200 {
		return weather.Weather{}, errors.Errorf("openWeatherMap: request failed with message: %v", r.Status)
	}
	defer r.Body.Close()
	return p.toWeather(r.Body)
}

func (p *openWeatherMapWeatherProvider) toWeather(data io.Reader) (weather.Weather, error) {
	var jsonData interface{}
	err := json.NewDecoder(data).Decode(&jsonData)
	if err != nil {
		return weather.Weather{}, errors.Wrap(err, "openWeatherMap: failed to unmarshal json response")
	}
	windSpeed, err := jsonpath.JsonPathLookup(jsonData, "$.wind.speed")
	if err != nil {
		return weather.Weather{}, errors.Wrapf(err, "openWeatherMap: failed to extract wind speed from %v", jsonData)
	}
	temperatureDegrees, err := jsonpath.JsonPathLookup(jsonData, "$.main.temp")
	if err != nil {
		return weather.Weather{}, errors.Wrapf(err, "openWeatherMap: failed to extract temperature degrees from %v", jsonData)
	}
	w := weather.Weather{
		WindSpeed:          int(math.Round(windSpeed.(float64))),
		TemperatureDegrees: int(math.Round(temperatureDegrees.(float64))),
	}
	log.WithField("weather", w).
		WithField("provider", "openWeatherMap").
		Debug("got weather data")
	return w, nil
}

package providers

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"weather-reporter/internal/weather"
)

func Test_Should_Build_Yahoo_Api_Url(t *testing.T) {
	city := "test-city"
	assertFunc := func(req *http.Request) {
		assert.Contains(t, req.URL.RawQuery, city)
	}
	client := NewClientStub("", 0, nil, assertFunc)
	provider := NewYahooWeatherProvider(client)
	_, _ = provider.Get(city)
}

func Test_Should_Return_Error_From_Yahoo_Http_Client(t *testing.T) {
	msg := "test error message"
	client := NewClientStub("", 0, errors.New(msg))
	provider := NewYahooWeatherProvider(client)
	_, err := provider.Get("test")
	assert.Contains(t, err.Error(), msg)
}

func Test_Should_Return_Error_When_Yahoo_Response_Status_Is_Not_OK(t *testing.T) {
	client := NewClientStub("", 500, nil)
	provider := NewYahooWeatherProvider(client)
	_, err := provider.Get("test")
	assert.Contains(t, err.Error(), "request failed")
}

func Test_Should_Return_Error_When_Yahoo_Response_Is_Not_Valid_Json(t *testing.T) {
	client := NewClientStub("TEST", 200, nil)
	provider := NewYahooWeatherProvider(client)
	_, err := provider.Get("test")
	assert.Contains(t, err.Error(), "failed to unmarshal json")
}

func Test_Should_Return_Error_When_Yahoo_Response_Has_No_Wind_Speed(t *testing.T) {
	client := NewClientStub(`{"query":{"results":{"channel":{"item":{"condition":{"temp":"0"}}}}}}`, 200, nil)
	provider := NewYahooWeatherProvider(client)
	_, err := provider.Get("test")
	assert.Contains(t, err.Error(), "failed to extract wind speed")
}

func Test_Should_Return_Error_When_Yahoo_Response_Has_No_Temperature(t *testing.T) {
	client := NewClientStub(`{"query":{"results":{"channel":{"wind":{"speed":"0"}}}}}`, 200, nil)
	provider := NewYahooWeatherProvider(client)
	_, err := provider.Get("test")
	assert.Contains(t, err.Error(), "failed to extract temperature")
}

func Test_Should_Return_Weather_From_Yahoo_Response(t *testing.T) {
	client := NewClientStub(`{"query":{"results":{"channel":{"wind":{"speed":"2"},"item":{"condition":{"temp":"33"}}}}}}`, 200, nil)
	provider := NewYahooWeatherProvider(client)
	w, err := provider.Get("test")
	assert.NoError(t, err)
	assert.Equal(t, w, weather.Weather{TemperatureDegrees: 1, WindSpeed: 2})
}

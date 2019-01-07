package providers

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"weather-reporter/internal/weather"
)

func Test_Should_Build_OWM_Api_Url(t *testing.T) {
	appID := "test-id"
	city := "test-city"
	assertFunc := func(req *http.Request) {
		assert.Contains(t, req.URL.RawQuery, fmt.Sprintf("appid=%v&q=%v&units=metric", appID, city))
	}
	client := NewClientStub("", 0, nil, assertFunc)
	provider := NewOpenWeatherMapWeatherProvider(client, appID)
	_, _ = provider.Get(city)
}

func Test_Should_Return_Error_From_OWM_Http_Client(t *testing.T) {
	msg := "test error message"
	client := NewClientStub("", 0, errors.New(msg))
	provider := NewOpenWeatherMapWeatherProvider(client, "")
	_, err := provider.Get("test")
	assert.Contains(t, err.Error(), msg)
}

func Test_Should_Return_Error_When_OWM_Response_Status_Is_Not_OK(t *testing.T) {
	client := NewClientStub("", 500, nil)
	provider := NewOpenWeatherMapWeatherProvider(client, "")
	_, err := provider.Get("test")
	assert.Contains(t, err.Error(), "request failed")
}

func Test_Should_Return_Error_When_OWM_Response_Is_Not_Valid_Json(t *testing.T) {
	client := NewClientStub("TEST", 200, nil)
	provider := NewOpenWeatherMapWeatherProvider(client, "")
	_, err := provider.Get("test")
	assert.Contains(t, err.Error(), "failed to unmarshal json")
}

func Test_Should_Return_Error_When_OWM_Response_Has_No_Wind_Speed(t *testing.T) {
	client := NewClientStub(`{"main":{"temp":0}}`, 200, nil)
	provider := NewOpenWeatherMapWeatherProvider(client, "")
	_, err := provider.Get("test")
	assert.Contains(t, err.Error(), "failed to extract wind speed")
}

func Test_Should_Return_Error_When_OWM_Response_Has_No_Temperature(t *testing.T) {
	client := NewClientStub(`{"wind":{"speed":0}}`, 200, nil)
	provider := NewOpenWeatherMapWeatherProvider(client, "")
	_, err := provider.Get("test")
	assert.Contains(t, err.Error(), "failed to extract temperature")
}

func Test_Should_Return_Weather_From_OWM_Response(t *testing.T) {
	client := NewClientStub(`{"main":{"temp":1},"wind":{"speed":2}}`, 200, nil)
	provider := NewOpenWeatherMapWeatherProvider(client, "")
	w, err := provider.Get("test")
	assert.NoError(t, err)
	assert.Equal(t, w, weather.Weather{TemperatureDegrees: 1, WindSpeed: 2})
}

package weather

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_Should_Return_Error_When_No_Providers_Configured(t *testing.T) {
	cache := new(cacheMock)
	cache.On("Get", mock.Anything).Return(Weather{}, false)
	service := NewWeatherService(cache)
	_, err := service.GetCurrentWeather("")
	assert.Contains(t, err.Error(), "no providers configured")
}

func Test_Should_Return_Last_Error_When_All_Providers_Failed(t *testing.T) {
	cache := new(cacheMock)
	cache.On("Get", mock.Anything).Return(Weather{}, false)
	p1 := provider(func(city string) (weather Weather, e error) {
		return Weather{}, errors.New("error-1")
	})
	p2 := provider(func(city string) (weather Weather, e error) {
		return Weather{}, errors.New("error-2")
	})
	service := NewWeatherService(cache, p1, p2)
	_, err := service.GetCurrentWeather("")
	assert.Contains(t, err.Error(), "error-2")
}

func Test_Should_Return_Weather_From_Cache_When_All_Providers_Failed(t *testing.T) {
	weather := Weather{TemperatureDegrees: 1}
	cache := new(cacheMock)
	cache.On("Get", "test").Return(weather, true)
	p1 := provider(func(city string) (weather Weather, e error) {
		return Weather{}, errors.New("error-1")
	})
	p2 := provider(func(city string) (weather Weather, e error) {
		return Weather{}, errors.New("error-2")
	})
	service := NewWeatherService(cache, p1, p2)
	actualWeather, err := service.GetCurrentWeather("test")
	assert.NoError(t, err)
	assert.Equal(t, weather, actualWeather)
}

func Test_Should_Put_Weather_From_Provider_To_Cache(t *testing.T) {
	city := "test"
	weather := Weather{TemperatureDegrees: 1}
	p := provider(func(city string) (Weather, error) {
		return weather, nil
	})
	cache := new(cacheMock)
	cache.On("Put", city, weather).Once()
	service := NewWeatherService(cache, p)
	_, _ = service.GetCurrentWeather(city)
}

func Test_Should_Return_Weather_From_Provider(t *testing.T) {
	cache := new(cacheMock)
	cache.On("Get", mock.Anything).Return(Weather{}, false)
	cache.On("Put", mock.Anything, mock.Anything).Maybe()
	city := "test"
	weather := Weather{TemperatureDegrees: 1}
	p := provider(func(c string) (Weather, error) {
		if c == city {
			return weather, nil
		}
		return Weather{}, errors.New("unexpected error")
	})
	service := NewWeatherService(cache, p)
	actualWeather, err := service.GetCurrentWeather(city)
	assert.NoError(t, err)
	assert.Equal(t, weather, actualWeather)
}

type cacheMock struct {
	mock.Mock
}

func (c *cacheMock) Get(city string) (Weather, bool) {
	args := c.Called(city)
	return args.Get(0).(Weather), args.Bool(1)
}

func (c *cacheMock) Put(city string, weather Weather) {
	c.Called(city, weather)
}

type providerStub struct {
	handler func(city string) (Weather, error)
}

func (p *providerStub) Get(city string) (Weather, error) {
	return p.handler(city)
}

func provider(handler func(city string) (Weather, error)) Provider {
	return &providerStub{handler: handler}
}

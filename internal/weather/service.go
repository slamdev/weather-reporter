package weather

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Service interface {
	GetCurrentWeather(city string) (Weather, error)
}

func NewWeatherService(cache Cache, weatherProviders ...Provider) Service {
	return &service{
		weatherProviders: weatherProviders,
		cache:            cache,
	}
}

type service struct {
	weatherProviders []Provider
	cache            Cache
}

func (s *service) GetCurrentWeather(city string) (Weather, error) {
	log.WithField("city", city).Debug("searching for weather")
	weather, err := s.getWeatherFromProvider(city)
	if err == nil {
		return weather, nil
	}
	err = errors.Wrapf(err, "failed to get %v weather from providers", city)
	weather, found := s.cache.Get(city)
	if found {
		log.WithField("city", city).
			WithField("error", fmt.Sprintf("%+v", err)).
			Warn("failed to get weather from provider; cached result will be returned")
		return weather, nil
	}
	return Weather{}, err
}

func (s *service) getWeatherFromProvider(city string) (Weather, error) {
	if len(s.weatherProviders) == 0 {
		return Weather{}, errors.New("no providers configured")
	}
	var lastError error
	for _, currentProvider := range s.weatherProviders {
		weather, err := currentProvider.Get(city)
		if err == nil {
			s.cache.Put(city, weather)
			return weather, nil
		}
		log.WithField("city", city).
			WithField("error", err).
			Warn("failed to get weather from provider")
		lastError = errors.Wrapf(err, "failed to get %v weather from provider", city)
	}
	return Weather{}, lastError
}

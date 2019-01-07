package http

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type WeatherHandler func(string) (interface{}, error)

func CreateWeatherHttpRouter(handler WeatherHandler) Router {
	return Router{
		Method:  "GET",
		Path:    "/v1/weather",
		Queries: []string{"city", "{city:[A-z]+}"},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handleRequest(w, r, handler)
		}),
	}
}

func handleRequest(writer http.ResponseWriter, request *http.Request, handler WeatherHandler) {
	routeVars := mux.Vars(request)
	city := routeVars["city"]
	data, err := handler(city)
	if err != nil {
		sendErrorResponse(writer, errors.Wrap(err, "failed to retrieve data"))
		return
	}
	sendAsJson(writer, data)
}

func sendAsJson(writer http.ResponseWriter, data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		sendErrorResponse(writer, errors.Wrap(err, "failed to marshal json"))
		return
	}
	_, err = writer.Write(response)
	if err != nil {
		err = errors.Wrap(err, "failed to write response")
		log.WithField("error", fmt.Sprintf("%+v", err)).Error()
	}
}

func sendErrorResponse(writer http.ResponseWriter, err error) {
	log.WithField("error", fmt.Sprintf("%+v", err)).Error()
	writer.WriteHeader(http.StatusInternalServerError)
	_, writeError := writer.Write([]byte(err.Error()))
	if writeError != nil {
		writeError = errors.Wrap(writeError, "failed to write response")
		log.WithField("error", fmt.Sprintf("%+v", writeError)).Error()
	}
}

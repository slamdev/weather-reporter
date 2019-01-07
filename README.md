# Weather Reporter

## AIM
Create a HTTP Service that reports on sydney weather. This service will source its information from the either of the below providers:
1. [Yahoo](https://developer.yahoo.com/weather/) (primary):
```bash
curl "https://query.yahooapis.com/v1/public/yql?q=select%20item.condition%2C%20wind%20from%20weather.forecast%20where%20woeid%20%3D%201105779&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys"
```
2. [OpenWeatherMap](https://openweathermap.org/current) (failover):
```bash
curl "http://api.openweathermap.org/data/2.5/weather?q=sydney,AU&appid=2326504fb9b100bee21400190e4dbe6d"
```

## Specs
- The service can hard-code Sydney as a city.
- The service should return a JSON payload with a unified response containing temperature in degrees Celsius and wind speed.
- If one of the provider goes down, service can quickly failover to a different provider without affecting customers.
- Have scalability and reliability in mind when designing the solution.
- Weather results are fine to be cached for up to 3 seconds on the server in normal behaviour to prevent hitting weather providers. Those results must be served as stale if all weather providers are down.

## Expected Output
Calling the service via `curl "http://localhost:8080/v1/weather?city=sydney"` should output the following JSON payload:
```json
{
  "wind_speed": 20,
  "temperature_degrees": 29
}
```

## Running

```bash
docker build -t weather-reporter .
docker run -p 8080:8080 weather-reporter
curl http://localhost:8080/v1/weather?city=sydney
```

## Monitoring

Service exposes `/health` endpoint for general health monitoring.

Service exposes `/metrics` endpoint in prometheus format. The following metrics are collected:
- requests count
- response times
- request\responses to weather providers
- cache hits and misses

## Logging

Logs format is configured via `LOG_FORMAT` env variable. It should be set to `json` for better integration with log collectors.

package weather

type Provider interface {
	Get(city string) (Weather, error)
}

type Weather struct {
	WindSpeed          int `json:"wind_speed"`
	TemperatureDegrees int `json:"temperature_degrees"`
}

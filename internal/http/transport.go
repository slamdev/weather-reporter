package http

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func InstrumentHttpTransport(name string, transport http.RoundTripper) http.RoundTripper {
	return promhttp.InstrumentRoundTripperInFlight(registerInFlightGaugeMetric(name),
		promhttp.InstrumentRoundTripperCounter(registerCounterMetric(name),
			promhttp.InstrumentRoundTripperDuration(registerHistVecMetric(name), transport),
		),
	)
}

func registerHistVecMetric(name string) prometheus.ObserverVec {
	metric := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:      name + "_request_duration_seconds",
			Namespace: "weather_reporter",
			Help:      "A histogram of request latencies.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method"},
	)
	prometheus.MustRegister(metric)
	return metric
}

func registerCounterMetric(name string) *prometheus.CounterVec {
	metric := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      name + "_client_api_requests_total",
			Namespace: "weather_reporter",
			Help:      "A counter for requests from the wrapped client.",
		},
		[]string{"code", "method"},
	)
	prometheus.MustRegister(metric)
	return metric
}

func registerInFlightGaugeMetric(name string) prometheus.Gauge {
	metric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      name + "_client_in_flight_requests",
		Namespace: "weather_reporter",
		Help:      "A gauge of in-flight requests for the wrapped client.",
	})
	prometheus.MustRegister(metric)
	return metric
}

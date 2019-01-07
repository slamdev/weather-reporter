package weather

import (
	impl "github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type Cache interface {
	Get(city string) (Weather, bool)
	Put(city string, weather Weather)
}

func NewWeatherCache(expiration time.Duration) Cache {
	return &cache{
		cache:  impl.New(expiration, expiration/2),
		metric: registerCacheMetric(),
	}
}

type cache struct {
	cache  *impl.Cache
	metric *prometheus.CounterVec
}

func (c *cache) Get(city string) (Weather, bool) {
	weather, found := c.cache.Get(city)
	if found {
		c.metric.WithLabelValues("hit").Inc()
		return weather.(Weather), true
	}
	c.metric.WithLabelValues("miss").Inc()
	return Weather{}, false
}

func (c *cache) Put(city string, weather Weather) {
	c.cache.Set(city, weather, impl.DefaultExpiration)
}

func registerCacheMetric() *prometheus.CounterVec {
	metric := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "weather_reporter",
		Name:      "cache",
		Help:      "Counter of cache hits or misses.",
	}, []string{"state"})
	prometheus.MustRegister(metric)
	return metric
}

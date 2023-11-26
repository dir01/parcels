package metrics

import (
	"github.com/dir01/parcels/service"
	"github.com/prometheus/client_golang/prometheus"
)

func NewPrometheus(apiNames []service.APIName) service.Metrics {
	apiLabels := make([]string, len(apiNames))
	for i, name := range apiNames {
		apiLabels[i] = string(name)
	}

	parcelDeliveredCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "parcel_delivered_total",
		Help: "Parcel have been delivered, no requests will hit APIs",
	})
	prometheus.MustRegister(parcelDeliveredCounter)

	fetchedChanged := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "fetched_changed_total",
		Help: "API have been fetched and API response does not match cached one",
	}, apiLabels)
	prometheus.MustRegister(fetchedChanged)

	fetchedFirst := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "fetched_first_total",
		Help: "API have been fetched for the first time",
	}, apiLabels)
	prometheus.MustRegister(fetchedFirst)

	fetchedUnchanged := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "fetched_unchanged_total",
		Help: "API have been fetched and API response is the same as cached one",
	}, apiLabels)
	prometheus.MustRegister(fetchedUnchanged)

	apiHit := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_hit_total",
		Help: "Total amount of API hits",
	}, apiLabels)
	prometheus.MustRegister(apiHit)

	apiParseError := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "api_parse_error_total",
		Help: "API response fetched successfully, but we failed to parse its response",
	}, apiLabels)
	prometheus.MustRegister(apiParseError)

	cacheBustAfterSuccess := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cache_bust_after_success_total",
		Help: "After inspecting cached response that was successful, we decided to bust it and refetch API",
	}, apiLabels)
	prometheus.MustRegister(cacheBustAfterSuccess)

	cacheHitAfterSuccess := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cache_hit_after_success_total",
		Help: "After inspecting cached response that was successful, we decided to use it and not refetch API",
	}, apiLabels)
	prometheus.MustRegister(cacheHitAfterSuccess)

	cacheBustAfterUnknownError := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cache_bust_after_unknown_error_total",
		Help: "After inspecting cached response that had status 'unknown error', we decided to bust it and refetch API",
	}, apiLabels)
	prometheus.MustRegister(cacheBustAfterUnknownError)

	cacheHitAfterUnknownError := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cache_hit_after_unknown_error_total",
		Help: "After inspecting cached response that had status 'unknown error', we decided to use it and not refetch API",
	}, apiLabels)
	prometheus.MustRegister(cacheHitAfterUnknownError)

	cacheBustAfterNotFoundError := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cache_bust_after_not_found_error_total",
		Help: "After inspecting cached response that had status 'not found', we decided to bust it and refetch API",
	}, apiLabels)
	prometheus.MustRegister(cacheBustAfterNotFoundError)

	cacheHitAfterNotFoundError := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cache_hit_after_not_found_error_total",
		Help: "After inspecting cached response that had status 'not found', we decided to use it and not refetch API",
	}, apiLabels)
	prometheus.MustRegister(cacheHitAfterNotFoundError)

	return &PrometheusMetrics{
		parcelDeliveredCounter:      parcelDeliveredCounter,
		fetchedChangedCounter:       fetchedChanged,
		fetchedFirstCounter:         fetchedFirst,
		fetchedUnchanged:            fetchedUnchanged,
		apiHit:                      apiHit,
		apiParseError:               apiParseError,
		cacheBustAfterSuccess:       cacheBustAfterSuccess,
		cacheHitAfterSuccess:        cacheHitAfterSuccess,
		cacheBustAfterUnknownError:  cacheBustAfterUnknownError,
		cacheHitAfterUnknownError:   cacheHitAfterUnknownError,
		cacheBustAfterNotFoundError: cacheBustAfterNotFoundError,
		cacheHitAfterNotFoundError:  cacheHitAfterNotFoundError,
	}
}

type PrometheusMetrics struct {
	parcelDeliveredCounter      prometheus.Counter
	fetchedChangedCounter       *prometheus.CounterVec
	fetchedFirstCounter         *prometheus.CounterVec
	fetchedUnchanged            *prometheus.CounterVec
	apiHit                      *prometheus.CounterVec
	apiParseError               *prometheus.CounterVec
	cacheBustAfterSuccess       *prometheus.CounterVec
	cacheHitAfterSuccess        *prometheus.CounterVec
	cacheBustAfterUnknownError  *prometheus.CounterVec
	cacheHitAfterUnknownError   *prometheus.CounterVec
	cacheBustAfterNotFoundError *prometheus.CounterVec
	cacheHitAfterNotFoundError  *prometheus.CounterVec
}

func (p *PrometheusMetrics) ParcelDelivered() {
	p.parcelDeliveredCounter.Inc()
}

func (p *PrometheusMetrics) FetchedChanged(apiName service.APIName) {
	p.fetchedChangedCounter.WithLabelValues(string(apiName)).Inc()
}

func (p *PrometheusMetrics) FetchedFirst(apiName service.APIName) {
	p.fetchedFirstCounter.WithLabelValues(string(apiName)).Inc()
}

func (p *PrometheusMetrics) FetchedUnchanged(apiName service.APIName) {
	p.fetchedUnchanged.WithLabelValues(string(apiName)).Inc()
}

func (p *PrometheusMetrics) APIHit(name service.APIName) {
	p.apiHit.WithLabelValues(string(name)).Inc()
}

func (p *PrometheusMetrics) APIParseError(apiName service.APIName) {
	p.apiParseError.WithLabelValues(string(apiName)).Inc()
}

func (p *PrometheusMetrics) CacheBustAfterSuccess(apiName service.APIName, willRefetch bool) {
	if willRefetch {
		p.cacheBustAfterSuccess.WithLabelValues(string(apiName)).Inc()
	} else {
		p.cacheHitAfterSuccess.WithLabelValues(string(apiName)).Inc()
	}
}

func (p *PrometheusMetrics) CacheBustAfterUnknownError(apiName service.APIName, willRefetch bool) {
	if willRefetch {
		p.cacheBustAfterUnknownError.WithLabelValues(string(apiName)).Inc()
	} else {
		p.cacheHitAfterUnknownError.WithLabelValues(string(apiName)).Inc()
	}
}

func (p *PrometheusMetrics) CacheBustAfterNotFoundError(apiName service.APIName, willRefetch bool) {
	if willRefetch {
		p.cacheBustAfterNotFoundError.WithLabelValues(string(apiName)).Inc()
	} else {
		p.cacheHitAfterNotFoundError.WithLabelValues(string(apiName)).Inc()
	}
}

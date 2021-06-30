package msclient

import . "github.com/prometheus/client_golang/prometheus"

type metrics struct {
	requests            *CounterVec
	requestErrors       *CounterVec
	statusCodeErrors    *CounterVec
	responseReadErrors  *CounterVec
	responseParseErrors *CounterVec
	successfulResponses *CounterVec

	successfulResponseTimes *HistogramVec
}

const labelURI = "uri"
const labelMethod = "method"

var labels = []string{labelURI, labelMethod}

func newMetrics(r *Registry, mxSubsystem string) *metrics {
	cVec := func(name string) *CounterVec {
		return NewCounterVec(CounterOpts{Subsystem: mxSubsystem, Name: name}, labels)
	}

	mx := &metrics{
		requests:            cVec("requests"),
		requestErrors:       cVec("requestErrors"),
		statusCodeErrors:    cVec("statusCodeErrors"),
		responseReadErrors:  cVec("responseReadErrors"),
		responseParseErrors: cVec("responseParseErrors"),
		successfulResponses: cVec("successfulResponses"),

		successfulResponseTimes: NewHistogramVec(HistogramOpts{
			Subsystem:   mxSubsystem,
			Name:        "successfulResponseTimes",
			Buckets:     []float64{0.01, 0.1, 0.5, 1.0, 3.0, 5.0, 10.0, 15.0, 20.0},
		}, labels),
	}
	r.MustRegister(mx.requests, mx.requestErrors, mx.statusCodeErrors, mx.responseReadErrors, mx.responseParseErrors,
		mx.successfulResponses, mx.successfulResponseTimes)
	return mx
}

package transport

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/valyala/fasthttp"
)

var RequestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "request_counter",
	Help: "Total number of requests",
}, []string{"method", "status"})

var TimeCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "time_request",
	Help: "Total",
}, []string{"method", "status"})

func metrics(f func(ctx *fasthttp.RequestCtx), methodName string) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		Now := time.Now()
		f(ctx)
		TimeWorkF := time.Now().Sub(Now)
		RequestCounter.WithLabelValues(methodName, strconv.Itoa(ctx.Response.StatusCode())).Inc()
		TimeCounter.WithLabelValues(methodName, strconv.Itoa(ctx.Response.StatusCode())).Add(float64(TimeWorkF))
	}

}

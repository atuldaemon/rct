package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"os/signal"
	"syscall"

	"github.com/atuldaemon/rct/booking"
	"github.com/atuldaemon/rct/parking"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	fieldKeys := []string{"method"}

	parkingStore, err := parking.NewInMemParkingStore()
	if err != nil {
		panic(err)
	}
	var p parking.Service
	{
		p = parking.NewService(parkingStore)
		p = parking.LoggingMiddleware(logger)(p)
		p = parking.NewInstrumentingService(
			kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
				Namespace: "api",
				Subsystem: "parking_service",
				Name:      "request_count",
				Help:      "Number of requests received.",
			}, fieldKeys),
			kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
				Namespace: "api",
				Subsystem: "parking_service",
				Name:      "request_latency_microseconds",
				Help:      "Total duration of requests in microseconds.",
			}, fieldKeys),
			p)
	}

	bookingStore, err := booking.NewInMemBookingStore()
	if err != nil {
		panic(err)
	}
	var b booking.Service
	{
		b = booking.NewService(bookingStore, p)
		b = booking.LoggingMiddleware(logger)(b)
		b = booking.NewInstrumentingService(
			kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
				Namespace: "api",
				Subsystem: "booking_service",
				Name:      "request_count",
				Help:      "Number of requests received.",
			}, fieldKeys),
			kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
				Namespace: "api",
				Subsystem: "booking_service",
				Name:      "request_latency_microseconds",
				Help:      "Total duration of requests in microseconds.",
			}, fieldKeys),
			b)
	}

	mux := http.NewServeMux()

	mux.Handle("/parking/v1/", parking.MakeHTTPHandler(p, log.With(logger, "component", "HTTP")))
	mux.Handle("/booking/v1/", booking.MakeHTTPHandler(b, log.With(logger, "component", "HTTP")))

	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())

	errs := make(chan error, 2)

	go func() {
		logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)

}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

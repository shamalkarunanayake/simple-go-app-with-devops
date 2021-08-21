package main

import (
	"fmt"
	"html"
	"math/rand"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	counter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "golang",
			Name:      "my_counter",
			Help:      "This is my counter",
		})

	gauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "golang",
			Name:      "my_gauge",
			Help:      "This is my gauge",
		})

	histogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "golang",
			Name:      "my_histogram",
			Help:      "This is my histogram",
		})

	summary = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace: "golang",
			Name:      "my_summary",
			Help:      "This is my summary",
		})
)

func main() {

	log := logrus.New()
	log.SetFormatter(&ecslogrus.Formatter{})

	epoch := time.Unix(0, 0).UTC()
	log.WithTime(epoch).WithField("custom", "foo").Info("hello")

	logrus.Print("Logging in Go!")

	for a := 0; a < 1000; a++ {
		log.Print("Logging in Go!%d\n", a)
	}

	logrus.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	log.WithFields(logrus.Fields{
		"omg":    false,
		"number": 11,
	}).Warn("The group's number increased tremendously!")

	//log.WithFields(log.Fields{
	//	"omg":    false,
	//	"number": 102,
	//}).Fatal("The ice breaks!")

	// A common pattern is to re-use fields between logging statements by re-using
	// the logrus.Entry returned from WithFields()
	contextLogger := log.WithFields(logrus.Fields{
		"common": "this is a common field",
		"other":  "I also should be logged always",
	})

	contextLogger.Info("I'll be logged with common and other field")
	contextLogger.Info("Me too")
	for b := 0; b < 1000; b++ {
		log.Print("Counter logs !%d\n", b)
	}
	prometheus.MustRegister(counter)
	prometheus.MustRegister(gauge)
	prometheus.MustRegister(histogram)
	prometheus.MustRegister(summary)

	go func() {
		for {
			counter.Add(rand.Float64() * 5)
			gauge.Add(rand.Float64()*15 - 5)
			histogram.Observe(rand.Float64() * 10)
			summary.Observe(rand.Float64() * 10)

			time.Sleep(time.Second)
		}
	}()

	standardFields := logrus.Fields{
		"hostname": "staging-1",
		"appname":  "go-app",
		"session":  "1ce3f6v",
	}

	log.WithFields(standardFields).WithFields(logrus.Fields{"string": "foo", "int": 1, "float": 1.1}).Info("My first ssl event from Golang")

	histogramVec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "prom_request_time",
		Help: "Time it has taken to retrieve the metrics",
	}, []string{"@timestamp"})

	prometheus.Register(histogramVec)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi there, I love Devops %s!", html.EscapeString(r.URL.Path))
	})
	http.Handle("/metrics", newHandlerWithHistogram(promhttp.Handler(), histogramVec))

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func newHandlerWithHistogram(handler http.Handler, histogram *prometheus.HistogramVec) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		status := http.StatusOK

		defer func() {
			histogram.WithLabelValues(fmt.Sprintf("%d", status)).Observe(time.Since(start).Seconds())
		}()

		if req.Method == http.MethodGet {
			handler.ServeHTTP(w, req)
			return
		}
		status = http.StatusBadRequest

		w.WriteHeader(status)
	})
}

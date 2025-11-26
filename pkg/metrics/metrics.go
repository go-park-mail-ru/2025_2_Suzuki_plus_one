package metrics

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Serve(log logger.Logger, serveString string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Info("Metrics server is starting at " + serveString + "/metrics")
	if err := http.ListenAndServe(serveString, nil); err != nil {
		log.Fatal("Failed to start metrics server: " + err.Error())
	}
}

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"net/http"
	"strconv"
)

var (
	prometheusPort     = 9090
	warningGaugeGroup  *prometheus.GaugeVec
	criticalGaugeGroup *prometheus.GaugeVec
)

func init() {
	var labelNames = []string{"coin", "name", "token_id", "provider", "human_name"}
	var warningOpts = prometheus.GaugeOpts{
		Name: "market_time_lag",
		Help: "Minute value. Tickers no updated last N min.",
	}
	var criticalOpts = prometheus.GaugeOpts{
		Name: "market_critical_time_lag",
		Help: "Minute value. Critical tickers no updated last N min.",
	}
	warningGaugeGroup = prometheus.NewGaugeVec(warningOpts, labelNames)
	criticalGaugeGroup = prometheus.NewGaugeVec(criticalOpts, labelNames)

	prometheus.MustRegister(warningGaugeGroup, criticalGaugeGroup)
	http.Handle("/metrics", promhttp.Handler())
	var err = http.ListenAndServe(strconv.Itoa(prometheusPort), nil)
	if err != nil {
		logger.Error("Prometheus metrics not started")
	}
}

func observeTickers(tickers []watchmarket.Ticker) {
	for i := 0; i < len(tickers); i++ {

	}
}

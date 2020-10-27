package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/db/postgres"
	"golang.org/x/net/context"
	"math"
	"net/http"
	"time"
)

var (
	prometheusPort     = 9090
	warningGaugeGroup  *prometheus.GaugeVec
	criticalGaugeGroup *prometheus.GaugeVec
	database postgres.Instance
)

func Init(instance postgres.Instance) {
	database  = instance
	var labelNames = []string{"coin", "name", "token_id", "provider", "human_name"}
	var warningOpts = prometheus.GaugeOpts{
		Name: "market_time_lag",
		Help: "Minute value. Tickers not updated last N min.",
	}
	var criticalOpts = prometheus.GaugeOpts{
		Name: "market_critical_time_lag",
		Help: "Minute value. Critical tickers not updated last N min.",
	}
	warningGaugeGroup = prometheus.NewGaugeVec(warningOpts, labelNames)
	criticalGaugeGroup = prometheus.NewGaugeVec(criticalOpts, labelNames)

	prometheus.MustRegister(warningGaugeGroup, criticalGaugeGroup)
	http.Handle("/metrics", promhttp.Handler())
	var err = http.ListenAndServe(fmt.Sprintf(":%d", prometheusPort), nil)
	if err != nil {
		logger.Error("Prometheus metrics not started")
	}
}

func RefreshMetrics() {
	var result, err = database.GetAllTickers(context.Background())
	if err != nil {
		return
	}

	for i := 0; i < len(result); i++ {
		var item = &result[i]
		var isCritical = isCriticalCoin(item)
		if isCritical {
			updateMetricsFor(criticalGaugeGroup, item)
		}
		if !isCritical {
			updateMetricsFor(warningGaugeGroup, item)
		}
	}
	return
}

func isCriticalCoin(ticker *models.Ticker) bool {
	return ticker != nil
}

func updateMetricsFor(gauge *prometheus.GaugeVec, ticker *models.Ticker) {
	//"coin", "name", "token_id", "provider", "human_name"
	var timeWithoutUpdates = time.Now().Sub(ticker.LastUpdated).Minutes()
	var roundedTime = math.Ceil(timeWithoutUpdates)
	gauge.WithLabelValues(fmt.Sprint(ticker.Coin), ticker.CoinName, ticker.TokenId, ticker.Provider, ticker.CoinName).Set(roundedTime)
}

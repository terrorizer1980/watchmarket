package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/db/postgres"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"golang.org/x/net/context"
	"math"
	"net/http"
	"time"
)

type Instance struct {
	warningGaugeGroup   *prometheus.GaugeVec
	criticalGaugeGroup  *prometheus.GaugeVec
	tickersCounterGroup *prometheus.CounterVec
	database            postgres.Instance
}

//todo: Remove before merge
var (
	TempCurrent *Instance
)

const (
	prometheusPort = 9090
)

func Init(db postgres.Instance) *Instance {
	labelNames := []string{"coin", "name", "token_id", "provider", "human_name"}
	warningOpts := prometheus.GaugeOpts{
		Name: "market_time_lag",
		Help: "Minute value. Tickers not updated last N min.",
	}
	criticalOpts := prometheus.GaugeOpts{
		Name: "market_critical_time_lag",
		Help: "Minute value. Critical tickers not updated last N min.",
	}
	tickersOpts := prometheus.CounterOpts{
		Name: "market_ticker_requests",
		Help: "Minute value. Critical tickers not updated last N min.",
	}

	warningGaugeGroup := prometheus.NewGaugeVec(warningOpts, labelNames)
	criticalGaugeGroup := prometheus.NewGaugeVec(criticalOpts, labelNames)
	tickersCounterGroup := prometheus.NewCounterVec(tickersOpts, labelNames)

	prometheus.MustRegister(warningGaugeGroup, criticalGaugeGroup, tickersCounterGroup)
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%d", prometheusPort), nil)

	return &Instance{
		warningGaugeGroup:   warningGaugeGroup,
		criticalGaugeGroup:  criticalGaugeGroup,
		tickersCounterGroup: tickersCounterGroup,
		database:            db,
	}
}

func (i *Instance) RefreshMetrics() {
	var result, err = i.database.GetAllTickers(context.Background())
	if err != nil {
		return
	}

	for _, item := range result {
		isCritical := isCriticalCoin(&item)
		if isCritical {
			updateGaugeMetricsFor(i.criticalGaugeGroup, &item)
		}
		if !isCritical {
			updateGaugeMetricsFor(i.warningGaugeGroup, &item)
		}
	}
	return
}

func isCriticalCoin(ticker *models.Ticker) bool {
	return ticker != nil
}

func updateGaugeMetricsFor(gauge *prometheus.GaugeVec, ticker *models.Ticker) {
	//"coin", "name", "token_id", "provider", "human_name"
	timeWithoutUpdates := time.Now().Sub(ticker.LastUpdated).Minutes()
	roundedTime := math.Ceil(timeWithoutUpdates)
	gauge.WithLabelValues(fmt.Sprint(ticker.Coin), ticker.CoinName, ticker.TokenId, ticker.Provider, ticker.CoinName).Set(roundedTime)
}

func (i *Instance) UpdateTickersCounterMetricsFor(ticker *watchmarket.Ticker) {
	//"coin", "name", "token_id", "provider", "human_name"
	i.tickersCounterGroup.WithLabelValues(fmt.Sprint(ticker.Coin), ticker.CoinName, ticker.TokenId, "", ticker.CoinName).Inc()
}

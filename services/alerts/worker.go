package alerts

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db"
	"github.com/trustwallet/watchmarket/db/models"
)

type Worker struct {
	database      db.Instance
	configuration config.Configuration
}

func Init(db db.Instance, cfg config.Configuration) *Worker {
	return &Worker{
		database:      db,
		configuration: cfg,
	}
}

func (w Worker) Run() {
	ctx := context.Background()
	assets, err := w.database.GetSubscribedAssets(ctx)
	if err != nil {
		log.Error(err)
		return
	}
	tickers, err := w.database.GetTickersByAssetIDs(assets, ctx)
	if err != nil {
		log.Error(err)
		return
	}
	assetsToTickersMap := w.getAssetsToTickersMap(assets, tickers)
	for asset, ticker := range assetsToTickersMap {
		subscribedUsers, err := w.database.GetSubscribedUsers(asset, ticker.Value, true, ctx)
		if err != nil {
			log.Error(err)
			continue
		}
		err = Notify(subscribedUsers, asset, ticker.Value, true)
		if err != nil {
			log.Error(err)
			continue
		}
	}
}

func (w Worker) getAssetsToTickersMap(assets []string, tickers []models.Ticker) map[string]models.Ticker {
	result := make(map[string]models.Ticker)
	for _, ticker := range tickers {
		current, ok := result[ticker.ID]
		if !ok {
			result[ticker.ID] = ticker
		} else {
			if ticker.ShowOption == models.AlwaysShow {
				result[ticker.ID] = ticker
				continue
			}
			if isBetterProvider(current.Provider, ticker.Provider, w.configuration.Markets.Priority.CoinInfo) {
				result[ticker.ID] = ticker
				continue
			}
		}
	}
	return result
}

func isBetterProvider(oldProvider, currentProvider string, providers []string) bool {
	for _, p := range providers {
		switch {
		case p == oldProvider && p != currentProvider:
			return false
		case p == currentProvider && p != oldProvider:
			return true
		case p == currentProvider && p == oldProvider:
			return false
		default:
			continue
		}
	}
	return false
}

func Notify(users []string, assetID string, price float64, higher bool) error {
	log.WithFields(log.Fields{"users": users, "asset": assetID, "price": price, "higher": higher}).
		Info("Notified")
	return nil
}

package db

import (
	"context"
	"github.com/trustwallet/watchmarket/db/models"
)

type (
	Instance interface {
		GetRates(currency string, ctx context.Context) ([]models.Rate, error)
		GetAllRates(ctx context.Context) ([]models.Rate, error)
		AddRates(rates []models.Rate, batchLimit uint, ctx context.Context) error

		AddTickers(tickers []models.Ticker, batchLimit uint, ctx context.Context) error
		GetTickers(coin uint, tokenId string, ctx context.Context) ([]models.Ticker, error)
		GetAllTickers(ctx context.Context) ([]models.Ticker, error)
		GetTickersByQueries(tickerQueries []models.TickerQuery, ctx context.Context) ([]models.Ticker, error)
		GetTickersByAssetIDs(assetIDs []string, ctx context.Context) ([]models.Ticker, error)

		GetSubscribedAssets(ctx context.Context) ([]string, error)
		GetSubscribedUsers(assetID string, price float64, higher bool, ctx context.Context) ([]string, error)
	}
)

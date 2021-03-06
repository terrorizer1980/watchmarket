package coinmarketcap

import (
	"context"
	"github.com/trustwallet/watchmarket/pkg/watchmarket"
	"strings"
)

func (p Provider) GetRates(ctx context.Context) (rates watchmarket.Rates, err error) {
	prices, err := p.client.fetchPrices(p.currency, ctx)
	if err != nil {
		return
	}
	rates = normalizeRates(prices, p.id)
	return
}

func normalizeRates(prices CoinPrices, provider string) watchmarket.Rates {
	var (
		result        watchmarket.Rates
		emptyPlatform Platform
	)

	for _, price := range prices.Data {
		if price.Platform != emptyPlatform {
			continue
		}
		result = append(result, watchmarket.Rate{
			Currency:         strings.ToUpper(price.Symbol),
			Rate:             watchmarket.TruncateWithPrecision(price.Quote.USD.Price, watchmarket.DefaultPrecision),
			Timestamp:        price.LastUpdated.Unix(),
			PercentChange24h: price.Quote.USD.PercentChange24h,
			Provider:         provider,
		})
	}
	return result
}

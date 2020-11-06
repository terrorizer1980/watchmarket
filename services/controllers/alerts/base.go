package alertscontroller

import (
	"context"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/services/controllers"
)

type Controller struct {
	database      db.Instance
	configuration config.Configuration
}

func NewController(
	database db.Instance,
	configuration config.Configuration,
) Controller {
	return Controller{
		database,
		configuration,
	}
}

func (c Controller) HandleSubscriptionsRequest(sr controllers.SubscriptionsRequest, ctx context.Context) error {
	subs := toSubscriptionsModel(sr.Subscriptions)
	err := c.database.AddPriceSubscriptions(subs, ctx)
	if err != nil {
		return err
	}
	return nil
}

func toSubscriptionsModel(subs []controllers.Subscription) []models.PriceSubscription {
	result := make([]models.PriceSubscription, 0, len(subs))
	for _, s := range subs {
		result = append(result, models.PriceSubscription{
			UserID:   s.UserID,
			Price:    s.Price,
			IsHigher: s.Higher,
			AssetID:  s.AssetID,
		})
	}
	return result
}

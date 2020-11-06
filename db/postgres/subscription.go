package postgres

import (
	"context"
	"github.com/trustwallet/watchmarket/db/models"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmgorm"
)

func (i *Instance) GetSubscribedAssets(ctx context.Context) ([]string, error) {
	g := apmgorm.WithContext(ctx, i.Gorm)
	span, _ := apm.StartSpan(ctx, "GetSubscribedAssets", "postgresql")
	defer span.End()

	var result []string
	g.Model(&models.PriceSubscription{})
	//err := g.Table("price_subscriptions").
	//	Where("asset_id like ?", "distinct%").
	//	Order("asset_id").
	//	Pluck("asset_id", &result).Error
	//if err != nil {
	//	return nil, err
	//}

	return result, nil
}

func (i *Instance) GetSubscribedUsers(assetID string, price float64, higher bool, ctx context.Context) ([]string, error) {
	return []string{"1", "2", "3"}, nil
}

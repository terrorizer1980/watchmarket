package postgres

import (
	"context"
	"github.com/trustwallet/watchmarket/db/models"
	"gorm.io/gorm/clause"
)

func (i *Instance) GetSubscribedAssets(ctx context.Context) ([]string, error) {
	var result []string
	err := i.Gorm.Model(&models.PriceSubscription{}).Distinct("asset_id").Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (i *Instance) GetSubscribedUsers(assetID string, price float64, higher bool, ctx context.Context) ([]string, error) {
	var subs []models.PriceSubscription
	err := i.Gorm.Model(&models.PriceSubscription{}).
		Where("asset_id = ? and price <= ? and is_higher = ?", assetID, price, higher).
		Find(&subs).Error
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(subs))
	for _, sub := range subs {
		result = append(result, sub.UserID)
	}
	return result, nil
}

func (i *Instance) AddPriceSubscription(sub models.PriceSubscription, ctx context.Context) error {
	return i.Gorm.Create(sub).Error
}

func (i *Instance) AddPriceSubscriptions(subs []models.PriceSubscription, ctx context.Context) error {
	return i.Gorm.Clauses(clause.OnConflict{DoNothing: true}).Create(&subs).Error
}

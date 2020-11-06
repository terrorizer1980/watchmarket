package postgres

import (
	"context"
)

func (i *Instance) GetSubscribedAssets(ctx context.Context) ([]string, error) {
	//g := apmgorm.WithContext(ctx, i.Gorm)
	//span, _ := apm.StartSpan(ctx, "GetSubscribedAssets", "postgresql")
	//defer span.End()

	return nil, nil
}

func (i *Instance) GetSubscribedUsers(assetID string, price float64, higher bool, ctx context.Context) ([]string, error) {
	return nil, nil
}

package alerts

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db"
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
	for _, asset := range assets {
		subscribedUsers, err := w.database.GetSubscribedUsers(asset, 0, true, ctx)
		if err != nil {
			log.Error(err)
			continue
		}
		err = Notify(subscribedUsers, asset, 0, true)
		if err != nil {
			log.Error(err)
			continue
		}
	}
}

func Notify(users []string, assetID string, price float64, higher bool) error {
	return nil
}
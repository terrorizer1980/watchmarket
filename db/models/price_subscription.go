package models

type PriceSubscription struct {
	UserID   string
	Price    float64
	IsHigher bool
	AssetID  string
}

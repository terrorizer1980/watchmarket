package models

type PriceSubscription struct {
	UserID   string `gorm:"primaryKey; index:,"`
	Price    float64
	IsHigher bool
	AssetID  string `gorm:"primaryKey; index:,"`
}

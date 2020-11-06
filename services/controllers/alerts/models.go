package alertscontroller

type SubscriptionsRequest struct {
	Subscriptions []Subscription `json:"subscriptions"`
}

type Subscription struct {
	AssetID string  `json:"asset_id"`
	Price   float64 `json:"price"`
	Higher  bool    `json:"higher"`
	UserID  string  `json:"user_id"`
}

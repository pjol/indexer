package indexer

import (
	"time"
)

type Listener struct {
	Owner        string    `json:"listener_owner"`
	OwnerId      string    `json:"owner_id"`
	LocationId   string    `json:"location_id"`
	Contract     string    `json:"contract"`
	TokenName    string    `json:"token_name"`
	Address      string    `json:"address"`
	Service      string    `json:"service"`
	Secret       string    `json:"secret"`
	Value        int       `json:"value"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expires_at"`
}

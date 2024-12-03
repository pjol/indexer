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

type DeleteRequest struct {
	OwnerId    string `json:"owner_id"`
	LocationId string `json:"location_id"`
	Service    string `json:"service"`
	Contract   string `json:"contract"`
	Secret     string `json:"secret"`
}

type KeyRequest struct {
	Id       string `json:"id"`
	Key      string `json:"key"`
	Contract string `json:"contract"`
	Name     string `json:"name"`
	Owner    string `json:"owner"`
}

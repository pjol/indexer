package listeners

import (
	"net/http"

	"github.com/citizenwallet/indexer/internal/services/db"
	"github.com/citizenwallet/indexer/pkg/indexer"
)

type Service struct {
	db *db.DB
}

func NewService(db *db.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) AddListener(w http.ResponseWriter, r *http.Request) {
	var l indexer.Listener
	_ = l
}

package listeners

import (
	"encoding/json"
	"net/http"

	"github.com/citizenwallet/indexer/internal/common"
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

//TODO: Add auth to AddListener

func (s *Service) AddListener(w http.ResponseWriter, r *http.Request) {
	var l indexer.Listener
	err := json.NewDecoder(r.Body).Decode(&l)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	suffix, err := s.db.TableNameSuffix(l.Contract)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lsdb, ok := s.db.ListenersDB[suffix]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = lsdb.AddListener(&l)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = common.Body(w, l, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

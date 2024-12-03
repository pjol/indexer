package listeners

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/citizenwallet/indexer/internal/common"
	"github.com/citizenwallet/indexer/internal/services/db"
	"github.com/citizenwallet/indexer/pkg/indexer"
	"github.com/google/uuid"
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
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if l.Service == "ZAPIER" {
		l.Value = 1
	}

	fmt.Println(l)
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

func (s *Service) DeleteListener(w http.ResponseWriter, r *http.Request) {
	var d indexer.DeleteRequest

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch d.Service {
	case "ZAPIER":
		key := r.Header["X-Api-Key"][0]
		suffix, err := s.db.TableNameSuffix(d.Contract)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		lsdb, ok := s.db.ListenersDB[suffix]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		valid, err := lsdb.GetKeyExists(key)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		err = lsdb.RemoveListener(&d)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func (s *Service) NewKey(w http.ResponseWriter, r *http.Request) {
	var k indexer.KeyRequest

	err := json.NewDecoder(r.Body).Decode(&k)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(k)

	k.Id = uuid.New().String()

	suffix, err := s.db.TableNameSuffix(k.Contract)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lsdb, ok := s.db.ListenersDB[suffix]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	key, err := lsdb.MakeKey(&k)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf(`{"key": "%s"}`, key)))
	return
}

func (s *Service) DeleteKey(w http.ResponseWriter, r *http.Request) {
	var k indexer.KeyRequest

	err := json.NewDecoder(r.Body).Decode(&k)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	suffix, err := s.db.TableNameSuffix(k.Contract)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lsdb, ok := s.db.ListenersDB[suffix]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = lsdb.RemoveKey(k.Key)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	return
}

func (s *Service) ExampleZapierResponse(w http.ResponseWriter, r *http.Request) {

	body := []byte(`[{
    "amount": 1,
    "txhash": "0x0",
    "from": "0x0"
  }]`)

	w.Write(body)
}

func (s *Service) AuthTest(w http.ResponseWriter, r *http.Request) {
	contract := r.URL.Query()["contract"][0]
	fmt.Println(contract)
	fmt.Println(r.Header)
	key := r.Header["X-Api-Key"][0]

	if contract == "" || key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	suffix, err := s.db.TableNameSuffix(contract)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lsdb, ok := s.db.ListenersDB[suffix]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	valid, err := lsdb.GetKeyExists(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "success"}`))
}

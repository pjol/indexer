package db

import (
	"database/sql"
	"fmt"

	"github.com/citizenwallet/indexer/pkg/indexer"
)

type ListenersDB struct {
	suffix string
	db     *sql.DB
	rdb    *sql.DB
}

func NewListenersDB(db, rdb *sql.DB, name string) (*ListenersDB, error) {
	lsdb := &ListenersDB{
		suffix: name,
		db:     db,
		rdb:    rdb,
	}
	return lsdb, nil
}

func (db *ListenersDB) Close() error {
	return db.db.Close()
}

func (db *ListenersDB) CloseR() error {
	return db.rdb.Close()
}

func (db *ListenersDB) CreateListenersTable() error {
	_, err := db.db.Exec(fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS t_listeners_%s(
		listener_owner text NOT NULL,
		owner_id text NOT NULL,
		location_id text,
		contract text NOT NULL,
		token_name text NOT NULL,
		address text NOT NULL,
		service text NOT NULL,
		secret text NOT NULL,
		value int NOT NULL,
		refresh_token,
		expires_at timestamp,
		created_at timestamp NOT NULL DEFAULT current_timestamp,
		updated_at timestamp NOT NULL DEFAULT current_timestamp,
		UNIQUE (address, service)
	);
	`, db.suffix))

	return err
}

func (db *ListenersDB) CreateListenersAuthTable() error {
	_, err := db.db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS t_listeners_auth_%s(
      id text PRIMARY KEY,
      owner text NOT NULL,
			name text NOT NULL,
      key text NOT NULL,
      UNIQUE(key)
		);
	`, db.suffix))

	return err
}

func (db *ListenersDB) CreateListenersTableIndexes() error {
	_, err := db.db.Exec(fmt.Sprintf(`
    CREATE INDEX IF NOT EXISTS idx_listeners_%s_owner ON t_listeners_%s (listener_owner);
    `, db.suffix, db.suffix))

	if err != nil {
		return err
	}

	_, err = db.db.Exec(fmt.Sprintf(`
    CREATE INDEX IF NOT EXISTS idx_listeners_%s_address ON t_listeners_%s (address);
    `, db.suffix, db.suffix))

	if err != nil {
		return err
	}

	return nil
}

func (db *ListenersDB) GetListenerDetails(address string) ([]*indexer.Listener, error) {
	rows, err := db.rdb.Query(fmt.Sprintf(`
		SELECT listener_owner, location_id, contract, token_name, address, service, secret, value
		FROM t_listeners_%s
		WHERE address = $1
	`, db.suffix), address)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	listeners := []*indexer.Listener{}
	for rows.Next() {
		var listener indexer.Listener
		err = rows.Scan(&listener.Owner, &listener.LocationId, &listener.Contract, &listener.TokenName, &listener.Address, &listener.Service, &listener.Secret, &listener.Value)
		if err != nil {
			return nil, err
		}

		listeners = append(listeners, &listener)
	}

	return listeners, nil
}

func (db *ListenersDB) AddListener(l *indexer.Listener) error {
	_, err := db.db.Exec(fmt.Sprintf(`
		INSERT INTO t_listeners_%s (listener_owner, owner_id, location_id, contract, token_name, address, service, secret, value, refresh_token, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, db.suffix), l.Owner, l.OwnerId, l.LocationId, l.Contract, l.TokenName, l.Address, l.Service, l.Secret, l.Value, l.RefreshToken, l.Expiry)

	if err != nil {
		return err
	}

	return nil
}

func (db *ListenersDB) RemoveListener(l *indexer.DeleteRequest) error {
	_, err := db.db.Exec(fmt.Sprintf(`
		DELETE FROM t_listeners_%s WHERE (owner_id = $1 AND secret = $2 AND service = $3 AND location_id = $4);
	`, db.suffix), l.OwnerId, l.Secret, l.Service, l.LocationId)

	if err != nil {
		return err
	}

	return nil
}

func (db *ListenersDB) MakeKey(k *indexer.KeyRequest) (string, error) {
	key := indexer.RandomString(32)

	_, err := db.db.Exec(fmt.Sprintf(`
		INSERT INTO t_listeners_auth_%s (id, owner, key, name) VALUES ($1, $2, $3, $4);
	`, db.suffix), k.Id, k.Owner, key, k.Name)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return key, nil
}

func (db *ListenersDB) GetKeyExists(key string) (bool, error) {
	row := db.db.QueryRow(fmt.Sprintf(`
		SELECT key FROM t_listeners_auth_%s WHERE key = $1;
	`, db.suffix), key)

	var exists string
	err := row.Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (db *ListenersDB) RemoveKey(key string) error {
	_, err := db.db.Exec(fmt.Sprintf(`
		DELETE FROM t_listeners_auth_%s WHERE key = $1;
	`, db.suffix), key)
	if err != nil {
		return err
	}

	db.db.Exec(fmt.Sprintf(`
		DELETE FROM t_listeners_%s WHERE secret = $1;
	`, db.suffix), key)

	return nil
}

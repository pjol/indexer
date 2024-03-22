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
		contract text NOT NULL,
		address text NOT NULL,
		service text NOT NULL,
		secret text NOT NULL,
		value int NOT NULL,
		created_at timestamp NOT NULL DEFAULT current_timestamp,
		updated_at timestamp NOT NULL DEFAULT current_timestamp,
		UNIQUE (address, service)
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
		SELECT listener_owner, contract, address, service, secret, value
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
		err = rows.Scan(&listener.Owner, &listener.Contract, &listener.Address, &listener.Service, &listener.Secret, &listener.Value)
		if err != nil {
			return nil, err
		}

		listeners = append(listeners, &listener)
	}

	return listeners, nil
}

func (db *ListenersDB) AddListener(l *indexer.Listener) error {
	_, err := db.db.Exec(fmt.Sprintf(`
		INSERT INTO t_listeners_%s (listener_owner, contract, address, service, secret, value)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, db.suffix), l.Owner, l.Contract, l.Address, l.Service, l.Secret, l.Value)

	if err != nil {
		return err
	}

	return nil
}

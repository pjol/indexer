package db

import (
	"database/sql"
	"fmt"
)

type SquareDB struct {
	suffix string
	db     *sql.DB
	rdb    *sql.DB
}

func NewSquareDB(db, rdb *sql.DB, name string) (*SquareDB, error) {
	sqrdb := &SquareDB{
		suffix: name,
		db:     db,
		rdb:    rdb,
	}
	return sqrdb, nil
}

func (db *SquareDB) Close() error {
	return db.db.Close()
}

func (db *SquareDB) CloseR() error {
	return db.rdb.Close()
}

func (db *SquareDB) CreateSquareTable(suffix string) error {
	_, err := db.db.Exec(fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS t_square_%s(
		merchant_name text NOT NULL
		address text NOT NULL
		merchant_secret text NOT NULL
		created_at timestamp NOT NULL DEFAULT current_timestamp,
		updated_at timestamp NOT NULL DEFAULT current_timestamp,
		UNIQUE(address)
	);
	`, suffix))

	return err
}

func (db *SquareDB) CreateSquareTableIndexes(suffix string) error {
	_, err := db.db.Exec(fmt.Sprintf(`
    CREATE INDEX IF NOT EXISTS idx_square_%s_merchant ON t_square_%s (merchant_name);
    `, suffix, suffix))

	if err != nil {
		return err
	}

	return nil
}

func (db *SquareDB) GetMerchantSecret(address string) (string, error) {
	var secret string
	err := db.rdb.QueryRow(fmt.Sprintf(`
		SELECT merchant_secret
		FROM t_square_%s
		WHERE address = $1
	`, db.suffix), address).Scan(secret)

	if err != nil {
		return "", err
	}

	return secret, nil
}

func (db *SquareDB) AddMerchant(name string, address string, secret string) error {
	_, err := db.db.Exec(fmt.Sprintf(`
		INSERT INTO t_square_%s (merchant_name, address, merchant_secret)
		VALUES ($1, $2, $3)
	`, db.suffix), name, address, secret)

	if err != nil {
		return err
	}

	return nil
}

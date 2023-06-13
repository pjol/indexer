package db

import (
	"fmt"
	"log"
	"math/big"

	"github.com/citizenwallet/node/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

const (
	dbBaseFolder = ".cw"
)

type DB struct {
	EventDB    *EventDB
	TransferDB map[string]*TransferDB
}

// NewDB instantiates a new DB
func NewDB(chainID *big.Int) (*DB, error) {
	basePath := storage.GetUserHomeDir()
	folderPath := fmt.Sprintf("%s/%s", basePath, dbBaseFolder)
	path := fmt.Sprintf("%s/events_%v.db", folderPath, chainID)

	// check if directory exists
	if !storage.Exists(folderPath) {
		err := storage.CreateDir(folderPath)
		if err != nil {
			return nil, err
		}
	}

	eventDB, err := NewEventDB(path)
	if err != nil {
		return nil, err
	}

	txdb := map[string]*TransferDB{}

	evs, err := eventDB.GetEvents()
	if err != nil {
		return nil, err
	}

	for _, ev := range evs {
		name := TransferName(chainID, ev.Contract)
		log.Default().Println("creating transfer db for: ", name)
		txdb[name], err = NewTransferDB(name)
		if err != nil {
			return nil, err
		}
	}

	return &DB{
		EventDB:    eventDB,
		TransferDB: txdb,
	}, nil
}

// TransferName returns the name of the transfer db for the given contract
func TransferName(chainID *big.Int, contract string) string {
	return fmt.Sprintf("%v_%s", chainID, contract)
}

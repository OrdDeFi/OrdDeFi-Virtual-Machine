package db_utils

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"log"
)

const CurrentDBVersion = "1"

const (
	AvailableSubAccount    = "a"
	TransferableSubAccount = "t"
)

type OrdDB struct {
	db *leveldb.DB
}

func OpenDB(path string) (*OrdDB, error) {
	levelDB, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	ordDB := OrdDB{
		db: levelDB,
	}

	return &ordDB, nil
}

func CloseDB(db *OrdDB) {
	err := db.db.Close()
	if err != nil {
		log.Println("Failed to close DB:", err)
	}
}

/*
Keeps string-in-string-out principle in all public data functions.
Cast from []byte to string automatically make a copy on []byte, to get avoid of stack memory issue and data interrupted.
*/

func (db OrdDB) Store(key string, value string) error {
	err := db.db.Put([]byte(key), []byte(value), nil)
	return err
}

func (db OrdDB) Read(key string) (*string, error) {
	byteRes, err := db.db.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}
	res := string(byteRes)
	return &res, nil
}

func (db OrdDB) ReadAllPrefix(prefixString string) (map[string]string, error) {
	prefix := []byte(prefixString)

	iter := db.db.NewIterator(util.BytesPrefix(prefix), nil)
	defer iter.Release()

	var result map[string]string
	result = make(map[string]string)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		result[string(key)] = string(value)
	}
	if err := iter.Error(); err != nil {
		return nil, err
	}
	return result, nil
}

func (db OrdDB) StoreKeyValues(keyValues map[string]string) error {
	batch := new(leveldb.Batch)
	for key, value := range keyValues {
		if value == "" {
			batch.Delete([]byte(key))
		} else {
			batch.Put([]byte(key), []byte(value))
		}
	}
	err := db.db.Write(batch, nil)
	return err
}

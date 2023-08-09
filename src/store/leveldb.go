package store

import (
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LevelDB struct {
	driver *leveldb.DB
}

func NewLevelDB() (*LevelDB, error) {
	db, err := leveldb.OpenFile("data/leveldb", nil)
	if err != nil {
		return nil, err
	}
	return &LevelDB{
		driver: db,
	}, nil
}

func (db *LevelDB) Put(key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return db.driver.Put([]byte(key), bytes, nil)
}

func (db *LevelDB) Get(key string, value interface{}) error {
	bytes, err := db.driver.Get([]byte(key), nil)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, value)
}

func (db *LevelDB) Search(prefix string) []string {
	var items = make([]string, 0)
	iter := db.driver.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer iter.Release()

	for iter.Next() {
		items = append(items, string(iter.Value()))
	}
	return items
}

func (db *LevelDB) Delete(key string) error {
	return db.driver.Delete([]byte(key), nil)
}

// Close release resources
func (db *LevelDB) Close() error {
	return db.driver.Close()
}

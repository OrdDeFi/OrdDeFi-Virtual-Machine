package memory_read

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
	"errors"
)

/*
CoinMeta read coin metadata from db
*/
func CoinMeta(db *db_utils.OrdDB, coinName string) (*memory_const.CoinMeta, error) {
	if coinName == "" {
		return nil, errors.New("read CoinMeta error: coin name is empty")
	}
	key := memory_const.CoinMetadataTable + ":" + coinName
	value, err := db.Read(key)
	if err != nil {
		if err.Error() == "leveldb: not found" {
			return nil, nil
		} else {
			return nil, err
		}
	}
	if value == nil {
		return nil, errors.New("read CoinMeta error: coin not found in db")
	}
	coinMeta, err := memory_const.CoinMetaFromJsonString(*value)
	if err != nil {
		return nil, err
	}
	return coinMeta, nil

}

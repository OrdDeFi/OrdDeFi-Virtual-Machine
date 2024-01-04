package memory_read

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
	"errors"
)

/*
Across all versions.
*/

func ODFIMeta() *memory_const.CoinMeta {
	result := new(memory_const.CoinMeta)
	result.Max = safe_number.SafeNumFromString("21000000")
	result.Lim = safe_number.SafeNumFromString("1000")
	result.AddrLim = safe_number.SafeNumFromString("1000")
	result.Desc = "The protocol major coin of OrdDeFi."
	result.Icon = ""
	return result
}

func ODGVMeta() *memory_const.CoinMeta {
	result := new(memory_const.CoinMeta)
	result.Max = safe_number.SafeNumFromString("21000000")
	result.Lim = safe_number.SafeNumFromString("1000")
	result.AddrLim = safe_number.SafeNumFromString("1000")
	result.Desc = "The governance coin of OrdDeFi."
	result.Icon = ""
	return result
}

/*
CoinMeta read coin metadata from db
*/
func CoinMeta(db *db_utils.OrdDB, coinName string) (*memory_const.CoinMeta, error) {
	if coinName == "" {
		return nil, errors.New("read CoinMeta error: coin name is empty")
	} else if coinName == "odfi" {
		return ODFIMeta(), nil
	} else if coinName == "odgv" {
		return ODGVMeta(), nil
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

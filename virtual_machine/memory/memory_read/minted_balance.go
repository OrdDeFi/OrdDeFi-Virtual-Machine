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

/*
TotalMintedBalance read coin minted balance from db
*/
func TotalMintedBalance(db *db_utils.OrdDB, coinName string) (*safe_number.SafeNum, error) {
	key := memory_const.TotalMintedBalanceTable + ":" + coinName
	value, err := db.Read(key)
	if err != nil {
		if err.Error() == "leveldb: not found" {
			return nil, nil
		} else {
			return nil, err
		}
	}
	if value == nil {
		return nil, errors.New("read TotalMintedBalance error: coin not found in db")
	}
	minted := safe_number.SafeNumFromString(*value)
	if minted == nil {
		return safe_number.SafeNumFromString("0"), nil
	}
	return minted, nil
}

/*
AddressMintedBalance read coin minted balance from db
*/
func AddressMintedBalance(db *db_utils.OrdDB, coinName string, address string) (*safe_number.SafeNum, error) {
	key := memory_const.AddressMintedBalanceTable + ":" + coinName + ":" + address
	value, err := db.Read(key)
	if err != nil {
		if err.Error() == "leveldb: not found" {
			return nil, nil
		} else {
			return nil, err
		}
	}
	if value == nil {
		return nil, errors.New("read AddressMintedBalanceTable error: coin or address not found in db")
	}
	minted := safe_number.SafeNumFromString(*value)
	if minted == nil {
		return safe_number.SafeNumFromString("0"), nil
	}
	return minted, nil
}

package memory_read

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
)

/*
For a single version.
*/

func AvailableBalance(db *db_utils.OrdDB, coinName string, address string) (*safe_number.SafeNum, error) {
	balanceKey := memory_const.CoinAddressAvailablePath(coinName, address)
	balanceString, err := db.Read(balanceKey)
	if err != nil {
		if err.Error() == "leveldb: not found" {
			value := "0"
			balanceString = &value
		} else {
			return nil, err
		}
	}
	num := safe_number.SafeNumFromString(*balanceString)
	return num, nil
}

func TransferableBalance(db *db_utils.OrdDB, coinName string, address string) (*safe_number.SafeNum, error) {
	balanceKey := memory_const.CoinAddressTransferablePath(coinName, address)
	balanceString, err := db.Read(balanceKey)
	if err != nil {
		if err.Error() == "leveldb: not found" {
			value := "0"
			balanceString = &value
		} else {
			return nil, err
		}
	}
	num := safe_number.SafeNumFromString(*balanceString)
	return num, nil
}

func Balance(db *db_utils.OrdDB, coinName string, address string) (*safe_number.SafeNum, *safe_number.SafeNum, error) {
	available, err := AvailableBalance(db, coinName, address)
	if err != nil {
		return nil, nil, err
	}
	transferable, err := TransferableBalance(db, coinName, address)
	if err != nil {
		return nil, nil, err
	}
	return available, transferable, nil
}

func AllCoinBalanceForAddress(db *db_utils.OrdDB, address string) (map[string]string, error) {
	prefix := memory_const.AddressCoinPrefix(address)
	result, err := db.ReadAllPrefix(prefix)
	return result, err
}

func AllAddressBalanceForCoin(db *db_utils.OrdDB, coinName string) (map[string]string, error) {
	prefix := memory_const.CoinAddressPrefix(coinName)
	result, err := db.ReadAllPrefix(prefix)
	return result, err
}

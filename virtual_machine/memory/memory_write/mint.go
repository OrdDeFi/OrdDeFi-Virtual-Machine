package memory_write

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
)

func WriteMintInfo(
	db *db_utils.OrdDB,
	coinName string,
	address string,
	newTotalMintedString string,
	newAddressMintedString string,
	newBalanceString string) error {
	totalMintedKey := memory_const.TotalMintedBalanceTable + ":" + coinName
	addressMintedKey := memory_const.AddressMintedBalanceTable + ":" + coinName + ":" + address

	batchWriting := CoinBalanceDoubleWriteKV(coinName, address, newBalanceString, db_utils.AvailableSubAccount)
	batchWriting[totalMintedKey] = newTotalMintedString
	batchWriting[addressMintedKey] = newAddressMintedString
	err := db.StoreKeyValues(batchWriting)
	return err
}

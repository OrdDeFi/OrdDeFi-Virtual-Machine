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
	newBalanceString string,
	version string) error {
	totalMintedKey := memory_const.TotalMintedBalanceTable + ":" + coinName
	addressMintedKey := memory_const.AddressMintedBalanceTable + ":" + coinName + ":" + address
	balanceKey1 := memory_const.CoinBalanceTable + ":v" + version + ":" + coinName + ":" + address
	balanceKey2 := memory_const.AddressBalanceTable + ":v" + version + ":" + address + ":" + coinName
	// Generate batch writing map
	var batchWriting map[string]string
	batchWriting = make(map[string]string)
	batchWriting[totalMintedKey] = newTotalMintedString
	batchWriting[addressMintedKey] = newAddressMintedString
	batchWriting[balanceKey1] = newBalanceString
	batchWriting[balanceKey2] = newBalanceString
	err := db.StoreKeyValues(batchWriting)
	return err
}

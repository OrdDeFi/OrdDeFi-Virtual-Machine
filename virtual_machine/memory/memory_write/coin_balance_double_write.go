package memory_write

import (
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
)

func CoinBalanceDoubleWriteKVForAvailable(coinName string, address string, newBalanceString string) map[string]string {
	balanceKey1 := memory_const.CoinAddressAvailablePath(coinName, address)
	balanceKey2 := memory_const.AddressCoinAvailablePath(coinName, address)
	// Generate batch writing map
	var batchWriting map[string]string
	batchWriting = make(map[string]string)
	batchWriting[balanceKey1] = newBalanceString
	batchWriting[balanceKey2] = newBalanceString
	return batchWriting
}

func CoinBalanceDoubleWriteKVForTransferable(coinName string, address string, newBalanceString string) map[string]string {
	balanceKey1 := memory_const.CoinAddressTransferablePath(coinName, address)
	balanceKey2 := memory_const.AddressCoinTransferablePath(coinName, address)
	// Generate batch writing map
	var batchWriting map[string]string
	batchWriting = make(map[string]string)
	batchWriting[balanceKey1] = newBalanceString
	batchWriting[balanceKey2] = newBalanceString
	return batchWriting
}

package memory_write

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
)

func coinBalanceDoubleWriteKVForAvailable(coinName string, address string, newBalanceString string) map[string]string {
	balanceKey1 := memory_const.CoinAddressAvailablePath(coinName, address)
	balanceKey2 := memory_const.AddressCoinAvailablePath(coinName, address)
	// Generate batch writing map
	var batchWriting map[string]string
	batchWriting = make(map[string]string)
	batchWriting[balanceKey1] = newBalanceString
	batchWriting[balanceKey2] = newBalanceString
	return batchWriting
}

func coinBalanceDoubleWriteKVForTransferable(coinName string, address string, newBalanceString string) map[string]string {
	balanceKey1 := memory_const.CoinAddressTransferablePath(coinName, address)
	balanceKey2 := memory_const.AddressCoinTransferablePath(coinName, address)
	// Generate batch writing map
	var batchWriting map[string]string
	batchWriting = make(map[string]string)
	batchWriting[balanceKey1] = newBalanceString
	batchWriting[balanceKey2] = newBalanceString
	return batchWriting
}

func CoinBalanceDoubleWriteKV(coinName string, address string, newBalanceString string, subAccount string) map[string]string {
	if subAccount == db_utils.AvailableSubAccount {
		return coinBalanceDoubleWriteKVForAvailable(coinName, address, newBalanceString)
	} else if subAccount == db_utils.TransferableSubAccount {
		return coinBalanceDoubleWriteKVForTransferable(coinName, address, newBalanceString)
	} else {
		return nil
	}

}

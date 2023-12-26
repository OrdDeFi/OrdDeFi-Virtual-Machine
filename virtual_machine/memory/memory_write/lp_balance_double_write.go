package memory_write

import "OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"

func LPBalanceDoubleWriteKV(lCoin string, rCoin string, address string, newBalanceString string) map[string]string {
	balanceKey1 := memory_const.LPAddressPath(lCoin, rCoin, address)
	balanceKey2 := memory_const.AddressLPPath(lCoin, rCoin, address)
	// Generate batch writing map
	var batchWriting map[string]string
	batchWriting = make(map[string]string)
	batchWriting[balanceKey1] = newBalanceString
	batchWriting[balanceKey2] = newBalanceString
	return batchWriting
}

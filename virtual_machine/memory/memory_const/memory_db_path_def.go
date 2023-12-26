package memory_const

import "OrdDeFi-Virtual-Machine/db_utils"

func CoinAddressPrefix(coinName string) string {
	path := CoinAddressBalanceTable + ":v" + db_utils.CurrentDBVersion + ":" + coinName + ":"
	return path
}

func CoinAddressAvailablePath(coinName string, address string) string {
	path := CoinAddressPrefix(coinName) + address + ":" + db_utils.AvailableSubAccount
	return path
}

func CoinAddressTransferablePath(coinName string, address string) string {
	path := CoinAddressPrefix(coinName) + address + ":" + db_utils.TransferableSubAccount
	return path
}

func AddressCoinPrefix(address string) string {
	path := AddressCoinBalanceTable + ":v" + db_utils.CurrentDBVersion + ":" + address + ":"
	return path
}

func AddressCoinAvailablePath(coinName string, address string) string {
	path := AddressCoinPrefix(address) + coinName + ":" + db_utils.AvailableSubAccount
	return path
}

func AddressCoinTransferablePath(coinName string, address string) string {
	path := AddressCoinPrefix(address) + coinName + ":" + db_utils.TransferableSubAccount
	return path
}

func UTXOCarryingBalancePath(txId string) string {
	return UTXOCarryingBalanceTable + ":" + txId + ":0" /*indicates the output index, which is always 0 */
}

func LPAddressPrefix(lpName string) string {
	path := LPAddressBalanceTable + ":v" + db_utils.CurrentDBVersion + ":" + lpName + ":"
	return path
}

func LPAddressPath(lpName string, address string) string {
	path := LPAddressPrefix(lpName) + address
	return path
}

func AddressLPPrefix(address string) string {
	path := AddressLPBalanceTable + ":v" + db_utils.CurrentDBVersion + ":" + address + ":"
	return path
}

func AddressLPPath(lpName string, address string) string {
	path := AddressLPPrefix(address) + lpName
	return path
}

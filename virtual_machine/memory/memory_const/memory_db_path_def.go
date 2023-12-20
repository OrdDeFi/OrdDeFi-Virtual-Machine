package memory_const

import "OrdDeFi-Virtual-Machine/db_utils"

func CoinAddressPrefix(coinName string) string {
	path := CoinBalanceTable + ":v" + db_utils.CurrentDBVersion + ":" + coinName + ":"
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
	path := CoinBalanceTable + ":v" + db_utils.CurrentDBVersion + ":" + address + ":"
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

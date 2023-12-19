package memory_const

func CoinAddressPrefix(coinName string, version string) string {
	path := CoinBalanceTable + ":v" + version + ":" + coinName + ":"
	return path
}

func CoinAddressAvailablePath(coinName string, address string, version string) string {
	path := CoinAddressPrefix(coinName, version) + address + ":a"
	return path
}

func CoinAddressTransferablePath(coinName string, address string, version string) string {
	path := CoinAddressPrefix(coinName, version) + address + ":t"
	return path
}

func AddressCoinPrefix(address string, version string) string {
	path := CoinBalanceTable + ":v" + version + ":" + address + ":"
	return path
}

func AddressCoinAvailablePath(coinName string, address string, version string) string {
	path := AddressCoinPrefix(address, version) + coinName + ":a"
	return path
}

func AddressCoinTransferablePath(coinName string, address string, version string) string {
	path := AddressCoinPrefix(address, version) + coinName + ":t"
	return path
}

package main

import (
	"brc20defi_vm/bitcoin_cli_channel"
	"brc20defi_vm/updater"
)

func main() {
	println("The Times 03/Jan/2009 Chancellor on brink of second bailout for banks.")
	println("BRC-20-DEFI indexer start to work.")
	blockNumber := bitcoin_cli_channel.GetBlockCount()
	updater.UpdateBlockNumber(blockNumber)
}

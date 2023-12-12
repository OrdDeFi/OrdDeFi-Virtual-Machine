package main

import (
	"brc20defi_vm/bitcoin_cli_channel"
	"brc20defi_vm/inscription_parser"
	"brc20defi_vm/updater"
	"flag"
)

func updateIndex() {
	println("BRC-20-DEFI indexer start to work.")
	blockNumber := bitcoin_cli_channel.GetBlockCount()
	updater.UpdateBlockNumber(blockNumber)
}

func parseRawTransaction(parseRawTransactionString string) {
	r, err := inscription_parser.ParseRawTransactionToInscription(parseRawTransactionString)
	if err != nil {
		println("parserawtransaction error:", err)
	} else {
		println(*r)
	}
}

func main() {
	println("The Times 03/Jan/2009 Chancellor on brink of second bailout for banks.")

	var parseRawTransactionString string
	flag.StringVar(&parseRawTransactionString, "parserawtransaction", "", "Parse Raw Transaction")
	flag.Parse()

	if parseRawTransactionString != "" {
		parseRawTransaction(parseRawTransactionString)
	} else {
		updateIndex()
	}
}

package main

import (
	"OrdDefi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDefi-Virtual-Machine/inscription_parser"
	"OrdDefi-Virtual-Machine/updater"
	"errors"
	"flag"
)

func updateIndex() {
	println("OrdDefi indexer start to work.")
	blockNumber := bitcoin_cli_channel.GetBlockCount()
	if blockNumber == 0 {
		return
	}
	updater.UpdateBlockNumber(blockNumber)
}

func parseRawTransaction(parseRawTransactionString string) {
	contentType, content, err := inscription_parser.ParseRawTransactionToInscription(parseRawTransactionString)
	if err != nil {
		println("parserawtransaction error:", err)
	} else {
		println(*contentType, len(content))
		println(string(content))
	}
}

func parseTransaction(txId string) error {
	rawTx := bitcoin_cli_channel.GetRawTransaction(txId)
	if rawTx == nil {
		err := errors.New("GetRawTransaction Failed")
		return err
	}
	parseRawTransaction(*rawTx)
	return nil
}

func main() {
	println("The Times 03/Jan/2009 Chancellor on brink of second bailout for banks.")

	var parseTransactionString string
	flag.StringVar(&parseTransactionString, "parsetransaction", "", "OrdDefi-Virtual-Machine -parsetransaction [txid]")
	var parseRawTransactionString string
	flag.StringVar(&parseRawTransactionString, "parserawtransaction", "", "OrdDefi-Virtual-Machine -parserawtransaction [raw transaction string]")
	flag.Parse()
	if parseTransactionString != "" {
		err := parseTransaction(parseTransactionString)
		if err != nil {
			println("parseTransaction error:", err)
		}
	} else if parseRawTransactionString != "" {
		parseRawTransaction(parseRawTransactionString)
	} else {
		updateIndex()
	}
}

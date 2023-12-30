package main

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDeFi-Virtual-Machine/inscription_parser"
	"OrdDeFi-Virtual-Machine/updater"
	"errors"
	"flag"
	"os"
)

func updateIndex(dataDir string, logDir string) error {
	println("OrdDeFi indexer start to work.")
	blockNumber := bitcoin_cli_channel.GetBlockCount()
	if blockNumber == 0 {
		err := errors.New("updateIndex error: bitcoin-cli getblockcount failed")
		return err
	}
	err := updater.UpdateBlockNumber(blockNumber, dataDir, logDir)
	return err
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
	flag.StringVar(&parseTransactionString, "parsetransaction", "", "OrdDeFi-Virtual-Machine -parsetransaction [txid]")
	var parseRawTransactionString string
	flag.StringVar(&parseRawTransactionString, "parserawtransaction", "", "OrdDeFi-Virtual-Machine -parserawtransaction [raw transaction string]")
	var parseDataDir string
	flag.StringVar(&parseDataDir, "data-dir", "", "OrdDeFi-Virtual-Machine -data-dir /path/of/storage")
	var parseLogDir string
	flag.StringVar(&parseLogDir, "log-dir", "", "OrdDeFi-Virtual-Machine -log-dir /path/of/log")
	flag.Parse()
	if parseDataDir == parseLogDir && parseDataDir != "" {
		println("-data-dir and -log-dir should be different")
		os.Exit(2)
	}
	if parseTransactionString != "" {
		err := parseTransaction(parseTransactionString)
		if err != nil {
			println("parseTransaction error:", err)
		}
	} else if parseRawTransactionString != "" {
		parseRawTransaction(parseRawTransactionString)
	} else {
		err := updateIndex(parseDataDir, parseLogDir)
		if err != nil {
			os.Exit(1)
		}
	}
}

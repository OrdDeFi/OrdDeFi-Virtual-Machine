package main

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDeFi-Virtual-Machine/subcommands"
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

func main() {
	println("The Times 03/Jan/2009 Chancellor on brink of second bailout for banks.")

	// subcommands
	var parseTransactionParam string
	flag.StringVar(&parseTransactionParam, "parsetransaction", "", "OrdDeFi-Virtual-Machine -parsetransaction [txid]")
	var parseRawTransactionParam string
	flag.StringVar(&parseRawTransactionParam, "parserawtransaction", "", "OrdDeFi-Virtual-Machine -parserawtransaction [raw transaction string]")
	var executeResultParam string
	flag.StringVar(&executeResultParam, "executeresult", "", "OrdDeFi-Virtual-Machine -executeresult [txid]")
	var checkUTXOTransferParam string
	flag.StringVar(&checkUTXOTransferParam, "checkutxotransfer", "", "OrdDeFi-Virtual-Machine -checkutxotransfer [txid:0]")

	// updater params
	var parseDataDir string
	flag.StringVar(&parseDataDir, "data-dir", "", "OrdDeFi-Virtual-Machine -data-dir /path/of/storage")
	var parseLogDir string
	flag.StringVar(&parseLogDir, "log-dir", "", "OrdDeFi-Virtual-Machine -log-dir /path/of/log")
	flag.Parse()
	if parseDataDir == parseLogDir && parseDataDir != "" {
		println("-data-dir and -log-dir should be different")
		os.Exit(2)
	}
	if parseTransactionParam != "" {
		subcommands.ParseTransaction(parseTransactionParam)
	} else if parseRawTransactionParam != "" {
		subcommands.ParseRawTransaction(parseRawTransactionParam)
	} else if executeResultParam != "" {
		subcommands.CheckExecuteResult(executeResultParam)
	} else if checkUTXOTransferParam != "" {
		subcommands.CheckUTXOTransfer(checkUTXOTransferParam)
	} else {
		err := updateIndex(parseDataDir, parseLogDir)
		if err != nil {
			os.Exit(1)
		}
	}
}

package main

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDeFi-Virtual-Machine/subcommands"
	"OrdDeFi-Virtual-Machine/updater"
	"errors"
	"flag"
	"os"
	"strings"
)

func updateIndex(dataDir string, logDir string, verbose bool) error {
	println("OrdDeFi indexer start to work.")
	blockNumber := bitcoin_cli_channel.GetBlockCount()
	if blockNumber == 0 {
		err := errors.New("updateIndex error: bitcoin-cli getblockcount failed")
		return err
	}
	err := updater.UpdateBlockNumber(blockNumber, dataDir, logDir, verbose)
	return err
}

func main() {
	println("The Times 03/Jan/2009 Chancellor on brink of second bailout for banks.")

	// DB path
	var dataDirParam string
	flag.StringVar(&dataDirParam, "data-dir", "", "OrdDeFi-Virtual-Machine -data-dir /path/of/storage")
	var logDirParam string
	flag.StringVar(&logDirParam, "log-dir", "", "OrdDeFi-Virtual-Machine -log-dir /path/of/log")
	if dataDirParam == logDirParam && dataDirParam != "" {
		println("-data-dir and -log-dir should be different")
		os.Exit(2)
	}

	// verbose
	var verboseParam string
	flag.StringVar(&verboseParam, "verbose", "", "OrdDeFi-Virtual-Machine -verbose true")

	// subcommands
	var parseTransactionParam string
	flag.StringVar(&parseTransactionParam, "parsetransaction", "", "OrdDeFi-Virtual-Machine -parsetransaction [txid]")
	var parseRawTransactionParam string
	flag.StringVar(&parseRawTransactionParam, "parserawtransaction", "", "OrdDeFi-Virtual-Machine -parserawtransaction [raw transaction string]")
	var executeResultParam string
	flag.StringVar(&executeResultParam, "executeresult", "", "OrdDeFi-Virtual-Machine -executeresult [txid]")
	var checkUTXOTransferParam string
	flag.StringVar(&checkUTXOTransferParam, "checkutxotransfer", "", "OrdDeFi-Virtual-Machine -checkutxotransfer [txid:0]")
	flag.Parse()

	if parseTransactionParam != "" {
		subcommands.ParseTransaction(parseTransactionParam)
	} else if parseRawTransactionParam != "" {
		subcommands.ParseRawTransaction(parseRawTransactionParam)
	} else if executeResultParam != "" {
		subcommands.CheckExecuteResult(executeResultParam, logDirParam)
	} else if checkUTXOTransferParam != "" {
		subcommands.CheckUTXOTransfer(checkUTXOTransferParam, dataDirParam)
	} else {
		verboseBool := strings.ToLower(verboseParam) == "true"
		err := updateIndex(dataDirParam, logDirParam, verboseBool)
		if err != nil {
			os.Exit(1)
		}
	}
}

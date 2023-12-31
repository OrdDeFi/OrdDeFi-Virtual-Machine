package main

import (
	"OrdDeFi-Virtual-Machine/subcommands"
	"OrdDeFi-Virtual-Machine/updater"
	"flag"
	"os"
	"strings"
)

func main() {
	// DB path
	var dataDirParam string
	flag.StringVar(&dataDirParam, "data-dir", "", "OrdDeFi-Virtual-Machine -data-dir /path/of/storage")
	var logDirParam string
	flag.StringVar(&logDirParam, "log-dir", "", "OrdDeFi-Virtual-Machine -log-dir /path/of/log")
	if dataDirParam == "" {
		dataDirParam = "./OrdDeFi_storage"
	}
	if logDirParam == "" {
		logDirParam = "./OrdDeFi_log"
	}
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
	var getAddressBalanceParam string
	flag.StringVar(&getAddressBalanceParam, "getaddressbalance", "", "OrdDeFi-Virtual-Machine -getaddressbalance [address]")
	var getAddressLPBalanceParam string
	flag.StringVar(&getAddressLPBalanceParam, "getaddresslpbalance", "", "OrdDeFi-Virtual-Machine -getaddresslpbalance [address]")
	var getLPAddressBalanceParam string
	flag.StringVar(&getLPAddressBalanceParam, "getlpaddressbalance", "", "OrdDeFi-Virtual-Machine -getlpaddressbalance [coinA-coinB]")
	var getCoinHoldersParam string
	flag.StringVar(&getCoinHoldersParam, "getcoinholders", "", "OrdDeFi-Virtual-Machine -getcoinholders [coin]")
	var getCoinMetaParam string
	flag.StringVar(&getCoinMetaParam, "getcoinmeta", "", "OrdDeFi-Virtual-Machine -getcoinmeta [coinName]")
	var getLPMetaParam string
	flag.StringVar(&getLPMetaParam, "getlpmeta", "", "OrdDeFi-Virtual-Machine -getlpmeta [coinA-coinB]")
	var getAllCoinsParam string
	flag.StringVar(&getAllCoinsParam, "getallcoins", "", "OrdDeFi-Virtual-Machine -getallcoins true")
	var getAllLPsParam string
	flag.StringVar(&getAllLPsParam, "getalllps", "", "OrdDeFi-Virtual-Machine -getalllps true")
	flag.Parse()

	if parseTransactionParam != "" {
		subcommands.ParseTransaction(parseTransactionParam)
	} else if parseRawTransactionParam != "" {
		subcommands.ParseRawTransaction(parseRawTransactionParam)
	} else if executeResultParam != "" {
		subcommands.CheckExecuteResult(executeResultParam, logDirParam)
	} else if checkUTXOTransferParam != "" {
		subcommands.CheckUTXOTransfer(checkUTXOTransferParam, dataDirParam)
	} else if getAddressBalanceParam != "" {
		subcommands.GetAddressBalance(getAddressBalanceParam, dataDirParam)
	} else if getAddressLPBalanceParam != "" {
		subcommands.GetAddressLPBalance(getAddressLPBalanceParam, dataDirParam)
	} else if getLPAddressBalanceParam != "" {
		subcommands.GetLPAddressBalance(getLPAddressBalanceParam, dataDirParam)
	} else if getCoinHoldersParam != "" {
		subcommands.GetCoinHolders(getCoinHoldersParam, dataDirParam)
	} else if getCoinMetaParam != "" {
		subcommands.GetCoinMeta(getCoinMetaParam, dataDirParam)
	} else if getLPMetaParam != "" {
		subcommands.GetLPMeta(getLPMetaParam, dataDirParam)
	} else if getAllCoinsParam != "" {
		subcommands.GetAllCoins(dataDirParam)
	} else if getAllLPsParam != "" {
		subcommands.GetAllLPs(dataDirParam)
	} else {
		verboseBool := strings.ToLower(verboseParam) == "true"
		err := updater.UpdateIndex(dataDirParam, logDirParam, verboseBool)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
	}
}

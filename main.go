package main

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDeFi-Virtual-Machine/rpc_server"
	"OrdDeFi-Virtual-Machine/subcommands"
	"OrdDeFi-Virtual-Machine/updater"
	"flag"
	"fmt"
	"os"
	"time"
)

func updateIndexTask(dataDirParam string, logDirParam string, verboseParam bool) {
	err := updater.UpdateIndex(dataDirParam, logDirParam, verboseParam)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func main() {
	// Version
	var versionParam bool
	flag.BoolVar(&versionParam, "version", false, "OrdDeFi-Virtual-Machine -version")

	// DB path
	var dataDirParam string
	flag.StringVar(&dataDirParam, "data-dir", "", "OrdDeFi-Virtual-Machine -data-dir /path/of/storage")
	var logDirParam string
	flag.StringVar(&logDirParam, "log-dir", "", "OrdDeFi-Virtual-Machine -log-dir /path/of/log")

	// Bitcoin-cli path
	var bitcoinCliParamPath string
	flag.StringVar(&bitcoinCliParamPath, "bitcoin-cli-param-file", "", "OrdDeFi-Virtual-Machine -bitcoin-cli-param-file /path/of/bitcoin-cli-param-file")

	// verbose
	var verboseParam bool
	flag.BoolVar(&verboseParam, "verbose", false, "OrdDeFi-Virtual-Machine -verbose")
	// server
	var serverParam bool
	flag.BoolVar(&serverParam, "server", false, "OrdDeFi-Virtual-Machine -server")
	var portParam int
	flag.IntVar(&portParam, "port", 9332, "OrdDeFi-Virtual-Machine -server -port [port_number], default port is 9332")

	// subcommands
	var parseTransactionParam string
	flag.StringVar(&parseTransactionParam, "parsetransaction", "", "OrdDeFi-Virtual-Machine -parsetransaction [txid]")
	var parseRawTransactionParam string
	flag.StringVar(&parseRawTransactionParam, "parserawtransaction", "", "OrdDeFi-Virtual-Machine -parserawtransaction [raw transaction string]")
	var executeResultParam string
	flag.StringVar(&executeResultParam, "executeresult", "", "OrdDeFi-Virtual-Machine -executeresult [txid]")
	var allExecuteResultParam string
	flag.StringVar(&allExecuteResultParam, "allexecuteresult", "", "OrdDeFi-Virtual-Machine -allexecuteresult [block_number|all]")
	var checkUTXOTransferParam string
	flag.StringVar(&checkUTXOTransferParam, "checkutxotransfer", "", "OrdDeFi-Virtual-Machine -checkutxotransfer [txid:0]")
	var getAddressUTXOTransferListParam string
	flag.StringVar(&getAddressUTXOTransferListParam, "getaddressutxotransferlist", "", "OrdDeFi-Virtual-Machine -getaddressutxotransferlist [address]")
	var getUTXOTransferListParam string
	flag.StringVar(&getUTXOTransferListParam, "getutxotransferlist", "", "OrdDeFi-Virtual-Machine -getutxotransferlist [tick|all]")
	var getUTXOTransferHistoryParam string
	flag.StringVar(&getUTXOTransferHistoryParam, "getutxotransferhistory", "", "OrdDeFi-Virtual-Machine -getutxotransferhistory [tick|all]")
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

	if versionParam {
		fmt.Println("v1.6.0")
		return
	}

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

	glBitcoinCliParams := bitcoin_cli_channel.GlobalParams()
	glBitcoinCliParams.LoadConfigPath(bitcoinCliParamPath)

	if parseTransactionParam != "" {
		subcommands.ParseTransaction(parseTransactionParam)
	} else if parseRawTransactionParam != "" {
		subcommands.ParseRawTransaction(parseRawTransactionParam)
	} else if executeResultParam != "" {
		subcommands.CheckExecuteResult(executeResultParam, logDirParam)
	} else if allExecuteResultParam != "" {
		subcommands.GetAllExecuteResult(allExecuteResultParam, logDirParam)
	} else if checkUTXOTransferParam != "" {
		subcommands.CheckUTXOTransfer(checkUTXOTransferParam, dataDirParam)
	} else if getUTXOTransferListParam != "" {
		subcommands.GetUTXOTransferList(getUTXOTransferListParam, dataDirParam)
	} else if getAddressUTXOTransferListParam != "" {
		subcommands.GetAddressUTXOTransferList(getAddressUTXOTransferListParam, dataDirParam)
	} else if getUTXOTransferHistoryParam != "" {
		subcommands.GetUTXOTransferHistory(getUTXOTransferHistoryParam, dataDirParam)
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
		println("The Times 03/Jan/2009 Chancellor on brink of second bailout for banks.")
		println("OrdDeFi indexer begin to update.")
		if serverParam {
			err := rpc_server.Serve(portParam, dataDirParam)
			if err != nil {
				println(err)
				os.Exit(1)
			}
			for {
				updateIndexTask(dataDirParam, logDirParam, verboseParam)
				time.Sleep(5 * time.Second)
			}
		} else {
			updateIndexTask(dataDirParam, logDirParam, verboseParam)
		}
	}
}

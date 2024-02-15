package subcommands

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
	"os"
	"strings"
)

func CheckExecuteResult(txId string, logDir string) {
	logDB, err := db_utils.OpenDB(logDir)
	if err != nil {
		println("open logDB error:", err.Error())
		os.Exit(4)
	}
	defer db_utils.CloseDB(logDB)

	key := memory_const.LogQueryTxTable + ":" + txId
	res, err := logDB.Read(key)
	if err != nil {
		println("read logDB error:", err.Error())
		os.Exit(5)
	}
	println(txId, "execute result:")
	components := strings.Split(*res, ";;;;;")
	for _, log := range components {
		println(log)
	}
}

func GetAllExecuteResult(blockNumber string, logDir string) {
	logDB, err := db_utils.OpenDB(logDir)
	if err != nil {
		println("open logDB error:", err.Error())
		os.Exit(4)
	}
	defer db_utils.CloseDB(logDB)

	key := memory_const.LogMainTable + ":" + blockNumber
	if strings.ToLower(blockNumber) == "all" {
		key = memory_const.LogMainTable
	}
	res, err := logDB.ReadAllPrefix(key)
	if err != nil {
		println("read logDB error:", err.Error())
		os.Exit(5)
	}
	for k, v := range res {
		println(k, ":", v)
	}
}

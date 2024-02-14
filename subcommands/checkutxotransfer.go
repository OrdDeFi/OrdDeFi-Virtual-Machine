package subcommands

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"os"
	"strings"
)

func CheckUTXOTransfer(utxo string, dataDir string) {
	db, err := db_utils.OpenDB(dataDir)
	if err != nil {
		println("open db error:", err.Error())
		os.Exit(6)
	}
	defer db_utils.CloseDB(db)

	components := strings.Split(utxo, ":")
	address, tick, amount, err := memory_read.UTXOCarryingBalance(db, components[0])
	if err != nil {
		println("CheckUTXOTransfer read UTXOCarryingBalance error:", err.Error())
		os.Exit(7)
	}
	if address == nil || tick == nil || amount == nil {
		println("No assets in UTXO:", utxo)
		return
	}
	println("From address:", *address)
	println("Tick:", *tick)
	println("Amount:", amount.String())
}

func GetUTXOTransferList(tick string, dataDir string) {
	db, err := db_utils.OpenDB(dataDir)
	if err != nil {
		println("open db error:", err.Error())
		os.Exit(27)
	}
	defer db_utils.CloseDB(db)

	r, err := memory_read.AllUTXOTransferForCoin(db, tick)
	if err != nil {
		println("GetCoinHoldersParam read AllAddressBalanceForCoin error:", err.Error())
		os.Exit(28)
	}
	for k, v := range r {
		println(k, ":", v)
	}
}

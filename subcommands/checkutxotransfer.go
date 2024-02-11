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
	println("From address:", *address)
	println("Tick:", *tick)
	println("Amount:", amount.String())
}

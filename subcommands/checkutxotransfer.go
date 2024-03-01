package subcommands

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"fmt"
	"os"
	"strings"
)

func CheckUTXOTransferData(utxo string, dataDir string) (*string, *string, *safe_number.SafeNum, error) {
	db, err := db_utils.OpenDB(dataDir)
	if err != nil {
		println("open db error:", err.Error())
		return nil, nil, nil, err
	}
	defer db_utils.CloseDB(db)

	components := strings.Split(utxo, ":")
	address, tick, amount, err := memory_read.UTXOCarryingBalance(db, components[0])
	if err != nil {
		println("CheckUTXOTransfer read UTXOCarryingBalance error:", err.Error())
		return nil, nil, nil, err
	}
	if address == nil || tick == nil || amount == nil {
		return nil, nil, nil, nil
	}
	return address, tick, amount, nil
}

func CheckUTXOTransfer(utxo string, dataDir string) {
	addressPtr, tickPtr, amount, err := CheckUTXOTransferData(utxo, dataDir)
	if err != nil {
		println("CheckUTXOTransferData error:", err.Error())
		os.Exit(6)
	}

	if addressPtr == nil || tickPtr == nil || amount == nil {
		fmt.Println("No assets in UTXO:", utxo)
		return
	}
	fmt.Println("From address:", *addressPtr)
	fmt.Println("Tick:", *tickPtr)
	fmt.Println("Amount:", amount.String())
}

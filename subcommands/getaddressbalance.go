package subcommands

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"os"
)

func GetAddressBalance(address string, dataDir string) {
	db, err := db_utils.OpenDB(dataDir)
	if err != nil {
		println("open db error:", err.Error())
		os.Exit(8)
	}
	defer db_utils.CloseDB(db)

	r, err := memory_read.AllCoinBalanceForAddress(db, address)
	if err != nil {
		println("GetAddressBalanceParam read AllCoinBalanceForAddress error:", err.Error())
		os.Exit(9)
	}
	for k, v := range r {
		println(k, ":", v)
	}
}

package subcommands

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"os"
)

func GetAddressUTXOTransferList(address string, dataDir string) {
	db, err := db_utils.OpenDB(dataDir)
	if err != nil {
		println("open db error:", err.Error())
		os.Exit(32)
	}
	defer db_utils.CloseDB(db)

	r, err := memory_read.AllUTXOTransferForAddress(db, address)
	if err != nil {
		println("GetUTXOTransferList read AllAddressBalanceForCoin error:", err.Error())
		os.Exit(33)
	}
	for k, v := range r {
		println(k, ":", v)
	}
}

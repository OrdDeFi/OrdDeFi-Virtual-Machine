package subcommands

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"os"
)

func GetAllLPs(dataDir string) {
	db, err := db_utils.OpenDB(dataDir)
	if err != nil {
		println("open db error:", err.Error())
		os.Exit(17)
	}
	defer db_utils.CloseDB(db)

	allLPs, err := memory_read.AllLiquidityProviders(db)
	if err != nil {
		println("GetAllLPs read AllLiquidityProviders error:", err.Error())
		os.Exit(18)
	}
	for _, lpName := range allLPs {
		println(lpName)
	}
}

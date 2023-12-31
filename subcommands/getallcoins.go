package subcommands

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"os"
)

func GetAllCoins(dataDir string) {
	db, err := db_utils.OpenDB(dataDir)
	if err != nil {
		println("open db error:", err.Error())
		os.Exit(15)
	}
	defer db_utils.CloseDB(db)

	allCoins, err := memory_read.AllCoins(db)
	if err != nil {
		println("GetAllCoins read AllCoins error:", err.Error())
		os.Exit(16)
	}
	for index, coinName := range allCoins {
		println(index, ":", coinName)
	}
}

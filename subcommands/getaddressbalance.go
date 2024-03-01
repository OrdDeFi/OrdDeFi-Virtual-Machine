package subcommands

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"fmt"
	"os"
)

func GetAddressBalanceData(address string, dataDir string) (map[string]string, error) {
	db, err := db_utils.OpenDB(dataDir)
	if err != nil {
		println("open db error:", err.Error())
		return nil, err
	}
	defer db_utils.CloseDB(db)

	r, err := memory_read.AllCoinBalanceForAddress(db, address)
	if err != nil {
		println("GetAddressBalanceParam read AllCoinBalanceForAddress error:", err.Error())
		return nil, err
	}
	return r, err
}

func GetAddressBalance(address string, dataDir string) {
	r, err := GetAddressBalanceData(address, dataDir)
	if err != nil {
		os.Exit(8)
	}
	for k, v := range r {
		fmt.Println(k, ":", v)
	}
}

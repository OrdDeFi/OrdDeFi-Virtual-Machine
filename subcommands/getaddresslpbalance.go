package subcommands

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"os"
	"strings"
)

func GetAddressLPBalance(address string, dataDir string) {
	db, err := db_utils.OpenDB(dataDir)
	if err != nil {
		println("open db error:", err.Error())
		os.Exit(10)
	}
	defer db_utils.CloseDB(db)

	r, err := memory_read.AllLiquidityProviderBalanceForAddress(db, address)
	if err != nil {
		println("GetAddressLPBalance read AllLiquidityProviderBalanceForAddress error:", err.Error())
		os.Exit(11)
	}
	for k, v := range r {
		println(k, ":", v)
	}
}

func GetLPAddressBalance(lpName string, dataDir string) {
	db, err := db_utils.OpenDB(dataDir)
	if err != nil {
		println("open db error:", err.Error())
		os.Exit(12)
	}
	defer db_utils.CloseDB(db)

	coinComponents := strings.Split(lpName, "-")
	if len(coinComponents) != 2 {
		println("GetLPMeta LP name parse error: LP name should be in abcd-edfg format", err.Error())
		os.Exit(13)
	}

	r, err := memory_read.AllAddressBalanceForLiquidityProvider(db, coinComponents[0], coinComponents[1])
	if err != nil {
		println("GetLPAddressBalance read AllAddressBalanceForLiquidityProvider error:", err.Error())
		os.Exit(14)
	}
	for k, v := range r {
		println(k, ":", v)
	}
}

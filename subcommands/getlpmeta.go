package subcommands

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"os"
	"strings"
)

func GetLPMeta(lpName string, dataDir string) {
	db, err := db_utils.OpenDB(dataDir)
	if err != nil {
		println("open db error:", err.Error())
		os.Exit(15)
	}
	defer db_utils.CloseDB(db)

	coinComponents := strings.Split(lpName, "-")
	if len(coinComponents) != 2 {
		println("GetLPMeta LP name parse error: LP name should be in abcd-edfg format", err.Error())
		os.Exit(16)
	}

	lpMeta, err := memory_read.LiquidityProviderMetadata(db, coinComponents[0], coinComponents[1])
	if err != nil {
		println("GetLPMeta read LiquidityProviderMetadata error:", err.Error())
		os.Exit(17)
	}
	println(lpMeta.LTick, "amount:", lpMeta.LAmt.String())
	println(lpMeta.RTick, "amount:", lpMeta.RAmt.String())
	println("LP token total supply:", lpMeta.Total.String())
}

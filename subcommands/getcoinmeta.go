package subcommands

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"os"
)

func GetCoinMeta(coinName string, dataDir string) {
	db, err := db_utils.OpenDB(dataDir)
	if err != nil {
		println("open db error:", err.Error())
		os.Exit(21)
	}
	defer db_utils.CloseDB(db)

	coinMeta, err := memory_read.CoinMeta(db, coinName)
	if err != nil {
		println("GetCoinMeta read CoinMeta error:", err.Error())
		os.Exit(22)
	}
	minted, err := memory_read.TotalMintedBalance(db, coinName)
	if err != nil {
		println("GetCoinMeta read TotalMintedBalance error:", err.Error())
		os.Exit(23)
	}
	if coinMeta == nil {
		println("Coin not exist:", coinName)
	}
	println("Name:", coinName)
	println("Description:", coinMeta.Desc)
	println("Icon(Base64):", coinMeta.Icon)
	println("Max amount(hard cap):", coinMeta.Max.String())
	println("Minted:", minted.String())
	println("Mint limit per tx:", coinMeta.Lim.String())
	println("Mint limit per address:", coinMeta.AddrLim.String())
}

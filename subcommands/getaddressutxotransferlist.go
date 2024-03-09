package subcommands

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"os"
)

func GetAddressUTXOTransferListData(address string, dataDir string) (map[string]string, error) {
	db, err := db_utils.OpenDB(dataDir)
	if err != nil {
		println("open db error:", err.Error())
		return nil, err
	}
	defer db_utils.CloseDB(db)

	r, err := memory_read.AllUTXOTransferForAddress(db, address)
	if err != nil {
		println("GetUTXOTransferList read AllAddressBalanceForCoin error:", err.Error())
		return nil, err
	}
	return r, err

}

func GetAddressUTXOTransferList(address string, dataDir string) {
	r, err := GetAddressUTXOTransferListData(address, dataDir)
	if err != nil {
		os.Exit(32)
	}
	for k, v := range r {
		println(k, ":", v)
	}
}

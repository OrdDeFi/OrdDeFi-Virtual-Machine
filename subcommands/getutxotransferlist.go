package subcommands

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type utxoListSortingPair struct {
	Key   string
	Value float64
}

func GetUTXOTransferListData(tick string, dataDir string) (map[string]string, error) {
	db, err := db_utils.OpenDB(dataDir)
	if err != nil {
		println("open db error:", err.Error())
		return nil, err
	}
	defer db_utils.CloseDB(db)

	r, err := memory_read.AllUTXOTransferForCoin(db, tick)
	if err != nil {
		println("GetUTXOTransferList read AllAddressBalanceForCoin error:", err.Error())
		return nil, err
	}
	return r, nil
}

func GetUTXOTransferList(tick string, dataDir string) {
	r, err := GetUTXOTransferListData(tick, dataDir)
	if err != nil {
		os.Exit(27)
	}
	var pairs []utxoListSortingPair
	for k, v := range r {
		floatValue, err := strconv.ParseFloat(v, 64)
		if err != nil {
			println("GetUTXOTransferList convert transfer containing value to float64 error:", err.Error())
			os.Exit(29)
		}
		pairs = append(pairs, utxoListSortingPair{k, floatValue})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value > pairs[j].Value
	})
	for _, pair := range pairs {
		fmt.Println(pair.Key, pair.Value)
	}
}

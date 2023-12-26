package memory_read

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
	"encoding/json"
)

/*
AllLiquidityProviders
Read all liquidity pairs
*/
func AllLiquidityProviders(db *db_utils.OrdDB) ([]string, error) {
	r, err := db.Read(memory_const.LpListTable)
	if err != nil && err.Error() != "leveldb: not found" {
		return nil, err
	}
	var liquidityProviders []string
	if r != nil {
		var result []string
		err2 := json.Unmarshal([]byte(*r), &result)
		if err2 != nil {
			return nil, err
		}
		if result != nil && len(result) != 0 {
			liquidityProviders = append(liquidityProviders, result...)
		}
	}
	return liquidityProviders, nil
}

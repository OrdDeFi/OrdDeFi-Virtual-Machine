package memory_read

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
	"encoding/json"
)

// AllDeployedCoins all coins except for odfi and odgv
func AllDeployedCoins(db *db_utils.OrdDB) ([]string, error) {
	r, err := db.Read(memory_const.CoinListTable)
	if err != nil && err.Error() != "leveldb: not found" {
		return nil, err
	}
	var defaultCoins []string
	if r != nil {
		var result []string
		err2 := json.Unmarshal([]byte(*r), &result)
		if err2 != nil {
			return nil, err
		}
		if result != nil && len(result) != 0 {
			defaultCoins = append(defaultCoins, result...)
		}
	}
	return defaultCoins, nil
}

// AllCoins all coins including for odfi and odgv
func AllCoins(db *db_utils.OrdDB) ([]string, error) {
	result, err := AllDeployedCoins(db)
	if err != nil && err.Error() != "leveldb: not found" {
		return nil, err
	}
	defaultCoins := []string{"odfi", "odgv"}
	defaultCoins = append(defaultCoins, result...)
	return defaultCoins, nil
}

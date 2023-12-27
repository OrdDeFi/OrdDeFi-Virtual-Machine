package memory_read

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
)

/*
LiquidityProviderBalance
Read lp token amount owning by address
*/
func LiquidityProviderBalance(db *db_utils.OrdDB, lCoin string, rCoin string, address string) (*safe_number.SafeNum, error) {
	balanceKey := memory_const.LPAddressPath(lCoin, rCoin, address)
	balanceString, err := db.Read(balanceKey)
	if err != nil {
		if err.Error() == "leveldb: not found" {
			value := "0"
			balanceString = &value
		} else {
			return nil, err
		}
	}
	num := safe_number.SafeNumFromString(*balanceString)
	return num, nil
}

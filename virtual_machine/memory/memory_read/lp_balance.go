package memory_read

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
)

/*
LiquidityProviderBalance
Read lp token amount owning by address
*/
func LiquidityProviderBalance(db *db_utils.OrdDB, coinName string, address string) (*safe_number.SafeNum, error) {
	return nil, nil
}

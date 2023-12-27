package memory_read

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
)

/*
Read each lp containing coins
*/

/*
LiquidityProviderMetadata
Read lp token total amount, and all coins contained by this lp.
return lp_token_total_amount, all_coins_contained, error
If lp not exist, return nil, nil
*/
func LiquidityProviderMetadata(db *db_utils.OrdDB, lcoinName string, rcoinName string) (*memory_const.LPMeta, error) {
	lpName := lcoinName + "-" + rcoinName
	lpMetaKey := memory_const.LpMetadataTable + ":" + lpName
	r, err := db.Read(lpMetaKey)
	if err != nil && err.Error() != "leveldb: not found" {
		return nil, err
	}
	if r == nil {
		return nil, nil
	}
	lpMeta, err := memory_const.LPMetaFromJsonString(*r)
	if err != nil {
		return nil, err
	}
	return lpMeta, nil
}

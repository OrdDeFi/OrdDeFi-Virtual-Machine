package memory_read

import "OrdDeFi-Virtual-Machine/safe_number"

/*
Read each lp containing coins
*/

/*
LiquidityProviderMetadata
Read lp token total amount, and all coins contained by this lp.
return lp_token_total_amount, all_coins_contained, error
If lp not exist, return nil, nil
If it's odfi-odgv pair, the all_coins_contained map should contain other coins trading fee.
*/
func LiquidityProviderMetadata(lcoinName string, rcoinName string) (map[string]safe_number.SafeNum, error) {
	return nil, nil
}

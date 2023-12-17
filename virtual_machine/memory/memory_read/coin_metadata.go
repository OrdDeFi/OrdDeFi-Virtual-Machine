package memory_read

import "OrdDeFi-Virtual-Machine/safe_number"

func CoinMetadataSave(
	coinName string,
	max *safe_number.SafeNum,
	lim *safe_number.SafeNum,
	addrLim *safe_number.SafeNum,
	desc string,
	icon string) {
}

/*
CoinMetadataQuery
return max, lim, addrLim, desc, icon
*/
func CoinMetadataQuery(coinName string) (*safe_number.SafeNum, *safe_number.SafeNum, *safe_number.SafeNum, *string, *string) {
	return nil, nil, nil, nil, nil
}

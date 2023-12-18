package memory_const

// CoinListTable stores all coin names
const CoinListTable = "coinlist"

// CoinMetadataTable stores deployed coins' meta
const CoinMetadataTable = "coinmeta"

/*
CoinBalanceTable stores a single coin's {address:balance}.
Path: CoinBalanceTable:coin_name:address
Value: balance string
*/
const CoinBalanceTable = "coinbalance"

// LpListTable stores all lp names
const LpListTable = "lplist"

// LpMetadataTable stores deployed lp' meta, coin value, price
const LpMetadataTable = "lpmeta"

/*
LpBalanceTable stores a single lp's {address:balance}.
Path: LpBalanceTable:lp_name:address
Value: balance string
*/
const LpBalanceTable = "lpbalance"

/*
AddressBalanceTable stores a single address's {coin|lp:balance}.
Path: AddressBalanceTable:address:coin_name|lp_name
Value: balance string
*/
const AddressBalanceTable = "addrbalance"
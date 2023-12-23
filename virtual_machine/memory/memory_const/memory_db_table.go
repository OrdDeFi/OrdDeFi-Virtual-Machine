package memory_const

// CoinListTable stores all coin names
const CoinListTable = "coinlist"

// CoinMetadataTable stores deployed coins' meta
const CoinMetadataTable = "coinmeta"

/*
TotalMintedBalanceTable stores the value a coin has been already minted
Path: TotalMintedBalanceTable:coin_name
Value: balance string
*/
const TotalMintedBalanceTable = "totalminted"

/*
AddressMintedBalanceTable stores the value a coin has been already minted
Path: AddressMintedBalanceTable:coin_name:address
Value: balance string
*/
const AddressMintedBalanceTable = "addrminted"

/*
CoinBalanceTable stores a single coin's {address:balance}.
Path: CoinBalanceTable:version:coin_name|lp_name:address
Value: balance string
*/
const CoinBalanceTable = "coinbalance"

/*
UTXOCarryingBalanceTable stores token balance in UTXO.
The balance cast from available status to transferable by OpTransfer.
*/
const UTXOCarryingBalanceTable = "utxobalance"

/*
AddressBalanceTable stores a single address's {coin|lp:balance}.
Path: AddressBalanceTable:version:address:coin_name|lp_name
Value: balance string
*/
const AddressBalanceTable = "addrbalance"

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

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
CoinAddressBalanceTable stores a single coin's {address:balance}.
Path: CoinAddressBalanceTable:version:coin_name|lp_name:address
Value: balance string
*/
const CoinAddressBalanceTable = "coinaddrbalance"

/*
AddressCoinBalanceTable stores a single address's {coin:balance}.
Path: AddressBalanceTable:version:address:coin_name|lp_name
Value: balance string
*/
const AddressCoinBalanceTable = "addrcoinbalance"

/*
UTXOCarryingBalanceTable stores token balance in UTXO.
The balance cast from available status to transferable by OpTransfer.
*/
const UTXOCarryingBalanceTable = "utxobalance"

/*
UTXOCarryingListTable stores token balance in UTXO.
The balance cast from available status to transferable by OpTransfer.
Double write for UTXOCarryingBalanceTable, to query by tick and address.
*/
const UTXOCarryingListTable = "utxocarryinglist"

/*
UTXOTransferHistoryTable stores UTXO transfer history.
Format UTXOTransferHistoryTable:tick:sender_address:carrying_assets_UTXO
*/
const UTXOTransferHistoryTable = "utxotransferhistory"

// LpListTable stores all lp names
const LpListTable = "lplist"

// LpMetadataTable stores deployed lp' meta, coin value, price
const LpMetadataTable = "lpmeta"

/*
LPAddressBalanceTable stores a single lp's {address:balance}.
Path: LpBalanceTable:lp_name:address
Value: balance string
*/
const LPAddressBalanceTable = "lpaddrbalance"

/*
AddressLPBalanceTable stores a single address's {lp:balance}.
Path: LpBalanceTable:lp_name:address
Value: balance string
*/
const AddressLPBalanceTable = "addrlpbalance"

const LogMainTable = "log"
const LogQueryTxTable = "txlog"

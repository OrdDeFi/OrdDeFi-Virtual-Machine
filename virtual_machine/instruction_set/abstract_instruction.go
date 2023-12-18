package instruction_set

type AbstractInstruction struct {
	// TxId
	TxId string

	// Address
	TxInAddr  string
	TxOutAddr string

	// General key
	P  string `json:"p"`  // Protocol
	Op string `json:"op"` // Operator

	// Keys for deploy / mint / transfer
	Tick string `json:"tick"` // Coin name

	// Keys for mint / transfer / remove liquidity / swap
	Amt string `json:"amt"` // @required. Amount for mint / transfer / remove liquidity(lp amount) / swap(costing coin amount)

	// Keys for deploy
	Max     string `json:"max"`  // @required. Max amount in circulation
	Lim     string `json:"lim"`  // @required. Max amount to be minted in a single tx
	AddrLim string `json:"alim"` // @optional, default: infinite. Max amount to be minted in a single address
	Desc    string `json:"desc"` // @optional. Description for coin
	Icon    string `json:"icon"` // @optional. Icon for coin, in Base64 encoding

	// Keys for transfer
	To string `json:"to"` // @optional. When to address passed, only self to self tx allowed to execute OpTransfer.

	// Keys for add liquidity / remove liquidity / swap
	Ltick string `json:"lt"` // @required. Left coin at pair
	Rtick string `json:"rt"` // @required. Right coin at pair

	// Keys for add liquidity
	Lamt string `json:"lamt"` // @required. Left coin amount to adding into liquidity provider
	Ramt string `json:"ramt"` // @required. Right coin amount to adding into liquidity provider

	// Keys for add liquidity
	AllowSwap string `json:"as"` // @optional, default: 1. allow swap 1(true) / 0(false)

	// Keys for swap
	Spend     string `json:"spend"` // @required. Spend which coin at swapping
	Threshold string `json:"trhd"`  // @optional, default: 0.5%. Allowed threshold at swapping. If slippage > threshold, swap will be aborted. Both 0.005 or 0.5% format accepted

	// Keys for version control
	Ver     string `json:"ver"`     // @optional. Balance add to which version of VM. Default is v1 for !!ALL VERSIONS!! of VM.
	FromVer string `json:"fromver"` // @required. Which version to change coins from
	ToVer   string `json:"tover"`   // @required. Which version to change coins to
}

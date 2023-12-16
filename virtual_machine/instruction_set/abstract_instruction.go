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

	// Key for deploy / mint / transfer
	Tick string `json:"tick"` // Coin name

	// Key for mint / transfer / remove liquidity / swap
	Amt string `json:"amt"` // @required. Amount for mint / transfer / remove liquidity(lp amount) / swap(costing coin amount)

	// Key for deploy
	Max     string `json:"max"`  // @required. Max amount in circulation
	Lim     string `json:"lim"`  // @required. Max amount to be minted in a single tx
	AddrLim string `json:"alim"` // @optional, default: infinite. Max amount to be minted in a single address
	Icon    string `json:"icon"` // @optional. Icon for coin, in Base64 encoding

	// Key for transfer
	To string `json:"to"` // @optional. When to address passed, only self to self tx allowed to execute OpTransfer.

	// Key for add liquidity / remove liquidity / swap
	Ltick string `json:"lt"` // @required. Left coin at pair
	Rtick string `json:"rt"` // @required. Right coin at pair

	// Key for add liquidity
	Lamt string `json:"lamt"` // @required. Left coin amount to adding into liquidity provider
	Ramt string `json:"ramt"` // @required. Right coin amount to adding into liquidity provider

	// Key for add liquidity
	AllowSwap string `json:"as"` // @optional, default: 1. allow swap 1(true) / 0(false)

	// Key for swap
	Spend     string `json:"spend"` // @required. Spend which coin at swapping
	Threshold string `json:"trhd"`  // @optional, default: 0.5%. Allowed threshold at swapping. If slippage > threshold, swap will be aborted. Both 0.005 or 0.5% format accepted
}

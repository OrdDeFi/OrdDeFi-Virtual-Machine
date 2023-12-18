package instruction_set

const OpNameDeploy = "deploy"
const OpNameMint = "mint"
const OpNameTransfer = "transfer"
const OpNameAddLiquidityProvider = "addlp"
const OpNameRemoveLiquidityProvider = "rmlp"
const OpNameSwap = "swap"
const OpNameChangeVersion = "chver"

// OpDeployInstruction {"p":"orddefi","op":"deploy","tick":"odfi","max":"21000000","lim":"1000","alim":"1000","icon":""}
type OpDeployInstruction struct {
	Tick    string // @required. Coin name to deploy
	Max     string // @required. Max amount in circulation
	Lim     string // @required. Max amount to be minted in a single tx
	AddrLim string // @optional, default: infinite. Max amount to be minted in a single address
	Desc    string // @optional. Description for coin
	Icon    string // @optional. Icon for coin, in Base64 encoding
}

func compileOpDeployInstruction(instruction AbstractInstruction) *OpDeployInstruction {
	op := OpDeployInstruction{}
	op.Tick = instruction.Tick
	op.Max = instruction.Max
	op.Lim = instruction.Lim
	op.AddrLim = instruction.AddrLim
	op.Desc = instruction.Desc
	op.Icon = instruction.Icon
	return &op
}

// OpMintInstruction {"p":"orddefi","op":"mint","tick":"odfi","amt":"1000"}
type OpMintInstruction struct {
	TxOutAddr string
	Tick      string // @required. Coin name to mint
	Amt       string // @required. Amount to mint
	Ver       string // @optional. Balance add to which version of VM. Default is v1 for !!ALL VERSIONS!! of VM.
}

func compileOpMintInstruction(instruction AbstractInstruction) *OpMintInstruction {
	op := OpMintInstruction{}
	op.TxOutAddr = instruction.TxOutAddr
	op.Tick = instruction.Tick
	op.Amt = instruction.Amt
	op.Ver = instruction.Ver
	return &op
}

/*
OpTransferInstruction
1. {"p":"orddefi","op":"transfer","tick":"odfi","amt":"1000"} from any address, send UTXO to apply transfer
2. {"p":"orddefi","op":"transfer","tick":"odfi","amt":"1000", "to": "bc1p****"} from self address, apply immediately
*/
type OpTransferInstruction struct {
	TxInAddr  string
	TxOutAddr string
	TxId      string // When to is nil, TxOut:0 will become an UTXO holding transferable coins, coins moved at UTXO becomes TxIn.
	Tick      string // @required. Coin name to transfer
	Amt       string // @required. Amount to transfer
	To        string // @optional. When to address passed, only self to self tx allowed to execute OpTransfer.
}

func compileOpTransferInstruction(instruction AbstractInstruction) *OpTransferInstruction {
	op := OpTransferInstruction{}
	op.TxInAddr = instruction.TxInAddr
	op.TxOutAddr = instruction.TxOutAddr
	op.TxId = instruction.TxId
	op.Tick = instruction.Tick
	op.Amt = instruction.Amt
	op.To = instruction.To
	return &op
}

/*
OpAddLiquidityProviderInstruction Add liquidity provider. Only self to self tx allow to execute OpAddLiquidityProvider
1. {"p":"orddefi","op":"addlp","ltick":"odfi","rtick":"odgv","lamt":"1000","ramt":"1000"} Add LP for odfi-odgv 50:50, allow swap.
2. {"p":"orddefi","op":"addlp","ltick":"odfi","rtick":"odgv","lamt":"1000","ramt":"1000","as":0} Add LP for odfi-odgv 50:50, disable swap.
*/
type OpAddLiquidityProviderInstruction struct {
	TxInAddr  string
	TxOutAddr string
	Ltick     string // @required. Left coin at pair
	Rtick     string // @required. Right coin at pair
	Lamt      string // @required. Left coin amount to adding into liquidity provider
	Ramt      string // @required. Right coin amount to adding into liquidity provider
	AllowSwap string // @optional, default: 1. allow swap 1(true) / 0(false)
}

func compileOpAddLiquidityProviderInstruction(instruction AbstractInstruction) *OpAddLiquidityProviderInstruction {
	op := OpAddLiquidityProviderInstruction{}
	op.TxInAddr = instruction.TxInAddr
	op.TxOutAddr = instruction.TxOutAddr
	op.Ltick = instruction.Ltick
	op.Rtick = instruction.Rtick
	op.Lamt = instruction.Lamt
	op.Ramt = instruction.Ramt
	op.AllowSwap = instruction.AllowSwap
	return &op
}

/*
OpRemoveLiquidityProviderInstruction Remove liquidity provider. Only self to self tx allow to execute OpRemoveLiquidityProvider
{"p":"orddefi","op":"rmlp","ltick":"odfi","rtick":"odgv","amt":"10"} Remove odfi-odgv LP by 10, and get odfi and odgv coins.
*/
type OpRemoveLiquidityProviderInstruction struct {
	TxInAddr  string
	TxOutAddr string
	Ltick     string // @required. Left coin at pair
	Rtick     string // @required. Right coin at pair
	Amt       string // @required. Amount to remove LP
}

func compileOpRemoveLiquidityProviderInstruction(instruction AbstractInstruction) *OpRemoveLiquidityProviderInstruction {
	op := OpRemoveLiquidityProviderInstruction{}
	op.TxInAddr = instruction.TxInAddr
	op.TxOutAddr = instruction.TxOutAddr
	op.Ltick = instruction.Ltick
	op.Rtick = instruction.Rtick
	op.Amt = instruction.Amt
	return &op
}

/*
OpSwapInstruction Swap at liquidity provider. Only self to self tx allow to execute OpSwap
1. {"p":"orddefi","op":"swap","ltick":"odfi","rtick":"odgv","spend":"odgv","amt":"10"} Buying odfi with odgv by 10 odgv.
2. {"p":"orddefi","op":"swap","ltick":"odfi","rtick":"odgv","spend":"odgv","amt":"10","trhd":"0.8%"} Buying odfi with odgv by 10 odgv.
*/
type OpSwapInstruction struct {
	TxInAddr  string
	TxOutAddr string
	Ltick     string // @required. Left coin at pair
	Rtick     string // @required. Right coin at pair
	Spend     string // @required. Spend which coin at swapping. e.g. "odgv"
	Amt       string // @required. Amount to spend for the spending coin
	Threshold string // @optional, default: 0.005. Allowed threshold at swapping. If slippage > threshold, swap will be aborted
}

func compileOpSwapInstruction(instruction AbstractInstruction) *OpSwapInstruction {
	op := OpSwapInstruction{}
	op.TxInAddr = instruction.TxInAddr
	op.TxOutAddr = instruction.TxOutAddr
	op.Ltick = instruction.Ltick
	op.Rtick = instruction.Rtick
	op.Spend = instruction.Spend
	op.Amt = instruction.Amt
	op.Threshold = instruction.Threshold
	return &op
}

/*
OpChangeVersionInstruction Move coins from a version to another.
{"p":"orddefi","op":"chver","tick":"odfi","fromver":"1","tover":"2","amt":"10"} Move 10 odfi from v1 to v2.
*/
type OpChangeVersionInstruction struct {
	TxInAddr  string
	TxOutAddr string
	Tick      string // @required. Coin name to change version
	FromVer   string // @required. From version
	ToVer     string // @required. To version
	Amt       string // @required. Amount to change version
}

func compileOpChangeVersionInstruction(instruction AbstractInstruction) *OpChangeVersionInstruction {
	op := OpChangeVersionInstruction{}
	op.TxInAddr = instruction.TxInAddr
	op.TxOutAddr = instruction.TxOutAddr
	op.Tick = instruction.Tick
	op.FromVer = instruction.FromVer
	op.ToVer = instruction.ToVer
	op.Amt = instruction.Amt
	return &op
}

func CompileInstruction(abstractInstruction AbstractInstruction) *interface{} {
	op := abstractInstruction.Op
	var res interface{}
	switch op {
	case OpNameDeploy:
		opDeploy := compileOpDeployInstruction(abstractInstruction)
		if opDeploy != nil {
			res = *opDeploy
		}
	case OpNameMint:
		opMint := compileOpMintInstruction(abstractInstruction)
		if opMint != nil {
			res = *opMint
		}
	case OpNameTransfer:
		opTransfer := compileOpTransferInstruction(abstractInstruction)
		if opTransfer != nil {
			res = *opTransfer
		}
	case OpNameAddLiquidityProvider:
		opAddLiquidityProvider := compileOpAddLiquidityProviderInstruction(abstractInstruction)
		if opAddLiquidityProvider != nil {
			res = *opAddLiquidityProvider
		}
	case OpNameRemoveLiquidityProvider:
		opRemoveLiquidityProvider := compileOpRemoveLiquidityProviderInstruction(abstractInstruction)
		if opRemoveLiquidityProvider != nil {
			res = *opRemoveLiquidityProvider
		}
	case OpNameSwap:
		opSwap := compileOpSwapInstruction(abstractInstruction)
		if opSwap != nil {
			res = *opSwap
		}
	case OpNameChangeVersion:
		opChangeVersion := compileOpChangeVersionInstruction(abstractInstruction)
		if opChangeVersion != nil {
			res = *opChangeVersion
		}
	}
	if res != nil {
		return &res
	}
	return nil
}

package instruction_set

import (
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/tx_utils"
	"strings"
)

const OpNameDeploy = "deploy"
const OpNameMint = "mint"
const OpNameTransfer = "transfer"
const OpNameAddLiquidityProvider = "addlp"
const OpNameRemoveLiquidityProvider = "rmlp"
const OpNameSwap = "swap"

// OpDeployInstruction {"p":"orddefi","op":"deploy","tick":"odfi","max":"21000000","lim":"1000","alim":"1000","icon":""}
type OpDeployInstruction struct {
	Tick           string // @required. Coin name to deploy
	Max            string // @required. Max amount in circulation
	Lim            string // @required. Max amount to be minted in a single tx
	AddrLim        string // @optional, default: infinite. Max amount to be minted in a single address
	Desc           string // @optional. Description for coin
	Icon           string // @optional. Icon for coin, in Base64 encoding
	RawInstruction string // Raw JSON string of instruction
}

func compileOpDeployInstruction(instruction AbstractInstruction) *OpDeployInstruction {
	op := OpDeployInstruction{}
	op.Tick = instruction.Tick
	op.Max = instruction.Max
	op.Lim = instruction.Lim
	op.AddrLim = instruction.AddrLim
	op.Desc = instruction.Desc
	op.Icon = instruction.Icon
	rawJsonString := instruction.JsonString()
	if rawJsonString != nil {
		op.RawInstruction = *rawJsonString
		return &op
	}
	return nil
}

// OpMintInstruction {"p":"orddefi","op":"mint","tick":"odfi","amt":"1000"}
type OpMintInstruction struct {
	TxOutAddr           string
	PreviousOutputIndex int
	Tick                string // @required. Coin name to mint
	Amt                 string // @required. Amount to mint
	RawInstruction      string // Raw JSON string of instruction
}

func (opMint OpMintInstruction) IsValidOpMintInstruction() bool {
	valid := opMint.PreviousOutputIndex == 0
	return valid
}

func compileOpMintInstruction(instruction AbstractInstruction) *OpMintInstruction {
	op := OpMintInstruction{}
	op.TxOutAddr = instruction.TxOutAddr
	op.PreviousOutputIndex = instruction.PreviousOutputIndex
	op.Tick = instruction.Tick
	op.Amt = instruction.Amt
	rawJsonString := instruction.JsonString()
	if rawJsonString != nil {
		op.RawInstruction = *rawJsonString
		return &op
	}
	return nil
}

/*
OpTransferInstruction
1. {"p":"orddefi","op":"transfer","tick":"odfi","amt":"1000"} from any address, send UTXO to apply transfer
2. {"p":"orddefi","op":"transfer","tick":"odfi","amt":"1000", "to": "bc1p****"} from self address, apply immediately
*/
type OpTransferInstruction struct {
	TxInAddr       string
	TxOutAddr      string
	TxId           string // When to is nil, TxOut:0 will become an UTXO holding transferable coins, coins moved at UTXO becomes TxIn.
	Tick           string // @required. Coin name to transfer
	Amt            string // @required. Amount to transfer
	To             string // @optional. When to address passed, only self to self tx allowed to execute OpTransfer.
	RawInstruction string // Raw JSON string of instruction
}

func compileOpTransferInstruction(instruction AbstractInstruction) *OpTransferInstruction {
	op := OpTransferInstruction{}
	op.TxOutAddr = instruction.TxOutAddr
	op.TxId = instruction.TxId
	op.Tick = instruction.Tick
	op.Amt = instruction.Amt
	op.To = instruction.To
	rawJsonString := instruction.JsonString()
	if rawJsonString != nil {
		op.RawInstruction = *rawJsonString
		return &op
	}
	return nil
}

/*
OpAddLiquidityProviderInstruction Add liquidity provider. Only self to self tx allow to execute OpAddLiquidityProvider
{"p":"orddefi","op":"addlp","ltick":"odfi","rtick":"odgv","lamt":"1000","ramt":"1000"} Add LP for odfi-odgv 50:50.
*/
type OpAddLiquidityProviderInstruction struct {
	TxInAddr       string
	TxOutAddr      string
	Ltick          string // @required. Left coin at pair
	Rtick          string // @required. Right coin at pair
	Lamt           string // @required. Left coin amount to adding into liquidity provider
	Ramt           string // @required. Right coin amount to adding into liquidity provider
	RawInstruction string // Raw JSON string of instruction
}

/*
ExtractParams alphabetical compare ltick and rtick, make smaller be the actual left
return actualLtick, actualRtick, actualLamt, actualRamt
*/
func (op OpAddLiquidityProviderInstruction) ExtractParams() (*string, *string, *safe_number.SafeNum, *safe_number.SafeNum) {
	if op.Ltick == "" || op.Rtick == "" {
		return nil, nil, nil, nil
	}
	cmpRes := strings.Compare(op.Ltick, op.Rtick)
	if cmpRes == 0 {
		return nil, nil, nil, nil
	}
	ltick := op.Ltick
	rtick := op.Rtick
	lamt := safe_number.SafeNumFromString(op.Lamt)
	ramt := safe_number.SafeNumFromString(op.Ramt)
	if lamt == nil || ramt == nil {
		return nil, nil, nil, nil
	}
	if cmpRes < 0 {
		return &ltick, &rtick, lamt, ramt
	} else if cmpRes > 0 {
		return &rtick, &ltick, ramt, lamt
	}
	return nil, nil, nil, nil
}

func compileOpAddLiquidityProviderInstruction(instruction AbstractInstruction) *OpAddLiquidityProviderInstruction {
	op := OpAddLiquidityProviderInstruction{}
	op.TxOutAddr = instruction.TxOutAddr
	op.Ltick = instruction.Ltick
	op.Rtick = instruction.Rtick
	op.Lamt = instruction.Lamt
	op.Ramt = instruction.Ramt
	rawJsonString := instruction.JsonString()
	if rawJsonString != nil {
		op.RawInstruction = *rawJsonString
		return &op
	}
	return nil
}

/*
OpRemoveLiquidityProviderInstruction Remove liquidity provider. Only self to self tx allow to execute OpRemoveLiquidityProvider
{"p":"orddefi","op":"rmlp","ltick":"odfi","rtick":"odgv","amt":"10"} Remove odfi-odgv LP by 10, and get odfi and odgv coins.
*/
type OpRemoveLiquidityProviderInstruction struct {
	TxInAddr       string
	TxOutAddr      string
	Ltick          string // @required. Left coin at pair
	Rtick          string // @required. Right coin at pair
	Amt            string // @required. Amount to remove LP
	RawInstruction string // Raw JSON string of instruction
}

/*
ExtractParams alphabetical compare ltick and rtick, make smaller be the actual left
return actualLtick, actualRtick, amt
*/
func (op OpRemoveLiquidityProviderInstruction) ExtractParams() (*string, *string, *safe_number.SafeNum) {
	if op.Ltick == "" || op.Rtick == "" {
		return nil, nil, nil
	}
	cmpRes := strings.Compare(op.Ltick, op.Rtick)
	if cmpRes == 0 {
		return nil, nil, nil
	}
	ltick := op.Ltick
	rtick := op.Rtick
	amt := safe_number.SafeNumFromString(op.Amt)
	if amt == nil {
		return nil, nil, nil
	}
	if cmpRes < 0 {
		return &ltick, &rtick, amt
	} else if cmpRes > 0 {
		return &rtick, &ltick, amt
	}
	return nil, nil, nil
}

func compileOpRemoveLiquidityProviderInstruction(instruction AbstractInstruction) *OpRemoveLiquidityProviderInstruction {
	op := OpRemoveLiquidityProviderInstruction{}
	op.TxOutAddr = instruction.TxOutAddr
	op.Ltick = instruction.Ltick
	op.Rtick = instruction.Rtick
	op.Amt = instruction.Amt
	rawJsonString := instruction.JsonString()
	if rawJsonString != nil {
		op.RawInstruction = *rawJsonString
		return &op
	}
	return nil
}

/*
OpSwapInstruction Swap at liquidity provider. Only self to self tx allow to execute OpSwap
1. {"p":"orddefi","op":"swap","ltick":"odfi","rtick":"odgv","spend":"odgv","amt":"10"} Buying odfi with odgv by 10 odgv.
2. {"p":"orddefi","op":"swap","ltick":"odfi","rtick":"odgv","spend":"odgv","amt":"10","trhd":"0.8%"} Buying odfi with odgv by 10 odgv.
*/
type OpSwapInstruction struct {
	TxInAddr       string
	TxOutAddr      string
	Ltick          string // @required. Left coin at pair
	Rtick          string // @required. Right coin at pair
	Spend          string // @required. Spend which coin at swapping. e.g. "odgv"
	Amt            string // @required. Amount to spend for the spending coin
	Threshold      string // @optional, default: 0.005. Allowed threshold at swapping. If slippage > threshold, swap will be aborted
	RawInstruction string // Raw JSON string of instruction
}

/*
ExtractParams alphabetical compare ltick and rtick, make smaller be the actual left
return actualLtick, actualRtick, amt
*/
func (op OpSwapInstruction) ExtractParams() (*string, *string, *safe_number.SafeNum) {
	if op.Ltick == "" || op.Rtick == "" {
		return nil, nil, nil
	}
	if op.Ltick != op.Spend && op.Rtick != op.Spend {
		return nil, nil, nil
	}
	cmpRes := strings.Compare(op.Ltick, op.Rtick)
	if cmpRes == 0 {
		return nil, nil, nil
	}
	ltick := op.Ltick
	rtick := op.Rtick
	amt := safe_number.SafeNumFromString(op.Amt)
	if amt == nil {
		return nil, nil, nil
	}
	if cmpRes < 0 {
		return &ltick, &rtick, amt
	} else if cmpRes > 0 {
		return &rtick, &ltick, amt
	}
	return nil, nil, nil
}

func compileOpSwapInstruction(instruction AbstractInstruction) *OpSwapInstruction {
	op := OpSwapInstruction{}
	op.TxOutAddr = instruction.TxOutAddr
	op.Ltick = instruction.Ltick
	op.Rtick = instruction.Rtick
	op.Spend = instruction.Spend
	op.Amt = instruction.Amt
	op.Threshold = instruction.Threshold
	rawJsonString := instruction.JsonString()
	if rawJsonString != nil {
		op.RawInstruction = *rawJsonString
		return &op
	}
	return nil
}

func CheckTickLegal(tick string) bool {
	if strings.Contains(tick, "-") {
		return false
	} else if strings.Contains(tick, "_") {
		return false
	} else if strings.Contains(tick, ":") {
		return false
	} else if strings.HasPrefix(tick, "$") {
		return false
	} else if strings.HasPrefix(tick, "@") {
		return false
	} else if strings.HasPrefix(tick, "#") {
		return false
	} else if strings.HasPrefix(tick, "%") {
		return false
	}
	length := len(tick)
	return length == 4
}

func CompileInstruction(abstractInstruction AbstractInstruction) *interface{} {
	if abstractInstruction.Tick != "" && CheckTickLegal(abstractInstruction.Tick) == false {
		return nil
	}
	if abstractInstruction.Ltick != "" && CheckTickLegal(abstractInstruction.Ltick) == false {
		return nil
	}
	if abstractInstruction.Rtick != "" && CheckTickLegal(abstractInstruction.Rtick) == false {
		return nil
	}
	if abstractInstruction.To != "" && tx_utils.IsValidBitcoinAddress(abstractInstruction.To) == false {
		return nil
	}
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
	}
	if res != nil {
		return &res
	}
	return nil
}

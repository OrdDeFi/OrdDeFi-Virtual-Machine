package virtual_machine

import (
	"OrdDefi-Virtual-Machine/tx_utils"
	"encoding/json"
	"errors"
	"github.com/btcsuite/btcd/wire"
	"strings"
)

const POrdDefi = "orddefi"

func isValidContentType(contentType string) bool {
	trimmedContentType := strings.TrimSpace(contentType)
	parts := strings.Split(trimmedContentType, ";")
	primaryType := parts[0]
	primaryType = strings.ToLower(primaryType)
	if primaryType == "text/plain" {
		return true
	}
	return false
}

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

func onlySelfTxAllowed(instruction AbstractInstruction) bool {
	op := instruction.Op
	if op == OpNameAddLiquidityProvider || op == OpNameRemoveLiquidityProvider || op == OpNameSwap {
		return true
	}
	if op == OpNameTransfer && instruction.To != "" {
		return true
	}
	return false
}

/*
preCompileInstructions
1. Check content-type, only "text/plain" available as instructions;
2. Parse content JSON string into []AbstractInstruction.
*/
func preCompileInstructions(contentType string, content []byte) []AbstractInstruction {
	if isValidContentType(contentType) == false {
		return nil
	}
	var abstractInstruction AbstractInstruction
	err := json.Unmarshal(content, &abstractInstruction)
	if err == nil {
		res := []AbstractInstruction{abstractInstruction}
		return res
	}
	var instructions []AbstractInstruction
	err2 := json.Unmarshal(content, &instructions)
	if err2 == nil {
		return instructions
	}
	return nil
}

func filterAbstractInstructions(rawInstructions []AbstractInstruction, tx *wire.MsgTx, txId string) ([]interface{}, error) {
	var res []interface{}
	for _, abstractInstruction := range rawInstructions {
		abstractInstruction.P = strings.ToLower(abstractInstruction.P)
		abstractInstruction.Op = strings.ToLower(abstractInstruction.Op)
		abstractInstruction.Tick = strings.ToLower(abstractInstruction.Tick)
		abstractInstruction.Ltick = strings.ToLower(abstractInstruction.Ltick)
		abstractInstruction.Rtick = strings.ToLower(abstractInstruction.Rtick)
		abstractInstruction.Spend = strings.ToLower(abstractInstruction.Spend)
		if abstractInstruction.P == POrdDefi {
			// parse output address
			firstOutputAddress, err := tx_utils.ParseFirstOutputAddress(tx)
			if err != nil {
				return nil, err
			}
			if firstOutputAddress == nil {
				return nil, errors.New("filterAbstractInstructions ParseFirstOutputAddress got empty address")
			}
			abstractInstruction.TxOutAddr = *firstOutputAddress
			// parse input address if needed
			if onlySelfTxAllowed(abstractInstruction) {
				firstInputAddress, err := tx_utils.ParseFirstInputAddress(tx)
				if err != nil {
					return nil, err
				}
				if firstInputAddress == nil {
					return nil, errors.New("filterAbstractInstructions ParseFirstInputAddress got empty address")
				}
				abstractInstruction.TxInAddr = *firstInputAddress
			}
			// save txid to abstract instruction
			abstractInstruction.TxId = txId
			instruction := CompileInstruction(abstractInstruction)
			if instruction != nil {
				res = append(res, *instruction)
			}
		}
	}
	return res, nil
}

func CompileInstructions(contentType string, content []byte, tx *wire.MsgTx, txId string) ([]interface{}, error) {
	instructions := preCompileInstructions(contentType, content)
	filteredInstructions, err := filterAbstractInstructions(instructions, tx, txId)
	return filteredInstructions, err
}

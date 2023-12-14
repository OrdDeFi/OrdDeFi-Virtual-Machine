package virtual_machine

import (
	"encoding/json"
	"github.com/btcsuite/btcd/wire"
	"strings"
)

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

type OrdefiInstruction struct {
	// General key
	P  string `json:"p"`  // Protocol
	Op string `json:"op"` // Operator

	// Key for deploy / mint / transfer
	Tick string `json:"tick"` // Coin name

	// Key for deploy
	Max     interface{} `json:"max"`  // @required. Max amount in circulation
	Lim     interface{} `json:"lim"`  // @required. Max amount to be minted in a single tx
	AddrLim interface{} `json:"alim"` // @optional, default: infinite. Max amount to be minted in a single address
	Icon    string      `json:"icon"` // @optional. Icon for coin, in Base64 encoding

	// Key for add liquidity
	Ltick     string      `json:"lt"` // @required. [coin_name]:[amount] e.g. brcd:1000
	Rtick     string      `json:"rt"` // @required. [coin_name]:[amount] e.g. brcg:1000
	AllowSwap interface{} `json:"as"` // @optional, default: 1. allow swap 1(true) / 0(false)

	// Key for remove liquidity and swap
	Pair string `json:"pair"` // @required. Removing pair of liquidity, [left_coin_name]-[right_coin_name] e.g. brcd-brcg

	// Key for swap
	Spend     interface{} `json:"spend"` // @required. Spend which coin at swapping
	Threshold interface{} `json:"thld"`  // @optional, default: 0.5%. Allowed threshold at swapping. If slippage > threshold, swap will be aborted. Both 0.005 or 0.5% format accepted.

	// Key for mint / transfer / remove liquidity / swap
	Amt interface{} `json:"amt"` // @required. Amount for mint / transfer / remove liquidity(lp amount) / swap(costing coin amount)
}

func preCompileInstructions(contentType string, content []byte, tx *wire.MsgTx, txId string) []OrdefiInstruction {
	if isValidContentType(contentType) == false {
		return nil
	}
	var instruction OrdefiInstruction
	err := json.Unmarshal(content, &instruction)
	if err == nil {
		res := []OrdefiInstruction{instruction}
		return res
	}
	var instructions []OrdefiInstruction
	err2 := json.Unmarshal(content, &instructions)
	if err2 == nil {
		return instructions
	}
	return nil
}

func CompileInstructions(contentType string, content []byte, tx *wire.MsgTx, txId string) []OrdefiInstruction {
	res := preCompileInstructions(contentType, content, tx, txId)
	return res
}

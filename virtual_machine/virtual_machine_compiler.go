package virtual_machine

import (
	"OrdDeFi-Virtual-Machine/tx_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"encoding/json"
	"errors"
	"github.com/btcsuite/btcd/wire"
	"strings"
)

const POrdDeFi = "orddefi"

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

func onlySelfTxAllowed(instruction instruction_set.AbstractInstruction) bool {
	op := instruction.Op
	if op == instruction_set.OpNameAddLiquidityProvider ||
		op == instruction_set.OpNameRemoveLiquidityProvider ||
		op == instruction_set.OpNameSwap ||
		op == instruction_set.OpNameChangeVersion {
		return true
	}
	if op == instruction_set.OpNameTransfer && instruction.To != "" {
		return true
	}
	return false
}

/*
preCompileInstructions
1. Check content-type, only "text/plain" available as instructions;
2. Parse content JSON string into []AbstractInstruction;
3. All content will be parsed as UTF-8.
*/
func preCompileInstructions(contentType string, content []byte) []instruction_set.AbstractInstruction {
	if isValidContentType(contentType) == false {
		return nil
	}
	var abstractInstruction instruction_set.AbstractInstruction
	err := json.Unmarshal(content, &abstractInstruction)
	if err == nil {
		res := []instruction_set.AbstractInstruction{abstractInstruction}
		return res
	}
	var instructions []instruction_set.AbstractInstruction
	err2 := json.Unmarshal(content, &instructions)
	if err2 == nil {
		for _, eachInst := range instructions {
			lowerOp := strings.ToLower(eachInst.Op)
			if (lowerOp == "deploy" || lowerOp == "mint") && len(instructions) > 1 {
				// Bulk execute instructs doesn't allow "deploy" and "mint", otherwise all instructions in slice will be aborted.
				return nil
			}
		}
		return instructions
	}
	return nil
}

func filterAbstractInstructions(rawInstructions []instruction_set.AbstractInstruction, tx *wire.MsgTx, txId string) ([]interface{}, error) {
	var res []interface{}
	for _, abstractInstruction := range rawInstructions {
		abstractInstruction.P = strings.ToLower(abstractInstruction.P)
		abstractInstruction.Op = strings.ToLower(abstractInstruction.Op)
		abstractInstruction.Tick = strings.ToLower(abstractInstruction.Tick)
		abstractInstruction.Ltick = strings.ToLower(abstractInstruction.Ltick)
		abstractInstruction.Rtick = strings.ToLower(abstractInstruction.Rtick)
		abstractInstruction.Spend = strings.ToLower(abstractInstruction.Spend)
		if abstractInstruction.P == POrdDeFi {
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
			instruction := instruction_set.CompileInstruction(abstractInstruction)
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

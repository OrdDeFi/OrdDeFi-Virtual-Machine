package authentication

import (
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"github.com/btcsuite/btcd/wire"
)

func InstructionShouldBeAuthed(instruction instruction_set.AbstractInstruction) bool {
	op := instruction.Op
	if op == instruction_set.OpNameAddLiquidityProvider ||
		op == instruction_set.OpNameRemoveLiquidityProvider ||
		op == instruction_set.OpNameSwap {
		return true
	}
	if op == instruction_set.OpNameTransfer && instruction.To != "" {
		return true
	}
	return false
}

func InstructionAuthenticate(tx *wire.MsgTx) (*bool, error) {
	return nil, nil
}

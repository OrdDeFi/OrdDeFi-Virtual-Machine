package operations

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"errors"
)

func getDiscount(instruction instruction_set.OpSwapInstruction, db *db_utils.OrdDB) (*float64, error) {
	return nil, nil
}

func calculateDeltaY(deltaX *safe_number.SafeNum, X *safe_number.SafeNum, Y *safe_number.SafeNum) *safe_number.SafeNum {
	// alpha = deltaX / X
	// beta = 1 - 1 / (1 + alpha)
	// deltaY = beta * Y
	return nil
}

func ExecuteOpSwap(instruction instruction_set.OpSwapInstruction, db *db_utils.OrdDB) error {
	if instruction.TxInAddr != instruction.TxOutAddr {
		return errors.New("no privileges on cross-address swap")
	}
	return nil
}

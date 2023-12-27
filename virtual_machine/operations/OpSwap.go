package operations

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"errors"
)

func DiscountForODFIAmount(totalValue *safe_number.SafeNum) (*float64, error) {
	if totalValue == nil {
		return nil, errors.New("getDiscount error: calc odfi total balance failed")
	}
	if totalValue.Compare(safe_number.SafeNumFromString("21000")) >= 0 {
		discount := 0.3
		return &discount, nil
	} else if totalValue.Compare(safe_number.SafeNumFromString("2100")) >= 0 {
		discount := 0.6
		return &discount, nil
	} else if totalValue.Compare(safe_number.SafeNumFromString("210")) >= 0 {
		discount := 0.9
		return &discount, nil
	} else {
		discount := 1.0
		return &discount, nil
	}
}

func getDiscount(instruction instruction_set.OpSwapInstruction, db *db_utils.OrdDB) (*float64, error) {
	available, transferable, err := memory_read.Balance(db, "odfi", instruction.TxOutAddr)
	if err != nil {
		return nil, err
	}
	if available == nil || transferable == nil {
		return nil, errors.New("getDiscount error: read odfi balance failed")
	}
	totalValue := available.Add(transferable)
	return DiscountForODFIAmount(totalValue)
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

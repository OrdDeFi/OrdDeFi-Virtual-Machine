package operations

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"errors"
	"fmt"
)

func DiscountForODFIAmount(totalValue *safe_number.SafeNum) (*string, error) {
	if totalValue == nil {
		return nil, errors.New("getDiscount error: calc odfi total balance failed")
	}
	if totalValue.Compare(safe_number.SafeNumFromString("21000")) >= 0 {
		discount := "0.3"
		return &discount, nil
	} else if totalValue.Compare(safe_number.SafeNumFromString("2100")) >= 0 {
		discount := "0.6"
		return &discount, nil
	} else if totalValue.Compare(safe_number.SafeNumFromString("210")) >= 0 {
		discount := "0.9"
		return &discount, nil
	} else {
		discount := "1"
		return &discount, nil
	}
}

func getDiscount(instruction instruction_set.OpSwapInstruction, db *db_utils.OrdDB) (*string, error) {
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

/*
getLPTakerFee
Trading **** at any pair, this pair will take 0.18% of consumed ****.
Final charged fee is affected by discount value.
*/
func getLPTakerFee(instruction instruction_set.OpSwapInstruction, db *db_utils.OrdDB) (*safe_number.SafeNum, error) {
	discount, err := getDiscount(instruction, db)
	if err != nil {
		return nil, err
	}
	if discount == nil {
		return nil, errors.New("ExecuteOpSwap error: get discount failed")
	}
	return nil, nil
}

/*
getODFITakerFee
Trading **** at any pair, ODFI-**** will take 0.02% of consumed ****.
If **** is ODFI, ODFITakerFee will be free.
Final charged fee is affected by discount value.
*/
func getODFITakerFee(instruction instruction_set.OpSwapInstruction, db *db_utils.OrdDB) (*safe_number.SafeNum, error) {
	discount, err := getDiscount(instruction, db)
	if err != nil {
		return nil, err
	}
	if discount == nil {
		return nil, errors.New("ExecuteOpSwap error: get discount failed")
	}
	return nil, nil
}

func calculateDeltaY(deltaX *safe_number.SafeNum, X *safe_number.SafeNum, Y *safe_number.SafeNum) *safe_number.SafeNum {
	// alpha = deltaX / X
	// beta = 1 - 1 / (1 + alpha)
	// deltaY = beta * Y
	return nil
}

func performSwap(instruction instruction_set.OpSwapInstruction, db *db_utils.OrdDB) error {
	//address := instruction.TxOutAddr
	//tick := instruction.Spend
	//lTick, rTick, consumingAmt := instruction.ExtractParams()

	return nil
}

func ExecuteOpSwap(instruction instruction_set.OpSwapInstruction, db *db_utils.OrdDB) error {
	if instruction.TxInAddr != instruction.TxOutAddr {
		return errors.New("no privileges on cross-address swap")
	}
	address := instruction.TxOutAddr
	tick := instruction.Spend
	lTick, rTick, consumingAmt := instruction.ExtractParams()
	if lTick == nil || rTick == nil || consumingAmt == nil {
		return errors.New("ExecuteOpSwap error: params extracting error")
	}
	if consumingAmt.IsZero() {
		return errors.New("ExecuteOpSwap error: amt is 0")
	}
	available, err := memory_read.AvailableBalance(db, tick, address)
	if err != nil {
		return err
	}
	if available.Compare(consumingAmt) < 0 {
		return fmt.Errorf("consumingAmt not enough: %s < %s", available.String(), consumingAmt.String())
	}
	return performSwap(instruction, db)
}

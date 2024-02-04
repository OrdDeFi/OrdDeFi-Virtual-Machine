package operations

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_write"
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
	if totalValue == nil {
		return nil, errors.New("getDiscount failed: totalValue is nil")
	}
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
		return nil, errors.New("getLPTakerFee error: get discount failed")
	}
	discountNum := safe_number.SafeNumFromString(*discount)
	if discountNum == nil {
		return nil, errors.New("getLPTakerFee error: convert discount number failed")
	}
	lTick, rTick, consumingAmt := instruction.ExtractParams()
	if lTick == nil || rTick == nil || consumingAmt == nil {
		return nil, errors.New("getLPTakerFee error: params extracting error")
	}
	standardFeeRate := safe_number.SafeNumFromString("0.0018")
	if standardFeeRate == nil {
		return nil, errors.New("getLPTakerFee error: standardFeeRate is nil")
	}
	standardFee := consumingAmt.Multiply(standardFeeRate)
	if standardFee == nil {
		return nil, errors.New("getLPTakerFee calculating standardFee failed")
	}
	actualFee := standardFee.Multiply(discountNum)
	if actualFee == nil {
		return nil, errors.New("getLPTakerFee calculating actualFee failed")
	}
	return actualFee, nil
}

/*
getODFITakerFee
Trading **** at any pair, ODFI-**** will take 0.02% of consumed ****.
If **** is ODFI, ODFITakerFee will be free.
Final charged fee is affected by discount value.
*/
func getODFITakerFee(instruction instruction_set.OpSwapInstruction, db *db_utils.OrdDB) (*safe_number.SafeNum, error) {
	tick := instruction.Spend
	if tick == "odfi" {
		return safe_number.SafeNumFromString("0"), nil
	}
	odfiMeta, err := memory_read.LiquidityProviderMetadata(db, "odfi", tick) // this method automatically decide with tick is left or right
	if err != nil {
		return nil, err
	}
	if odfiMeta == nil {
		return safe_number.SafeNumFromString("0"), nil
	}
	discount, err := getDiscount(instruction, db)
	if err != nil {
		return nil, err
	}
	if discount == nil {
		return nil, errors.New("getODFITakerFee error: get discount failed")
	}
	discountNum := safe_number.SafeNumFromString(*discount)
	if discountNum == nil {
		return nil, errors.New("getODFITakerFee error: convert discount number failed")
	}
	lTick, rTick, consumingAmt := instruction.ExtractParams()
	if lTick == nil || rTick == nil || consumingAmt == nil {
		return nil, errors.New("getODFITakerFee error: params extracting error")
	}
	standardFeeRate := safe_number.SafeNumFromString("0.0002")
	if standardFeeRate == nil {
		return nil, errors.New("getODFITakerFee error: standardFeeRate is nil")
	}
	standardFee := consumingAmt.Multiply(standardFeeRate)
	if standardFee == nil {
		return nil, errors.New("getODFITakerFee calculating standardFee failed")
	}
	actualFee := standardFee.Multiply(discountNum)
	if actualFee == nil {
		return nil, errors.New("getODFITakerFee calculating actualFee failed")
	}
	return actualFee, nil
}

func calculateDeltaY(deltaX *safe_number.SafeNum, X *safe_number.SafeNum, Y *safe_number.SafeNum) (*safe_number.SafeNum, error) {
	if deltaX == nil || X == nil || Y == nil {
		return nil, errors.New("calculateDeltaY error: deltaX, X or Y is nil")
	}
	// alpha = deltaX / X
	alpha := deltaX.DivideBy(X)
	if alpha == nil {
		return nil, errors.New("calculateDeltaY error: calculate alpha failed")
	}
	// beta = 1 - 1 / (1 + alpha)
	beta0 := safe_number.SafeNumFromString("1").Add(alpha)
	if beta0 == nil {
		return nil, errors.New("calculateDeltaY error: calculate beta0 failed")
	}
	beta1 := safe_number.SafeNumFromString("1").DivideBy(beta0)
	if beta1 == nil {
		return nil, errors.New("calculateDeltaY error: calculate beta1 failed")
	}
	beta := safe_number.SafeNumFromString("1").Subtract(beta1)
	if beta == nil {
		return nil, errors.New("calculateDeltaY error: calculate beta failed")
	}
	// deltaY = beta * Y
	deltaY := beta.Multiply(Y)
	if deltaY == nil {
		return nil, errors.New("calculateDeltaY error: calculate deltaY failed")
	}
	return deltaY, nil
}

func performSwap(instruction instruction_set.OpSwapInstruction, db *db_utils.OrdDB, lpMeta *memory_const.LPMeta) error {
	address := instruction.TxOutAddr
	// 1. preparing X, Y, spendingTick and buyingTick
	spendingTick := instruction.Spend
	buyingTick := ""
	lTick, rTick, consumingAmt := instruction.ExtractParams()
	if lTick == nil || rTick == nil || consumingAmt == nil {
		return errors.New("performSwap error: params extracting error")
	}
	var X *safe_number.SafeNum
	var Y *safe_number.SafeNum
	if spendingTick == *lTick {
		X = lpMeta.LAmt
		Y = lpMeta.RAmt
		buyingTick = *rTick
	} else if spendingTick == *rTick {
		X = lpMeta.RAmt
		Y = lpMeta.LAmt
		buyingTick = *lTick
	} else {
		return errors.New("performSwap error: Neither lTick nor rTick is equal to tick")
	}
	// 2. calculate odfi taker fee, lp taker fee and deltaX
	odfiTakerFee, err := getODFITakerFee(instruction, db)
	if err != nil {
		return err
	}
	if odfiTakerFee == nil {
		return errors.New("performSwap error: odfiTakerFee is nil")
	}
	lpTakerFee, err := getLPTakerFee(instruction, db)
	if err != nil {
		return err
	}
	if lpTakerFee == nil {
		return errors.New("performSwap error: lpTakerFee is nil")
	}
	// 2.1 if odfi-spending LP is equal to current lp, add odfiTakerFee to lpTakerFee, and set odfiTakerFee to zero.
	if spendingTick != "odfi" { // odfiTakerFee is fee when spending odfi, jump over this logic if it's odfi
		odfiLPName := memory_const.LPNameByTicks("odfi", spendingTick)
		lpName := memory_const.LPNameByTicks(buyingTick, spendingTick)
		if odfiLPName == nil || lpName == nil {
			return errors.New("performSwap tick error: calculate odfiLPName failed or calculate lpName failed")
		}
		if *odfiLPName == *lpName {
			lpTakerFee = lpTakerFee.Add(odfiTakerFee)
			if lpTakerFee == nil {
				return errors.New("performSwap calculate lpTakerFee failed")
			}
			odfiTakerFee = safe_number.SafeNumFromString("0")
			if odfiTakerFee == nil {
				return errors.New("performSwap calculate odfiTakerFee failed")
			}
		}
	}
	deltaX0 := consumingAmt.Subtract(odfiTakerFee)
	if deltaX0 == nil {
		return errors.New("performSwap calculate deltaX0 failed")
	}
	deltaX := deltaX0.Subtract(lpTakerFee)
	if deltaX == nil {
		return errors.New("performSwap calculate deltaX failed")
	}
	// 3. calculate deltaY
	deltaY, err := calculateDeltaY(deltaX, X, Y)
	if err != nil {
		return err
	}
	if deltaY == nil {
		return errors.New("performSwap calculate deltaY failed")
	}
	// 3.1 calculate slippage
	if instruction.Threshold != "" {
		slippageThreshold := safe_number.SafeNumFromString(instruction.Threshold)
		if slippageThreshold == nil {
			return errors.New("performSwap failed: parse slippage threshold failed:" + instruction.Threshold)
		}
		deltaPrice := deltaY.DivideBy(deltaX)
		if deltaPrice == nil {
			return errors.New("performSwap aborted: calculate deltaPrice failed")
		}
		originalPrice := Y.DivideBy(X)
		if originalPrice == nil {
			return errors.New("performSwap aborted: calculate originalPrice failed")
		}
		slippage := deltaPrice.DivideBy(originalPrice)
		if slippage == nil {
			return errors.New("performSwap aborted: calculate slippage failed")
		}
		slippage = safe_number.SafeNumFromString("1").Subtract(slippage)
		if slippage == nil {
			return errors.New("performSwap aborted: calculate slippage failed")
		}
		if slippage.Compare(slippageThreshold) > 0 {
			return fmt.Errorf("performSwap aborted: slippage larger than threshold: %s>%s", slippage.String(), instruction.Threshold)
		}
	}
	/*
		4.
		user spendingTick : - consumingAmt         (double-write)
		lp   spendingTick : + deltaX + lpTakerFee  (LPMeta)
		odfi-spendingTick : + odfiTakerFee         (LPMeta)
		---------------------------------------------------------
		user buyingTick   : + deltaY               (double-write)
		lp   buyingTick   : - deltaY               (LPMeta)
	*/
	err = memory_write.WriteSwapInfo(
		db,
		spendingTick,
		buyingTick,
		consumingAmt,
		lpTakerFee,
		odfiTakerFee,
		deltaX,
		deltaY,
		lpMeta,
		address)
	return err
}

func ExecuteOpSwap(instruction instruction_set.OpSwapInstruction, db *db_utils.OrdDB) error {
	address := instruction.TxOutAddr
	tick := instruction.Spend
	lTick, rTick, consumingAmt := instruction.ExtractParams()
	if lTick == nil || rTick == nil || consumingAmt == nil {
		return errors.New("ExecuteOpSwap error: params extracting error")
	}
	if tick != *lTick && tick != *rTick {
		return errors.New("ExecuteOpSwap error: Neither lTick nor rTick is equal to tick")
	}
	if *lTick == *rTick {
		return errors.New("ExecuteOpSwap error: lTick is equal to tick")
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
	lpMeta, err := memory_read.LiquidityProviderMetadata(db, *lTick, *rTick)
	if err != nil {
		return err
	}
	if lpMeta == nil {
		return errors.New("performSwap error: get LPMeta failed")
	}
	return performSwap(instruction, db, lpMeta)
}

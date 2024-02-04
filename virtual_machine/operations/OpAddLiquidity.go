package operations

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_write"
	"errors"
	"fmt"
)

func createLP(instruction instruction_set.OpAddLiquidityProviderInstruction, db *db_utils.OrdDB) error {
	// extract params
	lTick, rTick, lAmt, rAmt := instruction.ExtractParams()
	if lTick == nil || rTick == nil || lAmt == nil || rAmt == nil {
		return errors.New("OpAddLiquidityProvider error: params extracting error")
	}
	address := instruction.TxOutAddr
	err := memory_write.WriteCreateLPInfo(db, *lTick, *rTick, lAmt, rAmt, address)
	return err
}

func addToExistingLP(instruction instruction_set.OpAddLiquidityProviderInstruction, db *db_utils.OrdDB, lpMeta *memory_const.LPMeta) error {
	address := instruction.TxOutAddr
	lTick, rTick, lAmt, rAmt := instruction.ExtractParams()
	if lTick == nil || rTick == nil || lAmt == nil || rAmt == nil {
		return errors.New("OpAddLiquidityProvider error: params extracting error")
	}
	if lpMeta == nil {
		return errors.New("addToExistingLP failed: lpMeta is nil")
	}
	x := lpMeta.LAmt
	if x == nil {
		return errors.New("addToExistingLP failed: lpMeta.LAmt is nil")
	}
	y := lpMeta.RAmt
	if y == nil {
		return errors.New("addToExistingLP failed: lpMeta.RAmt is nil")
	}
	addingRatio := lAmt.DivideBy(rAmt)
	if addingRatio == nil {
		return fmt.Errorf("calulate addingRatio error: %s / %s", lAmt.String(), rAmt.String())
	}
	lpRatio := x.DivideBy(y)
	if lpRatio == nil {
		return fmt.Errorf("calulate lpRatio error: %s / %s", x.String(), y.String())
	}
	consumingLAmt := lAmt
	consumingRAmt := rAmt
	cmpRes := addingRatio.Compare(lpRatio)
	if cmpRes > 0 {
		// addingX exceed mixed amount
		consumingLAmt = rAmt.Multiply(lpRatio)
		if consumingLAmt == nil {
			return fmt.Errorf("calculate consumingLAmt error: %s * %s", rAmt.String(), lpRatio.String())
		}
	} else if cmpRes < 0 {
		// addingY exceed mixed amount
		consumingRAmt = lAmt.DivideBy(lpRatio)
		if consumingRAmt == nil {
			return fmt.Errorf("calculate consumingRAmt error: %s / %s", lAmt.String(), lpRatio.String())
		}
	}
	// addingLPAmount = (consumingLAmt / x) * lpMeta.Total
	// to avoid accuracy issue, calc multiply first
	// addingLPAmount = lpMeta.Total * consumingLAmt / x
	addingLPAmount0 := lpMeta.Total.Multiply(consumingLAmt)
	if addingLPAmount0 == nil {
		return fmt.Errorf("calulate addingLPAmount0 error: %s * %s", lpMeta.Total.String(), consumingLAmt.String())
	}
	addingLPAmount := addingLPAmount0.DivideBy(x)
	if addingLPAmount == nil {
		return fmt.Errorf("calulate addingLPAmount error: %s / %s", addingLPAmount0.String(), x.String())
	}
	err := memory_write.WriteAddToExistingLPInfo(db, *lTick, *rTick, consumingLAmt, consumingRAmt, addingLPAmount, lpMeta, address)
	return err
}

func ExecuteOpAddLiquidityProvider(instruction instruction_set.OpAddLiquidityProviderInstruction, db *db_utils.OrdDB) error {
	lTick, rTick, lAmt, rAmt := instruction.ExtractParams()
	if lTick == nil || rTick == nil || lAmt == nil || rAmt == nil {
		return errors.New("OpAddLiquidityProvider error: params extracting error")
	}
	if lAmt.IsZero() || rAmt.IsZero() {
		return errors.New("createLP error: lamt or ramt is 0")
	}
	lpMeta, err := memory_read.LiquidityProviderMetadata(db, *lTick, *rTick)
	if err != nil {
		return err
	}
	if lpMeta == nil {
		return createLP(instruction, db)
	} else {
		return addToExistingLP(instruction, db, lpMeta)
	}
}

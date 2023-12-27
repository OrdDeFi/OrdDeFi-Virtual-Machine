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
	address := instruction.TxOutAddr
	if lAmt.IsZero() || rAmt.IsZero() {
		return errors.New("createLP error: lamt or ramt is 0")
	}
	err := memory_write.WriteCreateLPInfo(db, *lTick, *rTick, lAmt, rAmt, address)
	return err
}

func addToExistingLP(instruction instruction_set.OpAddLiquidityProviderInstruction, db *db_utils.OrdDB, lpMeta *memory_const.LPMeta) error {
	lTick, rTick, lAmt, rAmt := instruction.ExtractParams()
	x := lpMeta.LAmt
	y := lpMeta.RAmt
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
	} else if cmpRes < 0 {
		// addingY exceed mixed amount
		consumingRAmt = lAmt.DivideBy(lpRatio)
	}
	addingLPRatio := consumingLAmt.DivideBy(x)
	if addingLPRatio == nil {
		return fmt.Errorf("calulate addingLPRatio error: %s / %s", consumingLAmt.String(), x.String())
	}
	addingLPAmount := addingLPRatio.Multiply(lpMeta.Total)
	if addingLPAmount == nil {
		return fmt.Errorf("calulate addingLPAmount error: %s * %s", addingLPRatio.String(), lpMeta.Total.String())
	}
	println(lTick, rTick, x.String(), y.String(), lAmt.String(), rAmt.String(), consumingLAmt.String(), consumingRAmt.String(), addingLPAmount.String())
	return nil
}

func ExecuteOpAddLiquidityProvider(instruction instruction_set.OpAddLiquidityProviderInstruction, db *db_utils.OrdDB) error {
	if instruction.TxInAddr != instruction.TxOutAddr {
		return errors.New("no privileges on cross-address add liquidity provider")
	}
	lTick, rTick, lAmt, rAmt := instruction.ExtractParams()
	if lTick == nil || rTick == nil || lAmt == nil || rAmt == nil {
		return errors.New("OpAddLiquidityProvider error: params extracting error")
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

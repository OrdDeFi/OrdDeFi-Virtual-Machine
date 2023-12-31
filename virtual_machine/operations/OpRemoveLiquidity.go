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

func removeLP(instruction instruction_set.OpRemoveLiquidityProviderInstruction, db *db_utils.OrdDB, lpMeta *memory_const.LPMeta) error {
	address := instruction.TxOutAddr
	lTick, rTick, consumingLPAmount := instruction.ExtractParams()
	if lTick == nil || rTick == nil || consumingLPAmount == nil {
		return errors.New("removeLP error: params extracting error")
	}
	x := lpMeta.LAmt
	if x == nil {
		return errors.New("removeLP error: lpMeta.LAmt is nil")
	}
	y := lpMeta.RAmt
	if y == nil {
		return errors.New("removeLP error: lpMeta.RAmt is nil")
	}
	total := lpMeta.Total
	if total == nil {
		return errors.New("removeLP error: lpMeta.Total is nil")
	}
	lpRatio := consumingLPAmount.DivideBy(total)
	if lpRatio == nil {
		return fmt.Errorf("calulate lpRatio error: %s / %s", consumingLPAmount.String(), total.String())
	}
	addingLAmt := x.Multiply(lpRatio)
	if addingLAmt == nil {
		return fmt.Errorf("calulate addingLAmt error: %s * %s", x.String(), lpRatio.String())
	}
	addingRAmt := y.Multiply(lpRatio)
	if addingRAmt == nil {
		return fmt.Errorf("calulate addingRAmt error: %s * %s", y.String(), lpRatio.String())
	}
	err := memory_write.WriteRemoveLPInfo(db, *lTick, *rTick, addingLAmt, addingRAmt, consumingLPAmount, lpMeta, address)
	return err
}

func ExecuteOpRemoveLiquidityProvider(instruction instruction_set.OpRemoveLiquidityProviderInstruction, db *db_utils.OrdDB) error {
	if instruction.TxInAddr != instruction.TxOutAddr {
		return errors.New("no privileges on cross-address remove liquidity provider")
	}
	lTick, rTick, consumingLPAmount := instruction.ExtractParams()
	if lTick == nil || rTick == nil || consumingLPAmount == nil {
		return errors.New("ExecuteOpRemoveLiquidityProvider error: params extracting error")
	}
	if consumingLPAmount.IsZero() {
		return errors.New("ExecuteOpRemoveLiquidityProvider error: consumingLPAmount is zero")
	}
	lpMeta, err := memory_read.LiquidityProviderMetadata(db, *lTick, *rTick)
	if err != nil {
		return err
	}
	if lpMeta == nil {
		return errors.New("ExecuteOpRemoveLiquidityProvider error: LP not exist")
	}
	return removeLP(instruction, db, lpMeta)
}

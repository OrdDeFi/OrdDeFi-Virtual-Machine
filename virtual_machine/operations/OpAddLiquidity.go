package operations

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_write"
	"errors"
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
	println(lTick, rTick, x.String(), y.String(), lAmt.String(), rAmt.String())
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

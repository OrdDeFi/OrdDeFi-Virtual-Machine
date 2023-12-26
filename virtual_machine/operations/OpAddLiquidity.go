package operations

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
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

func addToExistingLP(instruction instruction_set.OpAddLiquidityProviderInstruction, db *db_utils.OrdDB, coinMap map[string]safe_number.SafeNum) error {
	lTick, rTick, lAmt, rAmt := instruction.ExtractParams()
	x := coinMap[*lTick]
	y := coinMap[*rTick]
	println(x.String(), y.String(), lAmt.String(), rAmt.String())
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
	coinMap, err := memory_read.LiquidityProviderMetadata(*lTick, *rTick)
	if err != nil {
		return err
	}
	if coinMap == nil {
		return createLP(instruction, db)
	} else {
		return addToExistingLP(instruction, db, coinMap)
	}
}

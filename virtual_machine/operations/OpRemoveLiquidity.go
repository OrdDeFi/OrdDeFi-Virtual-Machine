package operations

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"errors"
)

func ExecuteOpRemoveLiquidityProvider(instruction instruction_set.OpRemoveLiquidityProviderInstruction, db *db_utils.OrdDB) error {
	if instruction.TxInAddr != instruction.TxOutAddr {
		return errors.New("no privileges on cross-address remove liquidity provider")
	}
	return nil
}

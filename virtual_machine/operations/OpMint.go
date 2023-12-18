package operations

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"errors"
)

func ExecuteOpMint(instruction instruction_set.OpMintInstruction, db *db_utils.OrdDB) error {
	coinName := instruction.Tick
	coinMeta, err := memory_read.CoinMeta(db, coinName)
	if err != nil {
		return err
	}
	if coinMeta == nil {
		return errors.New("CoinMeta not found named " + coinName)
	}

	return nil
}

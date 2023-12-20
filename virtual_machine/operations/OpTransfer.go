package operations

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"errors"
	"github.com/btcsuite/btcd/wire"
)

func executeImmediateTransfer(instruction instruction_set.OpTransferInstruction, db *db_utils.OrdDB) error {
	if instruction.TxInAddr != instruction.TxOutAddr {
		return errors.New("no privileges on cross-address immediate transfer")
	}
	// remove from current address available, add to "to" address available
	return nil
}

func executeUTXOTransfer(instruction instruction_set.OpTransferInstruction, db *db_utils.OrdDB) error {
	// remove from current address available, add to current address transferable
	// save a record on instruction.TxId, content: coinName:amountString
	return nil
}

func ExecuteTransfer(instruction instruction_set.OpTransferInstruction, db *db_utils.OrdDB) error {
	if instruction.To != "" {
		return executeImmediateTransfer(instruction, db)
	} else {
		return executeUTXOTransfer(instruction, db)
	}
}

func ApplyUTXOTransfer(tx *wire.MsgTx) error {
	if tx == nil {
		return errors.New("tx is nil")
	}
	return nil
}

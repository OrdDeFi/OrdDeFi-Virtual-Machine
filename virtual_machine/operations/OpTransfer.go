package operations

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_write"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/wire"
)

func performTransferBatchWriteKV(
	db *db_utils.OrdDB, coinName string,
	fromAccount string, fromSubAccount string,
	toAccount string, toSubAccount string,
	amount *safe_number.SafeNum) (map[string]string, error) {
	if amount == nil {
		return nil, errors.New("performTransferBatchWriteKV Transfer amount is nil")
	}
	if amount.IsZero() {
		return nil, errors.New("performTransferBatchWriteKV Transfer amount is 0")
	}
	var fromBalance *safe_number.SafeNum
	var toBalance *safe_number.SafeNum
	var err error
	// read from balance
	if fromSubAccount == db_utils.AvailableSubAccount {
		fromBalance, err = memory_read.AvailableBalance(db, coinName, fromAccount)
	} else if fromSubAccount == db_utils.TransferableSubAccount {
		fromBalance, err = memory_read.TransferableBalance(db, coinName, fromAccount)
	} else {
		return nil, errors.New("performTransferBatchWriteKV Sub-account error: " + fromSubAccount)
	}
	if err != nil {
		return nil, err
	}
	// read to balance
	if toSubAccount == db_utils.AvailableSubAccount {
		toBalance, err = memory_read.AvailableBalance(db, coinName, toAccount)
	} else if toSubAccount == db_utils.TransferableSubAccount {
		toBalance, err = memory_read.TransferableBalance(db, coinName, toAccount)
	} else {
		return nil, errors.New("performTransferBatchWriteKV Sub-account error: " + toSubAccount)
	}
	if err != nil {
		return nil, err
	}
	fromBalanceUpdated := fromBalance.Subtract(amount)
	if fromBalanceUpdated == nil {
		return nil, fmt.Errorf("performTransferBatchWriteKV from address balance error: %s - %s", fromBalance.String(), amount.String())
	}
	if fromBalanceUpdated.IsNegative() {
		return nil, fmt.Errorf("performTransferBatchWriteKV from address balance error: negative %s", fromBalanceUpdated.String())
	}
	toBalanceUpdated := toBalance.Add(amount)
	if toBalanceUpdated == nil {
		return nil, fmt.Errorf("performTransferBatchWriteKV to address balance error: %s + %s", toBalance.String(), amount.String())
	}
	if toBalanceUpdated.IsNegative() {
		return nil, fmt.Errorf("performTransferBatchWriteKV from address balance error: negative %s", toBalanceUpdated.String())
	}
	if fromBalanceUpdated.Add(toBalanceUpdated).IsEqualTo(fromBalance.Add(toBalance)) == false {
		return nil, fmt.Errorf("performTransferBatchWriteKV before calculation and after are not equal")
	}
	updateFromKV := memory_write.CoinBalanceDoubleWriteKV(coinName, fromAccount, fromBalanceUpdated.String(), fromSubAccount)
	if updateFromKV == nil {
		return nil, errors.New("performTransferBatchWriteKV updateFromKV generating error")
	}
	updateToKV := memory_write.CoinBalanceDoubleWriteKV(coinName, toAccount, toBalanceUpdated.String(), toSubAccount)
	if updateToKV == nil {
		return nil, errors.New("performTransferBatchWriteKV updateToKV generating error")
	}
	for k, v := range updateToKV {
		updateFromKV[k] = v
	}
	return updateFromKV, nil
}

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

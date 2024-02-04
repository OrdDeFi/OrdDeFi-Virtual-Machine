package authentication

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDeFi-Virtual-Machine/tx_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"errors"
	"github.com/btcsuite/btcd/wire"
)

func InstructionShouldBeAuthed(instruction instruction_set.AbstractInstruction) bool {
	op := instruction.Op
	if op == instruction_set.OpNameAddLiquidityProvider ||
		op == instruction_set.OpNameRemoveLiquidityProvider ||
		op == instruction_set.OpNameSwap {
		return true
	}
	if op == instruction_set.OpNameTransfer && instruction.To != "" {
		return true
	}
	return false
}

func InstructionAuthenticate(tx *wire.MsgTx) (*bool, error) {
	firstInputAddress, err := tx_utils.ParseFirstInputAddress(tx)
	if err != nil {
		return nil, err
	}
	if firstInputAddress == nil {
		return nil, errors.New("InstructionAuthenticate ParseFirstInputAddress got empty address")
	}
	firstOutputAddress, err := tx_utils.ParseFirstOutputAddress(tx)
	if err != nil {
		return nil, err
	}
	if firstOutputAddress == nil {
		return nil, errors.New("InstructionAuthenticate ParseFirstOutputAddress got empty address")
	}
	if *firstInputAddress == *firstOutputAddress {
		result := true
		return &result, nil
	}
	prevTxId := tx.TxIn[0].PreviousOutPoint.Hash.String()
	rawPrevTx := bitcoin_cli_channel.GetRawTransaction(prevTxId)
	if rawPrevTx == nil {
		return nil, errors.New("InstructionAuthenticate GetRawTransaction got empty rawPrevTx")
	}
	prevTx := bitcoin_cli_channel.DecodeRawTransaction(*rawPrevTx)
	if prevTx == nil {
		return nil, errors.New("InstructionAuthenticate DecodeRawTransaction got empty prevTx")
	}
	prevInputAddress, err := tx_utils.ParseFirstInputAddress(prevTx)
	if err != nil {
		return nil, err
	}
	if prevInputAddress == nil {
		return nil, errors.New("InstructionAuthenticate ParseFirstInputAddress got empty prevInputAddress")
	}
	commitInputAddrEqualToRevealOutputAddr := *prevInputAddress == *firstOutputAddress
	commitTxLastOutputIsOrdDeFiAuth := false
	prevLastOutput := prevTx.TxOut[len(prevTx.TxOut)-1]
	if prevLastOutput.Value == 0 {
		if len(prevLastOutput.PkScript) >= 2 {
			data := prevLastOutput.PkScript[2:]
			dataStr := string(data)
			if dataStr == "orddefi:auth" {
				commitTxLastOutputIsOrdDeFiAuth = true
			}
		}
	}
	if commitInputAddrEqualToRevealOutputAddr && commitTxLastOutputIsOrdDeFiAuth {
		result := true
		return &result, nil
	}

	result := false
	return &result, nil
}

package updater

import (
	"OrdDefi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDefi-Virtual-Machine/inscription_parser"
	"OrdDefi-Virtual-Machine/virtual_machine"
	"errors"
)

func UpdateBlockNumber(blockNumber int) {
	var err error
	blockHash := bitcoin_cli_channel.GetBlockHash(blockNumber)
	if blockHash == nil {
		err = errors.New("UpdateBlockNumber GetBlockHash failed")
		return
	}
	block := bitcoin_cli_channel.GetBlock(*blockHash)
	for _, txId := range block.Tx {
		rawTx := bitcoin_cli_channel.GetRawTransaction(txId)
		if rawTx == nil {
			err = errors.New("GetRawTransaction Failed")
			break
		}
		tx := bitcoin_cli_channel.DecodeRawTransaction(*rawTx)
		if tx == nil {
			err = errors.New("ParseRawTransaction -> DecodeRawTransaction Failed")
			break
		}
		contentType, content, err := inscription_parser.ParseTransactionToInscription(*tx)
		if err != nil {
			break
		}
		if contentType != nil && content != nil {
			instructions, err := virtual_machine.CompileInstructions(*contentType, content, tx, txId)
			if err != nil {
				break
			}
			virtual_machine.ExecuteInstructions(instructions)
		}
	}
	if err != nil {
		println("Updating block got error:", err) // failing
		println("Aborting update blocks...")
	}
}

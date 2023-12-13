package updater

import (
	"brc20defi_vm/bitcoin_cli_channel"
	"brc20defi_vm/inscription_parser"
	"brc20defi_vm/tx_utils"
	"errors"
)

func UpdateBlockNumber(blockNumber int) {
	blockHash := bitcoin_cli_channel.GetBlockHash(blockNumber)
	block := bitcoin_cli_channel.GetBlock(blockHash)
	var err error
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
			println("txId", txId)
			firstInputAddress, err := tx_utils.ParseFirstInputAddress(tx)
			if err != nil || firstInputAddress == nil {
				break
			}
			println("input", *firstInputAddress)
			firstOutputAddress, err := tx_utils.ParseFirstOutputAddress(tx)
			if err != nil || firstOutputAddress == nil {
				break
			}
			println("output", *firstOutputAddress)
			println(*contentType, len(content))
			println(string(content))
		}
	}
	if err != nil {
		println("Updating block got error:", err) // failing
		println("Aborting update blocks...")
	}
}

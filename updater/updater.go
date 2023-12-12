package updater

import (
	"brc20defi_vm/bitcoin_cli_channel"
	"brc20defi_vm/inscription_parser"
)

func UpdateBlockNumber(blockNumber int) {
	blockHash := bitcoin_cli_channel.GetBlockHash(blockNumber)
	block := bitcoin_cli_channel.GetBlock(blockHash)
	err := ""
	for _, txId := range block.Tx {
		rawTx := bitcoin_cli_channel.GetRawTransaction(txId)
		if rawTx == nil {
			err = "GetRawTransaction Failed"
			break
		}
		inscriptionContent, err := inscription_parser.ParseRawTransactionToInscription(*rawTx)
		if err != nil {
			break
		}
		if inscriptionContent != nil {
			println(*inscriptionContent)
		}
	}
	if err != "" {
		println(err) // failing
	}
}

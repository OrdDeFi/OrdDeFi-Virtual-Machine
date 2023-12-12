package updater

import "brc20defi_vm/bitcoin_cli_channel"

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
		tx := bitcoin_cli_channel.DecodeRawTransaction(*rawTx)
		if tx == nil {
			err = "DecodeRawTransaction Failed"
			break
		} else {
			for _, eachWitNess := range tx.TxIn[0].Witness {
				println(string(eachWitNess))
			}
		}
	}
	if err != "" {
		println(err) // failing
	}
}

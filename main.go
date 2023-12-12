package main

import "brc20defi_vm/bitcoin_cli_channel"

func main() {
	println("The Times 03/Jan/2009 Chancellor on brink of second bailout for banks.")
	println("BRC-20-DEFI indexer start to work.")
	blockNumber := bitcoin_cli_channel.GetBlockCount()
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
			println(string(tx.TxIn[0].SignatureScript))
		}
	}
	if err != "" {
		println(err) // failing
	}
}

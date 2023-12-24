package test

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDeFi-Virtual-Machine/tx_utils"
	"testing"
)

func TestParseCoinbaseTxInput(t *testing.T) {
	txId := TestingTxPool()[0]
	rawTx := bitcoin_cli_channel.GetRawTransaction(txId)
	if rawTx == nil {
		t.Errorf("TestParseCoinbaseTxInput GetRawTransaction error")
		return
	}
	tx := bitcoin_cli_channel.DecodeRawTransaction(*rawTx)
	if tx == nil {
		t.Errorf("TestParseCoinbaseTxInput DecodeRawTransaction error")
		return
	}
	for _, input := range tx.TxIn {
		address, inputValue, err := tx_utils.ParseInputAddressAndValue(input)
		if err != nil {
			t.Errorf("TestParseCoinbaseTxInput ParseInputAddressAndValue error %s", err)
			return
		}
		if address != nil || inputValue != nil {
			t.Errorf("TestParseCoinbaseTxInput ParseInputAddressAndValue error: for coinbase tx address and inputValue should be nil")
		}
	}
}

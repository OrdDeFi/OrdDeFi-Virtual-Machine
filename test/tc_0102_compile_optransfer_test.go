package test

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDeFi-Virtual-Machine/tx_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine"
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
	for i, input := range tx.TxIn {
		address, inputValue, err := tx_utils.ParseInputAddressAndValue(input)
		if address != nil && inputValue != nil {
			println(i, *address, *inputValue)
		} else {
			println(i, "nil", "nil")
		}
		if err != nil {
			t.Errorf("TestParseCoinbaseTxInput ParseInputAddressAndValue error %s", err)
			return
		}
		if address != nil || inputValue != nil {
			t.Errorf("TestParseCoinbaseTxInput ParseInputAddressAndValue error: for coinbase tx address and inputValue should be nil")
		}
	}
}

func TestValidToParam(t *testing.T) {
	txId := TestingTxPool()[0]
	commands := `{"p":"orddefi","op":"mint","tick":"odfi","amt":"1000","to":"bc1qm34lsc65zpw79lxes69zkqmk6ee3ewf0j77s3h"}`
	rawTx := bitcoin_cli_channel.GetRawTransaction(txId)
	if rawTx == nil {
		t.Errorf("TestCommandParse GetRawTransaction error")
	}
	tx := bitcoin_cli_channel.DecodeRawTransaction(*rawTx)
	if tx == nil {
		t.Errorf("TestCommandParse DecodeRawTransaction error")
	}
	instructions, err := virtual_machine.CompileInstructions("text/plain", []byte(commands), tx, txId)
	if err != nil {
		t.Errorf("TestCommandParse CompileInstructions error: %s", err.Error())
	}
	if len(instructions) != 1 {
		t.Errorf("TestCommandParse instructions count is not 1")
	}
}

func TestInvalidToParam(t *testing.T) {
	txId := TestingTxPool()[0]
	commands := `{"p":"orddefi","op":"transfer","tick":"odfi","amt":"1000","to":"abcd"}`
	rawTx := bitcoin_cli_channel.GetRawTransaction(txId)
	if rawTx == nil {
		t.Errorf("TestCommandParse GetRawTransaction error")
		return
	}
	tx := bitcoin_cli_channel.DecodeRawTransaction(*rawTx)
	if tx == nil {
		t.Errorf("TestCommandParse DecodeRawTransaction error")
		return
	}
	instructions, err := virtual_machine.CompileInstructions("text/plain", []byte(commands), tx, txId)
	if err != nil {
		t.Errorf("TestCommandParse CompileInstructions error: %s", err.Error())
		return
	}
	if len(instructions) != 0 {
		t.Errorf("TestCommandParse instructions count is not 1")
		return
	}
}

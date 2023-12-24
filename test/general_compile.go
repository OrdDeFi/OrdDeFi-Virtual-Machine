package test

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDeFi-Virtual-Machine/virtual_machine"
	"testing"
)

func TestingCompile(t *testing.T, cmd string, txId string) []interface{} {
	rawTx := bitcoin_cli_channel.GetRawTransaction(txId)
	if rawTx == nil {
		t.Errorf("TestingCompile GetRawTransaction error")
		return nil
	}
	tx := bitcoin_cli_channel.DecodeRawTransaction(*rawTx)
	if tx == nil {
		t.Errorf("TestingCompile DecodeRawTransaction error")
		return nil
	}
	instructions, err := virtual_machine.CompileInstructions("text/plain", []byte(cmd), tx, txId)
	if err != nil {
		t.Errorf("TestingCompile CompileInstructions error: %s", err.Error())
		return nil
	}
	if len(instructions) != 1 {
		t.Errorf("TestingCompile instructions count is not 1")
		return nil
	}
	return instructions
}

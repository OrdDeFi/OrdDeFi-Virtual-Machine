package test

import (
	"OrdDefi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDefi-Virtual-Machine/virtual_machine"
	"OrdDefi-Virtual-Machine/virtual_machine/instruction_set"
	"testing"
)

func TestCompileMintInSingleCommand(t *testing.T) {
	commands := `[{"p":"orddefi","op":"mint","tick":"ODFI","amt":"1000"}]`
	txId := "a8d1df8510d5ac3ad1199ebd987464226e1900260ab5cb10a3d19f7dabd460bc"
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
	if len(instructions) != 1 {
		t.Errorf("TestCommandParse CompileInstructions error: instructions length should be 1")
		return
	}
	for _, instruction := range instructions {
		switch instruction.(type) {
		case instruction_set.OpMintInstruction:
			println("succeed")
		default:
			t.Errorf("TestDeployInSingleCommand error: instruction type error, expected OpDeployInstruction")
		}
	}
}

func TestCompileMintInSingleSliceCommands(t *testing.T) {
	commands := `[
		{"p":"orddefi","op":"mint","tick":"ODFI","amt":"1000"}
	]`
	txId := "a8d1df8510d5ac3ad1199ebd987464226e1900260ab5cb10a3d19f7dabd460bc"
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
	if len(instructions) != 1 {
		t.Errorf("TestCommandParse CompileInstructions error: instructions length should be 1")
		return
	}
	for _, instruction := range instructions {
		switch instruction.(type) {
		case instruction_set.OpMintInstruction:
			println("succeed")
		default:
			t.Errorf("TestDeployInSingleCommand error: instruction type error, expected OpDeployInstruction")
		}
	}
}

func TestCompileMintInBatchCommands(t *testing.T) {
	commands := `[
		{"p":"orddefi","op":"mint","tick":"ODFI","amt":"1000"},
		{"p":"orddefi","op":"mint","tick":"ODFI","amt":"1000"}
	]`
	txId := "a8d1df8510d5ac3ad1199ebd987464226e1900260ab5cb10a3d19f7dabd460bc"
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
		t.Errorf("TestCommandParse CompileInstructions error: instructions should be nil")
		return
	}
}

package test

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDeFi-Virtual-Machine/virtual_machine"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"fmt"
	"testing"
)

// single command test
func TestCompileMintInSingleCommand(t *testing.T) {
	commands := `{"p":"orddefi","op":"mint","tick":"ODFI","amt":"1000"}`
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

// Test a single command using a single slice
func TestCompileMintInSingleSliceCommands(t *testing.T) {
	_, err := TestingMintInSingleSliceCommands("odfi", "a8d1df8510d5ac3ad1199ebd987464226e1900260ab5cb10a3d19f7dabd460bc", "1000")
	if err != nil {
		t.Errorf("%s", err.Error())
	}
}

// Test multiple mint commands using json array
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
	if len(instructions) != 0 {
		t.Errorf("TestCommandParse CompileInstructions error: If there is more than one batch command, the deploy/mint/transfer command cannot be included in the batch command")
		return
	}

	if err != nil {
		t.Errorf("TestCommandParse CompileInstructions error: %s", err.Error())
		return
	}
	if len(instructions) != 0 {
		t.Errorf("TestCommandParse CompileInstructions error: instructions should be nil")
		return
	}
}

func TestingMintInSingleSliceCommands(tick string, txId string, amt string) (*instruction_set.OpMintInstruction, error) {
	commands := `[
		{"p":"orddefi","op":"mint","tick":"` + tick + `","amt":"` + amt + `"}
	]`
	rawTx := bitcoin_cli_channel.GetRawTransaction(txId)
	if rawTx == nil {
		return nil, fmt.Errorf("TestCommandParse GetRawTransaction error")
	}
	tx := bitcoin_cli_channel.DecodeRawTransaction(*rawTx)
	if tx == nil {
		return nil, fmt.Errorf("TestCommandParse DecodeRawTransaction error")
	}
	instructions, err := virtual_machine.CompileInstructions("text/plain", []byte(commands), tx, txId)
	if err != nil {
		return nil, fmt.Errorf("TestCommandParse CompileInstructions error: %s", err.Error())
	}
	if len(instructions) != 1 {
		return nil, fmt.Errorf("TestCommandParse CompileInstructions error: instructions length should be 1")
	}
	for _, instruction := range instructions {
		switch value := instruction.(type) {
		case instruction_set.OpMintInstruction:
			return &value, nil
		default:
			return nil, fmt.Errorf("TestDeployInSingleCommand error: instruction type error, expected OpDeployInstruction")
		}
	}
	return nil, fmt.Errorf("TestDeployInSingleCommand error: instruction type error: no instruction compiled")
}

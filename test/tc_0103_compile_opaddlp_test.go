package test

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDeFi-Virtual-Machine/virtual_machine"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"testing"
)

func TestCompileAddLP(t *testing.T) {
	commands := `[{"p":"orddefi","op":"addlp","ltick":"ODFI","lamt":"1000","rtick":"ODFI","rlamt":"1000"}]`
	txId := "61de96170018ce878b1adf287b8ac9cf0e4f0ad8c5a69af203cc25bbde72a13e"
	rawTx := bitcoin_cli_channel.GetRawTransaction(txId)
	if rawTx == nil {
		t.Errorf("TestCompileAddLP GetRawTransaction error")
		return
	}
	tx := bitcoin_cli_channel.DecodeRawTransaction(*rawTx)
	if tx == nil {
		t.Errorf("TestCompileAddLP DecodeRawTransaction error")
		return
	}
	instructions, err := virtual_machine.CompileInstructions("text/plain", []byte(commands), tx, txId)
	if err != nil {
		t.Errorf("TestCommandParse CompileInstructions error: %s", err.Error())
		return
	}
	if len(instructions) != 1 {
		t.Errorf("TestCompileAddLP CompileInstructions error: instructions length should be 1")
		return
	}
	for _, instruction := range instructions {
		switch instruction.(type) {
		case instruction_set.OpAddLiquidityProviderInstruction:
			println("succeed")
		default:
			t.Errorf("TestCompileAddLP error: instruction type error, expected OpAddLiquidityProviderInstruction")
		}
	}
}

package test

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"testing"
)

func testingCompileRemoveLP(t *testing.T, commands string) (*string, *string, *safe_number.SafeNum) {
	txId := "61de96170018ce878b1adf287b8ac9cf0e4f0ad8c5a69af203cc25bbde72a13e"
	rawTx := bitcoin_cli_channel.GetRawTransaction(txId)
	if rawTx == nil {
		t.Errorf("testingCompileRemoveLP GetRawTransaction error")
		return nil, nil, nil
	}
	tx := bitcoin_cli_channel.DecodeRawTransaction(*rawTx)
	if tx == nil {
		t.Errorf("testingCompileRemoveLP DecodeRawTransaction error")
		return nil, nil, nil
	}
	instructions, err := virtual_machine.CompileInstructions("text/plain", []byte(commands), tx, txId)
	if err != nil {
		t.Errorf("TestCommandParse CompileInstructions error: %s", err.Error())
		return nil, nil, nil
	}
	if len(instructions) != 1 {
		t.Errorf("testingCompileRemoveLP CompileInstructions error: instructions length should be 1")
		return nil, nil, nil
	}
	for _, instruction := range instructions {
		switch value := instruction.(type) {
		case instruction_set.OpRemoveLiquidityProviderInstruction:
			return value.ExtractParams()
		default:
			t.Errorf("testingCompileRemoveLP error: instruction type error, expected OpAddLiquidityProviderInstruction")
		}
	}
	return nil, nil, nil
}

func TestCompileInvalidRemoveLP(t *testing.T) {
	commands := `[{"p":"orddefi","op":"rmlp","ltick":"ODFI","rtick":"ODFI","amt":"1001"}]`
	ltick, rtick, amt := testingCompileRemoveLP(t, commands)
	if ltick != nil {
		t.Errorf("error: ltick should be nil")
	}
	if rtick != nil {
		t.Errorf("error: rtick should be nil")
	}
	if amt != nil {
		t.Errorf("error: lamt should be nil")
	}
}

func TestCompileValidRemoveLP(t *testing.T) {
	commands := `[{"p":"orddefi","op":"rmlp","ltick":"ODFI","rtick":"ODGV","amt":"1001"}]`
	ltick, rtick, amt := testingCompileRemoveLP(t, commands)
	if *ltick != "odfi" {
		t.Errorf("error: ltick should be nil")
	}
	if *rtick != "odgv" {
		t.Errorf("error: rtick should be nil")
	}
	if amt.String() != "1001" {
		t.Errorf("error: lamt should be nil")
	}
}

func TestCompileValidRemoveLP2(t *testing.T) {
	commands := `[{"p":"orddefi","op":"rmlp","ltick":"ODGV","rtick":"ODFI","amt":"1001"}]`
	ltick, rtick, amt := testingCompileRemoveLP(t, commands)
	if *ltick != "odfi" {
		t.Errorf("error: ltick should be nil")
	}
	if *rtick != "odgv" {
		t.Errorf("error: rtick should be nil")
	}
	if amt.String() != "1001" {
		t.Errorf("error: lamt should be nil")
	}
}

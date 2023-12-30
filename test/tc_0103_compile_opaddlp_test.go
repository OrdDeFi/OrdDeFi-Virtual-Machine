package test

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"testing"
)

func testingCompileAddLP(t *testing.T, commands string) (*string, *string, *safe_number.SafeNum, *safe_number.SafeNum) {
	txId := "61de96170018ce878b1adf287b8ac9cf0e4f0ad8c5a69af203cc25bbde72a13e"
	rawTx := bitcoin_cli_channel.GetRawTransaction(txId)
	if rawTx == nil {
		t.Errorf("testingCompileAddLP GetRawTransaction error")
		return nil, nil, nil, nil
	}
	tx := bitcoin_cli_channel.DecodeRawTransaction(*rawTx)
	if tx == nil {
		t.Errorf("testingCompileAddLP DecodeRawTransaction error")
		return nil, nil, nil, nil
	}
	instructions, err := virtual_machine.CompileInstructions("text/plain", []byte(commands), tx, txId)
	if err != nil {
		t.Errorf("TestCommandParse CompileInstructions error: %s", err.Error())
		return nil, nil, nil, nil
	}
	if len(instructions) != 1 {
		t.Errorf("testingCompileAddLP CompileInstructions error: instructions length should be 1")
		return nil, nil, nil, nil
	}
	for _, instruction := range instructions {
		switch value := instruction.(type) {
		case instruction_set.OpAddLiquidityProviderInstruction:
			return value.ExtractParams()
		default:
			t.Errorf("testingCompileAddLP error: instruction type error, expected OpAddLiquidityProviderInstruction")
		}
	}
	return nil, nil, nil, nil
}

func TestCompileInvalidAddLP(t *testing.T) {
	// same tick
	commands := `[{"p":"orddefi","op":"addlp","ltick":"ODFI","lamt":"1000","rtick":"ODFI","ramt":"1001"}]`
	ltick, rtick, lamt, ramt := testingCompileAddLP(t, commands)
	if ltick != nil {
		t.Errorf("error: ltick should be nil")
	}
	if rtick != nil {
		t.Errorf("error: rtick should be nil")
	}
	if lamt != nil {
		t.Errorf("error: lamt should be nil")
	}
	if ramt != nil {
		t.Errorf("error: ramt should be nil")
	}
}

func TestCompileValidAddLP(t *testing.T) {
	commands := `[{"p":"orddefi","op":"addlp","ltick":"ODFI","lamt":"1000","rtick":"ODGV","ramt":"1001"}]`
	ltick, rtick, lamt, ramt := testingCompileAddLP(t, commands)
	if *ltick != "odfi" {
		t.Errorf("error: ltick should be odfi")
	}
	if *rtick != "odgv" {
		t.Errorf("error: rtick should be odgv")
	}
	if lamt.String() != "1000" {
		t.Errorf("error: lamt should be 1000")
	}
	if ramt.String() != "1001" {
		t.Errorf("error: ramt should be 1001")
	}
}

func TestCompileValidAddLP2(t *testing.T) {
	commands := `[{"p":"orddefi","op":"addlp","rtick":"ODFI","ramt":"1000","ltick":"ODGV","lamt":"1001"}]`
	ltick, rtick, lamt, ramt := testingCompileAddLP(t, commands)
	if *ltick != "odfi" {
		t.Errorf("error: ltick should be odfi")
	}
	if *rtick != "odgv" {
		t.Errorf("error: rtick should be odgv")
	}
	if lamt.String() != "1000" {
		t.Errorf("error: lamt should be 1000")
	}
	if ramt.String() != "1001" {
		t.Errorf("error: ramt should be 1001")
	}
}

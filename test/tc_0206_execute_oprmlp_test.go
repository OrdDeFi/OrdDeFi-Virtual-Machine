package test

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"OrdDeFi-Virtual-Machine/virtual_machine/operations"
	"encoding/json"
	"fmt"
	"testing"
)

func removeLiquidityProviderInstruction(
	txId string,
	lTick string, rTick string,
	amt string) (*instruction_set.OpRemoveLiquidityProviderInstruction, error) {
	var paramsMap map[string]string
	paramsMap = make(map[string]string)
	paramsMap["p"] = "orddefi"
	paramsMap["op"] = "rmlp"
	paramsMap["ltick"] = lTick
	paramsMap["rtick"] = rTick
	paramsMap["amt"] = amt
	jsonData, err := json.Marshal(paramsMap)
	if err != nil {
		return nil, err
	}
	commands := string(jsonData)
	rawTx := bitcoin_cli_channel.GetRawTransaction(txId)
	if rawTx == nil {
		return nil, fmt.Errorf("TestingTransferInSingleSliceCommands GetRawTransaction error")
	}
	tx := bitcoin_cli_channel.DecodeRawTransaction(*rawTx)
	if tx == nil {
		return nil, fmt.Errorf("TestingTransferInSingleSliceCommands DecodeRawTransaction error")
	}
	instructions, err := virtual_machine.CompileInstructions("text/plain", []byte(commands), tx, txId)
	if err != nil {
		return nil, fmt.Errorf("TestCommandParse CompileInstructions error: %s", err.Error())
	}
	if len(instructions) != 1 {
		return nil, fmt.Errorf("TestingTransferInSingleSliceCommands CompileInstructions error: instructions length should be 1")
	}
	for _, instruction := range instructions {
		switch value := instruction.(type) {
		case instruction_set.OpRemoveLiquidityProviderInstruction:
			return &value, nil
		default:
			return nil, fmt.Errorf("TestingTransferInSingleSliceCommands error: instruction type error, expected OpDeployInstruction")
		}
	}
	return nil, fmt.Errorf("TestingTransferInSingleSliceCommands error: instruction type error: no instruction compiled")
}

func TestRemoveLP(t *testing.T) {
	// open db
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestExecuteMint OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")

	txId := "61de96170018ce878b1adf287b8ac9cf0e4f0ad8c5a69af203cc25bbde72a13e"
	lTick := "odfi"
	rTick := "odgv"
	instruction, err := removeLiquidityProviderInstruction(txId, lTick, rTick, "100")
	if err != nil {
		t.Errorf("compile instruction error: %s", err.Error())
		return
	}
	if instruction == nil {
		t.Errorf("compile instruction error: instruction is nil")
		return
	}
	err = operations.ExecuteOpRemoveLiquidityProvider(*instruction, db)
	if err != nil {
		t.Errorf("execute ExecuteOpRemoveLiquidityProvider error: %s", err.Error())
		return
	}
	checkUserBalance(t, db, "bc1q2f0tczgrukdxjrhhadpft2fehzpcrwrz549u90")
}

func TestRemoveLP2(t *testing.T) {
	// open db
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestExecuteMint OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")

	txId := "2ad883c0532fbbb69380b46f72d53b9426f89bdcb95a29805fc9a11e0cfd6997"
	lTick := "odfi"
	rTick := "odgv"
	instruction, err := removeLiquidityProviderInstruction(txId, lTick, rTick, "10")
	if err != nil {
		t.Errorf("compile instruction error: %s", err.Error())
		return
	}
	if instruction == nil {
		t.Errorf("compile instruction error: instruction is nil")
		return
	}
	err = operations.ExecuteOpRemoveLiquidityProvider(*instruction, db)
	if err != nil {
		t.Errorf("execute ExecuteOpRemoveLiquidityProvider error: %s", err.Error())
		return
	}
	checkUserBalance(t, db, "bc1qr35hws365juz5rtlsjtvmulu97957kqvr3zpw3")
}

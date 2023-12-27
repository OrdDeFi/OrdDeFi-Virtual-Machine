package test

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"OrdDeFi-Virtual-Machine/virtual_machine/operations"
	"encoding/json"
	"fmt"
	"testing"
)

func addLiquidityProviderInstruction(
	txId string,
	lTick string, rTick string,
	lAmt string, rAmt string) (*instruction_set.OpAddLiquidityProviderInstruction, error) {
	var paramsMap map[string]string
	paramsMap = make(map[string]string)
	paramsMap["p"] = "orddefi"
	paramsMap["op"] = "addlp"
	paramsMap["ltick"] = lTick
	paramsMap["rtick"] = rTick
	paramsMap["lamt"] = lAmt
	paramsMap["ramt"] = rAmt
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
		case instruction_set.OpAddLiquidityProviderInstruction:
			return &value, nil
		default:
			return nil, fmt.Errorf("TestingTransferInSingleSliceCommands error: instruction type error, expected OpDeployInstruction")
		}
	}
	return nil, fmt.Errorf("TestingTransferInSingleSliceCommands error: instruction type error: no instruction compiled")
}

func checkUserBalance(t *testing.T, db *db_utils.OrdDB, address string) {
	println(address, ":")
	odfiA, odfiT, err := memory_read.Balance(db, "odfi", address)
	if err != nil {
		t.Errorf("checkUserBalance OpenDB error: %s", err.Error())
	}
	println("ODFI a/t:", odfiA.String(), odfiT.String())
	odgvA, odgvT, err := memory_read.Balance(db, "odgv", address)
	if err != nil {
		t.Errorf("checkUserBalance OpenDB error: %s", err.Error())
	}
	println("ODGV a/t:", odgvA.String(), odgvT.String())
	lpAmt, err := memory_read.LiquidityProviderBalance(db, "odfi", "odgv", address)
	if err != nil {
		t.Errorf("checkUserBalance OpenDB error: %s", err.Error())
	}
	println("ODFI-ODGV user balance:", lpAmt.String())
	lpMeta, err := memory_read.LiquidityProviderMetadata(db, "odfi", "odgv")
	if err != nil {
		t.Errorf("checkUserBalance OpenDB error: %s", err.Error())
	}
	lpMetaJSON, err := lpMeta.JsonString()
	if err != nil {
		t.Errorf("checkUserBalance jpMeta convert JSON error: %s", err.Error())
	}
	println("LP Meta:", *lpMetaJSON)
}

func TestToFindLegalTestingAddress(t *testing.T) {
	for _, txId := range TestingTxPool() {
		lTick := "odfi"
		rTick := "odgv"
		inscription, err := addLiquidityProviderInstruction(txId, lTick, rTick, "50", "100")
		if err == nil {
			println(inscription.TxOutAddr, txId)
		}
	}
}

func TestAddLP(t *testing.T) {
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
	instruction, err := addLiquidityProviderInstruction(txId, lTick, rTick, "50", "100")
	if err != nil {
		t.Errorf("compile instruction error: %s", err.Error())
		return
	}
	if instruction == nil {
		t.Errorf("compile instruction error: instruction is nil")
		return
	}
	err = operations.ExecuteOpAddLiquidityProvider(*instruction, db)
	if err != nil {
		t.Errorf("execute OpAddLiquidityProvider error: %s", err.Error())
		return
	}
	checkUserBalance(t, db, "bc1q2f0tczgrukdxjrhhadpft2fehzpcrwrz549u90")
}

func TestAddLP2(t *testing.T) {
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
	instruction, err := addLiquidityProviderInstruction(txId, lTick, rTick, "10", "10")
	if err != nil {
		t.Errorf("compile instruction error: %s", err.Error())
		return
	}
	if instruction == nil {
		t.Errorf("compile instruction error: instruction is nil")
		return
	}
	err = operations.ExecuteOpAddLiquidityProvider(*instruction, db)
	if err != nil {
		t.Errorf("execute OpAddLiquidityProvider error: %s", err.Error())
		return
	}
	checkUserBalance(t, db, "bc1qr35hws365juz5rtlsjtvmulu97957kqvr3zpw3")
}

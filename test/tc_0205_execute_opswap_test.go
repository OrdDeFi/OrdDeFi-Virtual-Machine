package test

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"OrdDeFi-Virtual-Machine/virtual_machine/operations"
	"encoding/json"
	"fmt"
	"testing"
)

func swapInstruction(
	txId string,
	lTick string, rTick string,
	spend string, amt string) (*instruction_set.OpSwapInstruction, error) {
	var paramsMap map[string]string
	paramsMap = make(map[string]string)
	paramsMap["p"] = "orddefi"
	paramsMap["op"] = "swap"
	paramsMap["ltick"] = lTick
	paramsMap["rtick"] = rTick
	paramsMap["spend"] = spend
	paramsMap["amt"] = amt
	jsonData, err := json.Marshal(paramsMap)
	if err != nil {
		return nil, err
	}
	commands := string(jsonData)
	rawTx := bitcoin_cli_channel.GetRawTransaction(txId)
	if rawTx == nil {
		return nil, fmt.Errorf("swapInstruction GetRawTransaction error")
	}
	tx := bitcoin_cli_channel.DecodeRawTransaction(*rawTx)
	if tx == nil {
		return nil, fmt.Errorf("swapInstruction DecodeRawTransaction error")
	}
	instructions, err := virtual_machine.CompileInstructions("text/plain", []byte(commands), tx, txId)
	if err != nil {
		return nil, fmt.Errorf("swapInstruction CompileInstructions error: %s", err.Error())
	}
	if len(instructions) != 1 {
		return nil, fmt.Errorf("swapInstruction CompileInstructions error: instructions length should be 1")
	}
	for _, instruction := range instructions {
		switch value := instruction.(type) {
		case instruction_set.OpSwapInstruction:
			return &value, nil
		default:
			return nil, fmt.Errorf("swapInstruction error: instruction type error, expected OpDeployInstruction")
		}
	}
	return nil, fmt.Errorf("swapInstruction error: instruction type error: no instruction compiled")
}

func checkLPMeta(t *testing.T, db *db_utils.OrdDB, lTick string, rTick string) {
	lpMeta, err := memory_read.LiquidityProviderMetadata(db, lTick, rTick)
	if err != nil {
		t.Errorf("checkLPMeta error: %s", err.Error())
		return
	}
	if lpMeta != nil {
		println("LPMeta:", lpMeta.LTick, lpMeta.LAmt.String(), lpMeta.RTick, lpMeta.RAmt.String())
	}
}

func checkStatusForSwap(t *testing.T, db *db_utils.OrdDB, address string, lTick string, rTick string, spendingTick string) {
	var odfiLPName *string
	odfiLPName = nil
	if spendingTick != "odfi" {
		odfiLPName = memory_const.LPNameByTicks("odfi", spendingTick)
		if odfiLPName == nil {
			t.Errorf("checkStatusForSwap tick error: calculate odfiLPName failed")
			return
		}
	}

	lpName := memory_const.LPNameByTicks(lTick, rTick)
	if lpName == nil {
		t.Errorf("checkStatusForSwap tick error: calculate lpName failed")
		return
	}
	checkLPMeta(t, db, lTick, rTick)
	if odfiLPName != nil && *odfiLPName != *lpName {
		checkLPMeta(t, db, "odfi", spendingTick)
	}
	lTickAvailableAmt, err := memory_read.AvailableBalance(db, lTick, address)
	if err != nil {
		t.Errorf("checkStatusForSwap tick error: read AvailableBalance failed %s", err.Error())
		return
	}
	rTickAvailableAmt, err := memory_read.AvailableBalance(db, rTick, address)
	if err != nil {
		t.Errorf("checkStatusForSwap tick error: read AvailableBalance failed %s", err.Error())
		return
	}
	println("User available", lTick, lTickAvailableAmt.String(), rTick, rTickAvailableAmt.String())
}

func testSwapForParams(t *testing.T, db *db_utils.OrdDB, txId string, address string, lTick string, rTick string, spendingTick string, amt string) {
	swap, err := swapInstruction(txId, lTick, rTick, spendingTick, amt)
	if err != nil {
		t.Errorf("generate swap instruction failed, error: %s", err.Error())
	}

	println("swap", lTick, rTick, "by", spendingTick, "for", amt)
	println("status before swap:")
	checkStatusForSwap(t, db, address, lTick, rTick, spendingTick)

	err = operations.ExecuteOpSwap(*swap, db)
	if err != nil {
		t.Errorf("generate swap instruction failed, error: %s", err.Error())
	}

	println("status after swap:")
	checkStatusForSwap(t, db, address, lTick, rTick, spendingTick)
}

func TestSwapForODFILP(t *testing.T) {
	// open db
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestExecuteMint OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)

	txId := "61de96170018ce878b1adf287b8ac9cf0e4f0ad8c5a69af203cc25bbde72a13e"
	address := "bc1q2f0tczgrukdxjrhhadpft2fehzpcrwrz549u90"
	println("Test 1")
	testSwapForParams(t, db, txId, address, "odfi", "odgv", "odgv", "10")
	println("Test 2")
	testSwapForParams(t, db, txId, address, "odfi", "odgv", "odfi", "10")
	println("Test 3")
	testSwapForParams(t, db, txId, address, "odgv", "odfi", "odgv", "10")
	println("Test 4")
	testSwapForParams(t, db, txId, address, "odgv", "odfi", "odfi", "10")
}

func TestSwapForNoneODFILP(t *testing.T) {
	// open db
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestExecuteMint OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)

	txId := "61de96170018ce878b1adf287b8ac9cf0e4f0ad8c5a69af203cc25bbde72a13e"
	address := "bc1q2f0tczgrukdxjrhhadpft2fehzpcrwrz549u90"
	println("Test 5")
	testSwapForParams(t, db, txId, address, "half", "odgv", "odgv", "10")
	println("Test 6")
	testSwapForParams(t, db, txId, address, "half", "odgv", "half", "10")
	println("Test 7")
	testSwapForParams(t, db, txId, address, "odgv", "half", "odgv", "10")
	println("Test 8")
	testSwapForParams(t, db, txId, address, "odgv", "half", "half", "10")
}

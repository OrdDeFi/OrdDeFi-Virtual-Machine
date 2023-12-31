package test

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"OrdDeFi-Virtual-Machine/virtual_machine/operations"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

func checkLP(t *testing.T, db *db_utils.OrdDB, address string, lTick string, rTick string) {
	lpAmt, err := memory_read.LiquidityProviderBalance(db, lTick, rTick, address)
	if err != nil {
		t.Errorf("checkUserBalance OpenDB error: %s", err.Error())
	}
	fmt.Printf("%s-%s user balance: %s\n", lTick, rTick, lpAmt.String())
	lpMeta, err := memory_read.LiquidityProviderMetadata(db, lTick, rTick)
	if err != nil {
		t.Errorf("checkLP OpenDB error: %s", err.Error())
		return
	}
	if lpMeta == nil {
		return
	}
	lpMetaJSON, err := lpMeta.JsonString()
	if err != nil {
		t.Errorf("checkUserBalance jpMeta convert JSON error: %s", err.Error())
	}
	println("LP Meta:", *lpMetaJSON)
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
	halfA, halfT, err := memory_read.Balance(db, "half", address)
	if err != nil {
		t.Errorf("checkUserBalance OpenDB error: %s", err.Error())
	}
	println("HALF a/t:", halfA.String(), halfT.String())
	checkLP(t, db, address, "odfi", "odgv")
	checkLP(t, db, address, "odgv", "half")
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
	db, err := db_utils.OpenDB(testDBPath)
	if err != nil {
		t.Errorf("TestExecuteMint OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")

	txId := "61de96170018ce878b1adf287b8ac9cf0e4f0ad8c5a69af203cc25bbde72a13e"
	txIDInputAdd := "bc1q2f0tczgrukdxjrhhadpft2fehzpcrwrz549u90"

	lTick := "odfi"
	rTick := "odgv"
	//1. check user balance
	lTickInitBalance, rTickInitBalance, lpInitAmt := getBalanceData(db, lTick, rTick, txIDInputAdd)

	//2. mint odfi and odgv
	if lTickInitBalance.IsZero() {
		TestingMintForParam(t, db, lTick, txId, "1000")
	}
	if rTickInitBalance.IsZero() {
		TestingMintForParam(t, db, rTick, txId, "1000")
	}
	//3. check user balance
	lTickInitBalance, rTickInitBalance, lpInitAmt = getBalanceData(db, lTick, rTick, txIDInputAdd)
	println("lTickInitBlance:", lTickInitBalance.String())
	println("rTickInitBlance:", rTickInitBalance.String())
	println("lpInitAmt:", lpInitAmt.String())

	//4. check lp meta
	lTickLiquidity := "50"
	rTickLiquidity := "100"
	lpMeta, _ := memory_read.LiquidityProviderMetadata(db, lTick, rTick)
	var lTickLiquidityExceptUsed, rTickLiquidityExceptUsed, usedRatio *safe_number.SafeNum
	if lpMeta != nil {
		lTickLiquidityExceptUsed, rTickLiquidityExceptUsed, usedRatio = genAddLiquidity(lTickLiquidity, rTickLiquidity, lpMeta.LAmt, lpMeta.RAmt)
	} else {
		lTickLiquidityExceptUsed = safe_number.SafeNumFromString(lTickLiquidity)
		rTickLiquidityExceptUsed = safe_number.SafeNumFromString(rTickLiquidity)
		// usedRatio = safe_number.SafeNumFromString("1")
	}

	instruction, err := addLiquidityProviderInstruction(txId, lTick, rTick, lTickLiquidity, rTickLiquidity)
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
	// lpAmt, err := memory_read.LiquidityProviderBalance(db, lTick, rTick, txIDInputAdd)
	lTickEndBalance, rTickEndBalance, lpEndAmt := getBalanceData(db, lTick, rTick, txIDInputAdd)

	println("lTickEndBlance:", lTickEndBalance.String())
	println("rTickEndBlance:", rTickEndBalance.String())
	println("lpEndAmt:", lpEndAmt.String())
	assert.True(t, lTickEndBalance.IsEqualTo(lTickInitBalance.Subtract(lTickLiquidityExceptUsed)))
	assert.True(t, rTickEndBalance.IsEqualTo(rTickInitBalance.Subtract(rTickLiquidityExceptUsed)))
	if lpMeta != nil {
		println("usedRatio:", usedRatio.String())
		println("expected lpEndAmt:", lpMeta.Total.Multiply(usedRatio).String())
		println("lpEndAmt:", lpEndAmt.String())
		assert.True(t, lpEndAmt.Subtract(lpInitAmt).IsEqualTo(lpMeta.Total.Multiply(usedRatio)))
	} else {
		assert.True(t, lpEndAmt.IsEqualTo(safe_number.SafeNumFromString("1000")))
	}

	checkUserBalance(t, db, "bc1q2f0tczgrukdxjrhhadpft2fehzpcrwrz549u90")
}

func TestAddLP2(t *testing.T) {
	// open db
	db, err := db_utils.OpenDB(testDBPath)
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

func TestAddLP3(t *testing.T) {
	// open db
	db, err := db_utils.OpenDB(testDBPath)
	if err != nil {
		t.Errorf("TestExecuteMint OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")

	txId := "61de96170018ce878b1adf287b8ac9cf0e4f0ad8c5a69af203cc25bbde72a13e"
	lTick := "half"
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

func genAddLiquidity(lTickLiquidity string, rTickLiquidity string, lExistLiquid *safe_number.SafeNum, rExistLiquid *safe_number.SafeNum) (*safe_number.SafeNum, *safe_number.SafeNum, *safe_number.SafeNum) {
	lTickLiquidNum := safe_number.SafeNumFromString(lTickLiquidity)
	rTickLiquidNum := safe_number.SafeNumFromString(rTickLiquidity)

	originalRatio := lExistLiquid.DivideBy(rExistLiquid)

	newLExistLiquid := lExistLiquid.Add(lTickLiquidNum)

	expectedRExistLiquid := newLExistLiquid.DivideBy(originalRatio)

	actualRExistLiquidIncrease := expectedRExistLiquid.Subtract(rExistLiquid)

	if actualRExistLiquidIncrease.IsGreaterThan(rTickLiquidNum) {
		returnedRTickLiquidity := rTickLiquidNum
		returnLTickLiquidity := rExistLiquid.Add(returnedRTickLiquidity).Multiply(originalRatio).Subtract(lExistLiquid)
		return returnLTickLiquidity, returnedRTickLiquidity, returnLTickLiquidity.DivideBy(lExistLiquid)
	} else {
		returnedRTickLiquidity := actualRExistLiquidIncrease
		returnLTickLiquidity := lTickLiquidNum
		return returnLTickLiquidity, returnedRTickLiquidity, returnLTickLiquidity.DivideBy(lExistLiquid)
	}
}

func getBalanceData(db *db_utils.OrdDB, lTick string, rTick string, txIDInputAdd string) (*safe_number.SafeNum, *safe_number.SafeNum, *safe_number.SafeNum) {
	lTickBalance, _ := memory_read.AvailableBalance(db, lTick, txIDInputAdd)
	rTickBalance, _ := memory_read.AvailableBalance(db, rTick, txIDInputAdd)
	lpAmt, _ := memory_read.LiquidityProviderBalance(db, lTick, rTick, txIDInputAdd)
	return lTickBalance, rTickBalance, lpAmt
}

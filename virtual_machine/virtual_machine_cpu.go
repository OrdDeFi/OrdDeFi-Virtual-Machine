package virtual_machine

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
	"OrdDeFi-Virtual-Machine/virtual_machine/operations"
	"os"
	"strconv"
)

func recordLog(logDB *db_utils.OrdDB, err error, instruction interface{}, blockNumber int, txIndex int, txId string) {
	rawInstruction := ""
	switch value := instruction.(type) {
	case instruction_set.OpDeployInstruction:
		rawInstruction = value.RawInstruction
	case instruction_set.OpMintInstruction:
		rawInstruction = value.RawInstruction
	case instruction_set.OpTransferInstruction:
		rawInstruction = value.RawInstruction
	case instruction_set.OpAddLiquidityProviderInstruction:
		rawInstruction = value.RawInstruction
	case instruction_set.OpRemoveLiquidityProviderInstruction:
		rawInstruction = value.RawInstruction
	case instruction_set.OpSwapInstruction:
		rawInstruction = value.RawInstruction
	}
	result := "succeed"
	if err != nil {
		result = "error: " + err.Error()
	}
	var batchKV map[string]string
	batchKV = make(map[string]string)
	key := memory_const.LogMainTable + ":" + strconv.Itoa(blockNumber) + ":" + strconv.Itoa(txIndex) + ":" + txId
	key2 := memory_const.LogQueryTxTable + ":" + txId
	value := result + ";;;;;" + rawInstruction
	batchKV[key] = value
	batchKV[key2] = value
	storeLogErr := logDB.StoreKeyValues(batchKV)
	if storeLogErr != nil {
		println("recordLog got error:", storeLogErr.Error())
		os.Exit(3)
	}
}

func executeInstruction(instruction interface{}, db *db_utils.OrdDB, logDB *db_utils.OrdDB, blockNumber int, txIndex int, txId string, verbose bool) {
	var err error
	switch value := instruction.(type) {
	case instruction_set.OpDeployInstruction:
		err = operations.ExecuteOpDeploy(value, db)
	case instruction_set.OpMintInstruction:
		err = operations.ExecuteOpMint(value, db)
	case instruction_set.OpTransferInstruction:
		err = operations.ExecuteTransfer(value, db)
	case instruction_set.OpAddLiquidityProviderInstruction:
		err = operations.ExecuteOpAddLiquidityProvider(value, db)
	case instruction_set.OpRemoveLiquidityProviderInstruction:
		err = operations.ExecuteOpRemoveLiquidityProvider(value, db)
	case instruction_set.OpSwapInstruction:
		err = operations.ExecuteOpSwap(value, db)
	}
	if verbose {
		if err != nil {
			println(blockNumber, txIndex, txId, "error:", err.Error())
		} else {
			println(blockNumber, txIndex, txId, "succeed.")
		}
	}
	recordLog(logDB, err, instruction, blockNumber, txIndex, txId)
}

func ExecuteInstructions(instructions []interface{}, db *db_utils.OrdDB, logDB *db_utils.OrdDB, blockNumber int, txIndex int, txId string, verbose bool) {
	for _, eachInstruction := range instructions {
		executeInstruction(eachInstruction, db, logDB, blockNumber, txIndex, txId, verbose)
	}
}

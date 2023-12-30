package virtual_machine

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"OrdDeFi-Virtual-Machine/virtual_machine/operations"
)

func recordLog(logDB *db_utils.OrdDB, err error) {

}

func executeInstruction(instruction interface{}, db *db_utils.OrdDB, logDB *db_utils.OrdDB) {
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
	recordLog(logDB, err)
}

func ExecuteInstructions(instructions []interface{}, db *db_utils.OrdDB, logDB *db_utils.OrdDB) {
	for _, eachInstruction := range instructions {
		executeInstruction(eachInstruction, db, logDB)
	}
}

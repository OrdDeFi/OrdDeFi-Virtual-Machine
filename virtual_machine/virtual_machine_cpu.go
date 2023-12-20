package virtual_machine

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"OrdDeFi-Virtual-Machine/virtual_machine/operations"
)

func executeInstruction(instruction interface{}, db *db_utils.OrdDB) {
	switch value := instruction.(type) {
	case instruction_set.OpDeployInstruction:
		operations.ExecuteOpDeploy(value, db)
	case instruction_set.OpMintInstruction:
		operations.ExecuteOpMint(value, db)
	case instruction_set.OpTransferInstruction:
		operations.ExecuteTransfer(value, db)
	case instruction_set.OpAddLiquidityProviderInstruction:
		operations.ExecuteOpAddLiquidityProvider(value, db)
	case instruction_set.OpRemoveLiquidityProviderInstruction:
		operations.ExecuteOpRemoveLiquidityProvider(value, db)
	case instruction_set.OpSwapInstruction:
		operations.ExecuteOpSwap(value, db)
	}
}

func ExecuteInstructions(instructions []interface{}, db *db_utils.OrdDB) {
	for _, eachInstruction := range instructions {
		executeInstruction(eachInstruction, db)
	}
}

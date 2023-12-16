package virtual_machine

import (
	"OrdDefi-Virtual-Machine/virtual_machine/instruction_set"
	"OrdDefi-Virtual-Machine/virtual_machine/operations"
)

func executeInstruction(instruction interface{}) {
	switch value := instruction.(type) {
	case instruction_set.OpDeployInstruction:
		operations.ExecuteOpDeploy(value)
	case instruction_set.OpMintInstruction:
		operations.ExecuteOpMint(value)
	case instruction_set.OpTransferInstruction:
		operations.ExecuteTransfer(value)
	case instruction_set.OpAddLiquidityProviderInstruction:
		operations.ExecuteOpAddLiquidityProvider(value)
	case instruction_set.OpRemoveLiquidityProviderInstruction:
		operations.ExecuteOpRemoveLiquidityProvider(value)
	case instruction_set.OpSwapInstruction:
		operations.ExecuteSwap(value)
	}
}

func ExecuteInstructions(instructions []interface{}) {
	for _, eachInstruction := range instructions {
		executeInstruction(eachInstruction)
	}
}

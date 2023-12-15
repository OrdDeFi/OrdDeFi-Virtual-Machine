package virtual_machine

func executeInstruction(instruction interface{}) {
	switch value := instruction.(type) {
	case OpDeployInstruction:
		ExecuteOpDeploy(value)
	case OpMintInstruction:
		ExecuteOpMint(value)
	case OpTransferInstruction:
		ExecuteTransfer(value)
	case OpAddLiquidityProviderInstruction:
		ExecuteOpAddLiquidityProvider(value)
	case OpRemoveLiquidityProviderInstruction:
		ExecuteOpRemoveLiquidityProvider(value)
	case OpSwapInstruction:
		ExecuteSwap(value)
	}
}

func ExecuteInstructions(instructions []interface{}) {
	for _, eachInstruction := range instructions {
		executeInstruction(eachInstruction)
	}
}

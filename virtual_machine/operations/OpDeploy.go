package operations

import (
	"OrdDefi-Virtual-Machine/safe_number"
	"OrdDefi-Virtual-Machine/virtual_machine/instruction_set"
)

func ExecuteOpDeploy(instruction instruction_set.OpDeployInstruction) {
	max := safe_number.SafeNumFromString(instruction.Max)
	println(max.String())
	return
}

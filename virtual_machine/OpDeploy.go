package virtual_machine

import (
	"OrdDefi-Virtual-Machine/safe_number"
)

func ExecuteOpDeploy(instruction OpDeployInstruction) {
	max := safe_number.SafeNumFromString(instruction.Max)
	println(max.String())
	return
}

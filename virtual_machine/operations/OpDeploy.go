package operations

import (
	"OrdDefi-Virtual-Machine/safe_number"
	"OrdDefi-Virtual-Machine/virtual_machine/instruction_set"
	"OrdDefi-Virtual-Machine/virtual_machine/memory"
	"strings"
)

func checkTickLegal(tick string) bool {
	if strings.Contains(tick, "-") {
		return false
	} else if strings.Contains(tick, "_") {
		return false
	} else if strings.HasPrefix(tick, "$") {
		return false
	} else if strings.HasPrefix(tick, "@") {
		return false
	} else if strings.HasPrefix(tick, "#") {
		return false
	} else if strings.HasPrefix(tick, "%") {
		return false
	}
	length := len(tick)
	return length == 4
}

func ExecuteOpDeploy(instruction instruction_set.OpDeployInstruction) {
	tick := instruction.Tick
	if checkTickLegal(tick) == false {
		return
	}
	maxValue := safe_number.SafeNumFromString(instruction.Max)
	lim := safe_number.SafeNumFromString(instruction.Lim)
	addrLim := safe_number.SafeNumFromString(instruction.AddrLim)
	desc := instruction.Desc
	icon := instruction.Icon
	if maxValue != nil && lim != nil {
		memory.CoinListSave(tick)
		memory.CoinMetadataSave(tick, maxValue, lim, addrLim, desc, icon)
	}
}

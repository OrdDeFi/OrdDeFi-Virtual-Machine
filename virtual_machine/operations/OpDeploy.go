package operations

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_write"
	"errors"
)

func ExecuteOpDeploy(instruction instruction_set.OpDeployInstruction, db *db_utils.OrdDB) error {
	tick := instruction.Tick
	maxValue := safe_number.SafeNumFromString(instruction.Max)
	lim := safe_number.SafeNumFromString(instruction.Lim)
	addrLim := safe_number.SafeNumFromString(instruction.AddrLim)
	desc := instruction.Desc
	icon := instruction.Icon
	if maxValue != nil && lim != nil {
		coinMeta, err := memory_read.CoinMeta(db, tick)
		if err != nil {
			return err
		}
		if coinMeta != nil {
			return errors.New("Coin exist for name: " + tick)
		}
		err = memory_write.WriteDeployInfo(db, tick, maxValue, lim, addrLim, desc, icon)
		if err != nil {
			return err
		}
	}
	return nil
}

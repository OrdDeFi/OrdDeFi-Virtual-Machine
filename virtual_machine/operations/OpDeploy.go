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
	if tick == "odfi" {
		return errors.New("ExecuteOpMint error: tick cannot be odfi")
	} else if tick == "odgv" {
		return errors.New("ExecuteOpMint error: tick cannot be odgv")
	}
	maxValue := safe_number.SafeNumFromString(instruction.Max)
	if maxValue == nil {
		return errors.New("ExecuteOpMint error: maxValue is nil")
	}
	lim := safe_number.SafeNumFromString(instruction.Lim)
	if lim == nil {
		return errors.New("ExecuteOpMint error: lim is nil")
	}
	addrLim := safe_number.SafeNumFromString(instruction.AddrLim)
	if addrLim == nil {
		addrLim = lim
	}
	desc := instruction.Desc
	icon := instruction.Icon
	coinMeta, err := memory_read.CoinMeta(db, tick)
	if err != nil {
		return err
	}
	if coinMeta != nil {
		return errors.New("Coin exist for name: " + tick)
	}
	err = memory_write.WriteDeployInfo(db, tick, maxValue, lim, addrLim, desc, icon)
	return err
}

package memory_read

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
	"errors"
	"strings"
)

/*
UTXOCarryingBalance
:return &address, &coinName, coinAmount, error
*/
func UTXOCarryingBalance(db *db_utils.OrdDB, txId string) (*string, *string, *safe_number.SafeNum, error) {
	key := memory_const.UTXOCarryingBalancePath(txId)
	tickAndValue, err := db.Read(key) // tickAndValue format: "fromAddress:tick:value"
	if err != nil {
		if err.Error() == "leveldb: not found" {
			return nil, nil, nil, nil
		} else {
			return nil, nil, nil, err
		}
	}
	if tickAndValue == nil {
		return nil, nil, nil, errors.New("UTXOCarryingBalance error: tickAndValue is nil")
	}
	comps := strings.Split(*tickAndValue, ":")
	if len(comps) != 3 {
		return nil, nil, nil, errors.New("UTXOCarryingBalance error: tickAndValue parse error")
	}
	address := comps[0]
	tick := comps[1]
	value := comps[2]
	if instruction_set.CheckTickLegal(tick) == false {
		return nil, nil, nil, errors.New("UTXOCarryingBalance error: tick is not legal: " + tick)
	}
	valueNumber := safe_number.SafeNumFromString(value)
	if valueNumber == nil {
		return nil, nil, nil, errors.New("UTXOCarryingBalance error: value is not legal number: " + value)
	}
	return &address, &tick, valueNumber, nil
}

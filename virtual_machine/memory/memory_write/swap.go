package memory_write

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"errors"
)

func ODFISpendingTickLPName(spendingTick string) *string {
	return memory_const.LPNameByTicks("odfi", spendingTick)
}

func WriteSwapInfo(
	db *db_utils.OrdDB,
	spendingTick string,
	buyingTick string,
	consumingAmt *safe_number.SafeNum,
	lpTakerFee *safe_number.SafeNum,
	odfiTakerFee *safe_number.SafeNum,
	deltaX *safe_number.SafeNum,
	deltaY *safe_number.SafeNum,
	lpMeta *memory_const.LPMeta,
	address string,
) error {
	/*
		4.
		1. user spendingTick : - consumingAmt         (double-write)
		2. user buyingTick   : + deltaY               (double-write)
		3. odfi-spendingTick : + odfiTakerFee         (if LPMeta exist for odfi-spendingTick)
		4. lp   spendingTick : + deltaX + lpTakerFee  (LPMeta)
		5. lp   buyingTick   : - deltaY               (LPMeta)
	*/
	var batchKV map[string]string
	batchKV = make(map[string]string)
	// 1. user spendingTick : - consumingAmt
	spendingAvailable, err := memory_read.AvailableBalance(db, spendingTick, address)
	if err != nil {
		return err
	}
	updatedSpendingAvailable := spendingAvailable.Subtract(consumingAmt)
	if updatedSpendingAvailable == nil {
		return errors.New("WriteSwapInfo failed: calculate updatedSpendingAvailable failed")
	}
	spendingDoubleWriteKV := CoinBalanceDoubleWriteKV(spendingTick, address, updatedSpendingAvailable.String(), db_utils.AvailableSubAccount)
	for k, v := range spendingDoubleWriteKV {
		batchKV[k] = v
	}
	// 3. odfi-spendingTick : + odfiTakerFee
	if odfiTakerFee.IsZero() == false {
		odfiLPMeta, err := memory_read.LiquidityProviderMetadata(db, "odfi", spendingTick)
		if err != nil {
			return err
		}
		if spendingTick == odfiLPMeta.LTick {
			odfiLPMeta.LAmt = odfiLPMeta.LAmt.Add(odfiTakerFee)
			if odfiLPMeta.LAmt == nil {
				return errors.New("WriteSwapInfo failed: calculate odfiLPMeta.LAmt failed")
			}
		} else if spendingTick == odfiLPMeta.RTick {
			odfiLPMeta.RAmt = odfiLPMeta.RAmt.Add(odfiTakerFee)
			if odfiLPMeta.RAmt == nil {
				return errors.New("WriteSwapInfo failed: calculate odfiLPMeta.RAmt failed")
			}
		} else {
			return errors.New("WriteSwapInfo failed: spending tick error")
		}
		odfiLPJson, err := odfiLPMeta.JsonString()
		if err != nil {
			return err
		}
		if odfiLPJson == nil {
			return errors.New("WriteSwapInfo failed: generate odfiLPJson failed")
		}
		lpName := memory_const.LPNameByTicks("odfi", spendingTick)
		if lpName == nil {
			return errors.New("WriteSwapInfo failed: calculate odfi-spendingTick lpName failed")
		}
		odfiLPMetaKey := memory_const.LPMetaDBPath(*lpName)
		batchKV[odfiLPMetaKey] = *odfiLPJson
	}
	err = db.StoreKeyValues(batchKV)
	return err
}

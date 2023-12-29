package memory_write

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"errors"
	"fmt"
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
		1. user spendingTick : - consumingAmt         (double-write)
		2. user buyingTick   : + deltaY               (double-write)
		3. odfi-spendingTick : + odfiTakerFee         (if LPMeta exist for odfi-spendingTick)
		4. lp   spendingTick : + deltaX + lpTakerFee  (LPMeta)
		5. lp   buyingTick   : - deltaY               (LPMeta)
		-------------------------------------------------------------------------------------
		3 & 4 won't be the same LP (protected by the callee `performSwap 2.1`):
			* if odfi-spending LP is equal to current lp, add odfiTakerFee to lpTakerFee, and set odfiTakerFee to zero.
	*/
	var batchKV map[string]string
	batchKV = make(map[string]string)
	// 1. user spendingTick : - consumingAmt
	{ // stack param isolation
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
	}
	// 2. user buyingTick   : + deltaY
	{ // stack param isolation
		buyingAvailable, err := memory_read.AvailableBalance(db, buyingTick, address)
		if err != nil {
			return err
		}
		updatedBuyingAvailable := buyingAvailable.Add(deltaY)
		if updatedBuyingAvailable == nil {
			return errors.New("WriteSwapInfo failed: calculate updatedBuyingAvailable failed")
		}
		buyingDoubleWriteKV := CoinBalanceDoubleWriteKV(buyingTick, address, updatedBuyingAvailable.String(), db_utils.AvailableSubAccount)
		for k, v := range buyingDoubleWriteKV {
			batchKV[k] = v
		}
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
		odfiLPName := memory_const.LPNameByTicks("odfi", spendingTick)
		if odfiLPName == nil {
			return errors.New("WriteSwapInfo failed: calculate odfiLPName failed")
		}
		odfiLPMetaKey := memory_const.LPMetaDBPath(*odfiLPName)
		batchKV[odfiLPMetaKey] = *odfiLPJson
	}
	// 4. lp   spendingTick : + deltaX + lpTakerFee  (LPMeta)
	// 5. lp   buyingTick   : - deltaY               (LPMeta)
	{ // stack param isolation
		lpMeta, err := memory_read.LiquidityProviderMetadata(db, buyingTick, spendingTick)
		if err != nil {
			return err
		}
		if spendingTick == lpMeta.LTick {
			lpMeta.LAmt = lpMeta.LAmt.Add(deltaX)
			if lpMeta.LAmt == nil {
				return fmt.Errorf("WriteSwapInfo failed: calculate lpMeta.LAmt error (+deltaX %s)", deltaX.String())
			}
			lpMeta.LAmt = lpMeta.LAmt.Add(lpTakerFee)
			if lpMeta.LAmt == nil {
				return fmt.Errorf("WriteSwapInfo failed: calculate lpMeta.LAmt error (+lpTakerFee %s)", lpTakerFee.String())
			}
			lpMeta.RAmt = lpMeta.RAmt.Subtract(deltaY)
			if lpMeta.RAmt == nil {
				return fmt.Errorf("WriteSwapInfo failed: calculate lpMeta.RAmt error (-deltaY %s)", deltaY.String())
			}
		} else if spendingTick == lpMeta.RTick {
			lpMeta.RAmt = lpMeta.RAmt.Add(deltaX)
			if lpMeta.RAmt == nil {
				return fmt.Errorf("WriteSwapInfo failed: calculate lpMeta.RAmt error (+deltaX %s)", deltaX.String())
			}
			lpMeta.RAmt = lpMeta.RAmt.Add(lpTakerFee)
			if lpMeta.RAmt == nil {
				return fmt.Errorf("WriteSwapInfo failed: calculate lpMeta.RAmt error (+lpTakerFee %s)", lpTakerFee.String())
			}
			lpMeta.LAmt = lpMeta.LAmt.Subtract(deltaY)
			if lpMeta.LAmt == nil {
				return fmt.Errorf("WriteSwapInfo failed: calculate lpMeta.LAmt error (-deltaY %s)", deltaY.String())
			}

		} else {
			return errors.New("WriteSwapInfo failed: spending tick error")
		}
		lpJson, err := lpMeta.JsonString()
		if err != nil {
			return err
		}
		if lpJson == nil {
			return errors.New("WriteSwapInfo failed: generate lpJson failed")
		}
		lpName := memory_const.LPNameByTicks(buyingTick, spendingTick)
		if lpName == nil {
			return errors.New("WriteSwapInfo failed: calculate lpName failed")
		}
		odfiLPMetaKey := memory_const.LPMetaDBPath(*lpName)
		batchKV[odfiLPMetaKey] = *lpJson
	}
	err := db.StoreKeyValues(batchKV)
	return err
}

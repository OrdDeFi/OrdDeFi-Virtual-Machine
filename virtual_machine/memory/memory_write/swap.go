package memory_write

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
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
		user spendingTick : - consumingAmt         (double-write)
		lp   spendingTick : + deltaX + lpTakerFee  (LPMeta)
		odfi-spendingTick : + odfiTakerFee         (if LPMeta exist for odfi-spendingTick)
		----------------------------------------------------------------------------------
		user buyingTick   : + deltaY               (double-write)
		lp   buyingTick   : - deltaY               (LPMeta)
	*/
	return nil
}

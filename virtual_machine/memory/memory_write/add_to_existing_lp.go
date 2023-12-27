package memory_write

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"errors"
)

func WriteAddToExistingLPInfo(
	db *db_utils.OrdDB,
	lTick string,
	rTick string,
	consumingLAmt *safe_number.SafeNum,
	consumingRAmt *safe_number.SafeNum,
	addingLPAmount *safe_number.SafeNum,
	lpMeta *memory_const.LPMeta,
	address string,
) error {
	lpName := lTick + "-" + rTick
	// 1. LP token add to user's wallet
	currentLPAmount, err := memory_read.LiquidityProviderBalance(db, lTick, rTick, address)
	if err != nil {
		return err
	}
	updatedLPAmount := currentLPAmount.Add(addingLPAmount)
	if updatedLPAmount == nil {
		return errors.New("WriteAddToExistingLPInfo calculate updatedLPAmount error")
	}
	batchKV := LPBalanceDoubleWriteKV(lTick, rTick, address, updatedLPAmount.String())
	// 2. Remove users token (left)
	leftAvailable, err := memory_read.AvailableBalance(db, lTick, address)
	if err != nil {
		return err
	}
	updatedLeftAvailable := leftAvailable.Subtract(consumingLAmt)
	if updatedLeftAvailable == nil {
		return errors.New("create LP error: " + lTick + " not enough: " + leftAvailable.String() + "-" + consumingLAmt.String())
	}
	leftCoinBatchKV := CoinBalanceDoubleWriteKV(lTick, address, updatedLeftAvailable.String(), db_utils.AvailableSubAccount)
	for k, v := range leftCoinBatchKV {
		batchKV[k] = v
	}
	// 3. Remove users token (right)
	rightAvailable, err := memory_read.AvailableBalance(db, rTick, address)
	if err != nil {
		return err
	}
	updatedRightAvailable := rightAvailable.Subtract(consumingRAmt)
	if updatedRightAvailable == nil {
		return errors.New("create LP error: " + rTick + " not enough: " + rightAvailable.String() + "-" + consumingRAmt.String())
	}
	rightCoinBatchKV := CoinBalanceDoubleWriteKV(rTick, address, updatedRightAvailable.String(), db_utils.AvailableSubAccount)
	for k, v := range rightCoinBatchKV {
		batchKV[k] = v
	}
	// 4. LP meta update
	if consumingLAmt == nil || consumingRAmt == nil {
		return errors.New("WriteAddToExistingLPInfo error: consumingLAmt or consumingRAmt is nil")
	}
	lpMeta.LAmt = lpMeta.LAmt.Add(consumingLAmt)
	lpMeta.RAmt = lpMeta.RAmt.Add(consumingRAmt)
	lpMeta.LTick = lTick
	lpMeta.RTick = rTick
	lpMeta.Total = lpMeta.Total.Add(addingLPAmount)
	if lpMeta.LAmt == nil {
		return errors.New("WriteAddToExistingLPInfo create LPMeta error: LAmt is nil")
	}
	if lpMeta.RAmt == nil {
		return errors.New("WriteAddToExistingLPInfo create LPMeta error: RAmt is nil")
	}
	if lpMeta.Total == nil {
		return errors.New("WriteAddToExistingLPInfo create LPMeta error: Total is nil")
	}
	lpMetaJsonString, err := lpMeta.JsonString()
	if err != nil {
		return errors.New("WriteAddToExistingLPInfo create LPMeta JSON string error")
	}
	lpMetaKey := memory_const.LpMetadataTable + ":" + lpName
	batchKV[lpMetaKey] = *lpMetaJsonString
	// write to DB
	err = db.StoreKeyValues(batchKV)
	return err
}

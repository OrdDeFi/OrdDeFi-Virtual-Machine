package memory_write

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"errors"
	"fmt"
)

func WriteRemoveLPInfo(
	db *db_utils.OrdDB,
	lTick string,
	rTick string,
	addingLAmt *safe_number.SafeNum,
	addingRAmt *safe_number.SafeNum,
	consumingLPAmount *safe_number.SafeNum,
	lpMeta *memory_const.LPMeta,
	address string,
) error {
	if addingLAmt == nil || addingRAmt == nil {
		return errors.New("WriteRemoveLPInfo error: addingLAmt or addingRAmt is nil")
	}
	if consumingLPAmount == nil {
		return errors.New("WriteRemoveLPInfo error: consumingLPAmount is nil")
	}
	lpName := lTick + "-" + rTick
	// 1. LP token add to user's wallet
	currentLPAmount, err := memory_read.LiquidityProviderBalance(db, lTick, rTick, address)
	if err != nil {
		return err
	}
	updatedLPAmount := currentLPAmount.Subtract(consumingLPAmount)
	if updatedLPAmount == nil {
		return fmt.Errorf("WriteRemoveLPInfo error: LP token is not enough: %s - %s", currentLPAmount.String(), consumingLPAmount.String())
	}
	batchKV := LPBalanceDoubleWriteKV(lTick, rTick, address, updatedLPAmount.String())
	// 2. Remove users token (left)
	leftAvailable, err := memory_read.AvailableBalance(db, lTick, address)
	if err != nil {
		return err
	}
	updatedLeftAvailable := leftAvailable.Add(addingLAmt)
	if updatedLeftAvailable == nil {
		return errors.New("WriteRemoveLPInfo calculate updatedLeftAvailable error")
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
	updatedRightAvailable := rightAvailable.Add(addingRAmt)
	if updatedRightAvailable == nil {
		return errors.New("WriteRemoveLPInfo calculate updatedRightAvailable error")
	}
	rightCoinBatchKV := CoinBalanceDoubleWriteKV(rTick, address, updatedRightAvailable.String(), db_utils.AvailableSubAccount)
	for k, v := range rightCoinBatchKV {
		batchKV[k] = v
	}
	// 4. LP meta update
	lpMeta.LAmt = lpMeta.LAmt.Subtract(addingLAmt)
	lpMeta.RAmt = lpMeta.RAmt.Subtract(addingRAmt)
	lpMeta.LTick = lTick
	lpMeta.RTick = rTick
	lpMeta.Total = lpMeta.Total.Subtract(consumingLPAmount)
	if lpMeta.LAmt == nil {
		return errors.New("WriteRemoveLPInfo calculate LPMeta error: LAmt is nil")
	}
	if lpMeta.RAmt == nil {
		return errors.New("WriteRemoveLPInfo calculate LPMeta error: RAmt is nil")
	}
	if lpMeta.Total == nil {
		return errors.New("WriteRemoveLPInfo calculate LPMeta error: Total is nil")
	}
	lpMetaJsonString, err := lpMeta.JsonString()
	if err != nil {
		return errors.New("WriteRemoveLPInfo create LPMeta JSON string error")
	}
	lpMetaKey := memory_const.LpMetadataTable + ":" + lpName
	batchKV[lpMetaKey] = *lpMetaJsonString
	// write to DB
	err = db.StoreKeyValues(batchKV)
	return err
}

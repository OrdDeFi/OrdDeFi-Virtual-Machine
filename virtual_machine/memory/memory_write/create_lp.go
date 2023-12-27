package memory_write

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"encoding/json"
	"errors"
)

func WriteCreateLPInfo(
	db *db_utils.OrdDB,
	lTick string,
	rTick string,
	lAmt *safe_number.SafeNum,
	rAmt *safe_number.SafeNum,
	address string,
) error {
	// 1. LP list update
	lpName := lTick + "-" + rTick
	allLPs, err := memory_read.AllLiquidityProviders(db)
	if err != nil {
		return err
	}
	allLPs = append(allLPs, lpName)
	allLPsJsonData, err := json.Marshal(allLPs)
	if err != nil {
		return err
	}
	allLPsValue := string(allLPsJsonData)
	// 2. LP meta update
	lpMeta := new(memory_const.LPMeta)
	lpMeta.LAmt = lAmt
	lpMeta.RAmt = rAmt
	lpMeta.LTick = lTick
	lpMeta.RTick = rTick
	lpMeta.Total = safe_number.SafeNumFromString("1000")
	lpMetaJsonString, err := lpMeta.JsonString()
	if err != nil {
		return errors.New("WriteCreateLPInfo create LPMeta JSON string error")
	}
	// 3. LP token add to user's wallet
	batchKV := LPBalanceDoubleWriteKV(lTick, rTick, address, "1000")
	// 4. Remove users token (left)
	leftAvailable, err := memory_read.AvailableBalance(db, lTick, address)
	if err != nil {
		return err
	}
	updatedLeftAvailable := leftAvailable.Subtract(lAmt)
	if updatedLeftAvailable == nil {
		return errors.New("create LP error: " + lTick + " not enough: " + leftAvailable.String() + "-" + lAmt.String())
	}
	leftCoinBatchKV := CoinBalanceDoubleWriteKV(lTick, address, updatedLeftAvailable.String(), db_utils.AvailableSubAccount)
	for k, v := range leftCoinBatchKV {
		batchKV[k] = v
	}
	// 5. Remove users token (right)
	rightAvailable, err := memory_read.AvailableBalance(db, rTick, address)
	if err != nil {
		return err
	}
	updatedRightAvailable := rightAvailable.Subtract(rAmt)
	if updatedRightAvailable == nil {
		return errors.New("create LP error: " + rTick + " not enough: " + rightAvailable.String() + "-" + rAmt.String())
	}
	rightCoinBatchKV := CoinBalanceDoubleWriteKV(rTick, address, updatedRightAvailable.String(), db_utils.AvailableSubAccount)
	for k, v := range rightCoinBatchKV {
		batchKV[k] = v
	}
	// 6. Combine KV
	allLPsKey := memory_const.LpListTable
	batchKV[allLPsKey] = allLPsValue
	lpMetaKey := memory_const.LpMetadataTable + ":" + lpName
	batchKV[lpMetaKey] = *lpMetaJsonString
	err = db.StoreKeyValues(batchKV)
	return err
}

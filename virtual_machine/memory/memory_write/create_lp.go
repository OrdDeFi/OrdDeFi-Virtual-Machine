package memory_write

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"encoding/json"
)

func WriteCreateLPInfo(
	db *db_utils.OrdDB,
	lTick string,
	rTick string,
	lAmt *safe_number.SafeNum,
	rAmt *safe_number.SafeNum,
	address string,

) error {
	// LP list update
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
	// LP meta update
	lpMeta := new(memory_const.LPMeta)
	lpMeta.LAmt = lAmt
	lpMeta.RAmt = rAmt
	lpMeta.LTick = lTick
	lpMeta.RTick = rTick
	lpMetaJsonString, err := lpMeta.JsonString()
	if err != nil {
		return nil
	}
	batchKV := LPBalanceDoubleWriteKV(lTick, rTick, address, "1000")
	allLPsKey := memory_const.LpListTable
	batchKV[allLPsKey] = allLPsValue
	lpMetaKey := memory_const.LpMetadataTable + ":" + lpName
	batchKV[lpMetaKey] = *lpMetaJsonString
	err = db.StoreKeyValues(batchKV)
	return err
}

package memory_write

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"encoding/json"
	"errors"
)

func WriteDeployInfo(
	db *db_utils.OrdDB,
	coinName string,
	max *safe_number.SafeNum,
	lim *safe_number.SafeNum,
	addrLim *safe_number.SafeNum,
	desc string,
	icon string) error {
	coinMeta := new(memory_const.CoinMeta)
	coinMeta.Max = max
	coinMeta.Lim = lim
	coinMeta.AddrLim = addrLim
	coinMeta.Desc = desc
	coinMeta.Icon = icon
	// Generate coin meta key-value
	coinMetaJsonString, err := coinMeta.JsonString()
	if err != nil {
		return err
	}
	if coinMetaJsonString == nil {
		return errors.New("WriteDeployInfo parse coin meta JSON error: result string is empty")
	}
	coinMetaKey := memory_const.CoinMetadataTable + ":" + coinName
	// Generate coin list key-value
	allCoins, err := memory_read.AllDeployedCoins(db)
	if err != nil {
		return err
	}
	allCoins = append(allCoins, coinName)
	allCoinsJsonData, err := json.Marshal(allCoins)
	if err != nil {
		return err
	}
	allCoinsValue := string(allCoinsJsonData)
	allCoinsKey := memory_const.CoinListTable
	// Generate batch writing map
	var batchWriting map[string]string
	batchWriting = make(map[string]string)
	batchWriting[coinMetaKey] = *coinMetaJsonString
	batchWriting[allCoinsKey] = allCoinsValue
	err = db.StoreKeyValues(batchWriting)
	return err
}

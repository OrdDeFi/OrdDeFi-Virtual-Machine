package db_utils

import (
	"errors"
	"strconv"
)

const lastUpdateBlockKey = "lastUpdatedBlock"

/*
GetLastUpdatedBlock
return:
1. nil, nil if no block updated
2. number, nil if any blockNumber updated
3. nil, err if read error or convert string to int error
*/
func GetLastUpdatedBlock(controlDB *OrdDB) (*int, error) {
	r, err := controlDB.Read(lastUpdateBlockKey)
	if err != nil {
		if err.Error() == "leveldb: not found" {
			return nil, nil
		} else {
			return nil, err
		}
	}
	if r == nil {
		return nil, nil
	}
	result, err := strconv.Atoi(*r)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

/*
GetUpdatedBlockHash
return:
1. nil, nil if blockNumber not exist
2. hash, nil if blockNumber exist
3. nil, err if read error
*/
func GetUpdatedBlockHash(controlDB *OrdDB, blockNumber int) (*string, error) {
	r, err := controlDB.Read(strconv.Itoa(blockNumber))
	if err != nil {
		if err.Error() == "leveldb: not found" {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return r, nil
}

func SetLastUpdatedBlock(controlDB *OrdDB, blockNumber int, blockHash string) error {
	if blockHash == "" {
		return errors.New("SetLastUpdatedBlock error: block hash is empty")
	}
	blockNumberStr := strconv.Itoa(blockNumber)
	var batchKV map[string]string
	batchKV = make(map[string]string)
	batchKV[lastUpdateBlockKey] = blockNumberStr
	batchKV[blockNumberStr] = blockHash
	err := controlDB.StoreKeyValues(batchKV)
	return err
}

func ResetLastUpdatedBlockTo(controlDB *OrdDB, toBlockNumber *int, lastUpdatedBlockNumber *int) error {
	var batchKV map[string]string
	batchKV = make(map[string]string)
	if toBlockNumber == nil || lastUpdatedBlockNumber == nil {
		// remove all
		batchKV[lastUpdateBlockKey] = ""
	} else {
		// remove from toBlockNumber to lastUpdatedBlockNumber
		batchKV[lastUpdateBlockKey] = strconv.Itoa(*toBlockNumber)
		for i := *toBlockNumber + 1; i <= *lastUpdatedBlockNumber; i++ {
			blockNumberStr := strconv.Itoa(i)
			batchKV[blockNumberStr] = ""
		}
	}
	err := controlDB.StoreKeyValues(batchKV)
	return err
}

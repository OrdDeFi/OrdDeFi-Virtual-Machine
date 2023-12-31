package test

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_const"
	"fmt"
	"strconv"
	"testing"
)

func TestLPPath(t *testing.T) {
	expectedPath := "lpaddrbalance:v1:odfi-odgv:bc1qm34lsc65zpw79lxes69zkqmk6ee3ewf0j77s3h"
	path := memory_const.LPAddressPath("odfi", "odgv", "bc1qm34lsc65zpw79lxes69zkqmk6ee3ewf0j77s3h")
	if path != expectedPath {
		t.Errorf("DB path error, got %s, expected %s", path, expectedPath)
	}

	expectedPath = "addrlpbalance:v1:bc1qm34lsc65zpw79lxes69zkqmk6ee3ewf0j77s3h:odfi-odgv"
	path = memory_const.AddressLPPath("odfi", "odgv", "bc1qm34lsc65zpw79lxes69zkqmk6ee3ewf0j77s3h")
	if path != expectedPath {
		t.Errorf("DB path error, got %s, expected %s", path, expectedPath)
	}
}

func TestDB(t *testing.T) {
	db, err := db_utils.OpenDB(testDBPath)
	if err != nil {
		t.Errorf("TestDB OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")

	key := "test:coin_list"
	valueString := `["odfi","odgv"]`
	err = db.Store(key, valueString)
	if err != nil {
		t.Errorf("TestDB error: %s", err.Error())
	}
	resultString, err := db.Read(key)
	if err != nil {
		t.Errorf("TestDB error: %s", err.Error())
	}
	if *resultString != valueString {
		t.Errorf("TestDB error: result %s, expected %s", *resultString, valueString)
	}
}

func TestDBRead(t *testing.T) {
	db, err := db_utils.OpenDB(testDBPath)
	if err != nil {
		t.Errorf("TestDBRead OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")

	key := "test:coin_list"
	valueString := `["odfi","odgv"]`

	resultString, err := db.Read(key)
	if err != nil {
		t.Errorf("TestDB error: %s", err.Error())
	}
	if *resultString != valueString {
		t.Errorf("TestDB error: result %s, expected %s", *resultString, valueString)
	}
}

func TestDBReadPrefix(t *testing.T) {
	db, err := db_utils.OpenDB(testDBPath)
	if err != nil {
		t.Errorf("TestDBReadPrefix OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")
	prefix := "test:coin:odfi:"
	for i := 0; i < 20; i++ {
		key := prefix + strconv.Itoa(i)
		value := key + ":value"
		err = db.Store(key, value)
		if err != nil {
			t.Errorf("TestDBReadPrefix Store error: %s", err.Error())
		}
	}
	res, err := db.ReadAllPrefix(prefix)
	if err != nil {
		t.Errorf("TestDBReadPrefix ReadAllPrefix error: %s", err.Error())
	}
	for key, valueString := range res {
		if valueString != key+":value" {
			t.Errorf("TestDBReadPrefix ReadAllPrefix value: %s, expected %s", valueString, key+":value")
		}
	}
}

func TestDBBatchStore(t *testing.T) {
	db, err := db_utils.OpenDB(testDBPath)
	if err != nil {
		t.Errorf("TestDBBatchStore OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")

	var keyValues map[string]string
	keyValues = make(map[string]string)
	for i := 0; i < 20; i++ {
		key := "test:key" + strconv.Itoa(i)
		value := key + ":value"
		keyValues[key] = value
	}
	err = db.StoreKeyValues(keyValues)
	if err != nil {
		t.Errorf("TestDBBatchStore StoreKeyValues error: %s", err.Error())
	}
	for i := 0; i < 20; i++ {
		key := "test:key" + strconv.Itoa(i)
		expectedValue := key + ":value"
		value, err := db.Read(key)
		if err != nil {
			t.Errorf("TestDBBatchStore Read error: %s", err.Error())
		}
		if value == nil {
			t.Errorf("TestDBBatchStore Read error: value empty")
		}
		if expectedValue != *value {
			t.Errorf("TestDBBatchStore value error: %s, expected %s", *value, expectedValue)
		}
	}
}

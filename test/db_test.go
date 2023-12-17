package test

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"fmt"
	"strconv"
	"testing"
)

func TestDB(t *testing.T) {
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestDB OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")

	key := "coin_list"
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
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestDBRead OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")

	key := "coin_list"
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
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestDBReadPrefix OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")
	prefix := "coin:odfi:"
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
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestDBBatchStore OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")

	var keyValues map[string]string
	keyValues = make(map[string]string)
	for i := 0; i < 20; i++ {
		key := "key" + strconv.Itoa(i)
		value := key + ":value"
		keyValues[key] = value
	}
	err = db.StoreKeyValues(keyValues)
	if err != nil {
		t.Errorf("TestDBBatchStore StoreKeyValues error: %s", err.Error())
	}
	for i := 0; i < 20; i++ {
		key := "key" + strconv.Itoa(i)
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

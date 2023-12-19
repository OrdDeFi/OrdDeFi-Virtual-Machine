package test

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"OrdDeFi-Virtual-Machine/virtual_machine/operations"
	"fmt"
	"testing"
)

func testForParam(t *testing.T, db *db_utils.OrdDB, tick string, ver string, txId string) {
	// 1. compile instruction
	instruction, err := TestingMintInSingleSliceCommands(tick, ver, txId)
	if err != nil {
		t.Errorf("TestExecuteMint error: %s", err.Error())
	}
	if instruction == nil {
		t.Errorf("TestExecuteMint error: deploy instruction is nil")
	}

	// 2. execute deploy op
	err = operations.ExecuteOpMint(*instruction, db)
	if err != nil {
		t.Errorf("TestExecuteMint error: execute deploy error %s", err)
	}
}

func TestExecuteMint(t *testing.T) {
	// open db
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestExecuteMint OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")

	for _, txId := range TestingTxPool() {
		testForParam(t, db, "odfi", "", txId)
	}
}

func TestExecuteMintVer1(t *testing.T) {
	// open db
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestExecuteMint OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")

	for _, txId := range TestingTxPool() {
		testForParam(t, db, "odgv", "1", txId)
	}
}

func TestExecuteMintVer2(t *testing.T) {
	// open db
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestExecuteMint OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")

	for _, txId := range TestingTxPool() {
		testForParam(t, db, "odfi", "2", txId)
	}
}

func testReadCoin(t *testing.T, coinName string) {
	// 1. open db
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestExecuteMintVer2 OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")
	// 2. read balance
	num, err := memory_read.Balance(db, "odfi", "bc1pq89nvjf7fd0kkyu8z825vyg48gupgmf9ngm5g9zk3hp8cyltd9nqr0fhj5", "1")
	if err != nil {
		t.Errorf("TestReadBalance error: execute deploy error %s", err)
	}
	if num == nil {
		t.Errorf("TestReadBalance error: num should not be nil")
	}
	// 3. read coin balance
	coinRes, err := memory_read.AllAddressBalanceForCoin(db, coinName, "")
	if err != nil {
		t.Errorf("TestReadBalance AllAddressBalanceForCoin error: %s", err.Error())
	}
	for k, v := range coinRes {
		println(k, v)
	}
}

func TestReadODFIBalance(t *testing.T) {
	testReadCoin(t, "odfi")
}

func TestReadODGVBalance(t *testing.T) {
	testReadCoin(t, "odgv")
}

func testReadAddress(t *testing.T, address string) {
	// 1. open db
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestExecuteMintVer2 OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")

	// 2. read address balance
	coinRes, err := memory_read.AllCoinBalanceForAddress(db, address, "")
	if err != nil {
		t.Errorf("TestReadBalance AllAddressBalanceForCoin error: %s", err.Error())
	}
	for k, v := range coinRes {
		println(k, v)
	}
}

func TestReadAddress1(t *testing.T) {
	testReadAddress(t, "bc1qsl5psn0f6kk9rep0gejmnj9zgnts7gr40sv5xu")
}

func TestReadAddress2(t *testing.T) {
	testReadAddress(t, "bc1qvp0m6efzkawm5ywnymket459dhtx326f64en6x")
}

func TestRemainingToMint(t *testing.T) {
	// open db
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestExecuteMintVer2 OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")
	r, err := memory_read.TotalMintedBalance(db, "odfi")
	if err != nil {
		t.Errorf("TestReadBalance TotalMintedBalance error: %s", err.Error())
	}
	println("Minted", r.String())
}

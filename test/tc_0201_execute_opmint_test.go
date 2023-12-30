package test

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"OrdDeFi-Virtual-Machine/virtual_machine/operations"
	"fmt"
	"strings"
	"testing"
)

func TestingMintForParam(t *testing.T, db *db_utils.OrdDB, tick string, txId string, amt string) {
	// 1. compile instruction
	instruction, err := TestingMintInSingleSliceCommands(tick, txId, amt)
	if err != nil {
		t.Errorf("TestExecuteMint error: %s", err.Error())
	}
	if instruction == nil {
		t.Errorf("TestExecuteMint error: deploy instruction is nil")
	}

	// 2. execute deploy op
	err = operations.ExecuteOpMint(*instruction, db)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Mint ended for ") || strings.HasPrefix(err.Error(), "Address reached limit for ") {
		} else {
			if strings.HasPrefix(err.Error(), "repeat mint disabled") == false {
				t.Errorf("TestExecuteMint error: execute OpMint error %s", err)
			}
		}
	}
}

func TestExecuteMintODFI(t *testing.T) {
	// open db
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestExecuteMint OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")
	coinName := "odfi"

	for _, txId := range TestingTxPool() {
		TestingMintForParam(t, db, coinName, txId, "1000")
	}
	for _, address := range validTestingAddressPool() {
		balance, _ := memory_read.AvailableBalance(db, coinName, address)
		if !balance.IsEqualTo(safe_number.SafeNumFromString("1000")) {
			t.Errorf("%s minted balance should be 1000, now is %s", address, balance.String())
		}
	}
}

func TestExecuteMintODGV(t *testing.T) {
	// open db
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestExecuteMint OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")

	for _, txId := range TestingTxPool() {
		TestingMintForParam(t, db, "odgv", txId, "1000")
	}
}

func TestExecuteMintHALF(t *testing.T) {
	// open db
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestExecuteMint OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")

	for _, txId := range TestingTxPool() {
		TestingMintForParam(t, db, "half", txId, "1000")
	}
}

func TestingReadCoin(t *testing.T, coinName string) {
	// 1. open db
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestExecuteMintVer2 OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")
	// 2. read balance
	available, transferable, err := memory_read.Balance(db, "odfi", "bc1pq89nvjf7fd0kkyu8z825vyg48gupgmf9ngm5g9zk3hp8cyltd9nqr0fhj5")
	if err != nil {
		t.Errorf("TestReadBalance error: execute deploy error %s", err)
	}
	if available == nil {
		t.Errorf("TestReadBalance error: num should not be nil")
	}
	if transferable == nil {
		t.Errorf("TestReadBalance error: num should not be nil")
	}
	// 3. read coin balance
	coinRes, err := memory_read.AllAddressBalanceForCoin(db, coinName)
	if err != nil {
		t.Errorf("TestReadBalance AllAddressBalanceForCoin error: %s", err.Error())
	}
	for k, v := range coinRes {
		println(k, v)
	}
}

func TestReadODFIBalance(t *testing.T) {
	TestingReadCoin(t, "odfi")
}

func TestReadODGVBalance(t *testing.T) {
	TestingReadCoin(t, "odgv")
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
	coinRes, err := memory_read.AllCoinBalanceForAddress(db, address)
	if err != nil {
		t.Errorf("TestReadBalance AllAddressBalanceForCoin error: %s", err.Error())
	}
	for k, v := range coinRes {
		println(k, v)
	}
}

func TestReadAddress1(t *testing.T) {
	testReadAddress(t, "39Vc3f9NsBoLPa2Qg7JFk34qhfQe8vRqaq")
}

func TestReadAddress2(t *testing.T) {
	testReadAddress(t, "bc1qrm3x757f5vp3skudtackrrac5qejp9h9xz403j")
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

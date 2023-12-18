package test

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/operations"
	"fmt"
	"testing"
)

func TestExecuteMint(t *testing.T) {
	// 1. compile instruction
	instruction, err := TestingMintInSingleSliceCommands("odfi", "")
	if err != nil {
		t.Errorf("TestExecuteMint error: %s", err.Error())
	}
	if instruction == nil {
		t.Errorf("TestExecuteMint error: deploy instruction is nil")
	}

	// 2. open db
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestExecuteMint OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")

	// 3. execute deploy op
	err = operations.ExecuteOpMint(*instruction, db)
	if err != nil {
		t.Errorf("TestExecuteMint error: execute deploy error %s", err)
	}
}

func TestExecuteMintVer1(t *testing.T) {
	// 1. compile instruction
	instruction, err := TestingMintInSingleSliceCommands("odfi", "1")
	if err != nil {
		t.Errorf("TestExecuteMintVer1 error: %s", err.Error())
	}
	if instruction == nil {
		t.Errorf("TestExecuteMintVer1 error: deploy instruction is nil")
	}

	// 2. open db
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestExecuteMintVer1 OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")

	// 3. execute deploy op
	err = operations.ExecuteOpMint(*instruction, db)
	if err != nil {
		t.Errorf("TestExecuteMintVer1 error: execute deploy error %s", err)
	}
}

func TestExecuteMintVer2(t *testing.T) {
	// 1. compile instruction
	instruction, err := TestingMintInSingleSliceCommands("odfi", "2")
	if err != nil {
		t.Errorf("TestExecuteMintVer2 error: %s", err.Error())
	}
	if instruction == nil {
		t.Errorf("TestExecuteMintVer2 error: deploy instruction is nil")
	}

	// 2. open db
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestExecuteMintVer2 OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")

	// 3. execute deploy op
	err = operations.ExecuteOpMint(*instruction, db)
	if err != nil {
		t.Errorf("TestExecuteMintVer2 error: execute deploy error %s", err)
	}
}

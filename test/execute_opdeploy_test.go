package test

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/operations"
	"fmt"
	"testing"
)

func TestExecuteMintInBatchCommands(t *testing.T) {
	instruction, err := TestingDeployInSingleSliceCommands()
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if instruction == nil {
		t.Errorf("TestExecuteMintInBatchCommands error: deploy instruction is nil")
	}

	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestDBReadPrefix OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")

	operations.ExecuteOpDeploy(*instruction, db)
}

package test

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"fmt"
	"testing"
)

func TestLPMeta(t *testing.T) {
	// open db
	db, err := db_utils.OpenDB(testDBPath)
	if err != nil {
		t.Errorf("TestExecuteMint OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")

	meta, err := memory_read.LiquidityProviderMetadata(db, "noth", "ing.")
	if err != nil {
		t.Errorf("TestLPMeta error: %s", err.Error())
	}
	if meta != nil {
		t.Errorf("LP meta is not nil, expected nil")
	}
}

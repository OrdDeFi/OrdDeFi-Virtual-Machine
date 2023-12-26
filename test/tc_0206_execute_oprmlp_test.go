package test

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"testing"
)

func TestRemoveLiquidityPair(t *testing.T) {
	// open db
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestExecuteMint OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
}

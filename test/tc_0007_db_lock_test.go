package test

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"testing"
)

const testLockDBPath = "test_lock_db"

func checkDBLockState(t *testing.T, db *db_utils.OrdDB) {
	state, err := db_utils.CheckControlDBLockState(db)
	if err != nil {
		t.Errorf("CheckControlDBLockState error: %s", err.Error())
	}
	if state == nil {
		t.Errorf("CheckControlDBLockState error: state is nil")
	}
	println("current lock state:", *state)
}

func TestDBLock(t *testing.T) {
	// open db
	db, err := db_utils.OpenDB(testLockDBPath)
	if err != nil {
		t.Errorf("TestDBLock OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)

	checkDBLockState(t, db)
	err = db_utils.LockControlDB(db)
	if err != nil {
		t.Errorf("TestDBLock LockControlDB error: %s", err.Error())
	}
	checkDBLockState(t, db)
	err = db_utils.ReleaseLockControlDB(db)
	if err != nil {
		t.Errorf("TestDBLock LockControlDB error: %s", err.Error())
	}
	checkDBLockState(t, db)
}

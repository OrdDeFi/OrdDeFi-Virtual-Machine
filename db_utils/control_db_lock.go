package db_utils

func LockControlDB(controlDB *OrdDB) error {
	err := controlDB.Store("lock", "true")
	return err
}

func ReleaseLockControlDB(controlDB *OrdDB) error {
	err := controlDB.Store("lock", "false")
	return err
}

func CheckControlDBLockState(controlDB *OrdDB) (*bool, error) {
	value, err := controlDB.Read("lock")
	if err != nil {
		if err.Error() == "leveldb: not found" {
			result := false
			return &result, nil
		} else {
			return nil, err
		}
	}
	if value == nil {
		result := false
		return &result, nil
	}
	if *value == "false" {
		result := false
		return &result, nil
	}
	result := true
	return &result, nil
}

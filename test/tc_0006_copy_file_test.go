package test

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/file_utils"
	"strings"
	"testing"
)

func TestCopyFile(t *testing.T) {
	err := file_utils.CopyFile("/Users/satoshi/large.file", "/Users/satoshi/large_copy.file")
	if err != nil {
		if strings.HasSuffix(err.Error(), "no such file or directory") == false {
			t.Errorf("Copy file error: %s", err.Error())
		}
	}
}

func TestCopyDir(t *testing.T) {
	err := file_utils.CopyDir("/Users/satoshi/dir", "/Users/satoshi/dir_copy")
	if err != nil {
		t.Errorf("Copy dir error: %s", err.Error())
	}
}

func TestRemoveDir(t *testing.T) {
	err := file_utils.RemoveDir("/Users/satoshi/B")
	if err != nil {
		t.Errorf("Remove dir error: %s", err.Error())
	}
}

func TestBackupPath(t *testing.T) {
	backupPath := db_utils.BackupPathForMainPath("/Users/satoshi/OrdDeFi_storage", 825100)
	println(backupPath)
	if backupPath != "/Users/satoshi/OrdDeFi_storage_backup_825100" {
		t.Errorf("TestBackupPath error: expected %s, got %s", "/Users/satoshi/OrdDeFi_storage_backup_825100", backupPath)
	}

	backupPath = db_utils.BackupPathForMainPath("/Users/satoshi/OrdDeFi_storage/", 825100)
	println(backupPath)
	if backupPath != "/Users/satoshi/OrdDeFi_storage_backup_825100" {
		t.Errorf("TestBackupPath error: expected %s, got %s", "/Users/satoshi/OrdDeFi_storage_backup_825100", backupPath)
	}
}

func TestRestoreNumber(t *testing.T) {
	r1 := db_utils.RestoringBlockNumber(825050)
	if r1 != 825000 {
		t.Errorf("TestRestoreNumber failed, expected %d, got %d", 825000, r1)
	}
	r2 := db_utils.RestoringBlockNumber(825049)
	if r2 != 825000 {
		t.Errorf("TestRestoreNumber failed, expected %d, got %d", 825000, r2)
	}
}

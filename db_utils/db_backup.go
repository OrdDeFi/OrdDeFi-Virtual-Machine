package db_utils

import (
	"OrdDeFi-Virtual-Machine/file_utils"
	"path/filepath"
	"strconv"
)

const BackupAlignment = 50

func BackupPathForMainPath(mainDBPath string, blockNumber int) string {
	backupSuffix := "_backup_" + strconv.Itoa(blockNumber)
	cleanPath := filepath.Clean(mainDBPath)
	parentDir := filepath.Dir(cleanPath)
	backupPath := filepath.Join(parentDir, filepath.Base(mainDBPath)+backupSuffix)
	return backupPath
}

func Backup(mainDBPath string, blockNumber int) error {
	backupPath := BackupPathForMainPath(mainDBPath, blockNumber)
	err := file_utils.CopyDir(mainDBPath, backupPath)
	return err
}

func Restore(mainDBPath string, blockNumber int) error {
	backupPath := BackupPathForMainPath(mainDBPath, blockNumber)
	err := file_utils.CopyDir(backupPath, mainDBPath)
	return err
}

func RestoringBlockNumber(lastUpdatedBlockNumber int) int {
	result := (lastUpdatedBlockNumber / BackupAlignment) * BackupAlignment
	if result == lastUpdatedBlockNumber {
		result -= BackupAlignment
	}
	return result
}

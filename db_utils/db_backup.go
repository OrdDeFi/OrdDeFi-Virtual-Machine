package db_utils

import (
	"OrdDeFi-Virtual-Machine/file_utils"
	"path/filepath"
	"strconv"
)

const BackupAlignment = 50
const GenesisBlockNumber = 829832

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
	if err != nil {
		return err
	}
	prevBackupPath := BackupPathForMainPath(mainDBPath, blockNumber-2*BackupAlignment)
	err = file_utils.RemoveDir(prevBackupPath)
	return err
}

func Restore(mainDBPath string, restoreBlockNumber int) error {
	backupPath := BackupPathForMainPath(mainDBPath, restoreBlockNumber)
	err := file_utils.CopyDir(backupPath, mainDBPath)
	return err
}

func RestoringBlockNumber(lastUpdatedBlockNumber int, removeCurrentBlock bool) int {
	result := (lastUpdatedBlockNumber / BackupAlignment) * BackupAlignment
	if result == lastUpdatedBlockNumber && removeCurrentBlock {
		result -= BackupAlignment
	}
	return result
}

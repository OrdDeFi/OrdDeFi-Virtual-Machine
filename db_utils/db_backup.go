package db_utils

import (
	"path/filepath"
	"strconv"
)

func BackupPathForMainPath(mainDBPath string, blockNumber int) string {
	backupSuffix := "_backup_" + strconv.Itoa(blockNumber)
	cleanPath := filepath.Clean(mainDBPath)
	parentDir := filepath.Dir(cleanPath)
	backupPath := filepath.Join(parentDir, filepath.Base(mainDBPath)+backupSuffix)
	return backupPath
}

func Backup(mainDBPath string, blockNumber int) error {

	return nil
}

func Restore(mainDBPath string, blockNumber int) error {

	return nil
}

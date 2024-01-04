package updater

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/file_utils"
	"errors"
)

const controlDBPath = "./OrdDeFi_control"

func restoreDB(controlDB *db_utils.OrdDB, dataDir string, logDir string, lastUpdatedBlockNumber *int, discardCurrentBlock bool) error {
	if lastUpdatedBlockNumber == nil {
		// if lastUpdatedBlockNumber is nil, delete the current data immediately, and index from begin
		err := file_utils.RemoveDir(dataDir)
		if err != nil {
			return err
		}
		err = file_utils.RemoveDir(logDir)
		if err != nil {
			return err
		}
		err = db_utils.ResetLastUpdatedBlockTo(controlDB, nil, nil)
		return err
	} else {
		restoreBlockNumber := db_utils.RestoringBlockNumber(*lastUpdatedBlockNumber, discardCurrentBlock)
		err := db_utils.Restore(dataDir, restoreBlockNumber)
		if err != nil {
			return err
		}
		err = db_utils.Restore(logDir, restoreBlockNumber)
		if err != nil {
			return err
		}
		err = db_utils.ResetLastUpdatedBlockTo(controlDB, &restoreBlockNumber, lastUpdatedBlockNumber)
		return err
	}
}

func UpdateIndex(dataDir string, logDir string, verbose bool) error {
	println("The Times 03/Jan/2009 Chancellor on brink of second bailout for banks.")
	println("OrdDeFi indexer begin to update.")

	// 0.0 Check bitcoin-cli requirements
	reachedMinRequirement, err := bitcoin_cli_channel.VersionGreaterThanMinRequirement()
	if err != nil {
		return err
	}
	if *reachedMinRequirement == false {
		return errors.New("bitcoin-cli version lower than 24.0.1")
	}
	// 0.1 Open control db
	controlDB, err := db_utils.OpenDB(controlDBPath)
	if err != nil {
		return err
	}
	defer db_utils.CloseDB(controlDB)

	// 1.1 Check current block number
	currentBlockNumber := bitcoin_cli_channel.GetBlockCount()
	if currentBlockNumber == 0 {
		err := errors.New("updateIndex error: bitcoin-cli getblockcount failed")
		return err
	}
	// 1.2 Check last updated block number
	lastUpdatedBlock, err := db_utils.GetLastUpdatedBlock(controlDB)
	if err != nil {
		return err
	}

	// 2.1 Check control DB lock state. Locking status indicates the VM was exited by accident. Restore before running new instructions.
	lockState, err := db_utils.CheckControlDBLockState(controlDB)
	if err != nil {
		return err
	}
	if *lockState == true {
		// lastUpdatedBlock is *int, which could be nil.
		// If it is nil, indicates that no block was updated succeed. Remove dataDir and logDir.
		// Otherwise, restore dataDir and logDir
		err = restoreDB(controlDB, dataDir, logDir, lastUpdatedBlock, false)
		if err != nil {
			return err
		}
		// refresh lastUpdatedBlock
		lastUpdatedBlock, err = db_utils.GetLastUpdatedBlock(controlDB)
		if err != nil {
			return err
		}
	}
	// 2.2 Check block hash of lastUpdatedBlock
	if lastUpdatedBlock != nil {
		storedBlockHash, err := db_utils.GetUpdatedBlockHash(controlDB, *lastUpdatedBlock)
		if err != nil {
			return err
		}
		if storedBlockHash != nil {
			bitcoinCliBlockHash := bitcoin_cli_channel.GetBlockHash(*lastUpdatedBlock)
			if bitcoinCliBlockHash == nil {
				return errors.New("UpdateIndex GetBlockHash failed: bitcoinCliBlockHash is nil")
			}
			if *storedBlockHash != *bitcoinCliBlockHash {
				err = restoreDB(controlDB, dataDir, logDir, lastUpdatedBlock, true)
				if err != nil {
					return err
				}
				// refresh lastUpdatedBlock
				lastUpdatedBlock, err = db_utils.GetLastUpdatedBlock(controlDB)
				if err != nil {
					return err
				}
			}
		}
	}

	// 3.0 Lock control DB
	err = db_utils.LockControlDB(controlDB)
	if err != nil {
		return err
	}
	beginBlockNumber := db_utils.GenesisBlockNumber
	if lastUpdatedBlock != nil {
		beginBlockNumber = *lastUpdatedBlock + 1
	}
	// 3.1 Running instructions
	for indexingBlockNumber := beginBlockNumber; indexingBlockNumber <= currentBlockNumber; indexingBlockNumber++ {
		// get block hash and all txIds in block
		println("indexing block", indexingBlockNumber)
		blockHash := bitcoin_cli_channel.GetBlockHash(indexingBlockNumber)
		if blockHash == nil {
			return errors.New("UpdateBlockNumber GetBlockHash failed")
		}
		err = UpdateBlockNumber(indexingBlockNumber, blockHash, dataDir, logDir, verbose)
		if err != nil {
			return err
		}
		err = db_utils.SetLastUpdatedBlock(controlDB, indexingBlockNumber, *blockHash)
		if err != nil {
			return err
		}
		if indexingBlockNumber%50 == 0 {
			err = db_utils.Backup(dataDir, indexingBlockNumber)
			if err != nil {
				return err
			}
			err = db_utils.Backup(logDir, indexingBlockNumber)
			if err != nil {
				return err
			}
		}
	}
	// 3.2 Release control lock
	err = db_utils.ReleaseLockControlDB(controlDB)
	println("OrdDeFi indexer update finished. Last block:", currentBlockNumber)
	return err
}

package updater

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDeFi-Virtual-Machine/concurrent"
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/inscription_parser"
	"OrdDeFi-Virtual-Machine/virtual_machine"
	"OrdDeFi-Virtual-Machine/virtual_machine/operations"
	"errors"
	"sync"
)

func prepareRawTransaction(blockTx []string) map[string]string {
	tasks := make(chan concurrent.Task, 50)
	var wg sync.WaitGroup
	storage := concurrent.NewResultStorage()

	concurrent.StartWorkerPool(tasks, 10, &wg, storage)
	for _, txId := range blockTx {
		txId := txId
		tasks <- func() (string, string) {
			rawTx := bitcoin_cli_channel.GetRawTransaction(txId)
			if rawTx == nil {
				return txId, ""
			}
			return txId, *rawTx
		}
	}
	close(tasks)
	wg.Wait()

	result := storage.Results
	if len(blockTx) != len(result) {
		println("prepareRawTransaction error: len(blockTx) != len(result), expected", len(blockTx), "got", len(result))
		return nil
	}
	return result
}

func UpdateBlockNumber(blockNumber int, blockHash *string, dataDir string, logDir string, verbose bool) error {
	var err error
	if blockHash == nil {
		return errors.New("UpdateBlockNumber failed: blockHash is nil")
	}
	block := bitcoin_cli_channel.GetBlock(*blockHash)
	if block == nil {
		err = errors.New("UpdateBlockNumber GetBlock failed")
		return err
	}
	// open data DB and log DB
	if dataDir == "" {
		dataDir = "./OrdDeFi_storage"
	}
	if logDir == "" {
		logDir = "./OrdDeFi_log"
	}
	if dataDir == logDir {
		return errors.New("-data-dir and -log-dir should be different")
	}
	db, err := db_utils.OpenDB(dataDir)
	if err != nil {
		return err
	}
	defer db_utils.CloseDB(db)
	logDB, err := db_utils.OpenDB(logDir)
	if err != nil {
		return err
	}
	defer db_utils.CloseDB(logDB)

	rawTransactions := prepareRawTransaction(block.Tx)
	// enum txId, execute operations if exist
	for txIndex, txId := range block.Tx {
		//rawTx := bitcoin_cli_channel.GetRawTransaction(txId)
		//if rawTx == nil {
		//	err = errors.New("GetRawTransaction Failed")
		//	break
		//}
		rawTxStr := rawTransactions[txId]
		if rawTxStr == "" {
			err = errors.New("GetRawTransaction Failed")
			break
		}
		tx := bitcoin_cli_channel.DecodeRawTransaction(rawTxStr)
		if tx == nil {
			err = errors.New("ParseRawTransaction -> DecodeRawTransaction Failed")
			break
		}
		utxoTransferApplied, err := operations.ApplyUTXOTransfer(db, tx, blockNumber)
		if err != nil {
			break
		}
		if utxoTransferApplied {
			// If UTXO transfer applied, stop executing any instruction in this tx to avoid security issue.
			continue
		}
		contentType, content, err := inscription_parser.ParseTransactionToInscription(*tx)
		if err != nil {
			break
		}
		if contentType != nil && content != nil {
			instructions, err := virtual_machine.CompileInstructions(*contentType, content, tx, txId)
			if err != nil {
				break
			}
			if len(instructions) != 0 {
				if verbose {
					println("executing op:", blockNumber, txIndex, txId)
				}
				virtual_machine.ExecuteInstructions(instructions, db, logDB, blockNumber, txIndex, txId, verbose)
			}
		}
	}
	if err != nil {
		println("Updating block got error:", err) // failing
		println("Aborting update blocks...")
		return err
	}
	return nil
}

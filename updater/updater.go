package updater

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/inscription_parser"
	"OrdDeFi-Virtual-Machine/virtual_machine"
	"OrdDeFi-Virtual-Machine/virtual_machine/operations"
	"errors"
)

func UpdateBlockNumber(blockNumber int, dataDir string, logDir string) error {
	var err error
	// get block hash and all txIds in block
	blockHash := bitcoin_cli_channel.GetBlockHash(blockNumber)
	if blockHash == nil {
		err = errors.New("UpdateBlockNumber GetBlockHash failed")
		return err
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
	// enum txId, execute operations if exist
	for txIndex, txId := range block.Tx {
		rawTx := bitcoin_cli_channel.GetRawTransaction(txId)
		if rawTx == nil {
			err = errors.New("GetRawTransaction Failed")
			break
		}
		tx := bitcoin_cli_channel.DecodeRawTransaction(*rawTx)
		if tx == nil {
			err = errors.New("ParseRawTransaction -> DecodeRawTransaction Failed")
			break
		}
		utxoTransferApplied, err := operations.ApplyUTXOTransfer(db, tx)
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
				virtual_machine.ExecuteInstructions(instructions, db, logDB, blockNumber, txIndex, txId)
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

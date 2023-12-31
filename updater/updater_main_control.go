package updater

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
	"errors"
)

func UpdateIndex(dataDir string, logDir string, verbose bool) error {
	println("The Times 03/Jan/2009 Chancellor on brink of second bailout for banks.")
	println("OrdDeFi indexer start to work.")
	
	blockNumber := bitcoin_cli_channel.GetBlockCount()
	if blockNumber == 0 {
		err := errors.New("updateIndex error: bitcoin-cli getblockcount failed")
		return err
	}
	err := UpdateBlockNumber(blockNumber, dataDir, logDir, verbose)
	return err
}

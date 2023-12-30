package subcommands

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
)

func ParseTransaction(txId string) {
	rawTx := bitcoin_cli_channel.GetRawTransaction(txId)
	if rawTx == nil {
		println("GetRawTransaction Failed")
		return
	}
	ParseRawTransaction(*rawTx)
}

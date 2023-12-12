package inscription_parser

import (
	"brc20defi_vm/bitcoin_cli_channel"
	"errors"
	"github.com/btcsuite/btcd/wire"
)

func parseTransactionToInscription(tx wire.MsgTx) (*string, error) {
	res := ""
	for _, eachWitNess := range tx.TxIn[0].Witness {
		res += string(eachWitNess)
	}
	return &res, nil
}

func ParseRawTransactionToInscription(rawTransaction string) (*string, error) {
	tx := bitcoin_cli_channel.DecodeRawTransaction(rawTransaction)
	if tx == nil {
		err := errors.New("ParseRawTransaction -> DecodeRawTransaction Failed")
		return nil, err
	}
	return parseTransactionToInscription(*tx)
}

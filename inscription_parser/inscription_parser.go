package inscription_parser

import (
	"brc20defi_vm/bitcoin_cli_channel"
	"errors"
	"github.com/btcsuite/btcd/wire"
)

const TaprootAnnexPrefix = 0x50

func parseScript(script []byte) {

}

func parseTransactionToInscription(tx wire.MsgTx) (*string, error) {
	witness := tx.TxIn[0].Witness
	if len(witness) == 0 {
		// EmptyWitness
		return nil, nil
	}
	if len(witness) == 1 {
		// KeyPathSpend
		return nil, nil
	}
	// annex define, see https://github.com/bitcoin/bips/blob/master/bip-0341.mediawiki
	// If there are at least two witness elements, and the first byte of the last element is 0x50,
	// this last element is called annex and is removed from the witness stack.
	// The annex (or the lack of thereof) is always covered by the signature and contributes to transaction weight,
	// but is otherwise ignored during taproot validation.
	var annex bool
	lastWitness := witness[len(witness)-1]
	annex = lastWitness[0] == 0x50
	if len(witness) == 2 && annex {
		// KeyPathSpend
		return nil, nil
	}
	var script []byte
	if annex {
		script = witness[len(witness)-1]
	} else {
		script = witness[len(witness)-2]
	}
	parseScript(script)
	scriptStr := string(script)
	return &scriptStr, nil
}

func ParseRawTransactionToInscription(rawTransaction string) (*string, error) {
	tx := bitcoin_cli_channel.DecodeRawTransaction(rawTransaction)
	if tx == nil {
		err := errors.New("ParseRawTransaction -> DecodeRawTransaction Failed")
		return nil, err
	}
	return parseTransactionToInscription(*tx)
}

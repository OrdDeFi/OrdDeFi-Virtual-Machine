package tx_utils

import (
	"brc20defi_vm/bitcoin_cli_channel"
	"errors"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func ParseFirstInputAddress(tx *wire.MsgTx) (*string, error) {
	if tx == nil {
		return nil, errors.New("ParseInputAddress error: transaction is nil")
	}
	if len(tx.TxIn) == 0 {
		return nil, errors.New("ParseFirstInputAddress error: TxIn is empty")
	}
	previousTxId := tx.TxIn[0].PreviousOutPoint.Hash.String()
	previousOutputIndex := tx.TxIn[0].PreviousOutPoint.Index
	previousRawTx := bitcoin_cli_channel.GetRawTransaction(previousTxId)
	if previousRawTx == nil {
		return nil, errors.New("ParseFirstInputAddress GetRawTransaction failed")
	}
	previousTx := bitcoin_cli_channel.DecodeRawTransaction(*previousRawTx)
	if previousTx == nil {
		return nil, errors.New("ParseFirstInputAddress DecodeRawTransaction failed")
	}
	if int(previousOutputIndex) >= len(previousTx.TxOut) {
		return nil, errors.New("ParseFirstInputAddress previousTxOut outbound")
	}
	previousOutput := previousTx.TxOut[previousOutputIndex]
	_, previousOutputAddress, _, err := txscript.ExtractPkScriptAddrs(previousOutput.PkScript, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}
	address := previousOutputAddress[0].EncodeAddress()
	if address == "" {
		return nil, errors.New("ParseFirstInputAddress error: parse address failed")
	}
	return &address, nil
}

func ParseFirstOutputAddress(tx *wire.MsgTx) (*string, error) {
	if tx == nil {
		return nil, errors.New("ParseFirstOutputAddress error: transaction is nil")
	}
	var err error
	address := ""
	for _, output := range tx.TxOut {
		scriptType, outputAddress, _, err := txscript.ExtractPkScriptAddrs(output.PkScript, &chaincfg.MainNetParams)
		if err != nil {
			break
		}
		if scriptType == txscript.NullDataTy {
			// First output is OpReturn, using next output
			continue
		}
		address = outputAddress[0].EncodeAddress()
		break
	}
	if err != nil {
		return nil, err
	}
	if address == "" {
		return nil, errors.New("ParseFirstOutputAddress error: parse address failed")
	}
	return &address, nil
}

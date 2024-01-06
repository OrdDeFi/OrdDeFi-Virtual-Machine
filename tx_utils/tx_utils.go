package tx_utils

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
	"errors"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"strings"
)

/*
ParseInputAddressAndValue
returns &inputAddress, &inputValue, err
1. If txIn is not coinbase, inputAddress and inputValue are not nil, err is nil
2. If txIn is coinbase, inputAddress and inputValue are nil, err is nil
There couldn't be a situation ((address is nil) && (inputValue is not nil)), or ((address is not nil) && (inputValue is nil))
*/
func ParseInputAddressAndValue(txIn *wire.TxIn) (*string, *int64, error) {
	const coinbaseOutputIndex uint32 = 4294967295
	const coinbaseTxId = "0000000000000000000000000000000000000000000000000000000000000000"
	previousTxId := txIn.PreviousOutPoint.Hash.String()
	previousOutputIndex := txIn.PreviousOutPoint.Index
	if previousTxId == coinbaseTxId && previousOutputIndex == coinbaseOutputIndex {
		return nil, nil, nil
	}
	previousRawTx := bitcoin_cli_channel.GetRawTransaction(previousTxId)
	if previousRawTx == nil {
		return nil, nil, errors.New("ParseInputAddressAndValue GetRawTransaction failed")
	}
	previousTx := bitcoin_cli_channel.DecodeRawTransaction(*previousRawTx)
	if previousTx == nil {
		return nil, nil, errors.New("ParseInputAddressAndValue DecodeRawTransaction failed")
	}
	if int(previousOutputIndex) >= len(previousTx.TxOut) {
		return nil, nil, errors.New("ParseInputAddressAndValue previousTxOut outbound")
	}
	previousOutput := previousTx.TxOut[previousOutputIndex]
	_, previousOutputAddress, _, err := txscript.ExtractPkScriptAddrs(previousOutput.PkScript, &chaincfg.MainNetParams)
	if err != nil {
		return nil, nil, err
	}
	address := previousOutputAddress[0].EncodeAddress()
	if address == "" {
		return nil, nil, errors.New("ParseInputAddressAndValue error: parse address failed")
	}
	value := previousOutput.Value
	return &address, &value, nil
}

func ParseFirstInputAddress(tx *wire.MsgTx) (*string, error) {
	if tx == nil {
		return nil, errors.New("ParseFirstInputAddress error: transaction is nil")
	}
	if len(tx.TxIn) == 0 {
		return nil, errors.New("ParseFirstInputAddress error: TxIn is empty")
	}
	addrPointer, _, err := ParseInputAddressAndValue(tx.TxIn[0])
	return addrPointer, err
}

func ParseOutputAddress(output *wire.TxOut) (*string, error) {
	scriptType, outputAddress, _, err := txscript.ExtractPkScriptAddrs(output.PkScript, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}
	if scriptType == txscript.NullDataTy {
		// output is OpReturn
		return nil, nil
	} else if scriptType == txscript.NonStandardTy {
		// output is non-standard OpReturn
		return nil, nil
	}
	address := outputAddress[0].EncodeAddress()
	return &address, nil
}

func ParseFirstOutputAddress(tx *wire.MsgTx) (*string, error) {
	if tx == nil {
		return nil, errors.New("ParseFirstOutputAddress error: transaction is nil")
	}
	var err error
	address := ""
	for _, output := range tx.TxOut {
		addressPointer, err := ParseOutputAddress(output)
		if err != nil {
			break
		}
		if addressPointer == nil {
			// kind of OpReturn
			continue
		}
		address = *addressPointer
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

func IsValidBitcoinAddress(address string) bool {
	if strings.ToLower(address) == "blackhole" || strings.ToLower(address) == "wormhole" {
		return true
	}
	_, err := btcutil.DecodeAddress(address, &chaincfg.MainNetParams)
	if err != nil {
		return false
	}
	return true
}

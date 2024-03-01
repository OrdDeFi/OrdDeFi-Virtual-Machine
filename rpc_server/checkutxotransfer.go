package rpc_server

import (
	"OrdDeFi-Virtual-Machine/subcommands"
	"encoding/json"
	"net/http"
)

/*
checkUTXOTransfer

	param: utxo
	return: {
	  "utxo": "8a****e8:0",
	  "address": "bc1p****abcd",
	  "tick": "odfi",
	  "amount": "1000"
	}
*/
func checkUTXOTransfer(w http.ResponseWriter, r *http.Request) {
	utxoStr := r.URL.Query().Get("utxo")
	w.Header().Set("Content-Type", "application/json")
	addressPtr, tickPtr, amount, err := subcommands.CheckUTXOTransferData(utxoStr, glDataDir)
	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if addressPtr == nil || tickPtr == nil || amount == nil {
		_ = json.NewEncoder(w).Encode(map[string]string{"utxo": utxoStr})
		return
	}
	var jsonRes map[string]interface{}
	jsonRes = make(map[string]interface{})
	jsonRes["utxo"] = utxoStr
	jsonRes["address"] = *addressPtr
	jsonRes["tick"] = *tickPtr
	jsonRes["amount"] = amount.String()
	_ = json.NewEncoder(w).Encode(jsonRes)
}

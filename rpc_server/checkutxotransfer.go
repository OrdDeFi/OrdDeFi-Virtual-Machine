package rpc_server

import (
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
	_ = json.NewEncoder(w).Encode(map[string]string{"utxo": utxoStr})
}

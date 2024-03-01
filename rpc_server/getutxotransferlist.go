package rpc_server

import (
	"encoding/json"
	"net/http"
)

/*
getUTXOTransferList

	param: tick
	return: {
	  "tick": "odfi",
	  "transferable_utxos": {
		"0d****ee:0": {
		  "address": "bc1p****abc0",
		  "amount": "210"
		},
		"8a****e8:0": {
		  "address": "bc1p****abc1",
		  "amount": "1000"
		},
	  }
	}
*/
func getUTXOTransferList(w http.ResponseWriter, r *http.Request) {
	tickStr := r.URL.Query().Get("tick")

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"tick": tickStr})
}

package rpc_server

import (
	"encoding/json"
	"net/http"
)

/*
getAddressUTXOTransferList

	param: address
	return: {
	  "address": "bc1p****abcd",
	  "transferable_utxos": {
	    "odfi": {
	      "0d****ee:0": "210",
	      "8a****e8:0": "1000"
	    },
		"odgv": {
	      "1d****fe:0": "500",
	      "33****68:0": "10000"
	    }
	  }
	}
*/
func getAddressUTXOTransferList(w http.ResponseWriter, r *http.Request) {
	addressStr := r.URL.Query().Get("address")

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"address": addressStr})
}

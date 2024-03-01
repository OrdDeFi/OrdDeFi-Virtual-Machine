package rpc_server

import (
	"encoding/json"
	"net/http"
)

/*
getAddressBalance

	param: address
	return: {
	  "address": "bc1p****abcd",
	  "assets": {
	    "odfi": {
	      "a": "10000",
	      "t": "10000"
	    },
		"odgv": {
	      "a": "10000",
	      "t": "10000"
	    }
	  }
	}
*/
func getAddressBalance(w http.ResponseWriter, r *http.Request) {
	addressStr := r.URL.Query().Get("address")

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"address": addressStr})
}

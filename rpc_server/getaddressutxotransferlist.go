package rpc_server

import (
	"OrdDeFi-Virtual-Machine/subcommands"
	"encoding/json"
	"net/http"
	"strings"
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
func getAddressUTXOTransferList(w http.ResponseWriter, req *http.Request) {
	addressStr := req.URL.Query().Get("address")
	w.Header().Set("Content-Type", "application/json")
	r, err := subcommands.GetAddressUTXOTransferListData(addressStr, glDataDir)
	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	var jsonRes map[string]interface{}
	var transferUTXOMap map[string]interface{}
	jsonRes = make(map[string]interface{})
	transferUTXOMap = make(map[string]interface{})
	jsonRes["address"] = addressStr
	jsonRes["transferable_utxos"] = transferUTXOMap
	for k, v := range r {
		kComps := strings.Split(k, ":")
		if len(kComps) != 5 {
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "GetAddressUTXOTransferListData result key parsing error"})
			return
		}
		tick := kComps[2]
		txid := kComps[3]
		if tickMap, ok := transferUTXOMap[tick]; ok {
			if castTickMap, ok2 := tickMap.(map[string]string); ok2 {
				castTickMap[txid+":0"] = v
			}
		} else {
			var newTickMap map[string]string
			newTickMap = make(map[string]string)
			newTickMap[txid+":0"] = v
			transferUTXOMap[tick] = newTickMap
		}
	}
	_ = json.NewEncoder(w).Encode(jsonRes)
}

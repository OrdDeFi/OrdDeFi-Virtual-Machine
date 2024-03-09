package rpc_server

import (
	"OrdDeFi-Virtual-Machine/subcommands"
	"encoding/json"
	"net/http"
	"strings"
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
func getUTXOTransferList(w http.ResponseWriter, req *http.Request) {
	tickStr := req.URL.Query().Get("tick")
	w.Header().Set("Content-Type", "application/json")

	r, err := subcommands.GetUTXOTransferListData(tickStr, glDataDir)
	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	var jsonRes map[string]interface{}
	var transferUTXOMap map[string]interface{}
	jsonRes = make(map[string]interface{})
	transferUTXOMap = make(map[string]interface{})
	jsonRes["tick"] = tickStr
	jsonRes["transferable_utxos"] = transferUTXOMap
	for k, v := range r {
		kComps := strings.Split(k, ":")
		if len(kComps) != 5 {
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "GetUTXOTransferListData result key parsing error"})
			return
		}
		address := kComps[2]
		txid := kComps[3]
		if addressMap, ok := transferUTXOMap[address]; ok {
			if castAddressMap, ok2 := addressMap.(map[string]string); ok2 {
				castAddressMap[txid+":0"] = v
			}
		} else {
			var newAddressMap map[string]string
			newAddressMap = make(map[string]string)
			newAddressMap[txid+":0"] = v
			transferUTXOMap[address] = newAddressMap
		}
	}
	_ = json.NewEncoder(w).Encode(jsonRes)
}

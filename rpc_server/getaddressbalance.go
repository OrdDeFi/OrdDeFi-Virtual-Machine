package rpc_server

import (
	"OrdDeFi-Virtual-Machine/subcommands"
	"encoding/json"
	"net/http"
	"strings"
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
func getAddressBalance(w http.ResponseWriter, req *http.Request) {
	addressStr := req.URL.Query().Get("address")
	w.Header().Set("Content-Type", "application/json")

	r, err := subcommands.GetAddressBalanceData(addressStr, glDataDir)
	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var jsonRes map[string]interface{}
	var assetsMap map[string]interface{}
	jsonRes = make(map[string]interface{})
	assetsMap = make(map[string]interface{})
	jsonRes["address"] = addressStr
	jsonRes["assets"] = assetsMap
	for k, v := range r {
		kComps := strings.Split(k, ":")
		if len(kComps) != 5 {
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "GetAddressBalanceData result key parsing error"})
			return
		}
		tick := kComps[3]
		subAccount := kComps[4]
		if tickMap, ok := assetsMap[tick]; ok {
			if castTickMap, ok2 := tickMap.(map[string]string); ok2 {
				castTickMap[subAccount] = v
			}
		} else {
			var newTickMap map[string]string
			newTickMap = make(map[string]string)
			newTickMap[subAccount] = v
			assetsMap[tick] = newTickMap
		}
	}
	_ = json.NewEncoder(w).Encode(jsonRes)
	return
}

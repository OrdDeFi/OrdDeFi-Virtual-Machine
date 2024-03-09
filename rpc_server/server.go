package rpc_server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
)

var glDataDir string

func Serve(port int, dataDir string) error {
	glDataDir = dataDir

	http.HandleFunc("/getaddressbalance", getAddressBalance)
	http.HandleFunc("/checkutxotransfer", checkUTXOTransfer)
	http.HandleFunc("/getaddressutxotransferlist", getAddressUTXOTransferList)
	http.HandleFunc("/getutxotransferlist", getUTXOTransferList)

	l, e := net.Listen("tcp", ":"+strconv.Itoa(port))
	if e != nil {
		return fmt.Errorf("listen error: %s", e)
	}

	go func() {
		err := http.Serve(l, nil)
		if err != nil {
			log.Printf("serve error: %s", err)
		}
	}()
	return nil
}

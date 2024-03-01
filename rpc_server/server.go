package rpc_server

import (
	"net"
	"net/http"
	"strconv"
)

func Serve(port int) error {
	http.HandleFunc("/getaddressbalance", getAddressBalance)
	http.HandleFunc("/checkutxotransfer", checkUTXOTransfer)
	http.HandleFunc("/getaddressutxotransferlist", getAddressUTXOTransferList)
	http.HandleFunc("/getutxotransferlist", getUTXOTransferList)

	l, e := net.Listen("tcp", ":"+strconv.Itoa(port))
	if e != nil {
		panic("listen error:" + e.Error())
	}
	err := http.Serve(l, nil)
	if err != nil {
		return err
	}
	return nil
}

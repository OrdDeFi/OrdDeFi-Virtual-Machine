package rpc_server

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
)

func httpAddHandler(w http.ResponseWriter, r *http.Request) {
	aStr := r.URL.Query().Get("a")
	bStr := r.URL.Query().Get("b")
	a, _ := strconv.Atoi(aStr)
	b, _ := strconv.Atoi(bStr)

	var result int
	result = a + b

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]int{"result": result})
}

func Serve(port int) error {
	http.HandleFunc("/add", httpAddHandler)
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

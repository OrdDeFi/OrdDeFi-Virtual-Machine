package test

func TestingMultiCommandsAndTxid() ([]byte, string) {
	command := `[
		{"p":"orddefi","op":"deploy","tick":"odfi","max":"21000000","lim":"1000","alim":"1000","icon":""},
		{"p":"orddefi","op":"deploy","tick":"test_long_decimal","max":"9900 0000 0000 0000 0000.0123 4567 8901 2345 6789","lim":"1000","alim":"1000","icon":""},
		{"p":"orddefi","op":"mint","tick":"odfi","amt":"1000"},
		{"p":"orddefi","op":"mint","tick":"odgv","amt":"1000"},
		{"p":"orddefi","op":"transfer","tick":"odfi","amt":"1000"},
		{"p":"orddefi","op":"transfer","tick":"odfi","amt":"1000", "to": "bc1qm34lsc65zpw79lxes69zkqmk6ee3ewf0j77s3h"},
		{"p":"orddefi","op":"addlp","ltick":"odfi","rtick":"odgv","lamt":"1000","ramt":"1000"},
		{"p":"orddefi","op":"addlp","ltick":"odfi","rtick":"odgv","lamt":"1000","ramt":"1000","as":"0"},
		{"p":"orddefi","op":"swap","ltick":"odfi","rtick":"odgv","spend":"odgv","amt":"10","trhd":"0.8%"},
		{"p":"orddefi","op":"rmlp","ltick":"odfi","rtick":"odgv","amt":"10"}
	]`
	txid := ""
	return []byte(command), txid
}

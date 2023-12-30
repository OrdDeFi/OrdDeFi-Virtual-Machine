package subcommands

import "OrdDeFi-Virtual-Machine/inscription_parser"

func ParseRawTransaction(parseRawTransactionString string) {
	contentType, content, err := inscription_parser.ParseRawTransactionToInscription(parseRawTransactionString)
	if err != nil {
		println("parserawtransaction error:", err)
	} else {
		println(*contentType, len(content))
		println(string(content))
	}
}

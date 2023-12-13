package bitcoin_cli_channel

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/wire"
	"os/exec"
	"strconv"
	"strings"
)

func GetBlockCount() int {
	cmd := exec.Command("bitcoin-cli", "getblockcount")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("GetBlockCount Error:", err)
		return 0
	}
	outputStr := string(output)
	outputStr = strings.Trim(outputStr, "\t\n")
	result, ok := strconv.Atoi(outputStr)
	if ok == nil {
		return result
	}
	return 0
}

func GetBlockHash(blockNumber int) *string {
	cmd := exec.Command("bitcoin-cli", "getblockhash", strconv.Itoa(blockNumber))
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("GetBlockHash Error:", err)
		return nil
	}
	outputStr := string(output)
	outputStr = strings.Trim(outputStr, "\t\n")
	return &outputStr
}

type BitcoinBlock struct {
	Hash              string   `json:"hash"`
	Confirmations     int      `json:"confirmations"`
	Height            int      `json:"height"`
	Version           int      `json:"version"`
	VersionHex        string   `json:"versionHex"`
	Merkleroot        string   `json:"merkleroot"`
	Time              int      `json:"time"`
	Mediantime        int      `json:"mediantime"`
	Nonce             int64    `json:"nonce"`
	Bits              string   `json:"bits"`
	Difficulty        float64  `json:"difficulty"`
	Chainwork         string   `json:"chainwork"`
	NTx               int      `json:"nTx"`
	Previousblockhash string   `json:"previousblockhash"`
	Strippedsize      int      `json:"strippedsize"`
	Size              int      `json:"size"`
	Weight            int      `json:"weight"`
	Tx                []string `json:"tx"`
}

func GetBlock(blockHash string) *BitcoinBlock {
	cmd := exec.Command("bitcoin-cli", "getblock", blockHash)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("GetBlockHash Error:", err)
		return nil
	}
	var res BitcoinBlock
	err = json.Unmarshal(output, &res)
	if err != nil {
		fmt.Println("GetBlock Parse JSON Error:", err)
		return nil
	}
	return &res
}

func GetRawTransaction(txId string) *string {
	cmd := exec.Command("bitcoin-cli", "getrawtransaction", txId)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("GetRawTransaction Error:", err)
		return nil
	}
	outputStr := string(output)
	outputStr = strings.Trim(outputStr, "\t\n")
	return &outputStr
}

func DecodeRawTransaction(rawTransactionString string) *wire.MsgTx {
	rawTxBin, err := hex.DecodeString(rawTransactionString)
	if err != nil {
		// Handle decoding error
		fmt.Println("Failed to decode raw transaction:", err)
		return nil
	}
	var tx wire.MsgTx
	err = tx.Deserialize(bytes.NewReader(rawTxBin))
	if err != nil {
		// Handle deserialization error
		fmt.Println("Failed to deserialize transaction:", err)
		return nil
	}
	return &tx
}

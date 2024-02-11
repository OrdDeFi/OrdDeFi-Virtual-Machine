package bitcoin_cli_channel

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/wire"
	"os/exec"
	"strconv"
	"strings"
)

func combinedCmdArgs(baseCmd []string) []string {
	params := GlobalParams().Params
	finalCmd := append(baseCmd[:1], append(params, baseCmd[1:]...)...)
	return finalCmd
}

func GetBlockCount() int {
	args := []string{"bitcoin-cli", "getblockcount"}
	args = combinedCmdArgs(args)
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("GetBlockCount Error:", err)
		return 0
	}
	outputStr := string(output)
	outputStr = strings.TrimSpace(outputStr)
	result, ok := strconv.Atoi(outputStr)
	if ok == nil {
		return result
	}
	return 0
}

func GetBlockHash(blockNumber int) *string {
	args := []string{"bitcoin-cli", "getblockhash", strconv.Itoa(blockNumber)}
	args = combinedCmdArgs(args)
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("GetBlockHash Error:", err)
		return nil
	}
	outputStr := string(output)
	outputStr = strings.TrimSpace(outputStr)
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
	args := []string{"bitcoin-cli", "getblock", blockHash}
	args = combinedCmdArgs(args)
	cmd := exec.Command(args[0], args[1:]...)
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
	args := []string{"bitcoin-cli", "getrawtransaction", txId}
	args = combinedCmdArgs(args)
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("GetRawTransaction Error:", err)
		return nil
	}
	outputStr := string(output)
	outputStr = strings.TrimSpace(outputStr)
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

func GetVersion() *string {
	args := []string{"bitcoin-cli", "--version"}
	args = combinedCmdArgs(args)
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("GetRawTransaction Error:", err)
		return nil
	}
	outputStr := string(output)
	outputStr = strings.TrimSpace(outputStr)
	lines := strings.Split(outputStr, "\n")
	outputStr = lines[0]
	components := strings.Split(outputStr, "version ")
	if len(components) != 2 {
		return nil
	}
	versionStr := components[1]
	return &versionStr
}

func VersionGreaterThanMinRequirement() (*bool, error) {
	version := GetVersion()
	if version == nil {
		return nil, errors.New("GetVersion error")
	}
	cmpRes := strings.Compare(*version, "v24.0.1")
	result := cmpRes >= 0
	return &result, nil
}

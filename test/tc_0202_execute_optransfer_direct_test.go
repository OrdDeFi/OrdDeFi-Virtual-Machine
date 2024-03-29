package test

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"OrdDeFi-Virtual-Machine/virtual_machine/operations"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestingTransferInSingleSliceCommands(tick string, txId string, amt string, to string) (*instruction_set.OpTransferInstruction, error) {
	commands := `[
		{"p":"orddefi","op":"transfer","tick":"` + tick + `","amt":"` + amt + `"}
	]`
	if to != "" {
		commands = `[
			{"p":"orddefi","op":"transfer","tick":"` + tick + `","amt":"` + amt + `", "to":"` + to + `"}
		]`
	}
	rawTx := bitcoin_cli_channel.GetRawTransaction(txId)
	if rawTx == nil {
		return nil, fmt.Errorf("TestingTransferInSingleSliceCommands GetRawTransaction error")
	}
	tx := bitcoin_cli_channel.DecodeRawTransaction(*rawTx)
	if tx == nil {
		return nil, fmt.Errorf("TestingTransferInSingleSliceCommands DecodeRawTransaction error")
	}
	instructions, err := virtual_machine.CompileInstructions("text/plain", []byte(commands), tx, txId)
	if err != nil {
		return nil, fmt.Errorf("TestCommandParse CompileInstructions error: %s", err.Error())
	}
	if len(instructions) != 1 {
		return nil, fmt.Errorf("TestingTransferInSingleSliceCommands CompileInstructions error: instructions length should be 1")
	}
	for _, instruction := range instructions {
		switch value := instruction.(type) {
		case instruction_set.OpTransferInstruction:
			return &value, nil
		default:
			return nil, fmt.Errorf("TestingTransferInSingleSliceCommands error: instruction type error, expected OpDeployInstruction")
		}
	}
	return nil, fmt.Errorf("TestingTransferInSingleSliceCommands error: instruction type error: no instruction compiled")
}

func testDirectTransferCommand(t *testing.T, db *db_utils.OrdDB, tick string, txId string, amt string, to string) {
	// 1. compile instruction
	instruction, err := TestingTransferInSingleSliceCommands(tick, txId, amt, to)

	if to != "" && err != nil {
		if len(tick) == 4 {
			if err.Error() != "TestCommandParse CompileInstructions error: no privileges on cross-address transfer" {
				t.Errorf("testDirectTransferCommand error: %s", err.Error())
			}
		}
		return
	}
	if len(tick) != 4 && err != nil {
		if err.Error() != "TestingTransferInSingleSliceCommands CompileInstructions error: instructions length should be 1" {
			t.Errorf("testDirectTransferCommand error: %s", err.Error())
		}
		return
	}
	if instruction == nil {
		t.Errorf("testDirectTransferCommand error: transfer instruction is nil")
		return
	}

	// 2. execute deploy op
	err = operations.ExecuteTransfer(*instruction, db)

	if err != nil {
		if strings.HasPrefix(err.Error(), "performTransferBatchWriteKV from address balance error") == false {
			t.Errorf("TestExecuteMint error: execute OpMint error %s", err)
		}
	}
}

func testDirectBalanceNotEnough(t *testing.T, db *db_utils.OrdDB, tick string) {
	txId := "61de96170018ce878b1adf287b8ac9cf0e4f0ad8c5a69af203cc25bbde72a13e"
	txInAddr := "bc1q2f0tczgrukdxjrhhadpft2fehzpcrwrz549u90"
	to := "bc1qr35hws365juz5rtlsjtvmulu97957kqvr3zpw3"
	receiverTickInitBalance, _ := memory_read.AvailableBalance(db, tick, to)
	senderTickInitBalance, _ := memory_read.AvailableBalance(db, tick, txInAddr)

	testDirectTransferCommand(t, db, tick, txId, "1001", "bc1qr35hws365juz5rtlsjtvmulu97957kqvr3zpw3")

	receiverTickEndBalance, _ := memory_read.AvailableBalance(db, tick, to)
	senderTickEndBalance, _ := memory_read.AvailableBalance(db, tick, txInAddr)

	assert.True(t, receiverTickInitBalance.IsEqualTo(receiverTickEndBalance))
	assert.True(t, senderTickInitBalance.IsEqualTo(senderTickEndBalance))
	assert.True(t, !receiverTickEndBalance.IsNegative())
	assert.True(t, !senderTickEndBalance.IsNegative())
}

func testTickIllegalTransfer(t *testing.T, db *db_utils.OrdDB, tick string) {
	txId := "61de96170018ce878b1adf287b8ac9cf0e4f0ad8c5a69af203cc25bbde72a13e"
	txInAddr := "bc1q2f0tczgrukdxjrhhadpft2fehzpcrwrz549u90"
	to := "bc1qr35hws365juz5rtlsjtvmulu97957kqvr3zpw3"
	receiverTickInitBalance, _ := memory_read.AvailableBalance(db, tick, to)
	senderTickInitBalance, _ := memory_read.AvailableBalance(db, tick, txInAddr)
	testDirectTransferCommand(t, db, tick, txId, "50", "bc1qr35hws365juz5rtlsjtvmulu97957kqvr3zpw3")

	receiverTickEndBalance, _ := memory_read.AvailableBalance(db, tick, to)
	senderTickEndBalance, _ := memory_read.AvailableBalance(db, tick, txInAddr)

	assert.True(t, receiverTickInitBalance.IsZero())
	assert.True(t, senderTickInitBalance.IsZero())
	assert.True(t, receiverTickEndBalance.IsZero())
	assert.True(t, senderTickEndBalance.IsZero())
}

func testDirectNormalTransfer(t *testing.T, db *db_utils.OrdDB, tick string) {
	txId := "61de96170018ce878b1adf287b8ac9cf0e4f0ad8c5a69af203cc25bbde72a13e"
	txInAddr := "bc1q2f0tczgrukdxjrhhadpft2fehzpcrwrz549u90"
	to := "bc1qr35hws365juz5rtlsjtvmulu97957kqvr3zpw3"
	receiverTickInitBalance, _ := memory_read.AvailableBalance(db, tick, to)
	senderTickInitBalance, _ := memory_read.AvailableBalance(db, tick, txInAddr)
	if senderTickInitBalance.IsZero() {
		TestingMintForParam(t, db, tick, txId, "1000")
		senderTickInitBalance, _ = memory_read.AvailableBalance(db, tick, txInAddr)

	}

	testDirectTransferCommand(t, db, tick, txId, "50", "bc1qr35hws365juz5rtlsjtvmulu97957kqvr3zpw3")

	receiverTickEndBalance, _ := memory_read.AvailableBalance(db, tick, to)
	senderTickEndBalance, _ := memory_read.AvailableBalance(db, tick, txInAddr)
	println(receiverTickInitBalance.String(), receiverTickEndBalance.String(), senderTickInitBalance.String(), senderTickEndBalance.String())

	assert.True(t, receiverTickEndBalance.IsEqualTo(receiverTickInitBalance.Add(safe_number.SafeNumFromString("50"))))
	assert.True(t, senderTickEndBalance.IsEqualTo(senderTickInitBalance.Subtract(safe_number.SafeNumFromString("50"))))
	assert.True(t, !receiverTickEndBalance.IsNegative())
	assert.True(t, !senderTickEndBalance.IsNegative())
}

func TestDirectTransfer(t *testing.T) {
	// open db
	db, err := db_utils.OpenDB(testDBPath)
	if err != nil {
		t.Errorf("TestExecuteMint OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")

	testTickIllegalTransfer(t, db, "shift")
	testDirectBalanceNotEnough(t, db, "odfi")
	testDirectNormalTransfer(t, db, "odfi")
}

func TestReadODFIBalanceAfterTransfer(t *testing.T) {
	TestingReadCoin(t, "odfi")
}

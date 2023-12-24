package test

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/operations"
	"strings"
	"testing"
)

func testUTXOTransferCommand(t *testing.T, db *db_utils.OrdDB, tick string, txId string, amt string) {
	// 1. compile instruction
	instruction, err := TestingTransferInSingleSliceCommands(tick, txId, amt, "")
	if len(tick) != 4 && err != nil {
		if err.Error() != "testUTXOTransferCommand CompileInstructions error: instructions length should be 1" {
			t.Errorf("testDirectTransferCommand error: %s", err.Error())
		}
		return
	}
	if instruction == nil {
		t.Errorf("testUTXOTransferCommand error: transfer instruction is nil")
		return
	}

	// 2. execute deploy op
	err = operations.ExecuteTransfer(*instruction, db)
	if err != nil {
		if strings.HasPrefix(err.Error(), "performTransferBatchWriteKV from address balance error") == false {
			t.Errorf("testUTXOTransferCommand error: execute OpMint error %s", err)
		}
	}
}

func testUTXOIllegalTick(t *testing.T, db *db_utils.OrdDB, tick string) {

}

func testUTXOBalanceNotEnough(t *testing.T, db *db_utils.OrdDB, tick string) {

}

func testUTXONormalTransfer(t *testing.T, db *db_utils.OrdDB, tick string) {

}

func TestUTXOTransfer(t *testing.T) {

}

func TestUTXOTransferArriving(t *testing.T) {

}

func TestReadODFIBalanceAfterUTXOTransfer(t *testing.T) {
	TestingReadCoin(t, "odfi")
}

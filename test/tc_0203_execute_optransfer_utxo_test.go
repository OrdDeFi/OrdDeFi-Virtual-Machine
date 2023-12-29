package test

import (
	"OrdDeFi-Virtual-Machine/bitcoin_cli_channel"
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"OrdDeFi-Virtual-Machine/virtual_machine/operations"
	"fmt"
	"testing"
)

func testMint(t *testing.T, db *db_utils.OrdDB, tick string, txId string) {
	TestingMintForParam(t, db, tick, txId, "1000")
}

func testCreateUTXOPacker(t *testing.T, db *db_utils.OrdDB, tick string, txId string) {
	cmd := `{"p":"orddefi","op":"transfer","tick":"odfi","amt":"50.1"}`
	instructions := TestingCompile(t, cmd, txId)
	instruction := instructions[0]
	switch value := instruction.(type) {
	case instruction_set.OpTransferInstruction:
		err := operations.ExecuteTransfer(value, db)
		if err != nil {
			t.Errorf("testCreateUTXOPacker error: %s", err)
			return
		}
	default:
		t.Errorf("testCreateUTXOPacker error: instruction is not OpTransfer")
		return
	}

}

func testApplyUTXOTransfer(t *testing.T, db *db_utils.OrdDB, txId string) {
	rawTx := bitcoin_cli_channel.GetRawTransaction(txId)
	if rawTx == nil {
		t.Errorf("TestingCompile GetRawTransaction error")
		return
	}
	tx := bitcoin_cli_channel.DecodeRawTransaction(*rawTx)
	if tx == nil {
		t.Errorf("TestingCompile DecodeRawTransaction error")
		return
	}
	operations.ApplyUTXOTransfer(db, tx)
}

func testReadBalanceAfterUTXOTransfer(t *testing.T, db *db_utils.OrdDB, tick string) {
	coinRes, err := memory_read.AllAddressBalanceForCoin(db, tick)
	if err != nil {
		t.Errorf("testReadBalanceAfterUTXOTransfer AllAddressBalanceForCoin error: %s", err.Error())
	}
	for k, v := range coinRes {
		println(k, v)
	}
}

func TestUTXOTransfer(t *testing.T) {
	// open db
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestExecuteMint OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)

	tick := "odfi"
	fmt.Println("DB opened successfully.")
	testMint(t, db, tick, "51600cbc3c61c1c42329ccd11713f39a41a145def64edfd80d5232a2a01a5964")
	testReadBalanceAfterUTXOTransfer(t, db, tick)
	testCreateUTXOPacker(t, db, tick, "adbcaaf04d8f8767a72cbd06f20e735122869758ab7dc3a4ef31cdcc39aa04b7")
	testReadBalanceAfterUTXOTransfer(t, db, tick)
	testApplyUTXOTransfer(t, db, "f347176e945684f19d9acfa52466c10c8ea62d1b5f2fb6a0382e5ea7f1171107")
	testReadBalanceAfterUTXOTransfer(t, db, tick)
}

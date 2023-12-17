package test

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"OrdDeFi-Virtual-Machine/virtual_machine/operations"
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func randomTick() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	const letters = "abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, 4)
	for i := 0; i < 4; i++ {
		result[i] = letters[r.Intn(len(letters))]
	}
	randomString := string(result)
	return randomString
}

func checkStringSlicesEqual(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := 0; i < len(slice1); i++ {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}

func TestExecuteDeployInBatchCommands(t *testing.T) {
	// 1. compile instruction
	randTick := randomTick()
	instruction, err := TestingDeployInSingleSliceCommands(randTick)
	if err != nil {
		t.Errorf("%s", err.Error())
	}
	if instruction == nil {
		t.Errorf("TestExecuteDeployInBatchCommands error: deploy instruction is nil")
	}

	// 2. open db
	db, err := db_utils.OpenDB("./test_db")
	if err != nil {
		t.Errorf("TestDBReadPrefix OpenDB error: %s", err.Error())
	}
	defer db_utils.CloseDB(db)
	fmt.Println("DB opened successfully.")

	// 3. record coin list
	allCoinsBefore, err := memory_read.AllCoins(db)
	if err != nil {
		t.Errorf("TestExecuteDeployInBatchCommands memory_read error: %s", err.Error())
	}
	allDeployedCoinsBefore, err := memory_read.AllDeployedCoins(db)
	if err != nil {
		t.Errorf("TestExecuteDeployInBatchCommands memory_read error: %s", err.Error())
	}
	// 4. execute deploy op
	err = operations.ExecuteOpDeploy(*instruction, db)
	if err != nil {
		t.Errorf("TestExecuteDeployInBatchCommands error: execute deploy error: %s", err.Error())
	}
	// 5. check coin list again
	allCoinsAfter, err := memory_read.AllCoins(db)
	if err != nil {
		t.Errorf("TestExecuteDeployInBatchCommands memory_read error: %s", err.Error())
	}
	allDeployedCoinsAfter, err := memory_read.AllDeployedCoins(db)
	if err != nil {
		t.Errorf("TestExecuteDeployInBatchCommands memory_read error: %s", err.Error())
	}
	// Check
	allCoinsBefore = append(allCoinsBefore, randTick)
	allDeployedCoinsBefore = append(allDeployedCoinsBefore, randTick)

	sort.Strings(allCoinsBefore)
	sort.Strings(allDeployedCoinsBefore)
	sort.Strings(allCoinsAfter)
	sort.Strings(allDeployedCoinsAfter)
	if checkStringSlicesEqual(allCoinsBefore, allCoinsAfter) == false {
		t.Errorf("TestExecuteDeployInBatchCommands all coins not matching: %s, expected %s", allCoinsAfter, allCoinsBefore)
	}
	if checkStringSlicesEqual(allDeployedCoinsBefore, allDeployedCoinsAfter) == false {
		t.Errorf("TestExecuteDeployInBatchCommands all deployed coins not matching: %s, expected %s", allDeployedCoinsAfter, allDeployedCoinsBefore)
	}
	coinMeta, err := memory_read.CoinMeta(db, randTick)
	if err != nil {
		t.Errorf("TestExecuteDeployInBatchCommands memory_read error: %s", err.Error())
	}
	deployedMaxString := coinMeta.Max.String()
	deployedLimString := coinMeta.Lim.String()
	deployedAddrLimString := coinMeta.AddrLim.String()
	deployedDescString := coinMeta.Desc
	deployedIconString := coinMeta.Icon
	if deployedMaxString != (*instruction).Max {
		t.Errorf("TestExecuteDeployInBatchCommands memory_read error: max not matching")
	}
	if deployedLimString != (*instruction).Lim {
		t.Errorf("TestExecuteDeployInBatchCommands memory_read error: lim not matching")
	}
	if deployedAddrLimString != (*instruction).AddrLim {
		t.Errorf("TestExecuteDeployInBatchCommands memory_read error: addr lim not matching")
	}
	if deployedDescString != (*instruction).Desc {
		t.Errorf("TestExecuteDeployInBatchCommands memory_read error: addr lim not matching")
	}
	if deployedIconString != (*instruction).Icon {
		t.Errorf("TestExecuteDeployInBatchCommands memory_read error: addr lim not matching")
	}
}

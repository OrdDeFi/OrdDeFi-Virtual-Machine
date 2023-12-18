package operations

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"errors"
)

func ExecuteOpMint(instruction instruction_set.OpMintInstruction, db *db_utils.OrdDB) error {
	// check if address is legal
	address := instruction.TxOutAddr
	if address == "" {
		return errors.New("address is nil for OpMint")
	}
	// check if coin name is legal
	coinName := instruction.Tick
	coinMeta, err := memory_read.CoinMeta(db, coinName)
	if err != nil {
		return err
	}
	if coinMeta == nil {
		return errors.New("CoinMeta not found named " + coinName)
	}
	// check if amt is legal
	amount := safe_number.SafeNumFromString(instruction.Amt)
	if amount == nil {
		return errors.New("Amount parse failed: " + instruction.Amt)
	}
	// query total minted value
	totalMinted, err := memory_read.TotalMintedBalance(db, coinName)
	if err != nil {
		return err
	}
	if totalMinted == nil {
		return errors.New("total minted returns nil: " + coinName)
	}
	// query address minted value
	addressMinted, err := memory_read.AddressMintedBalance(db, coinName, address)
	if err != nil {
		return err
	}
	if addressMinted == nil {
		return errors.New("address minted returns nil: " + coinName + " @" + address)
	}
	return nil
}

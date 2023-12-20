package operations

import (
	"OrdDeFi-Virtual-Machine/db_utils"
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/instruction_set"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_read"
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_write"
	"errors"
)

func ExecuteOpMint(instruction instruction_set.OpMintInstruction, db *db_utils.OrdDB) error {
	if instruction.IsValidOpMintInstruction() == false {
		return errors.New("repeat mint disabled")
	}
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
	commandAmount := safe_number.SafeNumFromString(instruction.Amt)
	if commandAmount == nil {
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
	// calculating params
	// 1. calculating amount
	remaining := coinMeta.Max.Subtract(totalMinted)
	addrRemaining := coinMeta.AddrLim.Subtract(addressMinted)
	minRemaining := remaining.Min(addrRemaining)
	mintingAmount := minRemaining.Min(commandAmount)

	if remaining.IsZero() {
		return errors.New("Mint ended for " + coinName)
	}
	if addrRemaining.IsZero() {
		return errors.New("Address reached limit for " + coinName)
	}

	// 2. calculating new total minted and address minted
	newTotalMinted := totalMinted.Add(mintingAmount)
	newTotalMintedString := newTotalMinted.String()
	newAddressMinted := addressMinted.Add(mintingAmount)
	newAddressMintedString := newAddressMinted.String()
	// 3. calculating new balance
	balance, err := memory_read.AvailableBalance(db, coinName, address)
	if err != nil {
		return err
	}
	newBalance := balance.Add(mintingAmount)
	newBalanceString := newBalance.String()
	err = memory_write.WriteMintInfo(db, coinName, address, newTotalMintedString, newAddressMintedString, newBalanceString)
	return err
}

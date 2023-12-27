package test

import (
	"OrdDeFi-Virtual-Machine/safe_number"
	"OrdDeFi-Virtual-Machine/virtual_machine/operations"
	"fmt"
	"testing"
)

func testingDiscountForAmount(t *testing.T, amountString string) {
	r, err := operations.DiscountForODFIAmount(safe_number.SafeNumFromString(amountString))
	if err != nil {
		t.Errorf("calculating discount error: %s", err.Error())
	}
	fmt.Printf("discount for %s is %.2f\n", amountString, *r)
}

func TestDiscount(t *testing.T) {
	testingDiscountForAmount(t, "21000.01")
	testingDiscountForAmount(t, "21000")
	testingDiscountForAmount(t, "2100.01")
	testingDiscountForAmount(t, "2100")
	testingDiscountForAmount(t, "210.01")
	testingDiscountForAmount(t, "210")
	testingDiscountForAmount(t, "209.999")
	testingDiscountForAmount(t, "21")
	testingDiscountForAmount(t, "0")
}

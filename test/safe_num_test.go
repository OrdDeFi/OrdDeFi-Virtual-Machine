package test

import (
	"OrdDefi-Virtual-Machine/safe_number"
	"testing"
)

func TestEmptyString(t *testing.T) {
	num := safe_number.SafeNumFromString("")
	if num != nil {
		t.Errorf("SafeNumFromString(\"\") = %s; expected nil", num)
	}
}

func TestNegative(t *testing.T) {
	num := safe_number.SafeNumFromString("-1")
	if num != nil {
		t.Errorf("SafeNumFromString(\"-1\") = %s; expected nil", num)
	}
}

func TestAdd(t *testing.T) {
	num1 := safe_number.SafeNumFromString("1")
	num2 := safe_number.SafeNumFromString("2")
	res := num1.Add(num2)
	if res.String() != "3" {
		t.Errorf("%s + %s = %s; expected 3", num1, num2, res)
	}
}

func TestAddZero(t *testing.T) {
	num1 := safe_number.SafeNumFromString("1")
	num2 := safe_number.SafeNumFromString("0")
	res := num1.Add(num2)
	if res.String() != "1" {
		t.Errorf("%s + %s = %s; expected 1", num1, num2, res)
	}
}

func TestSubtract(t *testing.T) {
	num1 := safe_number.SafeNumFromString("5")
	num2 := safe_number.SafeNumFromString("2")
	res := num1.Subtract(num2)
	if res.String() != "3" {
		t.Errorf("%s - %s = %s; expected 3", num1, num2, res)
	}
}

func TestOverSubtract(t *testing.T) {
	num1 := safe_number.SafeNumFromString("2")
	num2 := safe_number.SafeNumFromString("5")
	res := num1.Subtract(num2)
	if res != nil {
		t.Errorf("%s - %s = %s; expected nil", num1, num2, res)
	}
}

func TestSubtractZero(t *testing.T) {
	num1 := safe_number.SafeNumFromString("2")
	num2 := safe_number.SafeNumFromString("0")
	res := num1.Subtract(num2)
	if res.String() != "2" {
		t.Errorf("%s - %s = %s; expected nil", num1, num2, res)
	}
}

func TestSubtractFraction(t *testing.T) {
	num1 := safe_number.SafeNumFromString("2")
	num2 := safe_number.SafeNumFromString("0.01234567890123")
	res := num1.Subtract(num2)
	if res.String() != "1.98765432109877" {
		t.Errorf("%s - %s = %s; expected nil", num1, num2, res)
	}
}

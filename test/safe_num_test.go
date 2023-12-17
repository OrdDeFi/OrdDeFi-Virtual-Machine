package test

import (
	"OrdDeFi-Virtual-Machine/safe_number"
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

func TestSubtract2(t *testing.T) {
	num1 := safe_number.SafeNumFromString("5")
	num2 := safe_number.SafeNumFromString("4")
	res := num1.Subtract(num2)
	if res.String() != "1" {
		t.Errorf("%s - %s = %s; expected 1", num1, num2, res)
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
		t.Errorf("%s - %s = %s; expected 2", num1, num2, res)
	}
}

func TestSubtractFraction(t *testing.T) {
	num1 := safe_number.SafeNumFromString("2")
	num2 := safe_number.SafeNumFromString("0.01234567890123")
	res := num1.Subtract(num2)
	if res.String() != "1.98765432109877" {
		t.Errorf("%s - %s = %s; expected 1.98765432109877", num1, num2, res)
	}
}

func TestMul(t *testing.T) {
	num1 := safe_number.SafeNumFromString("2.22")
	num2 := safe_number.SafeNumFromString("3.33")
	res := num1.Multiply(num2)
	if res.String() != "7.3926" {
		t.Errorf("%s * %s = %s; expected 7.3926", num1, num2, res)
	}
}

func TestMulFraction(t *testing.T) {
	num1 := safe_number.SafeNumFromString("2.22")
	num2 := safe_number.SafeNumFromString("0.333")
	res := num1.Multiply(num2)
	if res.String() != "0.73926" {
		t.Errorf("%s * %s = %s; expected 0.73926", num1, num2, res)
	}
}

func TestDivideBy(t *testing.T) {
	num1 := safe_number.SafeNumFromString("2.22")
	num2 := safe_number.SafeNumFromString("2")
	res := num1.DivideBy(num2)
	if res.String() != "1.11" {
		t.Errorf("%s / %s = %s; expected 1.11", num1, num2, res)
	}
}

func TestDivideByFraction(t *testing.T) {
	num1 := safe_number.SafeNumFromString("2.22")
	num2 := safe_number.SafeNumFromString("0.2")
	res := num1.DivideBy(num2)
	if res.String() != "11.1" {
		t.Errorf("%s / %s = %s; expected 11.1", num1, num2, res)
	}
}

package safe_number

import (
	"math"
	"math/big"
	"strings"
)

func formalNumString(inputStr string) *string {
	trimmedStr := strings.Replace(inputStr, " ", "", -1)
	trimmedStr = strings.Replace(trimmedStr, "\n", "", -1)
	trimmedStr = strings.Replace(trimmedStr, "\t", "", -1)
	if trimmedStr == "" {
		return nil
	}
	components := strings.Split(trimmedStr, ".")
	componentsLength := len(components)
	numStr := ""
	if componentsLength == 1 {
		numStr = inputStr
	} else if componentsLength == 2 {
		integerPart := components[0]
		fractionalPart := components[1]
		if len(fractionalPart) > 18 {
			fractionalPart = fractionalPart[:18]
		}
		numStr = integerPart + "." + fractionalPart
	} else {
		return nil
	}
	return &numStr
}

func bigFloatToInt10Pow18(num *big.Float) *big.Int {
	if num == nil {
		return nil
	}
	delta := 0.5
	if num.Sign() < 0 {
		delta = -0.5
	}
	multiplier := uint64(math.Pow(10, 18))
	num.Mul(num, new(big.Float).SetUint64(multiplier))
	num.Add(num, new(big.Float).SetFloat64(delta))
	bigInt, _ := num.Int(nil)
	return bigInt
}

type SafeNum struct {
	decimal big.Int // 10 ** 18 * raw_num_string
}

func SafeNumFromString(inputStr string) *SafeNum {
	if inputStr == "" {
		return nil
	}
	numStr := formalNumString(inputStr)
	if numStr == nil {
		return nil
	}
	d := new(big.Float)
	if d == nil {
		return nil
	}
	d.SetPrec(2048) // Supporting 617 digits in decimal
	d.SetString(*numStr)
	if d.Sign() < 0 {
		// negative params are now allowed
		return nil
	}
	bigInt := bigFloatToInt10Pow18(d)

	safeNum := SafeNum{
		decimal: *bigInt,
	}
	return &safeNum
}

func trimFractionTail(fractionalPart string) string {
	res := fractionalPart
	for i := len(res) - 1; i >= 0; i-- {
		if res[i] == '0' {
			if i == 0 {
				res = ""
			}
			continue
		}
		res = res[:i+1]
		break
	}
	return res
}

func (num SafeNum) String() string {
	intString := num.decimal.String()
	strLen := len(intString)
	if strLen > 18 {
		integerPart := intString[:strLen-18]
		fractionalPart := intString[strLen-18:]
		fractionalPart = trimFractionTail(fractionalPart)
		if fractionalPart == "" {
			return integerPart
		}
		return integerPart + "." + fractionalPart
	} else {
		prePart := "0."
		for i := 0; i < 18-strLen; i++ {
			prePart += "0"
		}
		return prePart + trimFractionTail(intString)
	}
}

func (num SafeNum) Add(rightNumber *SafeNum) *SafeNum {
	if rightNumber == nil {
		return nil
	} else if rightNumber.decimal.Sign() == 0 {
		return &num
	} else if rightNumber.decimal.Sign() < 0 {
		return nil
	}
	resultInt := new(big.Int)
	resultInt.Add(&num.decimal, &rightNumber.decimal)

	numCmpRes := num.decimal.Cmp(resultInt)
	if numCmpRes != -1 {
		// O.G. value should < resultInt, to protect from overflow
		return nil
	}
	rightCmpRes := rightNumber.decimal.Cmp(resultInt)
	if rightCmpRes != -1 {
		// right value should < resultInt, to protect from overflow
		return nil
	}
	safeNum := SafeNum{
		decimal: *resultInt,
	}
	return &safeNum
}

func (num SafeNum) Subtract(rightNumber *SafeNum) *SafeNum {
	if rightNumber == nil {
		return nil
	} else if rightNumber.decimal.Sign() == 0 {
		return &num
	} else if rightNumber.decimal.Sign() < 0 {
		return nil
	}
	resultInt := new(big.Int)
	resultInt.Sub(&num.decimal, &rightNumber.decimal)
	if resultInt.Sign() < 0 {
		// negative sub result not allowed
		return nil
	}

	numCmpRes := num.decimal.Cmp(resultInt)
	if numCmpRes != 1 {
		// O.G. value should > resultInt, to protect from overflow
		return nil
	}
	safeNum := SafeNum{
		decimal: *resultInt,
	}
	return &safeNum
}

func (num SafeNum) Multiply(rightNumber *SafeNum) *SafeNum {
	if rightNumber == nil {
		return nil
	} else if rightNumber.decimal.Sign() == 0 {
		return &num
	} else if rightNumber.decimal.Sign() < 0 {
		return nil
	}

	resultInt := new(big.Int)
	resultInt.Mul(&num.decimal, &rightNumber.decimal)
	resultRightShift18 := new(big.Int)
	multiplier := uint64(math.Pow(10, 18))
	resultRightShift18.Div(resultInt, new(big.Int).SetUint64(multiplier))
	if resultRightShift18.Sign() < 0 {
		// negative mul result not allowed
		return nil
	}
	numCmpRes := num.decimal.Cmp(resultRightShift18)
	oneCmpRight := new(big.Int).SetUint64(multiplier).Cmp(&rightNumber.decimal)
	if numCmpRes != oneCmpRight {
		// O.G. value cmp resultInt, should be same as 1e18 cmp to right number, to protect from overflow
		return nil
	}
	safeNum := SafeNum{
		decimal: *resultRightShift18,
	}
	return &safeNum
}

func (num SafeNum) DivideBy(rightNumber *SafeNum) *SafeNum {
	if rightNumber == nil {
		return nil
	} else if rightNumber.decimal.Sign() == 0 {
		return &num
	} else if rightNumber.decimal.Sign() < 0 {
		return nil
	}

	resultLeftShift18 := new(big.Int)
	multiplier := uint64(math.Pow(10, 18))
	resultLeftShift18.Mul(&num.decimal, new(big.Int).SetUint64(multiplier))
	resultInt := new(big.Int)
	resultInt.Div(resultLeftShift18, &rightNumber.decimal)
	if resultInt.Sign() < 0 {
		// negative mul result not allowed
		return nil
	}
	numCmpRes := num.decimal.Cmp(resultInt)
	rightCmp1e18 := rightNumber.decimal.Cmp(new(big.Int).SetUint64(multiplier))
	if numCmpRes != rightCmp1e18 {
		// O.G. value cmp resultInt, should be same as right number cmp to 1e18, to protect from overflow
		return nil
	}
	safeNum := SafeNum{
		decimal: *resultInt,
	}
	return &safeNum
}

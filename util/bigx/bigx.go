package bigx

import "math/big"

func Mul(val *big.Int, multiplier float64) *big.Int {
	result, accuracy := new(big.Float).Mul(new(big.Float).SetInt(val), big.NewFloat(multiplier)).Int(nil)
	if accuracy < 0 {
		result.Add(result, big.NewInt(1))
	}
	return result
}

func FromString(val string) *big.Int {
	result, _ := new(big.Int).SetString(val, 10)
	return result
}

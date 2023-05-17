package app

// import (
// 	"fmt"
// 	"math/big"
// 	"testing"

// 	"github.com/ethereum/go-ethereum/common/hexutil"
// 	"github.com/ethereum/go-ethereum/rpc"
// )

// func TestCalcBlockReward(t *testing.T) {
// 	// Get gas limit and used
// 	gasLimit := blockData["gasLimit"].(hexutil.Uint64)
// 	gasUsed := blockData["gasUsed"].(hexutil.Uint64)

// 	// Calculate gas price
// 	gasPrice := big.NewInt(20000000000) // Example gas price of 20 Gwei

// 	// Calculate block reward
// 	reward := new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(uint64(gasLimit)))
// 	reward.Add(reward, big.NewInt(5))
// 	reward.Div(reward, new(big.Int).SetUint64(uint64(gasUsed)))

// 	// Check expected reward value
// 	expectedRewardStr := "2775258912097940429"
// 	if reward.String() != expectedRewardStr {
// 		t.Errorf("Expected reward value: %s, got: %s", expectedRewardStr, reward.String())
// 	}

// 	fmt.Println("Block Reward:", reward.String())
// }

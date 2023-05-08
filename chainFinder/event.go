package chainFinder

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Event struct {
	Address      common.Address         `json:"address"`
	ContractName string                 `json:"contractName"`
	ChainID      *big.Int               `json:"chainID"`
	Data         map[string]interface{} `json:"data"`
	BlockHash    common.Hash            `json:"blockHash"`
	BlockNumber  string                 `json:"blockNumber"`
	Name         string                 `json:"name"`
	TxHash       common.Hash            `json:"txHash"`
	TxIndex      string                 `json:"txIndex"`
}

type EventHandler func(*Event) error

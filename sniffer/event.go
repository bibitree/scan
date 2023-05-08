package sniffer

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
	Gas          uint64                 `json:"gas"`
	GasPrice     *big.Int               `json:"gasPrice"`
	GasTipCap    *big.Int               `json:"gasTipCap"`
	GasFeeCap    *big.Int               `json:"gasFeeCap"`
	Value        *big.Int               `json:"value"`
	Nonce        uint64                 `json:"nonce"`
	To           *common.Address        `json:"to"`
}

type EventHandler func(*Event) error

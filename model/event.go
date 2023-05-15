package model

import (
	"math/big"
)

type Event struct {
	Address      string                 `json:"address"`
	ContractName string                 `json:"contractName"`
	ChainID      *big.Int               `json:"chainID"`
	Data         map[string]interface{} `json:"data"`
	BlockHash    string                 `json:"blockHash"`
	BlockNumber  string                 `json:"blockNumber"`
	Name         string                 `json:"name"`
	TxHash       string                 `json:"txHash"`
	TxIndex      string                 `json:"txIndex"`
	Gas          uint64                 `json:"gas"`
	GasPrice     *big.Int               `json:"gasPrice"`
	GasTipCap    *big.Int               `json:"gasTipCap"`
	GasFeeCap    *big.Int               `json:"gasFeeCap"`
	Value        string                 `json:"value"`
	Nonce        uint64                 `json:"nonce"`
	ToAddress    string                 `json:"toToAddress"`
}

type EventHandler func(*Event) error

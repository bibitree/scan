package model

import (
	"math/big"
)

type Event struct {
	Address          string                 `json:"address"`
	ContractName     string                 `json:"contractName"`
	ChainID          *big.Int               `json:"chainID"`
	Data             map[string]interface{} `json:"data"`
	BlockHash        string                 `json:"blockHash"`
	BlockNumber      string                 `json:"blockNumber"`
	Name             string                 `json:"name"`
	TxHash           string                 `json:"txHash"`
	TxIndex          string                 `json:"txIndex"`
	Gas              uint64                 `json:"gas"`
	GasPrice         *big.Int               `json:"gasPrice"`
	GasTipCap        *big.Int               `json:"gasTipCap"`
	GasFeeCap        *big.Int               `json:"gasFeeCap"`
	Value            string                 `json:"value"`
	Nonce            uint64                 `json:"nonce"`
	ToAddress        string                 `json:"toToAddress"`
	Status           bool                   `json:"status"`
	Timestamp        uint64                 `json:"timestamp"`
	MinerAddress     string                 `json:"minerAddress"`
	Size             string                 `json:"size"`
	BlockReward      string                 `json:"blockReward"`
	AverageGasTipCap string                 `json:"averageGasTipCap"`
	BlockBeasReward  string                 `json:"blockBeasReward"`
	NewAddress       string                 `json:"NewAddress"`
	NewToAddress     string                 `json:"NewToAddress"`
}

type ContractData struct {
	ContractName string                 `json:"contractName"`
	Data         map[string]interface{} `json:"data"`
	Name         string                 `json:"name"`
	TxHash       string                 `json:"txHash"`
	Contrac      string                 `json:"to"`
}

type BlockData struct {
	BlockHash    string `json:"blockHash"`
	BlockNumber  string `json:"blockNumber"`
	BlockReward  string `json:"blockReward"`
	MinerAddress string `json:"minerAddress"`
	Size         string `json:"size"`
	Timestamp    uint64 `json:"timestamp"`
}

type EventData struct {
	Address      string   `json:"address"`
	ChainID      *big.Int `json:"chainID"`
	BlockHash    string   `json:"blockHash"`
	BlockNumber  string   `json:"blockNumber"`
	TxHash       string   `json:"txHash"`
	TxIndex      string   `json:"txIndex"`
	Gas          uint64   `json:"gas"`
	GasPrice     *big.Int `json:"gasPrice"`
	GasTipCap    *big.Int `json:"gasTipCap"`
	GasFeeCap    *big.Int `json:"gasFeeCap"`
	Value        string   `json:"value"`
	Nonce        uint64   `json:"nonce"`
	To           string   `json:"to"`
	Status       bool     `json:"status"`
	Timestamp    uint64   `json:"timestamp"`
	NewAddress   string   `json:"NewAddress"`
	NewToAddress string   `json:"NewToAddress"`
}

type ErcTop struct {
	ContractAddress    string `json:"contractAddress"`
	ContractName       string `json:"contractName"`
	Value              string `json:"value"`
	NewContractAddress string `json:"nonce"`
	ContractTxCount    string `json:"contractTxCount"`
}

type EventHandler func(*Event) error

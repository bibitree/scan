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
	EventName    string                 `json:"eventName"`
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
	GasLimit     uint64 `json:"gasLimit"`
}

type BlockData2 struct {
	BlockHash       string `json:"blockHash"`
	BlockNumber     string `json:"blockNumber"`
	BlockReward     string `json:"blockReward"`
	AllBlockReward  uint64 `json:"AllBlockReward"`
	MinerAddress    string `json:"minerAddress"`
	Size            string `json:"size"`
	Timestamp       uint64 `json:"timestamp"`
	BlockBeasReward string `json:"blockBeasReward"`
	Count           int    `json:"count"`
	GasLimit        uint64 `json:"gasLimit"`
}

type BlockData3 struct {
	BlockHash       string `json:"blockHash"`
	BlockNumber     uint64 `json:"blockNumber"`
	BlockReward     uint64 `json:"blockReward"`
	AllBlockReward  uint64 `json:"AllBlockReward"`
	MinerAddress    string `json:"minerAddress"`
	Size            string `json:"size"`
	Timestamp       int    `json:"timestamp"`
	BlockBeasReward string `json:"blockBeasReward"`
	Count           int    `json:"count"`
	GasLimit        int    `json:"gasLimit"`
}

type AddressData struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
	Count   uint64 `json:"count"`
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
	Value        *big.Int `json:"value"`
	Nonce        uint64   `json:"nonce"`
	To           string   `json:"to"`
	Status       bool     `json:"status"`
	Timestamp    uint64   `json:"timestamp"`
	NewAddress   string   `json:"NewAddress"`
	NewToAddress string   `json:"NewToAddress"`
}

type ErcTop struct {
	ContractAddress    string   `json:"contractAddress"`
	ContractName       string   `json:"contractName"`
	Value              *big.Int `json:"value"`
	NewContractAddress string   `json:"nonce"`
	ContractTxCount    string   `json:"contractTxCount"`
	Decimals           int      `json:"decimals"`
	Symbol             string   `json:"symbol"`
}

type EventHandler func(*Event) error

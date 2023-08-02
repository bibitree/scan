package sniffer

import (
	"encoding/json"
	"errors"
	"ethgo/model"
	"ethgo/util/bigx"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Event struct {
	Address          common.Address         `json:"address"`
	ContractName     string                 `json:"contractName"`
	ChainID          *big.Int               `json:"chainID"`
	Data             map[string]interface{} `json:"data"`
	BlockHash        common.Hash            `json:"blockHash"`
	BlockNumber      uint64                 `json:"blockNumber"`
	Name             string                 `json:"name"`
	TxHash           common.Hash            `json:"txHash"`
	TxIndex          string                 `json:"txIndex"`
	Gas              uint64                 `json:"gas"`
	GasPrice         *big.Int               `json:"gasPrice"`
	GasTipCap        *big.Int               `json:"gasTipCap"`
	GasFeeCap        *big.Int               `json:"gasFeeCap"`
	TransactionFee   *big.Int               `json:"transactionFee"`
	Value            string                 `json:"value"`
	Nonce            uint64                 `json:"nonce"`
	To               common.Address         `json:"to"`
	Status           bool                   `json:"status"`
	Timestamp        uint64                 `json:"timestamp"`
	MinerAddress     string                 `json:"minerAddress"`
	Size             string                 `json:"size"`
	BlockReward      string                 `json:"blockReward"`
	AverageGasTipCap string                 `json:"averageGasTipCap"`
	NewAddress       string                 `json:"newAddress"`
	NewToAddress     string                 `json:"newToAddress"`
	GasLimit         uint64                 `json:"gasLimit"`
	Bytecode         []byte                 `json:"bytecode"`
	ContractAddr     common.Address         `json:"contractAddr"`
}

func (e *Event) String() string {
	return e.TxHash.Hex()
}

type ContractData struct {
	ContractName string                 `json:"contractName"`
	EventName    string                 `json:"eventName"`
	Data         map[string]interface{} `json:"data"`
	Name         string                 `json:"name"`
	TxHash       common.Hash            `json:"txHash"`
	Contrac      common.Address         `json:"to"`
}

type CreateContractData struct {
	Bytecode     string `json:"bytecode"`
	ContractAddr string `json:"contractAddr"`
	Code         string `json:"Code"`
	Abi          string `json:"abi"`
	Icon         string `json:"icon"`
	Time         int    `json:"timestamp"`
}

type SetCreateContractData struct {
	Bytecode     []byte `json:"bytecode"`
	ContractAddr string `json:"contractAddr"`
	Code         string `json:"Code"`
	Abi          string `json:"abi"`
	Icon         string `json:"icon"`
	Time         int    `json:"timestamp"`
}

type GetCreateContractIconData struct {
	Bytecode       []byte `json:"bytecode"`
	ContractAddr   string `json:"contractAddr"`
	BytecodeString string `json:"bytecodeString"`
	Icon           string `json:"icon"`
	Time           int    `json:"timestamp"`
}

type AddressData struct {
	Address string   `json:"address"`
	Balance *big.Int `json:"balance"`
}

type BalanceResponse struct {
	// 钱包地址
	Address string `json:"address" example:"0x51E72BDbA3A6Fc6337251581CB95625fa3A7767F"`
	// 余额（WEI）
	Wei string `json:"wei" example:"49335849638413224831"`
}

type BlockData struct {
	BlockHash    common.Hash `json:"blockHash"`
	BlockNumber  *big.Int    `json:"blockNumber"`
	BlockReward  *big.Int    `json:"blockReward"`
	MinerAddress string      `json:"minerAddress"`
	Size         string      `json:"size"`
	Timestamp    int         `json:"timestamp"`
	GasLimit     int         `json:"gasLimit"`
}

type EventData struct {
	Address        common.Address `json:"address"`
	ChainID        int            `json:"chainID"`
	BlockHash      common.Hash    `json:"blockHash"`
	BlockNumber    *big.Int       `json:"blockNumber"`
	TxHash         common.Hash    `json:"txHash"`
	TxIndex        int            `json:"txIndex"`
	Gas            *big.Int       `json:"gas"`
	GasPrice       *big.Int       `json:"gasPrice"`
	GasTipCap      *big.Int       `json:"gasTipCap"`
	GasFeeCap      *big.Int       `json:"gasFeeCap"`
	TransactionFee *big.Int       `json:"transactionFee"`
	Value          string         `json:"value"`
	Nonce          *big.Int       `json:"nonce"`
	To             common.Address `json:"to"`
	Status         int            `json:"status"`
	Timestamp      int            `json:"timestamp"`
	NewAddress     string         `json:"NewAddress"`
	NewToAddress   string         `json:"NewToAddress"`
}

type Event2 struct {
	Address          common.Address `json:"address"`
	ContractName     string         `json:"contractName"`
	ChainID          *big.Int       `json:"chainID"`
	Data             []byte         `json:"data"`
	BlockHash        common.Hash    `json:"blockHash"`
	BlockNumber      uint64         `json:"blockNumber"`
	Name             string         `json:"name"`
	TxHash           common.Hash    `json:"txHash"`
	TxIndex          string         `json:"txIndex"`
	Gas              uint64         `json:"gas"`
	GasPrice         *big.Int       `json:"gasPrice"`
	GasTipCap        *big.Int       `json:"gasTipCap"`
	GasFeeCap        *big.Int       `json:"gasFeeCap"`
	TransactionFee   *big.Int       `json:"transactionFee"`
	Value            string         `json:"value"`
	Nonce            uint64         `json:"nonce"`
	To               common.Address `json:"to"`
	Status           bool           `json:"status"`
	Timestamp        uint64         `json:"timestamp"`
	MinerAddress     string         `json:"minerAddress"`
	Size             string         `json:"size"`
	BlockReward      string         `json:"blockReward"`
	AverageGasTipCap string         `json:"averageGasTipCap"`
	GasLimit         uint64         `json:"gasLimit"`
	Bytecode         []byte         `json:"bytecode"`
	ContractAddr     common.Address `json:"contractAddr"`
}

func (e *Event2) ToEventData() *model.EventData {
	return &model.EventData{
		Address:        e.Address.Hex(),
		ChainID:        e.ChainID,
		BlockHash:      e.BlockHash.Hex(),
		BlockNumber:    big.NewInt(0).SetUint64(e.BlockNumber).String(),
		TxHash:         e.TxHash.Hex(),
		TxIndex:        e.TxIndex,
		Gas:            e.Gas,
		GasPrice:       e.GasPrice,
		GasTipCap:      e.GasTipCap,
		GasFeeCap:      e.GasFeeCap,
		TransactionFee: e.TransactionFee,
		Value:          bigx.FromString(e.Value),
		Nonce:          e.Nonce,
		To:             e.To.Hex(),
		Status:         e.Status,
		Timestamp:      e.Timestamp,
		NewAddress:     e.Address.Hex(),
		NewToAddress:   e.To.Hex(),
	}
}

type Event3 struct {
	Address        common.Address `json:"address"`
	ContractName   string         `json:"contractName"`
	ChainID        *big.Int       `json:"chainID"`
	Data           []byte         `json:"data"`
	BlockHash      common.Hash    `json:"blockHash"`
	BlockNumber    string         `json:"blockNumber"`
	Name           string         `json:"name"`
	TxHash         common.Hash    `json:"txHash"`
	TxIndex        string         `json:"txIndex"`
	Gas            uint64         `json:"gas"`
	GasPrice       *big.Int       `json:"gasPrice"`
	GasTipCap      *big.Int       `json:"gasTipCap"`
	GasFeeCap      *big.Int       `json:"gasFeeCap"`
	TransactionFee *big.Int       `json:"transactionFee"`
	Value          string         `json:"value"`
	Nonce          uint64         `json:"nonce"`
	To             common.Address `json:"to"`
	Status         bool           `json:"status"`
	Timestamp      uint64         `json:"timestamp"`
	MinerAddress   string         `json:"minerAddress"`
	Size           string         `json:"size"`
	BlockReward    string         `json:"blockReward"`
	Time           int64          `json:"time"`
}

type ChainData struct {
	// 合约地址
	BlockRewards            string `json:"blockRewards" example:"0"`
	SuperNodes              uint64 `json:"superNodes" example:"100"`
	BlockHeight             string `json:"blockHeight" example:"49335849638413224831"`
	TotalBlockRewards       string `json:"totalBlockRewards" example:"49335849638413224831"`
	GasPriceGasPrice        string `json:"gasPriceGasPrice" example:"49335849638413224831"`
	TotalnumberofAddresses  string `json:"totalnumberofAddresses" example:"49335849638413224831"`
	NumberTransactions      string `json:"numberTransactions" example:"49335849638413224831"`
	NumberTransfers         string `json:"numberTransfers" example:"49335849638413224831"`
	NumberTransactionsIn24H string `json:"numberTransactionsIn24H" example:"49335849638413224831"`
	NumberaddressesIn24H    string `json:"numberaddressesIn24H" example:"49335849638413224831"`
}

type ErcTop struct {
	ContractAddress    string   `json:"contractAddress"`
	ContractName       string   `json:"contractName"`
	Value              *big.Int `json:"value"`
	NewContractAddress string   `json:"nonce"`
	ContractTxCount    int      `json:"contractTxCount"`
	Decimals           int      `json:"decimals"`
	Symbol             string   `json:"symbol"`
}

func (event *Event) IsEmpty() bool {
	if (event.Address != common.Address{}) ||
		(event.ContractName != "") ||
		(event.ChainID != nil) ||
		(len(event.Data) > 0) ||
		(event.BlockHash != common.Hash{}) ||
		(event.BlockNumber != 0) ||
		(event.Name != "") ||
		(event.TxIndex != "") ||
		(event.Gas != 0) ||
		(event.GasPrice != nil) ||
		(event.GasTipCap != nil) ||
		(event.GasFeeCap != nil) ||
		(event.Value != "") ||
		(event.Nonce != 0) ||
		(event.To != common.Address{}) {
		return false
	}
	return true
}

func SerializeMap(data map[string]interface{}) ([]byte, error) {
	for _, value := range data {
		switch value.(type) {
		case []byte, bool, float32, float64, int, int16, int32, int64, string, uint, uint16, uint32, uint64:
			// 支持序列化的数据类型
		default:
			return nil, errors.New("value type not support")
		}
	}
	return json.Marshal(data)
}

func DeserializeJsonToMap(jsonData []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type EventHandler func([]*Event) error

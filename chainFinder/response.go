package chainFinder

import "ethgo/model"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Contract struct {
	Contract string `json:"contract" example:"0x51E72BDbA3A6Fc6337251581CB95625fa3A7767F"`
}

type Call struct {
	Address string      `json:"address" example:"0xa19844250b2b37c8518cb837b58ffed67f2e915D"`
	Method  string      `json:"method" example:"getDNA"`
	Args    interface{} `json:"args" swaggertype:"object,string" example:"id:1020"`
}

type GasPrice struct {
}

type ContractTxCount struct {
	Contract string `json:"contract" example:"0x51E72BDbA3A6Fc6337251581CB95625fa3A7767F"`
	Count    string `json:"count" example:"49335849638413224831"`
}

type Balance struct {
	// 钱包地址
	Address string `json:"address" example:"0x51E72BDbA3A6Fc6337251581CB95625fa3A7767F"`
}

type Paginate struct {
	Event        []model.EventData    `json:"event" example:"100"`
	PageNumber   uint64               `json:"pageNumber" example:"49335849638413224831"`
	Decimals     string               `json:"decimals" example:"49335849638413224831"`
	CreationTime uint64               `json:"creationTime" example:"49335849638413224831"`
	Address      string               `json:"address" example:"49335849638413224831"`
	ContractData []model.ContractData `json:"contractData" example:"0"`
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

type EventData struct {
	// 合约地址
	ContractData []model.ContractData `json:"contractData" example:"0"`
	Event        []model.EventData    `json:"event" example:"100"`
	BlockData    []model.BlockData3   `json:"blockData" example:"100"`
	PageNumber   uint64               `json:"pageNumber" example:"49335849638413224831"`
}

type AddressData struct {
	// 合约地址
	AddressData []model.AddressData `json:"contractData" example:"0"`
	PageNumber  uint64              `json:"pageNumber" example:"49335849638413224831"`
}

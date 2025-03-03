package chainFinder

import (
	"ethgo/model"
	"ethgo/sniffer"

	"github.com/ethereum/go-ethereum/common"
)

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

type Balance2 struct {
	// 合约地址
	Address string `json:"address" example:"0xa19844250b2b37c8518cb837b58ffed67f2e915D"`
	// 方法名(大小写敏感)
	Balance string `json:"Balance" example:"getDNA"`
}

type Paginate struct {
	Event              []model.EventData            `json:"event" example:"100"`
	PageNumber         uint64                       `json:"pageNumber" example:"49335849638413224831"`
	Decimals           string                       `json:"decimals" example:"49335849638413224831"`
	CreationTime       uint64                       `json:"creationTime" example:"49335849638413224831"`
	Address            string                       `json:"address" example:"49335849638413224831"`
	ContractData       []model.ContractData         `json:"contractData" example:"0"`
	CreateContractData []sniffer.CreateContractData `json:"createContractData" example:"100"`
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
	ContractData    []model.ContractData `json:"contractData" example:"0"`
	ContractsData2  []model.ContractData `json:"contractData2" example:"0"`
	Event           []model.EventData    `json:"event" example:"100"`
	SupplementEvent []model.EventData    `json:"supplementEvent" example:"100"`
	BlockData       []model.BlockData3   `json:"blockData" example:"100"`
	AddressData     []model.AddressData  `json:"addressData" example:"0"`
	PageNumber      uint64               `json:"pageNumber" example:"49335849638413224831"`
	EventPageNumber uint64               `json:"eventPageNumber" example:"49335849638413224831"`
	Balance         []model.Balance2     `json:"balance" example:"0"`
	ETHBalance      string               `json:"eTHBalance" example:"0"`
	ContractNum     int                  `json:"contractNum" example:"0"`
}

type IsContractAddressResponse struct {
	// 钱包地址
	Address string `json:"address" example:"0x51E72BDbA3A6Fc6337251581CB95625fa3A7767F"`
	// 余额（WEI）
	IsContract bool `json:"isContract" example:"49335849638413224831"`
}

type CompareBytecodeAndSourceCode struct {
	BytecodeString   string `json:"bytecodeString" example:"0x51E72BDbA3A6Fc6337251581CB95625fa3A7767F"`
	Code             string `json:"code" example:"49335849638413224831"`
	SolcVersion      string `json:"solcVersion" example:"49335849638413224831"`
	OptimizationRuns int    `json:"optimizationRuns" example:"49335849638413224831"`
}

type ContractData struct {
	// 合约地址
	ErcTop             []model.ErcTop               `json:"ercTop" example:"0"`
	CreateContractData []sniffer.CreateContractData `json:"createContractData" example:"100"`
}

type CreateContractData struct {
	Bytecode       []byte         `json:"bytecode"`
	ContractAddr   common.Address `json:"contractAddr"`
	BytecodeString string         `json:"bytecodeString"`
	Icon           string         `json:"icon"`
}

type AddressData struct {
	// 合约地址
	AddressData []model.AddressData `json:"contractData" example:"0"`
	PageNumber  uint64              `json:"pageNumber" example:"49335849638413224831"`
}

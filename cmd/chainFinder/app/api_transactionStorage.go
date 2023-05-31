package app

import (
	"encoding/json"
	"ethgo/model/mysqlOrders"
	"ethgo/model/orders"
	"ethgo/proto"
	"ethgo/sniffer"
	"ethgo/util/ginx"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// Transact
// @Description 缓存数据
// @Accept application/json
// @Produce application/json
// @Param object body proto.Transact{args=object} true "请求参数"
// @Success 200 {object} proto.Response{data=object}
// @Router /tyche/api/transact [post]
func (app *App) AcceptTransactionStorage(c *ginx.Context) {
	var request = new(proto.Event)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}

	var data, err = json.MarshalIndent(request.Data, "", "  ")
	if err != nil {
		panic(err)
	}

	request.Address = strings.ToLower(request.Address)
	if !common.IsHexAddress(request.Address) {
		c.Failure(http.StatusBadRequest, "无效的参数: address", nil)
		return
	}

	if request.TxIndex == "" {
		request.TxIndex = "0"
	}

	var event = sniffer.Event2{
		Address:          common.HexToAddress(request.Address),
		ContractName:     request.ContractName,
		ChainID:          request.ChainID,
		Data:             data,
		BlockHash:        request.BlockHash,
		BlockNumber:      request.BlockNumber,
		Name:             request.Name,
		TxHash:           request.TxHash,
		TxIndex:          request.TxIndex,
		Gas:              request.Gas,
		GasPrice:         request.GasPrice,
		GasTipCap:        request.GasTipCap,
		GasFeeCap:        request.GasFeeCap,
		Value:            request.Value,
		Nonce:            request.Nonce,
		To:               common.HexToAddress(request.To),
		Status:           request.Status,
		Timestamp:        request.Timestamp,
		MinerAddress:     request.MinerAddress,
		Size:             request.Size,
		BlockReward:      request.BlockReward,
		AverageGasTipCap: request.AverageGasTipCap,
		GasLimit:         request.GasLimit,
	}

	orders.CreateTransactionStorage(event)

	txhash := request.TxHash.String()
	if txhash != "0x0000000000000000000000000000000000000000000000000000000000000000" {
		orders.CreateTransactionTOPStorage(event)
	}
	app.SetChainData(c)
	c.Success(http.StatusOK, "succ", event)
}

func (app *App) SetChainData(c *ginx.Context) {

	blockHeight, gasPriceGasPrice, err := mysqlOrders.GetLatestEvent()
	if err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	numberTransactions, numberTransfers, numberTransactionsIn24H, numberaddressesIn24H, totalnumberofAddresses, err := mysqlOrders.GetEventStatistics(blockHeight)
	if err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}

	totalBlockRewards, err := mysqlOrders.GetAllAddressesAndBlockRewardSum(blockHeight)
	if err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	paginate := sniffer.ChainData{
		BlockRewards:            "10",
		SuperNodes:              100,
		BlockHeight:             blockHeight,
		GasPriceGasPrice:        gasPriceGasPrice,
		NumberTransactions:      numberTransactions,
		NumberTransfers:         numberTransfers,
		NumberTransactionsIn24H: numberTransactionsIn24H,
		NumberaddressesIn24H:    numberaddressesIn24H,
		TotalBlockRewards:       totalBlockRewards,
		TotalnumberofAddresses:  totalnumberofAddresses,
	}
	//CreateChainDataStorag
	orders.CreateChainDataStorag(paginate)
}

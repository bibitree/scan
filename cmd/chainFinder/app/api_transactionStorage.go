package app

import (
	"encoding/json"
	"errors"
	"ethgo/chainFinder"
	"ethgo/model/mysqlOrders"
	"ethgo/model/orders"
	"ethgo/proto"
	"ethgo/sniffer"
	"ethgo/util"
	"ethgo/util/ginx"
	"fmt"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

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
		Bytecode:         request.Bytecode,
		ContractAddr:     request.ContractAddr,
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

	blockHeight, _, err := mysqlOrders.GetLatestEvent()
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
	gasPriceGasPrice, err := app.ProcessGasPrice(c)
	if err != nil {
		log.Error(err)
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

func (app *App) ProcessGasPrice(c *ginx.Context) (string, error) {

	var gasPrice = chainFinder.GasPrice{}
	body, err := util.Post(app.conf.ChainFinder.GetGasPrice, gasPrice)
	if err != nil {
		return "", err
	}

	var res Response
	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}

	log.Debugf("应答: %v", string(body))

	if res.Code != http.StatusOK {
		if res.Message == "" {
			res.Message = fmt.Sprintf("%v", res.Code)
		}
		return "", errors.New(res.Message)
	}

	return res.Data.(string), nil
}

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
	"time"
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
	// var events = new(proto.Evensts)
	data, err := c.GetRawData()
	if err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}

	log.Debug(string(data))
	var events []*proto.Event
	if err := json.Unmarshal(data, &events); err != nil {
		log.Error(err)
	} else {
		log.Debug(events)
	}
	// if err := c.BindJSONEx1(data, events); err != nil {
	// 	c.Failure(http.StatusBadRequest, err.Error(), nil)
	// 	return
	// }

	for _, ev := range events {
		var data, err = json.MarshalIndent(ev.Data, "", "  ")
		if err != nil {
			panic(err)
		}
		if ev.TxIndex == "" {
			ev.TxIndex = "0"
		}

		var event = sniffer.Event2{
			Address:          ev.Address,
			ContractName:     ev.ContractName,
			ChainID:          ev.ChainID,
			Data:             data,
			BlockHash:        ev.BlockHash,
			BlockNumber:      ev.BlockNumber,
			Name:             ev.Name,
			TxHash:           ev.TxHash,
			TxIndex:          ev.TxIndex,
			Gas:              ev.Gas,
			GasPrice:         ev.GasPrice,
			GasTipCap:        ev.GasTipCap,
			GasFeeCap:        ev.GasFeeCap,
			TransactionFee:   ev.TransactionFee,
			Value:            ev.Value,
			Nonce:            ev.Nonce,
			To:               ev.To,
			Status:           ev.Status,
			Timestamp:        ev.Timestamp,
			MinerAddress:     ev.MinerAddress,
			Size:             ev.Size,
			BlockReward:      ev.BlockReward,
			AverageGasTipCap: ev.AverageGasTipCap,
			GasLimit:         ev.GasLimit,
			Bytecode:         ev.Bytecode,
			ContractAddr:     ev.ContractAddr,
		}

		orders.CreateTransactionStorage(event)
		txhash := ev.TxHash.String()
		if txhash != "0x0000000000000000000000000000000000000000000000000000000000000000" {
			orders.CreateTransactionTOPStorage(event)
		}
	}
	c.Success(http.StatusOK, "succ", nil)
}

func (app *App) setChainDataLoop() {
	for {
		app.SetChainData()
		time.Sleep(time.Second * 5)
	}
}

func (app *App) SetChainData() {
	log.Debugf("ENTER @TransactionStorage SetChainData")
	blockHeight, _, err := mysqlOrders.GetLatestEvent()
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("ENTER @TransactionStorage SetChainData blockHeight")
	numberTransactions, numberTransfers, numberTransactionsIn24H, numberaddressesIn24H, totalnumberofAddresses, err := mysqlOrders.GetEventStatistics(blockHeight)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("ENTER @TransactionStorage SetChainData numberTransactions")
	totalBlockRewards, err := mysqlOrders.GetAllAddressesAndBlockRewardSum(blockHeight)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("ENTER @TransactionStorage SetChainData totalBlockRewards")
	gasPriceGasPrice, err := app.ProcessGasPrice()
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

func (app *App) ProcessGasPrice() (string, error) {

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

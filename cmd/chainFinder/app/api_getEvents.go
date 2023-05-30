package app

import (
	"encoding/json"
	"errors"
	"ethgo/chainFinder"
	"ethgo/model/mysqlOrders"
	"ethgo/model/orders"
	"ethgo/proto"
	"ethgo/util"
	"ethgo/util/ginx"
	"fmt"
	"net/http"
)

func (app *App) GetAllEvents(c *ginx.Context) {
	messageReader, err := orders.NewTransactionTOPStorageReader(TRANSACT_CONSUMER_GROUP_NAME, TRANSACT_CONSUMER_NAME)
	if err != nil {
		panic(err)
	}

	msgs, err := messageReader.Read()
	if err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
	}
	c.Success(http.StatusOK, "succ", msgs)
}

func (app *App) GetEventsByBlockNumbers(c *ginx.Context) {
	var request = new(proto.ByBlockNumbers)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	events, Contract, n, err := mysqlOrders.GetEventsBetweenBlockNumbers(uint64(request.Star), uint64(request.End), uint64(request.PageNo), uint64(request.PageSize))
	if err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
	}
	Contracts, err := mysqlOrders.GetEventsByTxHashes(Contract)
	if err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
	}
	paginate := chainFinder.EventData{
		Event:        events,
		PageNumber:   n,
		ContractData: Contracts,
	}
	c.Success(http.StatusOK, "succ", paginate)
}

func (app *App) GetEventsByBlockNumber(c *ginx.Context) {
	var request = new(proto.ByBlockNumber)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	events, Contract, err := mysqlOrders.GetEventByBlockNumber(uint64(request.BlockNumber))
	if err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
	}
	Contracts, err := mysqlOrders.GetEventsByTxHashes(Contract)
	if err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
	}
	eventData := chainFinder.EventData{
		ContractData: Contracts,
		Event:        events,
	}
	c.Success(http.StatusOK, "succ", eventData)
}

func (app *App) GetEventsByTxHash(c *ginx.Context) {
	var request = new(proto.ByTxHash)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	events, Contract, err := mysqlOrders.GetEventByTxHash(request.TxHash)
	if err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
	}
	Contracts, err := mysqlOrders.GetEventsByTxHashes(Contract)
	if err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
	}
	eventData := chainFinder.EventData{
		ContractData: Contracts,
		Event:        events,
	}
	c.Success(http.StatusOK, "succ", eventData)
}

func (app *App) GetEventsByBlockHash(c *ginx.Context) {
	var request = new(proto.ByBlockHash)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	events, Contract, err := mysqlOrders.GetEventByBlockHash(request.BlockHash)
	if err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
	}
	Contracts, err := mysqlOrders.GetEventsByTxHashes(Contract)
	if err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
	}
	eventData := chainFinder.EventData{
		ContractData: Contracts,
		Event:        events,
	}
	c.Success(http.StatusOK, "succ", eventData)
}

func (app *App) GetERCTop(c *ginx.Context) {
	var request = new(proto.AllEvents)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	events, err := mysqlOrders.GetErcTopData(uint64(request.N))
	if err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
	}
	c.Success(http.StatusOK, "succ", events)
}

func (app *App) GetEventsByContract(c *ginx.Context) {
	var request = new(proto.ByContract)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	events, Contract, n, err := mysqlOrders.GetEventsByToAddressAndBlockNumber(request.Contract, uint64(request.PageNo), uint64(request.PageSize))
	if err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
	}
	Contracts, err := mysqlOrders.GetEventsByTxHashes(Contract)
	if err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
	}
	decimals, err := app.ProcessCall(request.Contract, "decimals")
	if err != nil {
		c.Success(http.StatusOK, "succ", nil)
		return
	}
	if decimals == nil {
		c.Success(http.StatusOK, "succ", nil)
		return
	}
	decimalsData := decimals.([]interface{})
	decimalsFloat64 := decimalsData[0].(float64)
	decimalsString := fmt.Sprintf("%.0f", decimalsFloat64)

	address, time, err := mysqlOrders.GetEventAddressByToAddress(request.Contract)
	if err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
	}
	paginate := chainFinder.Paginate{
		Event:        events,
		PageNumber:   n,
		Decimals:     decimalsString,
		CreationTime: time,
		Address:      address,
		ContractData: Contracts,
	}

	c.Success(http.StatusOK, "succ", paginate)
}

func (app *App) ProcessCall(contract string, name string) (interface{}, error) {

	var call = chainFinder.Call{
		Address: contract,
		Method:  name,
	}
	body, err := util.Post(app.conf.ChainFinder.Callback, call)
	if err != nil {
		return nil, err
	}

	var res chainFinder.Response
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	if res.Code != http.StatusOK {
		if res.Message == "" {
			res.Message = fmt.Sprintf("%v", res.Code)
		}
		return nil, errors.New(res.Message)
	}

	return res.Data, nil
}

// func (app *App) ProcessERCContractTxCount(contract string) (interface{}, error) {
// 	var call = chainFinder.Contract{
// 		Contract: contract,
// 	}
// 	body, err := util.Post(app.conf.ChainFinder.ContractTxCount, call)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var res chainFinder.Response
// 	if err := json.Unmarshal(body, &res); err != nil {
// 		return nil, err
// 	}

// 	log.Debugf("应答: %v", string(body))

// 	if res.Code != http.StatusOK {
// 		if res.Message == "" {
// 			res.Message = fmt.Sprintf("%v", res.Code)
// 		}
// 		return nil, errors.New(res.Message)
// 	}

// 	return res.Data, nil
// }

func (app *App) ProcessContractCreationTime(contract string) (interface{}, error) {

	var call = chainFinder.Contract{
		Contract: contract,
	}
	body, err := util.Post(app.conf.ChainFinder.ContractCreationTime, call)
	if err != nil {
		return nil, err
	}

	var res chainFinder.Response
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	log.Debugf("应答: %v", string(body))

	if res.Code != http.StatusOK {
		if res.Message == "" {
			res.Message = fmt.Sprintf("%v", res.Code)
		}
		return nil, errors.New(res.Message)
	}

	return res.Data, nil
}

func (app *App) GetChainData(c *ginx.Context) {
	messageReader, err := orders.NewCreateChainDataStoragReader(TRANSACT_CONSUMER_GROUP_NAME, TRANSACT_CONSUMER_NAME)
	if err != nil {
		panic(err)
	}

	msgs, err := messageReader.Read()
	if err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
	}
	c.Success(http.StatusOK, "succ", msgs)
}

func (app *App) GetBlockNum(c *ginx.Context) {
	var request = new(proto.ByContract)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}

	blockHeight, _, err := mysqlOrders.GetLatestEvent()
	if err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	paginate := chainFinder.ChainData{
		BlockHeight: blockHeight,
	}

	c.Success(http.StatusOK, "succ", paginate)
}

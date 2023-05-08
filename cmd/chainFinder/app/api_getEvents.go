package app

import (
	"ethgo/model/mysqlOrders"
	"ethgo/proto"
	"ethgo/util/ginx"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
)

func (app *App) GetAllEvents(c *ginx.Context) {
	var request = new(proto.AllEvents)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	events, err := mysqlOrders.GetAllEvents()
	if err != nil {
		log.Fatal(err)
	}
	c.Success(http.StatusOK, "succ", events)
}

func (app *App) GetEventsByBlockNumbers(c *ginx.Context) {
	var request = new(proto.ByBlockNumbers)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	events, err := mysqlOrders.GetEventsBetweenBlockNumbers(uint64(request.Star), uint64(request.End))
	if err != nil {
		log.Fatal(err)
	}
	c.Success(http.StatusOK, "succ", events)
}

func (app *App) GetEventsByBlockNumber(c *ginx.Context) {
	var request = new(proto.ByBlockNumber)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	events, err := mysqlOrders.GetEventByBlockNumber(uint64(request.BlockNumber))
	if err != nil {
		log.Fatal(err)
	}
	c.Success(http.StatusOK, "succ", events)
}

func (app *App) GetEventsByTxHash(c *ginx.Context) {
	var request = new(proto.ByTxHash)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	events, err := mysqlOrders.GetEventByTxHash(common.HexToHash(request.TxHash))
	if err != nil {
		log.Fatal(err)
	}
	c.Success(http.StatusOK, "succ", events)
}

func (app *App) GetEventsByBlockHash(c *ginx.Context) {
	var request = new(proto.ByBlockHash)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	events, err := mysqlOrders.GetEventByBlockHash(common.HexToHash(request.BlockHash))
	if err != nil {
		log.Fatal(err)
	}
	c.Success(http.StatusOK, "succ", events)
}

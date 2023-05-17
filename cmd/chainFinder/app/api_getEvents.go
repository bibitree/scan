package app

import (
	"ethgo/model/mysqlOrders"
	"ethgo/proto"
	"ethgo/util/ginx"
	"net/http"
)

func (app *App) GetAllEvents(c *ginx.Context) {

	var request = new(proto.AllEvents)
	log.Info("请求内容", request.N)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	events, err := mysqlOrders.GetAllEvents(uint64(request.N))
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
	events, err := mysqlOrders.GetEventsBetweenBlockNumbers(uint64(request.Star), uint64(request.End), uint64(request.PageNo), uint64(request.PageSize))
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
	events, err := mysqlOrders.GetEventByTxHash(request.TxHash)
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
	events, err := mysqlOrders.GetEventByBlockHash(request.BlockHash)
	if err != nil {
		log.Fatal(err)
	}
	c.Success(http.StatusOK, "succ", events)
}

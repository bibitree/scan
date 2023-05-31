package app

import (
	"ethgo/util/ginx"

	"github.com/gin-gonic/gin"
)

func (app *App) Router(g *gin.Engine) {
	g.POST("/chainFinder/api/AcceptTransactionStorage", ginx.WrapHandler(app.AcceptTransactionStorage))
	g.POST("/chainFinder/api/GetAllEvents", ginx.WrapHandler(app.GetAllEvents))
	g.POST("/chainFinder/api/GetEventsByBlockNumbers", ginx.WrapHandler(app.GetEventsByBlockNumbers))
	g.POST("/chainFinder/api/GetEventsByBlockNumber", ginx.WrapHandler(app.GetEventsByBlockNumber))
	g.POST("/chainFinder/api/GetEventsByTxHash", ginx.WrapHandler(app.GetEventsByTxHash))
	g.POST("/chainFinder/api/GetEventsByBlockHash", ginx.WrapHandler(app.GetEventsByBlockHash))
	g.POST("/chainFinder/api/GetERCTop", ginx.WrapHandler(app.GetERCTop))
	g.POST("/chainFinder/api/GetEventsByContract", ginx.WrapHandler(app.GetEventsByContract))
	g.POST("/chainFinder/api/GetChainData", ginx.WrapHandler(app.GetChainData))
	g.POST("/chainFinder/api/GetBlockNum", ginx.WrapHandler(app.GetBlockNum))
	g.POST("/chainFinder/api/GetAddressList", ginx.WrapHandler(app.GetAddressList))
	g.POST("/chainFinder/api/GetEventsByAddress", ginx.WrapHandler(app.GetEventsByAddress))
}

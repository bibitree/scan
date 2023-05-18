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
}

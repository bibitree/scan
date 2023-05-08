package app

import (
	"ethgo/util/ginx"

	"github.com/gin-gonic/gin"
)

func (app *App) Router(g *gin.Engine) {
	g.POST("/tyche/api/AcceptTransactionStorage", ginx.WrapHandler(app.AcceptTransactionStorage))
	g.POST("/tyche/api/GetAllEvents", ginx.WrapHandler(app.GetAllEvents))
	g.POST("/tyche/api/GetEventsByBlockNumbers", ginx.WrapHandler(app.GetEventsByBlockNumbers))
	g.POST("/tyche/api/GetEventsByBlockNumber", ginx.WrapHandler(app.GetEventsByBlockNumber))
	g.POST("/tyche/api/GetEventsByTxHash", ginx.WrapHandler(app.GetEventsByTxHash))
	g.POST("/tyche/api/GetEventsByBlockHash", ginx.WrapHandler(app.GetEventsByBlockHash))
}

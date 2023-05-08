package app

import (
	"encoding/json"
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

	var bytes, err = json.MarshalIndent(request, "", "  ")
	if err != nil {
		panic(err)
	}
	log.Debugf("%v", string(bytes))

	request.Address = strings.ToLower(request.Address)
	if !common.IsHexAddress(request.Address) {
		c.Failure(http.StatusBadRequest, "无效的参数: address", nil)
		return
	}

	var event = sniffer.Event{
		Address:      common.HexToAddress(request.Address),
		ContractName: request.ContractName,
		ChainID:      request.ChainID,
		Data:         request.Data,
		BlockHash:    request.BlockHash,
		BlockNumber:  request.BlockNumber,
		Name:         request.Name,
		TxHash:       request.TxHash,
		TxIndex:      request.TxIndex,
	}

	orders.CreateTransactionStorage(event)

	// if err = orders.Pending(id, contract.Address.String(), hexutil.Encode(inputData)); err != nil {
	// 	log.Errorf("Failed to %v: %v, %v", c.Request.URL, err, transactor)
	// 	c.Failure(http.StatusInternalServerError, err.Error(), nil)
	// 	return
	// }

	// if err = app.base.Transact(context.Background(), request.OrderID, transactor); err != nil {
	// 	log.Errorf("Failed to %v: %v, %v", c.Request.URL, err, transactor)
	// 	c.Failure(http.StatusInternalServerError, err.Error(), nil)
	// 	return
	// }

	c.Success(http.StatusOK, "succ", event)
}

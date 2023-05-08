package chainFinder

import (
	"context"
	"encoding/json"
	"ethgo/model/mysqlOrders"
	"ethgo/model/orders"

	"ethgo/sniffer"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func (t *ChainFinder) WatchTransactionStorage(ctx context.Context) {
	messageReader, err := orders.NewTransactionStorageReader(TRANSACT_CONSUMER_GROUP_NAME, TRANSACT_CONSUMER_NAME)
	if err != nil {
		panic(err)
	}

	messageDispatcher := NewMessageDispatcher(messageReader, t.TransactionStorage, t.conf.SucceedNumberOfConcurrent)
	messageDispatcher.Run(ctx)
}

func (t *ChainFinder) TransactionStorage(ctx context.Context, message *orders.Message) AfterFunc {
	log.Debugf("ENTER @TransactionStorage 订单")
	defer log.Debugf("  LEAVE @TransactionStorage 订单")

	chainID, err := message.Int64("ChainID")
	if err != nil {
		// 发生错误，处理错误逻辑
		return After(t.conf.NetworkRetryInterval, message)
	}

	yourMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(message.String("Data")), &yourMap)
	if err != nil {
		// 处理错误
		return After(t.conf.NetworkRetryInterval, message)
	}

	var event = sniffer.Event{
		Address:      common.HexToAddress(message.String("Address")),
		ContractName: message.String("ContractName"),
		ChainID:      big.NewInt(chainID),
		Data:         yourMap,
		BlockHash:    common.HexToHash(message.String("BlockHash")),
		BlockNumber:  message.String("BlockNumber"),
		Name:         message.String("Name"),
		TxHash:       common.HexToHash(message.String("TxHash")),
		TxIndex:      message.String("TxIndex"),
	}
	log.Info(event)
	mysqlOrders.InsertEvent(event)
	mysqlOrders.InsertEventData(event)
	return t.ack(message)
}

package chainFinder

import (
	"context"
	"encoding/json"

	// "ethgo/model/mysqlOrders"
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
	gasPrice, err := message.Int64("GasPrice")
	if err != nil {
		// 发生错误，处理错误逻辑
		return After(t.conf.NetworkRetryInterval, message)
	}
	gasTipCap, err := message.Int64("GasTipCap")
	if err != nil {
		// 发生错误，处理错误逻辑
		return After(t.conf.NetworkRetryInterval, message)
	}
	gasFeeCap, err := message.Int64("GasFeeCap")
	if err != nil {
		// 发生错误，处理错误逻辑
		return After(t.conf.NetworkRetryInterval, message)
	}

	nonce, err := message.Uint64("Nonce")
	if err != nil {
		// 发生错误，处理错误逻辑
		return After(t.conf.NetworkRetryInterval, message)
	}

	gas, err := message.Uint64("Gas")
	if err != nil {
		// 发生错误，处理错误逻辑
		return After(t.conf.NetworkRetryInterval, message)
	}

	var yourMap map[string]interface{}
	data := message.Bytes("Data")
	if len(data) == 2 {
		yourMap = make(map[string]interface{})
	} else {
		err = json.Unmarshal(data, &yourMap)
		if err != nil {
			// 处理错误
			return After(t.conf.NetworkRetryInterval, message)
		}
	}

	timestamp, err := message.Uint64("Timestamp")
	if err != nil {
		// 发生错误，处理错误逻辑
		return After(t.conf.NetworkRetryInterval, message)
	}

	sizeStr := message.String("Size")

	var event = sniffer.Event{
		Address:          common.HexToAddress(message.String("Address")),
		ContractName:     message.String("ContractName"),
		ChainID:          big.NewInt(chainID),
		Data:             yourMap,
		BlockHash:        common.HexToHash(message.String("BlockHash")),
		BlockNumber:      message.String("BlockNumber"),
		Name:             message.String("Name"),
		TxHash:           common.HexToHash(message.String("TxHash")),
		TxIndex:          message.String("TxIndex"),
		Gas:              gas,
		GasPrice:         big.NewInt(gasPrice),
		GasTipCap:        big.NewInt(gasTipCap),
		GasFeeCap:        big.NewInt(gasFeeCap),
		Value:            message.String("Value"),
		Nonce:            nonce,
		To:               common.HexToAddress(message.String("To")),
		Status:           message.String("Status") == "true",
		Timestamp:        timestamp,
		MinerAddress:     message.String("MinerAddress"),
		Size:             sizeStr,
		BlockReward:      message.String("BlockReward"),
		AverageGasTipCap: message.String("AverageGasTipCap"),
	}
	log.Info(event)
	mysqlOrders.InsertEvent(event)
	// mysqlOrders.InsertEventData(event)
	return t.ack(message)
}

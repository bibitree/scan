package chainFinder

import (
	"context"
	"sync"
	"time"

	// "ethgo/model/mysqlOrders"
	"ethgo/model/mysqlOrders"

	"ethgo/model/orders"

	"ethgo/sniffer"
	"math/big"
)

func (t *ChainFinder) WatchBalanceStorage(ctx context.Context) {
	messageReader, err := orders.NewCreateBalanceReader(TRANSACT_CONSUMER_GROUP_NAME, TRANSACT_CONSUMER_NAME)
	if err != nil {
		panic(err)
	}

	messageDispatcher := NewMessageDispatcher(messageReader, t.BalanceStorage, t.conf.SucceedNumberOfConcurrent)
	messageDispatcher.Run(ctx)
}

var balanceUpCache sync.Map

func (t *ChainFinder) BalanceStorage(ctx context.Context, message *orders.Message) AfterFunc {
	log.Debugf("ENTER @BalanceStorage 订单")
	defer log.Debugf("  LEAVE @BalanceStorage 订单")
	log.Info(message.String("Address"))

	var balance = Balance{
		Address: message.String("Address"),
	}
	now := time.Now().Unix()
	if v, ok := balanceUpCache.Load(balance.Address); ok {
		lastUp := v.(int64)
		if now-lastUp < 10 {
			return t.ack(message)
		}
	}

	balanceSupply, err := t.ProcessBalance(balance)
	if err != nil {
		// log.Error(err)
		return nil
	}
	balanceSupplyDataMap := balanceSupply.(map[string]interface{})
	wei := balanceSupplyDataMap["wei"].(string)
	num := new(big.Int)
	num, _ = num.SetString(wei, 10)
	if err != nil {
		log.Error(err)
		return nil
	}

	var event = sniffer.AddressData{
		Address: message.String("Address"),
		Balance: num,
	}
	log.Info(event)

	mysqlOrders.InsertAddressData(event)
	balanceUpCache.Store(balance.Address, now)
	return t.ack(message)
}

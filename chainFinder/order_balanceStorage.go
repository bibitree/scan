package chainFinder

import (
	"context"

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

func (t *ChainFinder) BalanceStorage(ctx context.Context, message *orders.Message) AfterFunc {
	log.Debugf("ENTER @BalanceStorage 订单")
	defer log.Debugf("  LEAVE @BalanceStorage 订单")
	log.Info(message.String("Address"))

	var balance = Balance{
		Address: message.String("Address"),
	}

	balanceSupply, err := t.ProcessBalance(balance)
	if err != nil {
		log.Error(err)
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
	return t.ack(message)
}

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

func (t *ChainFinder) GetBalancesByUsers(userAddresses []string) (map[string]interface{}, error) {
	// 创建一个用于存储用户余额的 map
	balances := make(map[string]interface{})

	for _, address := range userAddresses {
		balance := Balance{
			Address: address,
		}

		// 调用 ProcessBalance 方法获取余额
		data, err := t.ProcessBalance(balance)
		if err != nil {
			log.Errorf("Failed to fetch balance for address %s: %v", address, err)
			continue // 跳过错误的地址，处理下一个
		}

		// 如果余额为空，跳过该地址
		if data == nil {
			log.Debugf("Address %s has no balance, skipping", address)
			continue
		}

		// 将余额存入 map
		balances[address] = data
	}

	return balances, nil
}

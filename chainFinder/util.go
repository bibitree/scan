package chainFinder

import (
	"ethgo/model/orders"
)

func (t *ChainFinder) ack(message *orders.Message) AfterFunc {
	if err := message.Ack(); err != nil {
		log.Errorf("Ack message failed: %v, %v", message.Source(), err)
		return After(t.conf.RedisRetryInterval, message)
	}
	return nil
}

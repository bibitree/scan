package tyche

import (
	"context"
	"ethgo/eth"
	"ethgo/model/orders"
	"ethgo/tyche/gaslimit"
	"ethgo/tyche/gasprice"
	"ethgo/util/bigx"
	"ethgo/util/ethx"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/garyburd/redigo/redis"
)

func (t *Tyche) watchPending(ctx context.Context) {
	messageReader, err := orders.NewPendingReader(PENDING_CONSUMER_GROUP_NAME, PENDING_CONSUMER_NAME)
	if err != nil {
		panic(err)
	}

	messageDispatcher := NewMessageDispatcher(messageReader, t.pending, t.conf.PendingNumberOfConcurrent)
	messageDispatcher.Run(ctx)
}

func (t *Tyche) pending(ctx context.Context, message *orders.Message) AfterFunc {
	var id = message.ID()
	var nonce = message.Nonce()
	var contractAddr = message.String("to")
	var inputData = message.String("inputData")
	if contractAddr == "" {
		panic("to field must not be empty")
	}
	if inputData == "" {
		panic("inputData field must not be empty")
	}

	distance := t.limiter.Distance(nonce)
	if distance >= t.conf.PendingNumberOfConcurrent*2 {
		// log.Debugf("The <PENDING> order will be suspended: %v, %v", id, nonce)
		return After(t.conf.PendingRetryInterval, message)
	}

	log.Debugf("ENTER @PENDING 订单: %v, %v, %v", id, nonce, t.limiter.max)
	defer log.Debugf("  LEAVE @PENDING 订单: %v, %v, %v", id, nonce, t.limiter.max)

	release := orders.Lock(id)
	defer release()

	order, err := orders.Get(id, orders.FIELD_STATUS, orders.FIELD_TX_HASH, orders.FIELD_TX_DATA)
	switch err {
	case nil:
		break
	case redis.ErrNil:
		panic(id)
	default:
		log.Errorf("Get order data failed: %v, %v, %v", id, nonce, err)
		return After(t.conf.RedisRetryInterval, message)
	}

	if order.Status != orders.PENDING_STATUS {
		log.Debugf("Not in pending status: %v, %v %v", id, nonce, order.Status)
		return t.ack(message)
	}

	var signedTx *types.Transaction
	if order.TxHash != "" {
		signedTx, err = ethx.Unmarshal(order.TxData)
		if err != nil {
			panic(err)
		}

		gasPrice := gasprice.Get()
		if bigx.Mul(signedTx.GasPrice(), 1.1).Cmp(gasPrice) >= 0 {
			log.Infof("The transaction will be reload: %v, %v, %v", id, nonce, signedTx.Hash())
		} else {
			log.Infof("The transaction will be bumping gas price: %v, %v, %v", id, nonce, signedTx.Hash())
			signedTx = nil
		}
	}

	if signedTx == nil {
		address := common.HexToAddress(contractAddr)
		baseTx := new(types.LegacyTx)
		baseTx.To = &address
		baseTx.Nonce = nonce
		baseTx.Gas = gaslimit.Get()
		baseTx.GasPrice = gasprice.Get()
		baseTx.Value = big.NewInt(0)
		baseTx.Data = hexutil.MustDecode(inputData)

		signedTx, err = t.signTx(types.NewTx(baseTx))
		if err != nil {
			panic(err)
		}

		bytes, err := signedTx.MarshalBinary()
		if err != nil {
			panic(err)
		}

		err = orders.Bind(id, signedTx.Hash().String(), hexutil.Encode(bytes))
		if err != nil {
			log.Errorf("bind transaction data failed: %v, %w", id, err)
			return After(t.conf.RedisRetryInterval, message)
		}
	}

SEND:
	if err = t.sendTransaction(ctx, id, signedTx); err != nil {
		if eth.IsRpcError(err) {
			// In this case, we will try again later
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(time.Second):
				goto SEND
			}
		}

		var modifier = orders.NewModifier(id, signedTx.Hash().String(), nonce)
		modifier.Set("address", contractAddr)
		modifier.Set("error", err.Error())
		modifier.Set("recovery", true)
		if err := orders.Error(modifier); err != nil {
			log.Errorf("Publish <ERROR> message failed: %v, %v", id, err)
			return After(t.conf.RedisRetryInterval, message)
		}

		return t.ack(message)
	}

	var modifier = orders.NewModifier(id, signedTx.Hash().String(), nonce)
	modifier.Set("sent", time.Now().Unix())
	if err := orders.Sent(modifier); err != nil {
		log.Errorf("Publish <SENT> message failed: %v, %v", id, err)
		return After(t.conf.RedisRetryInterval, message)
	}

	return t.ack(message)
}

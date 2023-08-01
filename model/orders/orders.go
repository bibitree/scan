package orders

import (
	"encoding/json"
	"errors"
	"ethgo/model"
	"ethgo/sniffer"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	FIELD_ID                = "id"
	FIELD_CREATED_AT        = "createdAt"
	FIELD_STATUS            = "status"
	FIELD_UPDATED_AT        = "updatedAt"
	FIELD_NUMBER_OF_RETRIES = "numberOfRetries"
	FIELD_NONCE             = "nonce"
	FIELD_TX_HASH           = "txHash"
	FIELD_TX_DATA           = "txData"
	FIELD_REASON            = "reason"
)

const (
	PENDING_STATUS = "pending"
	SENT_STATUS    = "sent"
	SUCCEED_STATUS = "succ"
	FAILED_STATUS  = "fail"
	ERROR_STATUS   = "error"
)

const (
	ERROR_ORDER_EXPIRED   = 15 * 60 * 60 * 24
	FAILED_ORDER_EXPIRED  = 15 * 60 * 60 * 24
	SUCCEED_ORDER_EXPIRED = 15 * 60 * 60 * 24
)

type Order struct {
	Id              string `redis:"id" json:"id"`
	Status          string `redis:"status" json:"status"`
	CreatedAt       int64  `redis:"createdAt" json:"createdAt"`
	UpdatedAt       int64  `redis:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	NumberOfRetries int64  `redis:"numberOfRetries,omitempty" json:"numberOfRetries,omitempty"`
	Nonce           uint64 `redis:"nonce" json:"nonce"`
	TxHash          string `redis:"txHash,omitempty" json:"txHash,omitempty"`
	TxData          string `redis:"txData,omitempty" json:"txData,omitempty"`
	Reason          string `redis:"reason,omitempty" json:"reason,omitempty"`
}

type SignerFunc func(order *Order) error

var nonceLock sync.Mutex

func Init(nonce uint64) error {
	var red = model.RedisPool.Get()
	defer red.Close()

	red.Send("SETNX", keys.NonceAt(), nonce)
	return red.Flush()
}

func Error(modifier *Modifier) error {
	var red = model.RedisPool.Get()
	defer red.Close()

	red.Send("MUTIL")
	{

		args := redis.Args{keys.Error()}
		args = args.Add("*")
		args = args.Add(modifier.Values()...)
		red.Send("XADD", args...)
	}

	{
		args := redis.Args{keys.Entity(modifier.ID())}
		args = args.Add(FIELD_STATUS, ERROR_STATUS)
		args = args.Add(FIELD_UPDATED_AT, time.Now().Unix())
		red.Send("HSET", args...)
	}

	red.Send("EXPIRE", keys.Entity(modifier.ID()), ERROR_ORDER_EXPIRED)
	red.Send("DEL", keys.Nonce(modifier.Nonce()))
	red.Send("EXEC")
	return red.Flush()
}

func Failed(modifier *Modifier) error {
	var red = model.RedisPool.Get()
	defer red.Close()

	red.Send("MUTIL")
	{
		args := redis.Args{keys.Failed()}
		args = args.Add("*")
		args = args.Add(modifier.Values()...)
		red.Send("XADD", args...)
	}

	{
		args := redis.Args{keys.Entity(modifier.ID())}
		args = args.Add(FIELD_STATUS, FAILED_STATUS)
		args = args.Add(FIELD_TX_HASH, modifier.Hash())
		args = args.Add(FIELD_UPDATED_AT, time.Now().Unix())
		red.Send("HSET", args...)
	}

	red.Send("EXPIRE", keys.Entity(modifier.ID()), FAILED_ORDER_EXPIRED)
	red.Send("DEL", keys.Nonce(modifier.Nonce()))
	red.Send("EXEC")
	return red.Flush()
}

func Succeed(modifier *Modifier) error {
	var red = model.RedisPool.Get()
	defer red.Close()

	red.Send("MUTIL")
	{
		args := redis.Args{keys.Succeed()}
		args = args.Add("*")
		args = args.Add(modifier.Values()...)
		red.Send("XADD", args...)
	}

	{
		args := redis.Args{keys.Entity(modifier.ID())}
		args = args.Add(FIELD_STATUS, SUCCEED_STATUS)
		args = args.Add(FIELD_TX_HASH, modifier.Hash())
		args = args.Add(FIELD_UPDATED_AT, time.Now().Unix())
		red.Send("HSET", args...)
	}

	red.Send("EXPIRE", keys.Entity(modifier.ID()), SUCCEED_ORDER_EXPIRED)
	red.Send("DEL", keys.Nonce(modifier.Nonce()))
	red.Send("EXEC")
	return red.Flush()
}

func Sent(modifier *Modifier) error {
	var red = model.RedisPool.Get()
	defer red.Close()

	red.Send("MUTIL")
	{
		args := redis.Args{keys.Sent()}
		args = args.Add("*")
		args = args.Add(modifier.Values()...)
		red.Send("XADD", args...)
	}

	{
		args := redis.Args{keys.Entity(modifier.ID())}
		args = args.Add(FIELD_STATUS, SENT_STATUS)
		args = args.Add(FIELD_UPDATED_AT, time.Now().Unix())
		red.Send("HSET", args...)
	}

	red.Send("EXEC")
	return red.Flush()
}

func CreateChainDataStorag(event sniffer.ChainData) error {
	red := model.RedisPool.Get()
	defer red.Close()

	chainDataKey := keys.CreateChainDataStorag()

	// check if the hash key exists
	exists, err := redis.Bool(red.Do("EXISTS", chainDataKey))
	if err != nil {
		return err
	}

	// if the hash key exists, update the data
	if exists {
		red.Send("HMSET", chainDataKey,
			"BlockRewards", event.BlockRewards,
			"SuperNodes", event.SuperNodes,
			"BlockHeight", event.BlockHeight,
			"TotalBlockRewards", event.TotalBlockRewards,
			"GasPriceGasPrice", event.GasPriceGasPrice,
			"TotalnumberofAddresses", event.TotalnumberofAddresses,
			"NumberTransactions", event.NumberTransactions,
			"NumberTransfers", event.NumberTransfers,
			"NumberTransactionsIn24H", event.NumberTransactionsIn24H,
			"NumberaddressesIn24H", event.NumberaddressesIn24H,
			"Time", time.Now().Unix(),
		)
	} else { // if not exists, create the hash key
		red.Send("HSET", chainDataKey, "BlockRewards", event.BlockRewards)
		red.Send("HSET", chainDataKey, "SuperNodes", event.SuperNodes)
		red.Send("HSET", chainDataKey, "BlockHeight", event.BlockHeight)
		red.Send("HSET", chainDataKey, "TotalBlockRewards", event.TotalBlockRewards)
		red.Send("HSET", chainDataKey, "GasPriceGasPrice", event.GasPriceGasPrice)
		red.Send("HSET", chainDataKey, "TotalnumberofAddresses", event.TotalnumberofAddresses)
		red.Send("HSET", chainDataKey, "NumberTransactions", event.NumberTransactions)
		red.Send("HSET", chainDataKey, "NumberTransfers", event.NumberTransfers)
		red.Send("HSET", chainDataKey, "NumberTransactionsIn24H", event.NumberTransactionsIn24H)
		red.Send("HSET", chainDataKey, "NumberaddressesIn24H", event.NumberaddressesIn24H)
		red.Send("HSET", chainDataKey, "Time", time.Now().Unix())
	}

	_, err = red.Do("EXEC")
	if err != nil {
		return err
	}

	return nil
}

func CreateTransactionStorage(event sniffer.Event2) error {
	var red = model.RedisPool.Get()
	defer red.Close()
	red.Send("MUTIL")
	{
		args := redis.Args{keys.CreateTransactionStorage()}
		args = args.Add("*")
		args = args.Add("Address", event.Address.String())
		args = args.Add("ContractName", event.ContractName)
		args = args.Add("ChainID", event.ChainID)
		args = args.Add("Data", event.Data)
		args = args.Add("BlockHash", event.BlockHash)
		args = args.Add("BlockNumber", event.BlockNumber)
		args = args.Add("Name", event.Name)
		args = args.Add("TxHash", event.TxHash)
		args = args.Add("TxIndex", event.TxIndex)
		args = args.Add("Gas", event.Gas)
		args = args.Add("GasPrice", event.GasPrice)
		args = args.Add("GasTipCap", event.GasTipCap)
		args = args.Add("GasFeeCap", event.GasFeeCap)
		args = args.Add("TransactionFee", event.TransactionFee)
		args = args.Add("Value", event.Value)
		args = args.Add("Nonce", event.Nonce)
		args = args.Add("To", event.To.String())
		args = args.Add("Status", event.Status)
		args = args.Add("Timestamp", event.Timestamp)
		args = args.Add("MinerAddress", event.MinerAddress)
		args = args.Add("Size", event.Size)
		args = args.Add("BlockReward", event.BlockReward)
		args = args.Add("AverageGasTipCap", event.AverageGasTipCap)
		args = args.Add("GasLimit", event.GasLimit)
		args = args.Add("GasLimit", event.GasLimit)
		args = args.Add("Bytecode", event.Bytecode)
		args = args.Add("ContractAddr", event.ContractAddr)
		args = args.Add("Time", time.Now().Unix())

		red.Send("XADD", args...)
	}
	red.Send("EXEC")
	return red.Flush()
}

func CreateBalance(event sniffer.Event2) error {
	var red = model.RedisPool.Get()
	defer red.Close()
	red.Send("MUTIL")
	{
		args := redis.Args{keys.CreateBalanceStorage()}
		args = args.Add("*")
		args = args.Add("Address", event.Address.String())
		args = args.Add("BlockNumber", event.BlockNumber)
		args = args.Add("TxIndex", event.TxIndex)
		args = args.Add("Time", time.Now().Unix())
		red.Send("XADD", args...)
	}
	red.Send("EXEC")
	return red.Flush()
}

func CreateTransactionTOPStorage(event sniffer.Event2) error {
	var red = model.RedisPool.Get()
	defer red.Close()
	red.Send("MULTI")

	// Construct args with keys and values for the event
	args := redis.Args{keys.CreateTransactionTOPStorage()}
	byteArray, _ := json.Marshal(event)
	args = args.Add("*")
	args = args.Add("event", byteArray)
	red.Send("XADD", args...)
	red.Send("EXEC")
	return red.Flush()
	// args = args.Add("*")
	// args = args.Add("Address", event.Address)
	// args = args.Add("ContractName", event.ContractName)
	// args = args.Add("ChainID", event.ChainID)
	// args = args.Add("Data", event.Data)
	// args = args.Add("BlockHash", event.BlockHash)
	// args = args.Add("BlockNumber", event.BlockNumber)
	// args = args.Add("Name", event.Name)
	// args = args.Add("TxHash", event.TxHash)
	// args = args.Add("TxIndex", event.TxIndex)
	// args = args.Add("Gas", event.Gas)
	// args = args.Add("GasPrice", event.GasPrice)
	// args = args.Add("GasTipCap", event.GasTipCap)
	// args = args.Add("GasFeeCap", event.GasFeeCap)
	// args = args.Add("Value", event.Value)
	// args = args.Add("Nonce", event.Nonce)
	// args = args.Add("To", event.To)
	// args = args.Add("Status", event.Status)
	// args = args.Add("Timestamp", event.Timestamp)
	// args = args.Add("MinerAddress", event.MinerAddress)
	// args = args.Add("Size", event.Size)
	// args = args.Add("BlockReward", event.BlockReward)
	// args = args.Add("AverageGasTipCap", event.AverageGasTipCap)
	// args = args.Add("Time", time.Now().Unix())

	// // Check if there are more than maxStorage transactions already stored
	// numStored, err := redis.Int(red.Do("XLEN", keys.CreateTransactionTOPStorage()))
	// if err != nil {
	// 	return err
	// }

	// if numStored >= maxStorage {
	// 	// Delete the oldest transaction
	// 	oldestID, err := redis.Strings(red.Do("XRANGE", keys.CreateTransactionTOPStorage(), "-", "+", "COUNT", 1))
	// 	if err != nil {
	// 		return err
	// 	}
	// 	red.Send("XDEL", keys.CreateTransactionTOPStorage(), oldestID[0])
	// }

	// Add the new transaction

}

func Pending(id, contractAddr, inputData string) error {
	release := Lock(id)
	defer release()

	if b, err := Exists(id); err != nil {
		return err
	} else if b {
		return errors.New("the ID already exists")
	}

	// 目前我们采用一个 Account 一个进程的方式来处理链上事务， 所以，
	// 这里我们只要保证当前进程的 Nonce 唯一即可
	nonceLock.Lock()
	defer nonceLock.Unlock()

	var red = model.RedisPool.Get()
	defer red.Close()

	// 取本地 Nonce 值
	nonce, err := redis.Uint64(red.Do("GET", keys.NonceAt()))
	if err != nil {
		return err
	}

	// 通过 Redis 的事务机制， 将下列操作原子化
	red.Send("MUTIL")
	{
		// 向 Pending 队列发送消息
		args := redis.Args{keys.Pending()}
		args = args.Add("*")
		args = args.Add("id", id)
		args = args.Add("nonce", nonce)
		args = args.Add("to", contractAddr)
		args = args.Add("inputData", inputData)
		red.Send("XADD", args...)
	}

	{
		// 更新本地 Nonce 值
		red.Send("INCR", keys.NonceAt())
	}

	{
		// 创建 Nonce 索引
		red.Send("SET", keys.Nonce(nonce), id)
	}

	{
		// 创建订单实体
		var order = new(Order)
		order.Id = id
		order.Status = PENDING_STATUS
		order.Nonce = nonce
		order.CreatedAt = time.Now().Unix()

		args := redis.Args{keys.Entity(order.Id)}
		args = args.AddFlat(order)
		red.Send("HSET", args...)
	}

	red.Send("EXEC")
	return red.Flush()
}

func Bind(id, txHash, txData string) error {
	var red = model.RedisPool.Get()
	defer red.Close()

	args := redis.Args{keys.Entity(id)}
	args = args.Add(FIELD_TX_HASH, txHash)
	args = args.Add(FIELD_TX_DATA, txData)
	args = args.Add(FIELD_UPDATED_AT, time.Now().Unix())

	red.Send("HSET", args...)
	return red.Flush()
}

func Replace(id, txHash, txData string, sent *Modifier) error {
	var red = model.RedisPool.Get()
	defer red.Close()

	red.Send("MUTIL")
	{
		args := redis.Args{keys.Sent()}
		args = args.Add("*")
		args = args.Add(sent.Values()...)

		red.Send("XADD", args...)
	}

	{
		args := redis.Args{keys.Entity(id)}
		args = args.Add(FIELD_TX_HASH, txHash)
		args = args.Add(FIELD_TX_DATA, txData)
		args = args.Add(FIELD_UPDATED_AT, time.Now().Unix())

		red.Send("HSET", args...)
	}

	red.Send("HINCRBY", keys.Entity(id), FIELD_NUMBER_OF_RETRIES, 1)
	red.Send("EXEC")
	return red.Flush()
}

func Lock(id string) func() {
	l := getLocker(id)
	l.Lock()
	return l.Unlock
}

func IsCompleted(id string) (bool, error) {

	status, err := Status(id)
	if err != nil {
		return false, err
	}

	switch status {
	case PENDING_STATUS, SENT_STATUS:
		return false, nil
	case SUCCEED_STATUS, FAILED_STATUS, ERROR_STATUS:
		return true, nil
	default:
		panic(status)
	}
}

func ID(nonce uint64) (string, error) {
	var red = model.RedisPool.Get()
	defer red.Close()

	return redis.String(red.Do("GET", keys.Nonce(nonce)))
}

func Cancel(nonce uint64) error {
	var red = model.RedisPool.Get()
	defer red.Close()

	red.Send("DEL", keys.Nonce(nonce))
	return red.Flush()
}

func Status(id string) (string, error) {
	var red = model.RedisPool.Get()
	defer red.Close()

	return redis.String(red.Do("HGET", keys.Entity(id), FIELD_STATUS))
}

func Set(id string, args ...interface{}) error {
	var red = model.RedisPool.Get()
	defer red.Close()

	argv := redis.Args{keys.Entity(id)}
	argv = argv.Add(args...)
	red.Send("HSET", argv...)
	return red.Flush()
}

func Get(id string, fields ...string) (*Order, error) {
	var red = model.RedisPool.Get()
	defer red.Close()

	if len(fields) > 0 {
		args := redis.Args{keys.Entity(id)}
		args = args.AddFlat(fields)
		values, err := redis.Values(red.Do("HMGET", args...))
		if err != nil {
			return nil, err
		}

		var result []interface{}
		for index, field := range fields {
			result = append(result, []byte(field), values[index])
		}

		var order = new(Order)
		if err := redis.ScanStruct(result, order); err != nil {
			return nil, err
		}

		return order, nil
	}

	values, err := redis.Values(red.Do("HGETALL", keys.Entity(id)))
	if err != nil {
		return nil, err
	}

	var order = new(Order)
	if err := redis.ScanStruct(values, order); err != nil {
		return nil, err
	}
	return order, nil
}

func TxHash(id string) (string, error) {
	var red = model.RedisPool.Get()
	defer red.Close()

	return redis.String(red.Do("HGET", keys.Entity(id), FIELD_TX_HASH))
}

func TxData(id string) (string, error) {
	var red = model.RedisPool.Get()
	defer red.Close()

	return redis.String(red.Do("HGET", keys.Entity(id), FIELD_TX_DATA))
}

func NumberOfRetries(id string) (int64, error) {
	var red = model.RedisPool.Get()
	defer red.Close()

	reply, err := redis.Int64(red.Do("HGET", keys.Entity(id), FIELD_NUMBER_OF_RETRIES))
	if err == redis.ErrNil {
		return 0, nil
	}
	return reply, err
}

func NonceAt() (uint64, error) {
	var red = model.RedisPool.Get()
	defer red.Close()

	return redis.Uint64(red.Do("GET", keys.NonceAt()))
}

func Exists(id string) (bool, error) {
	var red = model.RedisPool.Get()
	defer red.Close()

	reply, err := redis.Int(red.Do("EXISTS", keys.Entity(id)))
	return reply == 1, err
}

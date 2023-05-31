package orders

import (
	"ethgo/model"
	"ethgo/sniffer"
	"ethgo/util/redisx"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

type ReadOption interface {
	apply(*readOption)
}
type readOption struct {
	IsBlock  bool
	IsLimit  bool
	Count    int
	Timeout  int64
	MessagId string
}

type funcOption struct {
	f func(*readOption)
}

func (fdo *funcOption) apply(do *readOption) {
	fdo.f(do)
}

func WithNoBlock() ReadOption {
	return &funcOption{
		func(ro *readOption) {
			ro.IsBlock = false
			ro.Timeout = 0
		},
	}
}

func WithBlock(timeout int64) ReadOption {
	return &funcOption{
		func(ro *readOption) {
			ro.IsBlock = true
			ro.Timeout = timeout
		},
	}
}

func WithPending() ReadOption {
	return &funcOption{
		func(ro *readOption) {
			ro.MessagId = "0-0"
		},
	}
}

func WithNoLimit() ReadOption {
	return &funcOption{
		func(ro *readOption) {
			ro.IsLimit = false
			ro.Count = 0
		},
	}
}

func WithLimit(count int) ReadOption {
	return &funcOption{
		func(ro *readOption) {
			ro.IsLimit = true
			ro.Count = count
		},
	}
}

func WithLatest() ReadOption {
	return &funcOption{
		func(ro *readOption) {
			ro.MessagId = ">"
		},
	}
}

func defaultOptions() readOption {
	return readOption{
		IsBlock:  false,
		IsLimit:  false,
		Count:    0,
		Timeout:  0,
		MessagId: ">",
	}
}

type MessageReader struct {
	streamName string
	groupName  string
	readerName string
}

func (reader *MessageReader) StreamName() string {
	return reader.streamName
}

func (reader *MessageReader) GroupName() string {
	return reader.groupName
}

func (reader *MessageReader) Name() string {
	return reader.readerName
}

func (reader *MessageReader) Trim(timeout int64) (int64, error) {
	var red = model.RedisPool.Get()
	defer red.Close()

	reply, err := redis.Values(red.Do("XPENDING", reader.streamName, reader.groupName))
	if err != nil {
		return 0, err
	}

	pendingNumber, err := redis.Int64(reply[0], nil)
	if err != nil {
		panic(err)
	}

	// All entries that have an ID lower than messageId will be evicted:
	var messageId = time.Now().UnixMilli() - timeout*1000
	if pendingNumber > 0 {
		minID, err := redis.String(reply[1], nil)
		if err != nil {
			panic(err)
		}

		values := strings.SplitN(minID, "-", 2)
		lastMessageId, err := strconv.ParseInt(values[0], 10, 64)
		if err != nil {
			panic(err)
		}

		if lastMessageId < messageId {
			messageId = lastMessageId
		}
	}

	return redis.Int64(red.Do("XTRIM", reader.streamName, "MINID", messageId))
}

func (reader *MessageReader) Read(opts ...ReadOption) ([]*Message, error) {
	var red = model.RedisPool.Get()
	defer red.Close()

	var options = defaultOptions()
	for _, opt := range opts {
		opt.apply(&options)
	}

	var args = redis.Args{}
	args = args.Add("GROUP", reader.groupName, reader.readerName)
	if options.IsBlock {
		args = args.Add("BLOCK", options.Timeout)
	}
	if options.IsLimit {
		args = args.Add("COUNT", options.Count)
	}

	args = args.Add("STREAMS", reader.streamName, options.MessagId)

	streams, err := redisx.Streams(redis.Values(red.Do("XREADGROUP", args...)))
	if err != nil {
		return nil, err
	}

	var messages = []*Message{}
	for i := range streams[0].Value {
		var msg = streams[0].Value[i]
		messages = append(messages, &Message{
			streamName: reader.streamName,
			groupName:  reader.groupName,
			messageId:  msg.Id,
			data:       msg.Data,
		})
	}

	return messages, nil
}

func NewMessageReader(streamName, groupName, readerName string) (*MessageReader, error) {
	var red = model.RedisPool.Get()
	defer red.Close()

	_, err := red.Do("XGROUP", "CREATE", streamName, groupName, "0-0", "MKSTREAM")
	if err != nil {
		if err.Error() != "BUSYGROUP Consumer Group name already exists" {
			return nil, err
		}
	}

	return &MessageReader{
		streamName: streamName,
		groupName:  groupName,
		readerName: readerName,
	}, nil
}

func ReadTransactionTOPStorage() ([]sniffer.Event2, error) {
	var transactions []sniffer.Event2
	red := model.RedisPool.Get()
	defer red.Close()

	reply, err := redis.Values(red.Do("XRANGE", keys.CreateTransactionTOPStorage(), "-", "+"))
	if err != nil {
		return nil, err
	}

	for _, data := range reply {
		transaction, err := redis.Values(data, nil)
		if err != nil {
			return nil, err
		}

		var event sniffer.Event2
		if err = redis.ScanStruct(transaction, &event); err != nil {
			return nil, err
		}

		transactions = append(transactions, event)
	}

	return transactions, nil
}

func ReadChainDataStorag() (sniffer.ChainData, error) {
	var chainData sniffer.ChainData
	red := model.RedisPool.Get()
	defer red.Close()

	chainDataKey := keys.CreateChainDataStorag()

	// Check if the hash key exists
	exists, err := redis.Bool(red.Do("EXISTS", chainDataKey))
	if err != nil {
		return chainData, err
	}

	// If the hash key exists, get the data
	if exists {
		blockRewards, _ := redis.Float64(red.Do("HGET", chainDataKey, "BlockRewards"))
		superNodes, _ := redis.Uint64(red.Do("HGET", chainDataKey, "SuperNodes"))
		blockHeight, _ := redis.String(red.Do("HGET", chainDataKey, "BlockHeight"))
		totalBlockRewards, _ := redis.Float64(red.Do("HGET", chainDataKey, "TotalBlockRewards"))
		gasPrice, _ := redis.Float64(red.Do("HGET", chainDataKey, "GasPriceGasPrice"))
		totalAddrs, _ := redis.String(red.Do("HGET", chainDataKey, "TotalnumberofAddresses"))
		numTransfers, _ := redis.Int(red.Do("HGET", chainDataKey, "NumberTransactions"))
		numTxs, _ := redis.Int(red.Do("HGET", chainDataKey, "NumberTransfers"))
		numTxs24Hr, _ := redis.Int(red.Do("HGET", chainDataKey, "NumberTransactionsIn24H"))
		NumberaddressesIn24H, _ := redis.Int(red.Do("HGET", chainDataKey, "NumberaddressesIn24H"))

		chainData = sniffer.ChainData{
			BlockHeight:             blockHeight,
			BlockRewards:            strconv.FormatFloat(blockRewards, 'f', -1, 64),
			SuperNodes:              superNodes,
			TotalBlockRewards:       strconv.FormatFloat(totalBlockRewards, 'f', -1, 64),
			GasPriceGasPrice:        strconv.FormatFloat(gasPrice, 'f', -1, 64),
			TotalnumberofAddresses:  totalAddrs,
			NumberTransfers:         strconv.Itoa(numTransfers),
			NumberTransactions:      strconv.Itoa(numTxs),
			NumberTransactionsIn24H: strconv.Itoa(numTxs24Hr),
			NumberaddressesIn24H:    strconv.Itoa(NumberaddressesIn24H),
		}
	}

	return chainData, nil
}

func NewErrorReader(groupName string, readerName string) (*MessageReader, error) {
	return NewMessageReader(keys.Error(), groupName, readerName)
}

func NewFailedReader(groupName string, readerName string) (*MessageReader, error) {
	return NewMessageReader(keys.Failed(), groupName, readerName)
}

func NewPendingReader(groupName string, readerName string) (*MessageReader, error) {
	return NewMessageReader(keys.Pending(), groupName, readerName)
}

func NewSentReader(groupName string, readerName string) (*MessageReader, error) {
	return NewMessageReader(keys.Sent(), groupName, readerName)
}

func NewSucceedReader(groupName string, readerName string) (*MessageReader, error) {
	return NewMessageReader(keys.Succeed(), groupName, readerName)
}

func NewTransactionStorageReader(groupName string, readerName string) (*MessageReader, error) {
	return NewMessageReader(keys.CreateTransactionStorage(), groupName, readerName)
}

func NewTransactionTOPStorageReader(groupName string, readerName string) (*MessageReader, error) {
	return NewMessageReader(keys.CreateTransactionTOPStorage(), groupName, readerName)
}

func NewCreateChainDataStoragReader(groupName string, readerName string) (*MessageReader, error) {
	return NewMessageReader(keys.CreateChainDataStorag(), groupName, readerName)
}

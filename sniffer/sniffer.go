package sniffer

import (
	"context"
	"errors"
	"ethgo/eth"
	"ethgo/model/blocknumber"
	"ethgo/util/ethx"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Sniffer struct {
	chainID      *big.Int
	conf         *Config
	contracts    map[common.Address]*ethx.Contract
	handler      EventHandler
	addresses    []common.Address
	filterTopics []common.Hash
}

type TransactionInfo struct {
	From             common.Address
	TxIndex          uint              // index of the transaction within the block
	BlockNumber      uint64            // number of the block containing the transaction
	BlockHash        common.Hash       // hash of the block containing the transaction
	Tx               types.Transaction // the actual transaction object
	Status           bool
	Timestamp        uint64
	MinerAddress     string
	Size             string
	BlockReward      string
	AverageGasTipCap string
}

type Contract struct {
	*abi.ABI
	*bind.BoundContract

	Address common.Address
	Name    string
}

func defaultEventHandler(*Event) error {
	panic("请注册 EventHandler")
}

func New(conf *Config) (*Sniffer, error) {

	sf := &Sniffer{
		conf:         conf,
		handler:      defaultEventHandler,
		contracts:    make(map[common.Address]*ethx.Contract),
		addresses:    make([]common.Address, 0),
		filterTopics: make([]common.Hash, 0),
	}

	for _, v := range conf.Contracts {
		var address = common.HexToAddress(v.Addr)
		contract, err := ethx.NewContract(address, v.ABI)
		if err != nil {
			return nil, err
		}

		var eventIDs = make(map[string]common.Hash)
		for k, v := range contract.Events {
			eventIDs[strings.ToLower(k)] = v.ID
		}

		var filterTopics = make([]common.Hash, 0)
		for _, eventName := range v.Events {
			id, ok := eventIDs[strings.ToLower(eventName)]
			if !ok {
				return nil, ErrNoEvent
			}
			filterTopics = append(filterTopics, id)
		}

		sf.addresses = append(sf.addresses, address)
		sf.filterTopics = append(sf.filterTopics, filterTopics...)
		sf.contracts[address] = contract
	}

	return sf, nil
}

func (s *Sniffer) SetEventHandler(handler EventHandler) {
	if handler == nil {
		handler = defaultEventHandler
	}
	s.handler = handler
}

func (s *Sniffer) Run(ctx context.Context, backend eth.Backend) error {
	chainID, err := backend.ChainID(ctx)
	if err != nil {
		return err
	}

	latest, err := backend.BlockNumber(ctx)
	if err != nil {
		return err
	}

	if err := blocknumber.SetNX(latest); err != nil {
		return err
	}

	s.chainID = chainID

	s.run(ctx, backend)
	return nil
}

func (s *Sniffer) run(ctx context.Context, backend eth.Backend) {
	log.Info("开始侦听")
	defer log.Info("结束侦听")

	for {

		goto QUERY

	WAIT:
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
		}

	QUERY:
		// Beginning of the queried range.
		// 获取本地最新块
		fromBlockNumber, err := blocknumber.Get()
		if err != nil {
			log.With(err).Error("Failed to blocknumber.Get")
			goto WAIT
		}

		// End of the range.
		//获取当前安全块
		toBlockNumber, err := s.getSecurityBlockNumber(ctx, backend)
		if err != nil {
			log.With(err).Error("Failed to getSecurityBlockNumber")
			goto WAIT
		}

		if fromBlockNumber > toBlockNumber {
			// Is latest block number.
			// 如果本地块号大于安全块号，则表示本地已经是最新块，进入等待状态。
			goto WAIT
		}

		// Clipping block number range.
		// 如果查询的块数超过了配置文件中指定的块数，则将结束块号调整并限制查询块数。
		blockCnt := toBlockNumber - fromBlockNumber + 1
		if blockCnt > s.conf.NumberOfBlocks {
			toBlockNumber = fromBlockNumber + s.conf.NumberOfBlocks - 1
		}

		log.Debugf("起始块: %d, 结束块: %d", fromBlockNumber, toBlockNumber)

		// Executes a filter query.
		// 执行日志筛选操作，从区块中抽取感兴趣的日志信息。
		transaction, err := s.filterLogsAndTransactions(ctx, backend, fromBlockNumber, toBlockNumber)
		if err != nil {
			log.With(err).Error("Failed to filterLogs")
			goto WAIT
		}

		// log.Info(transaction)
		// Handle all logs.
		// 处理抽取到的日志信息，并在处理过程中出现错误则进入等待状态。
		if err := s.handleLogs(ctx, backend, transaction); err != nil {
			log.With(err).Error("Failed to handleLogs")
			goto WAIT
		}

		// Update local block number in redis.
		// 将本地块号更新为安全块号的下一个块号。
		blocknumber.Set(toBlockNumber + 1)
		goto QUERY
	}
}

func (s *Sniffer) getSecurityBlockNumber(ctx context.Context, backend eth.Backend) (uint64, error) {
	latestBlockNumber, err := backend.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}

	securityHeight := s.conf.SecrityHeight
	if latestBlockNumber < securityHeight {
		return 0, fmt.Errorf("no blocks")
	}

	return latestBlockNumber - securityHeight, nil
}

func (s *Sniffer) filterLogsAndTransactions(ctx context.Context, backend eth.Backend, fromBlockNumber uint64, toBlockNumber uint64) ([]TransactionInfo, error) {
	transactionsInfo, err := s.getTransactionsInBlocks(ctx, backend, fromBlockNumber, toBlockNumber)
	if err != nil {
		return nil, err
	}
	return transactionsInfo, nil
}

func (s *Sniffer) filterLogs(ctx context.Context, backend eth.Backend, blockNumber *big.Int, to common.Address) ([]types.Log, error) {
	query := ethereum.FilterQuery{
		FromBlock: blockNumber,
		ToBlock:   blockNumber,
		Addresses: []common.Address{to},
		Topics:    [][]common.Hash{},
	}
	return backend.FilterLogs(ctx, query)
}

func (s *Sniffer) getTransactionsInBlocks(ctx context.Context, backend eth.Backend, fromBlockNumber uint64, toBlockNumber uint64) ([]TransactionInfo, error) {
	transactions := make([]TransactionInfo, 0)

	for blockNumber := big.NewInt(int64(fromBlockNumber)); blockNumber.Cmp(new(big.Int).SetUint64(toBlockNumber)) <= 0; blockNumber.Add(blockNumber, big.NewInt(1)) {
		block, err := backend.BlockByNumber(ctx, blockNumber)
		if err != nil {
			log.Info("err_ getBlock", err)
			continue
		}
		blockReward, averageGasTipCap, err := getBlockRewar(block)
		if err != nil {
			log.Info("err_ getBlockRewar", err)
			continue
		}
		timestamp := block.Time()
		minerAddress := block.Coinbase().String()
		size := block.Size().String()
		// .Status
		txInfo := TransactionInfo{
			BlockNumber:      block.NumberU64(),
			BlockHash:        block.Hash(),
			Timestamp:        timestamp,
			MinerAddress:     minerAddress,
			Size:             size,
			BlockReward:      blockReward,
			AverageGasTipCap: averageGasTipCap,
		}
		transactions = append(transactions, txInfo)
		if err == nil {
			for txIndex, tx := range block.Transactions() {
				msg, err := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), big.NewInt(int64(block.NumberU64()))) // 获取交易对应的消息信息
				if err != nil {
					log.Info("err_ tx_Hash", tx.Hash().String())
					continue
				}
				receipt, err := backend.TransactionReceipt(ctx, tx.Hash())
				var status bool
				if err == nil && receipt != nil {
					status = receipt.Status == types.ReceiptStatusSuccessful
				}
				timestamp := block.Time()
				minerAddress := block.Coinbase().String()
				size := block.Size().String()
				// .Status
				txInfo := TransactionInfo{
					TxIndex:          uint(txIndex),
					BlockNumber:      block.NumberU64(),
					BlockHash:        block.Hash(),
					From:             msg.From(),
					Tx:               *tx,
					Status:           status,
					Timestamp:        timestamp,
					MinerAddress:     minerAddress,
					Size:             size,
					BlockReward:      blockReward,
					AverageGasTipCap: averageGasTipCap,
				}
				transactions = append(transactions, txInfo)
			}
		}
	}
	return transactions, nil
}

func getBlockRewar(block *types.Block) (string, string, error) {
	// 安全检测1：检查参数block是否为空
	if block == nil {
		return "", "", errors.New("block is nil")
	}

	// 获取所有的交易
	txs := block.Transactions()

	// 安全检测2：检查txs是否为空
	if len(txs) == 0 {
		return "", "", errors.New("no transactions in block")
	}

	gasTipCapSum := new(big.Int)
	gasFee := big.NewInt(0)
	for _, tx := range txs {
		// 安全检测3：检查交易tx是否为空
		if tx == nil {
			continue
		}
		var gasTipCap = tx.GasTipCap()
		var gas = tx.Gas()
		gasTipCapSum = gasTipCapSum.Add(gasTipCapSum, gasTipCap)
		gasFee = new(big.Int).Add(gasFee, new(big.Int).Mul(gasTipCap, new(big.Int).SetUint64(gas)))
	}

	// 安全检测4：检查len(txs)是否为0，防止除数为0
	if len(txs) == 0 {
		return "", "", errors.New("divide by zero")
	}
	averageGasTipCap := new(big.Int).Div(gasTipCapSum, big.NewInt(int64(len(txs))))
	var averageGasTipCapData = new(big.Float).Quo(new(big.Float).SetInt(averageGasTipCap), big.NewFloat(1e9))
	averageGasTipCapData.SetPrec(64)                             // 设置精度位数为64
	averageGasTipCapDataStr := averageGasTipCapData.Text('f', 6) // 将结果转换为字符串，保留18位小数
	var blockReward = new(big.Float).Quo(new(big.Float).SetInt(gasFee), big.NewFloat(1e18)).String()
	return blockReward, averageGasTipCapDataStr, nil
}

// var blockReward = new(big.Float).Quo(new(big.Float).SetInt(gasFee), new(big.Float).SetInt64(1e18)).String()
// return blockReward
// }

func (s *Sniffer) handleLogs(ctx context.Context, backend eth.Backend, txs []TransactionInfo) error {
	for _, tx := range txs {
		event := new(Event)
		// 解析交易数据成为事件对象
		if err := s.unpackTransaction(ctx, backend, &tx, event); err != nil {
			log.Panic(err)
		}
		log.Info("完成")
		// 处理反序列化后的事件
		if err := s.handleEvent(ctx, event); err != nil {
			return err
		}
		// 在应用程序关闭时，可以取消所有正在进行的处理任务
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
	return nil
}

func (s *Sniffer) unpackTransaction(ctx context.Context, backend eth.Backend, tx *TransactionInfo, out *Event) error {
	out.Name = ""                           // 设置Event结构中的事件名
	out.Data = make(map[string]interface{}) // 准备一个空的数据映射

	// 设置Event对象的其他属性
	out.Address = tx.From
	out.BlockHash = tx.BlockHash
	out.TxHash = tx.Tx.Hash()
	out.BlockNumber = strconv.FormatUint(tx.BlockNumber, 10)
	out.TxIndex = strconv.FormatUint(uint64(tx.TxIndex), 10)
	out.Gas = tx.Tx.Gas()
	out.GasPrice = tx.Tx.GasPrice()
	out.GasTipCap = tx.Tx.GasTipCap()
	out.GasFeeCap = tx.Tx.GasFeeCap()
	out.Value = tx.Tx.Value().String()
	out.Nonce = tx.Tx.Nonce()

	to := tx.Tx.To()
	out.Size = tx.Size
	out.Status = tx.Status
	out.MinerAddress = tx.MinerAddress
	out.Timestamp = tx.Timestamp
	out.BlockReward = tx.BlockReward
	out.AverageGasTipCap = tx.AverageGasTipCap
	if to == nil {
		out.To = common.Address{}
	} else {
		out.To = *to
	}
	out.ContractName = ""
	out.ChainID = s.chainID
	if len(tx.Tx.Data()) > 0 {

		txLogs, err := s.filterLogs(ctx, backend, big.NewInt(int64(tx.BlockNumber)), out.To)
		if err != nil {
			return err
		}
		if len(txLogs) == 0 {
			return nil
		}

		for _, log := range txLogs {
			if log.TxHash != tx.Tx.Hash() {
				return nil
			}
			if len(log.Topics) == 0 {
				continue
			}
			// if len(log.Topics) >= 0 {
			// 	continue
			// }

			// 遍历所有待匹配地址
			for _, address := range s.addresses {
				// 在嗅探器对象的合约映射中查找是否存在与地址匹配的合约对象
				contract := s.contracts[address]

				// 根据日志中第一个topic查找对应的事件
				event, err := contract.EventByID(log.Topics[0])
				if err == nil { // 如果找到了对应的事件
					out.ContractName = contract.Name        // 设置Event结构中的合约名
					out.ChainID = s.chainID                 // 设置Event结构中的链ID
					out.Name = event.Name                   // 设置Event结构中的事件名
					out.Data = make(map[string]interface{}) // 准备一个空的数据映射
					// 解压日志中的数据成为Event结构中的映射
					err := contract.UnpackLogIntoMap(out.Data, out.Name, log)

					if err != nil {
						fmt.Println("+++++++++++++++++++++++++", out.TxHash.String())
						continue
					}
					// 设置Event对象的其他属性
					out.Address = txLogs[0].Address
					out.BlockHash = txLogs[0].BlockHash
					out.TxHash = txLogs[0].TxHash
					out.BlockNumber = strconv.FormatUint(txLogs[0].BlockNumber, 10)
					out.TxIndex = strconv.FormatUint(uint64(txLogs[0].TxIndex), 10)
					return nil // 成功解析后结束函数
				}
			}
		}
		return nil
	}
	return nil
}

func (s *Sniffer) handleEvent(ctx context.Context, event *Event) error {
	log.Info(event)
	for {
		err := s.handler(event)
		if err == nil {
			return nil
		}

		log.Warn(err)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
		}
	}
}

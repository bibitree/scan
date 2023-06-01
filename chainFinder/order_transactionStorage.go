package chainFinder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	// "ethgo/model/mysqlOrders"
	"ethgo/model/mysqlOrders"
	"ethgo/model/orders"
	"ethgo/util"

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

	address := common.HexToAddress(message.String("Address")).String()
	if address != "" {
		address = t.conf.PrefixChain + address[2:]
	}

	txHash := common.HexToHash(message.String("TxHash")).String()
	if txHash == "0x0000000000000000000000000000000000000000000000000000000000000000" {
		return t.BlockStorage(ctx, message)
	}

	// var yourMap map[string]interface{}
	data := message.String("ContractName")
	if data != "" {
		err := t.ContractStorage(ctx, message)
		if err != nil {
			// 处理错误
			return After(t.conf.NetworkRetryInterval, message)
		}
	}

	Value := message.String("Value")
	if Value != "0" && Value != "" {
		err := t.AddressStorage(ctx, message, Value)
		if err != nil {
			// 处理错误
			return After(t.conf.NetworkRetryInterval, message)
		}
	}

	chainID, err := message.Int("ChainID")
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

	nonce, err := message.Int64("Nonce")
	if err != nil {
		// 发生错误，处理错误逻辑
		return After(t.conf.NetworkRetryInterval, message)
	}

	gas, err := message.Int64("Gas")
	if err != nil {
		// 发生错误，处理错误逻辑
		return After(t.conf.NetworkRetryInterval, message)
	}

	toAddress := common.HexToAddress(message.String("To")).String()
	if toAddress != "" {
		toAddress = t.conf.PrefixChain + toAddress[2:]
	}
	log.Info(toAddress)

	timestamp, err := message.Int("Timestamp")
	if err != nil {
		// 发生错误，处理错误逻辑
		return After(t.conf.NetworkRetryInterval, message)
	}

	blockNumber, err := message.Int64("BlockNumber")
	if err != nil {
		// 发生错误，处理错误逻辑
		return After(t.conf.NetworkRetryInterval, message)
	}

	txIndex, err := message.Int("TxIndex")
	if err != nil {
		// 发生错误，处理错误逻辑
		return After(t.conf.NetworkRetryInterval, message)
	}

	var Status int
	Status = 0
	if message.String("Status") == "1" {
		Status = 1
	}
	var event = sniffer.EventData{
		Address:      common.HexToAddress(message.String("Address")),
		ChainID:      chainID,
		BlockHash:    common.HexToHash(message.String("BlockHash")),
		BlockNumber:  big.NewInt(blockNumber),
		TxHash:       common.HexToHash(message.String("TxHash")),
		TxIndex:      txIndex,
		Gas:          big.NewInt(gas),
		GasPrice:     big.NewInt(gasPrice),
		GasTipCap:    big.NewInt(gasTipCap),
		GasFeeCap:    big.NewInt(gasFeeCap),
		Value:        Value,
		Nonce:        big.NewInt(nonce),
		To:           common.HexToAddress(message.String("To")),
		Status:       Status,
		Timestamp:    timestamp,
		NewAddress:   address,
		NewToAddress: toAddress,
	}

	log.Info(event)
	mysqlOrders.InsertEvent(event)
	return t.ack(message)
}

func (t *ChainFinder) AddressStorage(ctx context.Context, message *orders.Message, value string) AfterFunc {
	log.Debugf("ENTER @ContractStorage 订单")
	defer log.Debugf("  LEAVE @ContractStorage 订单")

	log.Info(message.String("Address"))

	i := new(big.Int)
	i.SetString(value, 10)
	var event = sniffer.AddressData{
		Address: message.String("To"),
		Balance: i,
	}

	log.Info(event)
	mysqlOrders.InsertAddressData(event)
	return t.ack(message)
}

func (t *ChainFinder) ContractStorage(ctx context.Context, message *orders.Message) AfterFunc {
	log.Debugf("ENTER @ContractStorage 订单")
	defer log.Debugf("  LEAVE @ContractStorage 订单")

	toAddress := common.HexToAddress(message.String("To")).String()
	if toAddress != "" {
		toAddress = t.conf.PrefixChain + toAddress[2:]
	}

	var yourMap map[string]interface{}
	data := message.Bytes("Data")
	if len(data) == 2 {
		yourMap = make(map[string]interface{})
	} else {
		err := json.Unmarshal(data, &yourMap)

		if err != nil {
			// 处理错误
			return After(t.conf.NetworkRetryInterval, message)
		}
	}
	log.Info(toAddress)
	var name string
	if message.String("ContractName") == "ERC20" {
		name, _ = t.StoreERCInfo(common.HexToAddress(message.String("To")).String())
	}

	var event = sniffer.ContractData{
		ContractName: name,
		EventName:    message.String("ContractName"),
		Data:         yourMap,
		Name:         message.String("Name"),
		TxHash:       common.HexToHash(message.String("TxHash")),
		Contrac:      common.HexToAddress(message.String("To")),
	}

	log.Info(event)
	mysqlOrders.InsertContractData(event)
	return t.ack(message)
}

func (t *ChainFinder) BlockStorage(ctx context.Context, message *orders.Message) AfterFunc {
	log.Debugf("ENTER @BlockStorage 订单")
	defer log.Debugf("  LEAVE @BlockStorage 订单")

	timestamp, err := message.Int("Timestamp")
	if err != nil {
		// 发生错误，处理错误逻辑
		return After(t.conf.NetworkRetryInterval, message)
	}

	gasLimit, err := message.Int("GasLimit")
	if err != nil {
		// 发生错误，处理错误逻辑
		return After(t.conf.NetworkRetryInterval, message)
	}

	blockNumber, err := message.Int64("BlockNumber")
	if err != nil {
		// 发生错误，处理错误逻辑
		return After(t.conf.NetworkRetryInterval, message)
	}

	blockReward, err := message.Float64("BlockReward")
	if err != nil {
		// 发生错误，处理错误逻辑
		return After(t.conf.NetworkRetryInterval, message)
	}
	a := int64(blockReward * 1e18)
	var event = sniffer.BlockData{
		BlockHash:    common.HexToHash(message.String("BlockHash")),
		BlockNumber:  big.NewInt(blockNumber),
		Timestamp:    timestamp,
		MinerAddress: message.String("MinerAddress"),
		Size:         message.String("Size"),
		BlockReward:  big.NewInt(a),
		GasLimit:     gasLimit,
	}
	log.Info(event)
	mysqlOrders.InsertBlock(event)
	return t.ack(message)
}

func (t *ChainFinder) StoreERCInfo(event string) (string, error) {

	var contract = Contract{
		Contract: event,
	}
	ercName, err := t.ProcessERCName(contract)
	if err != nil {
		log.Error(err)
		return "", err
	}
	ercTotalSupply, err := t.ProcessERCTotalSupply(contract)
	if err != nil {
		log.Error(err)
		return "", err
	}

	ethContractTxCount, err := mysqlOrders.GetEventsByContractAddress(event)
	if err != nil {
		log.Error(err)
		return "", err
	}
	fmt.Print(ethContractTxCount)

	address := event
	if address != "" {
		address = t.conf.PrefixChain + address[2:]
	}
	// dataMap, ok := ethContractTxCount.(map[string]interface{})
	// if !ok {
	// 	return "", errors.New("invalid data type")
	// }
	// fmt.Print(dataMap["count"].(int))
	// count, ok := ethContractTxCount.(ContractTxCount)

	// name, ok := ercName.(map[string]interface{})
	ercSlice := ercName.([]interface{})
	ercString := ercSlice[0].(string)
	fmt.Println(ercString)

	ercTotalSupply1 := ercTotalSupply.([]interface{})
	ercTotalSupplyFloat64 := ercTotalSupply1[0].(float64)
	ercTotalSupplyString := fmt.Sprintf("%.0f", ercTotalSupplyFloat64)
	fmt.Println(ercString)
	var ercTop = sniffer.ErcTop{
		ContractAddress:    event,
		ContractName:       ercString,
		Value:              ercTotalSupplyString,
		ContractTxCount:    len(ethContractTxCount),
		NewContractAddress: address,
	}
	mysqlOrders.InsertErcTop(ercTop)
	return ercString, nil
}

func (t *ChainFinder) ProcessERCName(contract Contract) (interface{}, error) {

	var call = Call{
		Address: contract.Contract,
		Method:  "name",
	}
	body, err := util.Post(t.conf.Callback, call)
	if err != nil {
		return nil, err
	}

	var res Response
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	log.Debugf("应答: %v", string(body))

	if res.Code != http.StatusOK {
		if res.Message == "" {
			res.Message = fmt.Sprintf("%v", res.Code)
		}
		return nil, errors.New(res.Message)
	}

	return res.Data, nil
}

func (t *ChainFinder) ProcessERCTotalSupply(contract Contract) (interface{}, error) {

	var call = Call{
		Address: contract.Contract,
		Method:  "totalSupply",
	}
	body, err := util.Post(t.conf.Callback, call)
	if err != nil {
		return nil, err
	}

	var res Response
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	log.Debugf("应答: %v", string(body))

	if res.Code != http.StatusOK {
		if res.Message == "" {
			res.Message = fmt.Sprintf("%v", res.Code)
		}
		return nil, errors.New(res.Message)
	}

	return res.Data, nil
}

func (t *ChainFinder) ProcessERCContractTxCount(contract Contract) (interface{}, error) {

	var call = Contract{
		Contract: contract.Contract,
	}
	body, err := util.Post(t.conf.ContractTxCount, call)
	if err != nil {
		return nil, err
	}

	var res Response
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	log.Debugf("应答: %v", string(body))

	if res.Code != http.StatusOK {
		if res.Message == "" {
			res.Message = fmt.Sprintf("%v", res.Code)
		}
		return nil, errors.New(res.Message)
	}

	return res.Data, nil
}

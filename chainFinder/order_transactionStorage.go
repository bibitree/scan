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

	address := common.HexToAddress(message.String("Address")).String()

	if address != "" {
		address = t.conf.PrefixChain + address[2:]
	}

	toAddress := common.HexToAddress(message.String("To")).String()
	if toAddress != "" {
		toAddress = t.conf.PrefixChain + toAddress[2:]
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

	log.Info(toAddress)

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
		NewAddress:       address,
		NewToAddress:     toAddress,
	}

	log.Info(event)

	if event.ContractName == "ERC20" {
		t.StoreERCInfo(common.HexToAddress(message.String("To")).String())
	}
	mysqlOrders.InsertEvent(event)
	// mysqlOrders.InsertEventData(event)
	return t.ack(message)
}

func (t *ChainFinder) StoreERCInfo(event string) error {

	var contract = Contract{
		Contract: event,
	}
	ercName, err := t.ProcessERCName(contract)
	if err != nil {
		log.Error(err)
		return err
	}
	ercTotalSupply, err := t.ProcessERCTotalSupply(contract)
	if err != nil {
		log.Error(err)
		return err
	}

	ethContractTxCount, err := t.ProcessERCContractTxCount(contract)
	if err != nil {
		log.Error(err)
		return err
	}
	fmt.Print(ethContractTxCount)

	address := event
	if address != "" {
		address = t.conf.PrefixChain + address[2:]
	}
	dataMap, ok := ethContractTxCount.(map[string]interface{})
	if !ok {
		return errors.New("invalid data type")
	}
	fmt.Print(dataMap["count"].(string))
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
		ContractTxCount:    dataMap["count"].(string),
		NewContractAddress: address,
	}
	mysqlOrders.InsertErcTop(ercTop)

	// mysqlOrders.InsertErcTop(ercTop)
	return nil
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

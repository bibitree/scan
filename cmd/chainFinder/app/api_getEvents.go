package app

import (
	"encoding/json"
	"errors"
	"ethgo/chainFinder"
	"ethgo/model"
	"ethgo/model/mysqlOrders"
	"ethgo/model/orders"
	"ethgo/proto"
	"ethgo/util"
	"ethgo/util/ginx"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (app *App) GetAllEvents(c *ginx.Context) {
	var request = new(proto.AllEvents)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	// ReadTransactionTOPStorage()
	messageReader, err := orders.GetLatestTransactionTOPStorage(request.N)
	if err != nil {
		panic(err)
	}
	c.Success(http.StatusOK, "succ", messageReader)
}

func (app *App) GetEventsByBlockNumbers(c *ginx.Context) {
	var request = new(proto.ByBlockNumbers)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	events, Contract, n, _, err := mysqlOrders.GetEventsBetweenBlockNumbers(uint64(request.Star), uint64(request.End), uint64(request.PageNo), uint64(request.PageSize))
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
	}
	Contracts, _, err := mysqlOrders.GetEventsByTxHashes(Contract)
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
	}
	BlockNumberS, err := mysqlOrders.GetBlockDataByBlockNumber2(uint64(request.Star), uint64(request.End), uint64(request.PageNo), uint64(request.PageSize))
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
	}
	paginate := chainFinder.EventData{
		Event:        events,
		PageNumber:   n,
		ContractData: Contracts,
		BlockData:    BlockNumberS,
	}
	c.Success(http.StatusOK, "succ", paginate)
}

func (app *App) GetEventsByAddress(c *ginx.Context) {
	var request = new(proto.ByAddress)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	events, Contract, page, err := mysqlOrders.GetEventsByAddress(request.Address, request.PageNo, request.PageSize)
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
	}
	Contracts, _, err := mysqlOrders.GetEventsByTxHashes(Contract)
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
	}

	eventData := chainFinder.EventData{
		ContractData: Contracts,
		Event:        events,
		PageNumber:   page,
	}
	c.Success(http.StatusOK, "succ", eventData)
}

func (app *App) GetAddres(c *ginx.Context) {
	var request = new(proto.ByAddress)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	addressData, err := mysqlOrders.GetTopAddress(request.Address)
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
	}
	addressDataList := []model.AddressData{
		{
			Address: addressData.Address,
			Balance: addressData.Balance,
			Count:   addressData.Count,
		},
	}
	c.Success(http.StatusOK, "succ", addressDataList)
}

func (app *App) GetAddressList(c *ginx.Context) {
	var request = new(proto.ByAddresss)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	events, pageNumber, err := mysqlOrders.GetTopAddresses(uint64(request.PageNo), uint64(request.PageSize))
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
	}
	addressData := chainFinder.AddressData{
		AddressData: events,
		PageNumber:  pageNumber,
	}
	c.Success(http.StatusOK, "succ", addressData)
}

func (app *App) GetEventsByBlockNumber(c *ginx.Context) {
	var request = new(proto.ByBlockNumber)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	events, Contract, _, err := mysqlOrders.GetEventByBlockNumber(uint64(request.BlockNumber))
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
	}
	Contracts, _, err := mysqlOrders.GetEventsByTxHashes(Contract)
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
	}

	BlockNumberS, err := mysqlOrders.GetBlockDataByBlockNumber([]string{strconv.Itoa(request.BlockNumber)})
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
	}

	eventData := chainFinder.EventData{
		ContractData: Contracts,
		Event:        events,
		BlockData:    BlockNumberS,
	}
	c.Success(http.StatusOK, "succ", eventData)
}

func (app *App) GetEventsByTxHash(c *ginx.Context) {
	var request = new(proto.ByTxHash)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	events, Contract, err := mysqlOrders.GetEventByTxHash(request.TxHash)
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
	}
	Contracts, _, err := mysqlOrders.GetEventsByTxHashes(Contract)
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
	}
	eventData := chainFinder.EventData{
		ContractData: Contracts,
		Event:        events,
	}
	c.Success(http.StatusOK, "succ", eventData)
}

func (app *App) GetEventsByBlockHash(c *ginx.Context) {
	var request = new(proto.ByBlockHash)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	events, Contract, blockNumber, err := mysqlOrders.GetEventByBlockHash(request.BlockHash)
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
	}
	Contracts, _, err := mysqlOrders.GetEventsByTxHashes(Contract)
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
	}
	BlockNumberS, err := mysqlOrders.GetBlockDataByBlockNumber(blockNumber)
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
	}

	eventData := chainFinder.EventData{
		ContractData: Contracts,
		Event:        events,
		BlockData:    BlockNumberS,
	}
	c.Success(http.StatusOK, "succ", eventData)
}

func (app *App) GetERCTop(c *ginx.Context) {
	var request = new(proto.AllEvents)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	events, addresss, err := mysqlOrders.GetErcTopData(uint64(request.N))
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
	}
	Contracts, err := mysqlOrders.GetCreateContractIconData(addresss)
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
	}

	contractData := chainFinder.ContractData{
		CreateContractData: Contracts,
		ErcTop:             events,
	}

	c.Success(http.StatusOK, "succ", contractData)
}

func (app *App) GetEventsByContract(c *ginx.Context) {
	var request = new(proto.ByContract)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	events, Contract, n, err := mysqlOrders.GetEventsByToAddressAndBlockNumber(request.Contract, uint64(request.PageNo), uint64(request.PageSize))
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
	}
	contracts, addresslist, err := mysqlOrders.GetEventsByTxHashes(Contract)
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
	}
	createContractData, err := mysqlOrders.GetCreateContractIconData(addresslist)
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
	}

	decimals, err := app.ProcessCall(request.Contract, "decimals")
	if err != nil {
		c.Success(http.StatusOK, "succ", nil)
		return
	}
	if decimals == nil {
		c.Success(http.StatusOK, "succ", nil)
		return
	}
	decimalsData := decimals.([]interface{})
	decimalsFloat64 := decimalsData[0].(float64)
	decimalsString := fmt.Sprintf("%.0f", decimalsFloat64)

	address, time, err := mysqlOrders.GetEventAddressByToAddress(request.Contract)
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
	}
	paginate := chainFinder.Paginate{
		Event:              events,
		PageNumber:         n,
		Decimals:           decimalsString,
		CreationTime:       time,
		Address:            address,
		ContractData:       contracts,
		CreateContractData: createContractData,
	}

	c.Success(http.StatusOK, "succ", paginate)
}

func (app *App) ProcessCall(contract string, name string) (interface{}, error) {

	var call = chainFinder.Call{
		Address: contract,
		Method:  name,
	}
	body, err := util.Post(app.conf.ChainFinder.Callback, call)
	if err != nil {
		return nil, err
	}

	var res chainFinder.Response
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	if res.Code != http.StatusOK {
		if res.Message == "" {
			res.Message = fmt.Sprintf("%v", res.Code)
		}
		return nil, errors.New(res.Message)
	}

	return res.Data, nil
}

// func (app *App) ProcessERCContractTxCount(contract string) (interface{}, error) {
// 	var call = chainFinder.Contract{
// 		Contract: contract,
// 	}
// 	body, err := util.Post(app.conf.ChainFinder.ContractTxCount, call)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var res chainFinder.Response
// 	if err := json.Unmarshal(body, &res); err != nil {
// 		return nil, err
// 	}

// 	log.Debugf("应答: %v", string(body))

// 	if res.Code != http.StatusOK {
// 		if res.Message == "" {
// 			res.Message = fmt.Sprintf("%v", res.Code)
// 		}
// 		return nil, errors.New(res.Message)
// 	}

// 	return res.Data, nil
// }

func (app *App) ProcessContractCreationTime(contract string) (interface{}, error) {

	var call = chainFinder.Contract{
		Contract: contract,
	}
	body, err := util.Post(app.conf.ChainFinder.ContractCreationTime, call)
	if err != nil {
		return nil, err
	}

	var res chainFinder.Response
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

func (app *App) GetChainData(c *ginx.Context) {
	var request = new(proto.ChainData)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	messageReader, err := orders.ReadChainDataStorag()
	if err != nil {
		panic(err)
	}

	// msgs, err := messageReader.Read()
	// if err != nil {
	// 	c.Failure(http.StatusBadRequest, err.Error(), nil)
	// }
	c.Success(http.StatusOK, "succ", messageReader)
}

func (app *App) GetBlockNum(c *ginx.Context) {
	var request = new(proto.ByContract)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}

	blockHeight, _, err := mysqlOrders.GetLatestEvent()
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
		return
	}
	paginate := chainFinder.ChainData{
		BlockHeight: blockHeight,
	}

	c.Success(http.StatusOK, "succ", paginate)
}

func (app *App) GetCreateContractData(c *ginx.Context) {
	var request = new(proto.GetCreateContractDataByContract)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	createContractData, err := mysqlOrders.GetCreateContractData(request.Contract)
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
		return
	}
	c.Success(http.StatusOK, "succ", createContractData)
}

func (app *App) IsContractAddress(c *ginx.Context) {
	var request = new(proto.ByAddress)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}
	createContractData, err := mysqlOrders.IsContractDataExists(request.Address)
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
		return
	}

	isContractAddressResponse := chainFinder.IsContractAddressResponse{
		Address:    request.Address,
		IsContract: createContractData,
	}

	c.Success(http.StatusOK, "succ", isContractAddressResponse)
}

func (app *App) CompareByIcon(c *ginx.Context) {
	// 获取合约地址
	address := c.PostForm("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing contracimgFilet address"})
		return
	}
	// 获取上传图片
	imgFile, err := c.FormFile("img")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/// 存储图片到指定路径
	imgPath := "img/" + address + filepath.Ext(imgFile.Filename)
	if err := os.MkdirAll("img", os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := c.SaveUploadedFile(imgFile, imgPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = mysqlOrders.UpdateCreateContractIconData(address, imgPath)
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
	}
	c.Success(http.StatusOK, "succ", imgPath)
}

func (app *App) CompareBytecodeAndSourceCode(c *ginx.Context) {
	// 获取上传文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取合约地址
	address := c.PostForm("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing contract address"})
		return
	}

	// 打开上传文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	// 读取文件内容并保存为string
	content, err := ioutil.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fileContentAsString := string(content)

	createContractData, err := mysqlOrders.GetCreateContractData(address)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	compareBytecodeAndSourceCode := chainFinder.CompareBytecodeAndSourceCode{
		Code:           fileContentAsString,
		BytecodeString: createContractData.BytecodeString,
	}

	decimals, err := app.ProcessCompareBytecodeAndSourceCode(compareBytecodeAndSourceCode)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	if decimals == nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	decimalsData := decimals.(interface{})
	decimalsFloat64 := decimalsData.(string)
	// decimalsString := fmt.Sprintf("%.0f", decimalsFloat64)
	err = mysqlOrders.UpdateCreateContractData(address, "", decimalsFloat64, fileContentAsString)
	if err != nil {
		c.Failure(http.StatusBadGateway, err.Error(), nil)
	}
	c.Success(http.StatusOK, "succ", decimalsFloat64)
	// c.JSON(http.StatusOK, gin.H{"status": "success", "data": })
}

func (app *App) ProcessCompareBytecodeAndSourceCode(compareBytecodeAndSourceCode chainFinder.CompareBytecodeAndSourceCode) (interface{}, error) {

	body, err := util.Post(app.conf.ChainFinder.CompareBytecodeAndSourceCode, compareBytecodeAndSourceCode)
	if err != nil {
		return "", err
	}

	var res Response
	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}

	log.Debugf("应答: %v", string(body))

	if res.Code != http.StatusOK {
		if res.Message == "" {
			res.Message = fmt.Sprintf("%v", res.Code)
		}
		return "", errors.New(res.Message)
	}

	return res.Data, nil
}

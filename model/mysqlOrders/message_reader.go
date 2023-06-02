package mysqlOrders

import (
	"context"
	"database/sql"
	"encoding/json"
	"ethgo/model"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
)

const BlockBeasReward = "10"

func GetAllEvents(n uint64) ([]model.EventData, []string, error) {
	// 声明SQL语句，限制返回数量为 n
	sqlStr := fmt.Sprintf("SELECT * FROM event ORDER BY id DESC LIMIT %d", n)

	// 查询所有数据
	rows, err := model.MysqlPool.Query(sqlStr)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	// 将查询结果遍历转化为EventData类型并返回
	events := make([]model.EventData, 0)
	txHashList := make([]string, 0)

	for rows.Next() {
		event := model.EventData{}
		var chainID, blockHashBytes, txHash []byte
		var id uint64
		var gasPrice, gasTipCap, gasFeeCap string // Modify the variable types for these fields
		err := rows.Scan(&id, &event.Address, &chainID,
			&blockHashBytes, &event.BlockNumber, &txHash, &event.TxIndex,
			&event.Gas, &gasPrice, &gasTipCap, &gasFeeCap,
			&event.Value, &event.Nonce, &event.To,
			&event.Status, &event.Timestamp, &event.NewAddress,
			&event.NewToAddress)
		if err != nil {
			log.Fatal(err)
		}
		gasPriceParsed, _ := new(big.Int).SetString(gasPrice, 10)
		event.GasPrice = gasPriceParsed

		gasTipCapParsed, _ := new(big.Int).SetString(gasTipCap, 10)
		event.GasTipCap = gasTipCapParsed

		gasFeeCapParsed, _ := new(big.Int).SetString(gasFeeCap, 10)
		event.GasFeeCap = gasFeeCapParsed

		event.ChainID = new(big.Int).SetBytes(chainID) // 将 []byte 转为 *big.Int
		event.BlockHash = string(blockHashBytes)
		event.TxHash = string(txHash)

		events = append(events, event)
		txHashList = append(txHashList, string(event.TxHash))
	}

	return events, txHashList, nil
}

// 声明SQL语句
func GetEventByTxHash(txHash string) ([]model.EventData, []string, error) {
	sqlStr := `SELECT * FROM event WHERE txHash = ?`
	// 查询匹配的数据
	rows, err := model.MysqlPool.Query(sqlStr, txHash)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	// 将查询结果遍历转化为Event类型并返回
	events := make([]model.EventData, 0)
	txHashList := make([]string, 0)

	for rows.Next() {
		event := model.EventData{}
		var chainID, blockHashBytes, txHash []byte
		var id uint64
		var gasPrice, gasTipCap, gasFeeCap string // Modify the variable types for these fields
		err := rows.Scan(&id, &event.Address, &chainID,
			&blockHashBytes, &event.BlockNumber, &txHash, &event.TxIndex,
			&event.Gas, &gasPrice, &gasTipCap, &gasFeeCap,
			&event.Value, &event.Nonce, &event.To,
			&event.Status, &event.Timestamp, &event.NewAddress,
			&event.NewToAddress)
		if err != nil {
			log.Fatal(err)
		}
		gasPriceParsed, _ := new(big.Int).SetString(gasPrice, 10)
		event.GasPrice = gasPriceParsed

		gasTipCapParsed, _ := new(big.Int).SetString(gasTipCap, 10)
		event.GasTipCap = gasTipCapParsed

		gasFeeCapParsed, _ := new(big.Int).SetString(gasFeeCap, 10)
		event.GasFeeCap = gasFeeCapParsed

		event.ChainID = new(big.Int).SetBytes(chainID) // 将 []byte 转为 *big.Int
		event.BlockHash = string(blockHashBytes)
		event.TxHash = string(txHash)

		events = append(events, event)
		txHashList = append(txHashList, string(event.TxHash))
	}
	return events, txHashList, nil
}

// 声明SQL语句
func GetEventsByAddress(address string) ([]model.EventData, []string, error) {
	sqlStr := `SELECT * FROM event WHERE address = ? OR toAddress = ? ORDER BY blockNumber DESC` // Add ORDER BY clause to sort by blockNumber in descending order
	// 查询匹配的数据
	rows, err := model.MysqlPool.Query(sqlStr, address, address)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	// 将查询结果遍历转化为Event类型并返回
	events := make([]model.EventData, 0)
	txHashList := make([]string, 0)
	for rows.Next() {
		event := model.EventData{}
		var chainID, blockHashBytes, txHash []byte
		var id uint64
		var gasPrice, gasTipCap, gasFeeCap string // Modify the variable types for these fields
		err := rows.Scan(&id, &event.Address, &chainID,
			&blockHashBytes, &event.BlockNumber, &txHash, &event.TxIndex,
			&event.Gas, &gasPrice, &gasTipCap, &gasFeeCap,
			&event.Value, &event.Nonce, &event.To,
			&event.Status, &event.Timestamp, &event.NewAddress,
			&event.NewToAddress)
		if err != nil {
			log.Fatal(err)
		}
		gasPriceParsed, _ := new(big.Int).SetString(gasPrice, 10)
		event.GasPrice = gasPriceParsed
		gasTipCapParsed, _ := new(big.Int).SetString(gasTipCap, 10)
		event.GasTipCap = gasTipCapParsed
		gasFeeCapParsed, _ := new(big.Int).SetString(gasFeeCap, 10)
		event.GasFeeCap = gasFeeCapParsed
		event.ChainID = new(big.Int).SetBytes(chainID) // 将 []byte 转为 *big.Int
		event.BlockHash = string(blockHashBytes)
		event.TxHash = string(txHash)
		events = append(events, event)
		txHashList = append(txHashList, string(event.TxHash))
	}
	return events, txHashList, nil
}

func GetEventByBlockHash(blockHash string) ([]model.EventData, []string, []string, error) {
	// 声明SQL语句
	sqlStr := `SELECT * FROM event WHERE blockHash = ? ORDER BY blockNumber DESC`

	// 查询匹配的数据
	rows, err := model.MysqlPool.Query(sqlStr, blockHash)
	if err != nil {
		return nil, nil, nil, err
	}

	defer rows.Close()

	events := make([]model.EventData, 0)
	txHashList := make([]string, 0)
	blockNumberSet := make(map[string]bool)

	for rows.Next() {
		event := model.EventData{}
		var chainID, blockHashBytes, txHash []byte
		var id uint64
		var gasPrice, gasTipCap, gasFeeCap string // Modify the variable types for these fields
		err := rows.Scan(&id, &event.Address, &chainID,
			&blockHashBytes, &event.BlockNumber, &txHash, &event.TxIndex,
			&event.Gas, &gasPrice, &gasTipCap, &gasFeeCap,
			&event.Value, &event.Nonce, &event.To,
			&event.Status, &event.Timestamp, &event.NewAddress,
			&event.NewToAddress)
		if err != nil {
			log.Fatal(err)
		}
		gasPriceParsed, _ := new(big.Int).SetString(gasPrice, 10)
		event.GasPrice = gasPriceParsed

		gasTipCapParsed, _ := new(big.Int).SetString(gasTipCap, 10)
		event.GasTipCap = gasTipCapParsed

		gasFeeCapParsed, _ := new(big.Int).SetString(gasFeeCap, 10)
		event.GasFeeCap = gasFeeCapParsed

		event.ChainID = new(big.Int).SetBytes(chainID) // 将 []byte 转为 *big.Int
		event.BlockHash = string(blockHashBytes)
		event.TxHash = string(txHash)

		events = append(events, event)
		txHashList = append(txHashList, string(event.TxHash))
		blockNumberSet[event.BlockNumber] = true
	}
	// 从 set 中提取所有不重复的 blockNumber，并转换为 string list
	blockNumberList := make([]string, 0, len(blockNumberSet))
	for blockNumber := range blockNumberSet {
		blockNumberList = append(blockNumberList, blockNumber)
	}
	return events, txHashList, blockNumberList, nil
}

func GetEventByBlockNumber(blockNumber uint64) ([]model.EventData, []string, []string, error) {
	// 声明SQL语句
	sqlStr := `SELECT * FROM event WHERE blockNumber = ? ORDER BY blockNumber DESC`

	// 查询匹配的数据
	rows, err := model.MysqlPool.Query(sqlStr, blockNumber)
	if err != nil {
		return nil, nil, nil, err
	}

	defer rows.Close()

	events := make([]model.EventData, 0)
	txHashList := make([]string, 0)
	blockNumberSet := make(map[string]bool)
	for rows.Next() {
		event := model.EventData{}
		var chainID, blockHashBytes, txHash []byte
		var id uint64
		var gasPrice, gasTipCap, gasFeeCap string // Modify the variable types for these fields
		err := rows.Scan(&id, &event.Address, &chainID,
			&blockHashBytes, &event.BlockNumber, &txHash, &event.TxIndex,
			&event.Gas, &gasPrice, &gasTipCap, &gasFeeCap,
			&event.Value, &event.Nonce, &event.To,
			&event.Status, &event.Timestamp, &event.NewAddress,
			&event.NewToAddress)
		if err != nil {
			log.Fatal(err)
		}
		gasPriceParsed, _ := new(big.Int).SetString(gasPrice, 10)
		event.GasPrice = gasPriceParsed

		gasTipCapParsed, _ := new(big.Int).SetString(gasTipCap, 10)
		event.GasTipCap = gasTipCapParsed

		gasFeeCapParsed, _ := new(big.Int).SetString(gasFeeCap, 10)
		event.GasFeeCap = gasFeeCapParsed

		event.ChainID = new(big.Int).SetBytes(chainID) // 将 []byte 转为 *big.Int
		event.BlockHash = string(blockHashBytes)
		event.TxHash = string(txHash)

		events = append(events, event)
		txHashList = append(txHashList, string(event.TxHash))
		blockNumberSet[event.BlockNumber] = true
	}
	// 从 set 中提取所有不重复的 blockNumber，并转换为 string list
	blockNumberList := make([]string, 0, len(blockNumberSet))
	for blockNumber := range blockNumberSet {
		blockNumberList = append(blockNumberList, blockNumber)
	}
	return events, txHashList, blockNumberList, nil
}

func GetBlockDataByBlockNumber2(start uint64, end uint64, pageNo uint64, pageSize uint64) ([]model.BlockData3, error) {
	events := make([]model.BlockData3, 0)

	// 构造查询block表的sql语句，按照blockNumber从大到小排序
	sqlStr := fmt.Sprintf("SELECT * FROM block WHERE blockNumber BETWEEN %d AND %d ORDER BY blockNumber DESC LIMIT %d, %d", start, end, (pageNo-1)*pageSize, pageSize)
	countSqlStr := "SELECT COUNT(*) AS Count FROM event WHERE blockNumber = ?"
	// 使用查询数据库，并且在查询时使用超时参数
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := model.MysqlPool.QueryContext(ctx, sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 遍历查询结果并将其转换为ContractData类型
	for rows.Next() {
		var event model.BlockData3
		var id uint64
		err = rows.Scan(&id, &event.BlockHash, &event.BlockNumber, &event.BlockReward, &event.MinerAddress, &event.Size, &event.Timestamp, &event.GasLimit)
		if err != nil {
			return nil, err
		}
		countRows, err := model.MysqlPool.QueryContext(ctx, countSqlStr, event.BlockNumber)
		if err != nil {
			return nil, err
		}
		defer countRows.Close()
		if countRows.Next() {
			var count int
			err = countRows.Scan(&count)
			if err != nil {
				return nil, err
			}
			event.Count = count
		}
		event.BlockBeasReward = BlockBeasReward
		events = append(events, event)
	}

	return events, nil
}

func GetEventsBetweenBlockNumbers(start uint64, end uint64, pageNo uint64, pageSize uint64) ([]model.EventData, []string, uint64, []string, error) {
	// 计算偏移量
	offset := (pageNo - 1) * pageSize

	// 声明SQL语句和参数
	sqlStr := `SELECT COUNT(*) FROM event WHERE blockNumber BETWEEN ? AND ?`
	countSqlStr := `SELECT * FROM event WHERE blockNumber BETWEEN ? AND ? ORDER BY blockNumber DESC LIMIT ?,?`
	args := []interface{}{start, end}

	// 查询匹配的数据
	row := model.MysqlPool.QueryRow(sqlStr, args...)
	var total int
	err := row.Scan(&total)
	if err != nil {
		log.Fatal(err)
	}

	// 计算总页数
	pageCount := uint64(math.Ceil(float64(total) / float64(pageSize)))

	args = append(args, offset, pageSize)

	rows, err := model.MysqlPool.Query(countSqlStr, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// 将查询结果遍历转化为Event类型并返回
	events := make([]model.EventData, 0)
	txHashList := make([]string, 0)
	blockNumberList := make([]string, 0)
	for rows.Next() {
		event := model.EventData{}
		var chainID, blockHashBytes, txHash []byte
		var id uint64
		var gasPrice, gasTipCap, gasFeeCap string // Modify the variable types for these fields
		err := rows.Scan(&id, &event.Address, &chainID,
			&blockHashBytes, &event.BlockNumber, &txHash, &event.TxIndex,
			&event.Gas, &gasPrice, &gasTipCap, &gasFeeCap,
			&event.Value, &event.Nonce, &event.To,
			&event.Status, &event.Timestamp, &event.NewAddress,
			&event.NewToAddress)
		if err != nil {
			log.Fatal(err)
		}
		gasPriceParsed, _ := new(big.Int).SetString(gasPrice, 10)
		event.GasPrice = gasPriceParsed

		gasTipCapParsed, _ := new(big.Int).SetString(gasTipCap, 10)
		event.GasTipCap = gasTipCapParsed

		gasFeeCapParsed, _ := new(big.Int).SetString(gasFeeCap, 10)
		event.GasFeeCap = gasFeeCapParsed

		event.ChainID = new(big.Int).SetBytes(chainID) // 将 []byte 转为 *big.Int
		event.BlockHash = string(blockHashBytes)
		event.TxHash = string(txHash)

		events = append(events, event)
		txHashList = append(txHashList, string(event.TxHash))
		blockNumberList = append(blockNumberList, event.BlockNumber)
	}
	blockNumberSet := make(map[string]bool) // 新建一个 set 来去除重复的 BlockNumber
	for _, blockNumber := range blockNumberList {
		blockNumberSet[blockNumber] = true
	}
	distinctBlockNumbers := make([]string, 0) // 新建一个列表来存储去重后的 BlockNumber
	for blockNumber := range blockNumberSet {
		distinctBlockNumbers = append(distinctBlockNumbers, blockNumber)
	}
	return events, txHashList, pageCount, distinctBlockNumbers, nil
}

func GetErcTopData(n uint64) ([]model.ErcTop, error) {
	// 声明SQL语句，限制返回数量为 n
	sqlStr := fmt.Sprintf("SELECT * FROM ercTop ORDER BY ContractTxCount DESC LIMIT %d", n)
	// 查询所有数据
	rows, err := model.MysqlPool.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var id uint64
	// 将查询结果遍历转化为ErcTop类型并返回
	ercTops := make([]model.ErcTop, 0)
	for rows.Next() {
		ercTop := model.ErcTop{}
		err := rows.Scan(&id, &ercTop.ContractAddress, &ercTop.ContractName, &ercTop.Value, &ercTop.NewContractAddress, &ercTop.ContractTxCount)
		if err != nil {
			log.Fatal(err)
		}
		ercTops = append(ercTops, ercTop)
	}

	return ercTops, nil
}

func GetEventsByToAddressAndBlockNumber(toAddress string, pageNo uint64, pageSize uint64) ([]model.EventData, []string, uint64, error) {
	// 查询匹配的总条数
	countSql := `SELECT count(*) FROM event WHERE toAddress = ?`
	countArgs := []interface{}{toAddress}
	var count uint64
	err := model.MysqlPool.QueryRow(countSql, countArgs...).Scan(&count)
	if err != nil {
		return nil, nil, 0, err
	}

	// 计算总页数
	totalPage := count / pageSize
	if count%pageSize != 0 {
		totalPage++
	}

	// 计算偏移量
	offset := (pageNo - 1) * pageSize

	// 声明SQL语句和参数
	sqlStr := `SELECT * FROM event WHERE toAddress = ? ORDER BY blockNumber DESC LIMIT ?,?`
	args := []interface{}{toAddress, offset, pageSize}

	// 查询匹配的数据
	rows, err := model.MysqlPool.Query(sqlStr, args...)
	if err != nil {
		return nil, nil, 0, err
	}
	defer rows.Close()

	// 将查询结果遍历转化为Event类型并返回
	events := make([]model.EventData, 0)
	txHashList := make([]string, 0)

	for rows.Next() {
		event := model.EventData{}
		var chainID, blockHashBytes, txHash []byte
		var id uint64
		var gasPrice, gasTipCap, gasFeeCap string // Modify the variable types for these fields
		err := rows.Scan(&id, &event.Address, &chainID,
			&blockHashBytes, &event.BlockNumber, &txHash, &event.TxIndex,
			&event.Gas, &gasPrice, &gasTipCap, &gasFeeCap,
			&event.Value, &event.Nonce, &event.To,
			&event.Status, &event.Timestamp, &event.NewAddress,
			&event.NewToAddress)
		if err != nil {
			log.Fatal(err)
		}
		gasPriceParsed, _ := new(big.Int).SetString(gasPrice, 10)
		event.GasPrice = gasPriceParsed

		gasTipCapParsed, _ := new(big.Int).SetString(gasTipCap, 10)
		event.GasTipCap = gasTipCapParsed

		gasFeeCapParsed, _ := new(big.Int).SetString(gasFeeCap, 10)
		event.GasFeeCap = gasFeeCapParsed

		event.ChainID = new(big.Int).SetBytes(chainID) // 将 []byte 转为 *big.Int
		event.BlockHash = string(blockHashBytes)
		event.TxHash = string(txHash)

		events = append(events, event)
		txHashList = append(txHashList, string(event.TxHash))
	}

	return events, txHashList, totalPage, nil
}

func GetEventAddressByToAddress(toAddress string) (string, uint64, error) {
	// 查询最小的TimeStamp
	sqlStr := `SELECT MIN(timeStamp) FROM event WHERE toAddress = ?`
	row := model.MysqlPool.QueryRow(sqlStr, toAddress)

	var timestamp uint64
	err := row.Scan(&timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", 0, fmt.Errorf("no events found for toAddress %s", toAddress)
		}
		return "", 0, err
	}

	// 根据TimeStamp查询完整的Event
	sqlStr = `SELECT address FROM event WHERE timeStamp = ? AND toAddress = ?`
	row = model.MysqlPool.QueryRow(sqlStr, timestamp, toAddress)

	var address string
	err = row.Scan(&address)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", 0, fmt.Errorf("no events found for timeStamp %d and toAddress %s", timestamp, toAddress)
		}
		return "", 0, err
	}

	return address, timestamp, nil
}

func GetLatestEvent() (string, string, error) {
	// 声明SQL语句
	sqlStr := `SELECT blockNumber FROM block ORDER BY cast(blockNumber as unsigned) DESC LIMIT 1`
	// 查询匹配的数据
	row := model.MysqlPool.QueryRow(sqlStr)

	var blockNumber int64

	// 绑定查询结果到对应变量
	err := row.Scan(&blockNumber)
	if err != nil {
		if err == sql.ErrNoRows { // 如果查询结果为空，则返回空字符串
			return "", "", nil
		}
		return "", "", err
	}

	// 将查询结果转化为string并返回
	blockNumberStr := strconv.FormatInt(blockNumber, 10)
	gasPriceStr := "1300000000"

	return blockNumberStr, gasPriceStr, nil
}

func GetEventStatistics(number string) (totalDataCount string, emptyContractNameCount string, last24HoursDataCount string, last24HoursUniqueAddressCount string, uniqueAddressCount string, err error) {
	// 获取数据库所有数据总共有多少条
	num, err := strconv.Atoi(number)
	if err != nil {
		fmt.Println("Atoi error:", err)
	} else {
		num += 1
	}
	var count int64
	err = model.MysqlPool.QueryRow("SELECT COUNT(*) FROM event").Scan(&count)
	if err != nil {
		return "", "", "", "", "", err
	}
	totalDataCount = strconv.Itoa(int(count))

	var count2 int64
	// 获取数据库中contractName为空的数据总共有多少条
	err = model.MysqlPool.QueryRow("SELECT COUNT(*) FROM ercevent").Scan(&count2)
	if err != nil {
		return totalDataCount, "", "", "", "", err
	}
	emptyContractNameCount = strconv.Itoa(int(count) - int(count2))

	// 获取距离现在24小时之内的数据有多少条
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	rows, err := model.MysqlPool.Query("SELECT COUNT(*) FROM event WHERE Timestamp >= ? AND Timestamp <= ? AND address != ?", yesterday.Unix(), now.Unix(), "0x0000000000000000000000000000000000000000")
	if err != nil {
		return totalDataCount, emptyContractNameCount, "", "", "", err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return totalDataCount, emptyContractNameCount, "", "", "", err
		}
	}
	last24HoursDataCount = strconv.Itoa(int(count))

	// 获取距离现在24小时之内的所有数据获得其全部的address以及toAddress，并去重之后数量
	rows, err = model.MysqlPool.Query("SELECT DISTINCT Address, ToAddress FROM event WHERE Timestamp >= ? AND Timestamp <= ? AND address != ?", yesterday.Unix(), now.Unix(), "0x0000000000000000000000000000000000000000")
	if err != nil {
		return totalDataCount, emptyContractNameCount, last24HoursDataCount, "", "", err
	}
	defer rows.Close()
	uniqueAddresses := make(map[string]bool)
	for rows.Next() {
		var address, toAddress string
		err = rows.Scan(&address, &toAddress)
		if err != nil {
			return totalDataCount, emptyContractNameCount, last24HoursDataCount, "", "", err
		}
		uniqueAddresses[address] = true
		uniqueAddresses[toAddress] = true
	}
	last24HoursUniqueAddressCount = strconv.Itoa(len(uniqueAddresses))

	// 获取距离现在24小时之内的所有数据获得其全部的address以及toAddress，并去重之后数量
	rows, err = model.MysqlPool.Query("SELECT DISTINCT Address, ToAddress FROM event ")
	if err != nil {
		return totalDataCount, emptyContractNameCount, last24HoursDataCount, last24HoursUniqueAddressCount, "", err
	}
	defer rows.Close()
	allUniqueAddresses := make(map[string]bool)
	for rows.Next() {
		var address, toAddress string
		err = rows.Scan(&address, &toAddress)
		if err != nil {
			return totalDataCount, emptyContractNameCount, last24HoursDataCount, last24HoursUniqueAddressCount, "", err
		}
		allUniqueAddresses[address] = true
		allUniqueAddresses[toAddress] = true
	}
	uniqueAddressCount = strconv.Itoa(len(allUniqueAddresses))

	return totalDataCount, emptyContractNameCount, last24HoursDataCount, last24HoursUniqueAddressCount, uniqueAddressCount, err
}

func GetAllAddressesAndBlockRewardSum(number string) (string, error) {
	sqlStr := `SELECT DISTINCT blockReward FROM block WHERE blockNumber = ?`
	rows, err := model.MysqlPool.Query(sqlStr, number)
	if err != nil && sql.ErrNoRows != err {
		return "", err
	}
	defer rows.Close()
	var blockRewardSum *big.Float = big.NewFloat(0)

	for rows.Next() {
		var blockReward string
		err := rows.Scan(&blockReward)
		if err != nil {
			log.Fatal(err)
		}
		blockRewardFloat, _ := new(big.Float).SetString(blockReward)
		if blockRewardSum == nil {
			blockRewardSum = big.NewFloat(0)
		}
		if blockRewardFloat == nil {
			blockRewardFloat = big.NewFloat(0)
		}
		blockRewardSum.Add(blockRewardSum, blockRewardFloat)
	}

	num, err := strconv.ParseInt(number, 10, 64)
	if err != nil {
		fmt.Println("Atoi error:", err)
	} else {
		num += 1
	}

	BlockBeasRewardNum, err := strconv.ParseInt(BlockBeasReward, 10, 64)
	if err != nil {
		fmt.Println("Atoi error:", err)
	}
	num = num * BlockBeasRewardNum * 1e18

	blockRewardSum = new(big.Float).Add(blockRewardSum, big.NewFloat(float64(num)))
	blockRewardSumStr := blockRewardSum.String()

	return blockRewardSumStr, nil
}

func GetEventsByTxHashes(txHashes []string) ([]model.ContractData, error) {
	events := make([]model.ContractData, 0)
	// 将传入的TxHash列表转换为字符串形式，以便查询数据库
	txHashStr := fmt.Sprintf("'%s'", strings.Join(txHashes, "','"))
	// 构造sql语句
	sqlStr := fmt.Sprintf("SELECT * FROM ercevent WHERE txHash in (%s)", txHashStr)
	// 使用QueryContext来查询数据库，并且在查询时使用超时参数
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := model.MysqlPool.QueryContext(ctx, sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// 遍历查询结果并将其转换为ContractData类型
	for rows.Next() {
		var event model.ContractData
		var data string
		// var id string
		var id uint64
		err = rows.Scan(&id, &event.ContractName, &event.EventName, &data, &event.Name, &event.TxHash, &event.Contrac)
		if err != nil {
			return nil, err
		}
		// 解析data字段
		json.Unmarshal([]byte(data), &event.Data)
		events = append(events, event)
	}
	return events, nil
}

func GetEventsByContractAddress(contractAddress string) ([]model.ContractData, error) {
	events := make([]model.ContractData, 0)
	// 构造sql语句，查询与指定合约地址匹配的数据
	sqlStr := fmt.Sprintf("SELECT * FROM ercevent WHERE contrac = '%s'", contractAddress)
	// 使用QueryContext来查询数据库，并且在查询时使用超时参数
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := model.MysqlPool.QueryContext(ctx, sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// 遍历查询结果并将其转换为ContractData类型
	for rows.Next() {
		var event model.ContractData
		var data string
		var id uint64
		err := rows.Scan(&id, &event.ContractName, &event.EventName, &data, &event.Name, &event.TxHash, &event.Contrac)
		if err != nil {
			return nil, err
		}
		// 解析data字段
		json.Unmarshal([]byte(data), &event.Data)
		events = append(events, event)
	}
	return events, nil
}

func GetBlockDataByBlockNumber(blockNumber []string) ([]model.BlockData3, error) {
	events := make([]model.BlockData3, 0)
	// 将传入的TxHash列表转换为字符串形式，以便查询数据库
	blockNumberStr := fmt.Sprintf("'%s'", strings.Join(blockNumber, "','"))
	// 构造查询event数据库中BlockNumber与本次循环中BlockNumber相同的交易数量的sql语句
	countSqlStr := "SELECT COUNT(*) AS Count FROM event WHERE blockNumber = ?"
	// 构造查询block表的sql语句
	sqlStr := fmt.Sprintf("SELECT * FROM block WHERE blockNumber in (%s)", blockNumberStr)
	// 使用查询数据库，并且在查询时使用超时参数
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := model.MysqlPool.QueryContext(ctx, sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// 遍历查询结果并将其转换为ContractData类型
	for rows.Next() {
		var event model.BlockData3
		var id uint64
		err = rows.Scan(&id, &event.BlockHash, &event.BlockNumber, &event.BlockReward, &event.MinerAddress, &event.Size, &event.Timestamp, &event.GasLimit)
		if err != nil {
			return nil, err
		}
		event.BlockBeasReward = BlockBeasReward
		// 查询交易数量并加入到BlockData2结构中
		countRows, err := model.MysqlPool.QueryContext(ctx, countSqlStr, event.BlockNumber)
		if err != nil {
			return nil, err
		}
		defer countRows.Close()
		if countRows.Next() {
			var count int
			err = countRows.Scan(&count)
			if err != nil {
				return nil, err
			}
			event.Count = count
		}
		events = append(events, event)
	}
	return events, nil
}

func GetTopAddresses(pageNo uint64, pageSize uint64) ([]model.AddressData, uint64, error) {
	var sqlStr string
	var params []interface{}
	sqlStr = "SELECT address, Balance, Count FROM addressTop ORDER BY Balance DESC"
	// 根据传入的参数设置limit值
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	start := (pageNo - 1) * pageSize
	sqlStr += " LIMIT ?,?"
	params = append(params, start, pageSize)
	rows, err := model.MysqlPool.Query(sqlStr, params...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var addresses []model.AddressData
	for rows.Next() {
		var addressData model.AddressData
		if err := rows.Scan(&addressData.Address, &addressData.Balance, &addressData.Count); err != nil {
			return nil, 0, err
		}
		addressData.Count, err = GetEventDataCountByAddress(addressData.Address)
		if err != nil {
			return nil, 0, err
		}
		addresses = append(addresses, addressData)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	totalRows, err := model.MysqlPool.Query("SELECT COUNT(*) FROM addressTop")
	if err != nil {
		return nil, 0, err
	}
	defer totalRows.Close()
	var total uint64
	for totalRows.Next() {
		if err := totalRows.Scan(&total); err != nil {
			return nil, 0, err
		}
	}
	totalPages := uint64(math.Ceil(float64(total) / float64(pageSize)))
	return addresses, totalPages, nil
}

func GetEventDataCountByAddress(address string) (uint64, error) {
	// 声明SQL语句，查询所有address和toaddress等于传入值的数据的数量
	sqlStr := "SELECT COUNT(*) FROM event WHERE address = ? OR `toAddress` = ?"
	// 查询数据数量
	row := model.MysqlPool.QueryRow(sqlStr, address, address)
	var count uint64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

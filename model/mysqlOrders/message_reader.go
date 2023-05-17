package mysqlOrders

import (
	"encoding/json"
	"ethgo/model"
	"fmt"
	"log"
	"math/big"
)

// 声明SQL语句
func GetEventByTxHash(txHash string) ([]model.Event, error) {
	sqlStr := `SELECT * FROM event WHERE txHash = ?`
	// 查询匹配的数据
	rows, err := model.MysqlPool.Query(sqlStr, txHash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// 将查询结果遍历转化为Event类型并返回
	events := make([]model.Event, 0)
	for rows.Next() {
		event := model.Event{}
		var data, chainID, blockHashBytes, txHash []byte
		var id uint64
		var gasPrice, gasTipCap, gasFeeCap string // Modify the variable types for these fields
		err := rows.Scan(&id, &event.Address, &event.ContractName, &chainID, &data,
			&blockHashBytes, &event.BlockNumber, &event.Name, &txHash, &event.TxIndex,
			&event.Gas, &gasPrice, &gasTipCap, &gasFeeCap,
			&event.Value, &event.Nonce, &event.ToAddress)

		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(data, &event.Data)
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
	}
	return events, nil
}

func GetEventByBlockHash(blockHash string) ([]model.Event, error) {
	// 声明SQL语句
	sqlStr := `SELECT * FROM event WHERE blockHash = ?`

	// 查询匹配的数据
	rows, err := model.MysqlPool.Query(sqlStr, blockHash)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// 将查询结果遍历转化为Event类型并返回
	events := make([]model.Event, 0)
	for rows.Next() {
		event := model.Event{}
		var data, chainID, blockHashBytes, txHash []byte
		var id uint64
		var gasPrice, gasTipCap, gasFeeCap string // Modify the variable types for these fields
		err := rows.Scan(&id, &event.Address, &event.ContractName, &chainID, &data,
			&blockHashBytes, &event.BlockNumber, &event.Name, &txHash, &event.TxIndex,
			&event.Gas, &gasPrice, &gasTipCap, &gasFeeCap,
			&event.Value, &event.Nonce, &event.ToAddress)

		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(data, &event.Data)
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
	}
	return events, nil
}

func GetEventByBlockNumber(blockNumber uint64) ([]model.Event, error) {
	// 声明SQL语句
	sqlStr := `SELECT * FROM event WHERE blockNumber = ?`

	// 查询匹配的数据
	rows, err := model.MysqlPool.Query(sqlStr, blockNumber)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// 将查询结果遍历转化为Event类型并返回
	events := make([]model.Event, 0)
	for rows.Next() {
		event := model.Event{}
		var data, chainID, blockHashBytes, txHash []byte
		var id uint64
		var gasPrice, gasTipCap, gasFeeCap string // Modify the variable types for these fields
		err := rows.Scan(&id, &event.Address, &event.ContractName, &chainID, &data,
			&blockHashBytes, &event.BlockNumber, &event.Name, &txHash, &event.TxIndex,
			&event.Gas, &gasPrice, &gasTipCap, &gasFeeCap,
			&event.Value, &event.Nonce, &event.ToAddress)

		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(data, &event.Data)
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
	}
	return events, nil
}

func GetEventsBetweenBlockNumbers(start uint64, end uint64, pageNo uint64, pageSize uint64) ([]model.Event, error) {
	// 计算偏移量
	offset := (pageNo - 1) * pageSize

	// 声明SQL语句和参数
	sqlStr := `SELECT * FROM event WHERE blockNumber BETWEEN ? AND ? LIMIT ?,?`
	args := []interface{}{start, end, offset, pageSize}

	// 查询匹配的数据
	rows, err := model.MysqlPool.Query(sqlStr, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// 将查询结果遍历转化为Event类型并返回
	events := make([]model.Event, 0)
	for rows.Next() {
		event := model.Event{}
		var data, chainID, blockHashBytes, txHash []byte
		var id uint64
		var gasPrice, gasTipCap, gasFeeCap string // Modify the variable types for these fields
		err := rows.Scan(&id, &event.Address, &event.ContractName, &chainID, &data,
			&blockHashBytes, &event.BlockNumber, &event.Name, &txHash, &event.TxIndex,
			&event.Gas, &gasPrice, &gasTipCap, &gasFeeCap,
			&event.Value, &event.Nonce, &event.ToAddress)

		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(data, &event.Data)
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
	}
	return events, nil
}

func GetAllEvents(n uint64) ([]model.Event, error) {
	// 声明SQL语句，限制返回数量为 n
	sqlStr := fmt.Sprintf("SELECT * FROM event ORDER BY id DESC LIMIT %d", n)
	// 查询所有数据
	rows, err := model.MysqlPool.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 将查询结果遍历转化为Event类型并返回
	events := make([]model.Event, 0)
	for rows.Next() {
		event := model.Event{}
		var data, chainID, blockHashBytes, txHash []byte
		var id uint64
		var gasPrice, gasTipCap, gasFeeCap string // Modify the variable types for these fields
		err := rows.Scan(&id, &event.Address, &event.ContractName, &chainID, &data,
			&blockHashBytes, &event.BlockNumber, &event.Name, &txHash, &event.TxIndex,
			&event.Gas, &gasPrice, &gasTipCap, &gasFeeCap,
			&event.Value, &event.Nonce, &event.ToAddress)

		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(data, &event.Data)
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
	}
	return events, nil
}

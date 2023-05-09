package mysqlOrders

import (
	"encoding/json"
	"ethgo/model"
	"ethgo/sniffer"
	"log"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
)

func GetEventByTxHash(txHash common.Hash) ([]sniffer.Event, error) {
	// 声明SQL语句
	sqlStr := `SELECT * FROM event WHERE txHash = ?`
	// 查询匹配的数据
	rows, err := model.MysqlPool.Query(sqlStr, txHash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// 将查询结果遍历转化为Event类型并返回
	events := make([]sniffer.Event, 0)
	for rows.Next() {
		event := sniffer.Event{}
		var data []byte
		err := rows.Scan(&event.Address, &event.ContractName, &event.ChainID, &data,
			&event.BlockHash, &event.BlockNumber, &event.Name, &event.TxHash, &event.TxIndex,
			&event.Gas, &event.GasPrice, &event.GasTipCap, &event.GasFeeCap,
			&event.Value, &event.Nonce, &event.To)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(data, &event.Data)
		if err != nil {
			log.Fatal(err)
		}
		events = append(events, event)
	}
	return events, nil
}

func GetEventByBlockHash(blockHash common.Hash) ([]sniffer.Event, error) {
	// 声明SQL语句
	sqlStr := `SELECT * FROM event WHERE blockHash = ?`

	// 查询匹配的数据
	rows, err := model.MysqlPool.Query(sqlStr, blockHash)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// 将查询结果遍历转化为Event类型并返回
	events := make([]sniffer.Event, 0)
	for rows.Next() {
		event := sniffer.Event{}
		var data []byte
		err := rows.Scan(&event.Address, &event.ContractName, &event.ChainID, &data,
			&event.BlockHash, &event.BlockNumber, &event.Name, &event.TxHash, &event.TxIndex,
			&event.Gas, &event.GasPrice, &event.GasTipCap, &event.GasFeeCap,
			&event.Value, &event.Nonce, &event.To)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(data, &event.Data)
		if err != nil {
			log.Fatal(err)
		}
		events = append(events, event)
	}

	return events, nil
}

func GetEventByBlockNumber(blockNumber uint64) ([]sniffer.Event, error) {
	// 声明SQL语句
	sqlStr := `SELECT * FROM event WHERE blockNumber = ?`

	// 查询匹配的数据
	rows, err := model.MysqlPool.Query(sqlStr, blockNumber)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// 将查询结果遍历转化为Event类型并返回
	events := make([]sniffer.Event, 0)
	for rows.Next() {
		event := sniffer.Event{}
		var data []byte
		err := rows.Scan(&event.Address, &event.ContractName, &event.ChainID, &data,
			&event.BlockHash, &event.BlockNumber, &event.Name, &event.TxHash, &event.TxIndex,
			&event.Gas, &event.GasPrice, &event.GasTipCap, &event.GasFeeCap,
			&event.Value, &event.Nonce, &event.To)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(data, &event.Data)
		if err != nil {
			log.Fatal(err)
		}
		events = append(events, event)
	}

	return events, nil
}

func GetEventsBetweenBlockNumbers(start uint64, end uint64) ([]sniffer.Event, error) {
	// 声明SQL语句
	sqlStr := `SELECT * FROM event WHERE blockNumber BETWEEN ? AND ?`
	// 查询匹配的数据
	rows, err := model.MysqlPool.Query(sqlStr, strconv.Itoa(int(start)), strconv.Itoa(int(end)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// 将查询结果遍历转化为Event类型并返回
	events := make([]sniffer.Event, 0)
	for rows.Next() {
		event := sniffer.Event{}
		var data []byte
		err := rows.Scan(&event.Address, &event.ContractName, &event.ChainID, &data,
			&event.BlockHash, &event.BlockNumber, &event.Name, &event.TxHash, &event.TxIndex,
			&event.Gas, &event.GasPrice, &event.GasTipCap, &event.GasFeeCap,
			&event.Value, &event.Nonce, &event.To)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(data, &event.Data)
		if err != nil {
			log.Fatal(err)
		}
		events = append(events, event)
	}
	return events, nil
}

func GetAllEvents() ([]sniffer.Event, error) {
	// 声明SQL语句
	sqlStr := `SELECT * FROM newtransaction`

	// 查询所有数据
	rows, err := model.MysqlPool.Query(sqlStr)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// 将查询结果遍历转化为Event类型并返回
	events := make([]sniffer.Event, 0)
	for rows.Next() {
		event := sniffer.Event{}
		var data []byte
		err := rows.Scan(&event.Address, &event.ContractName, &event.ChainID, &data,
			&event.BlockHash, &event.BlockNumber, &event.Name, &event.TxHash, &event.TxIndex,
			&event.Gas, &event.GasPrice, &event.GasTipCap, &event.GasFeeCap,
			&event.Value, &event.Nonce, &event.To)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(data, &event.Data)
		if err != nil {
			log.Fatal(err)
		}
		events = append(events, event)
	}

	return events, nil
}

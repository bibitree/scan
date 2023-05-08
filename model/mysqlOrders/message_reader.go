package mysqlOrders

import (
	"ethgo/model"
	"ethgo/sniffer"
	"fmt"
	"log"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
)

// 查询数据
func QueryData() error {
	// 声明SQL语句
	sqlStr := "SELECT name, age FROM user"
	// 执行查询操作，返回一个Rows结果集对象和错误信息
	rows, err := model.MysqlPool.Query(sqlStr)
	if err != nil {
		return err
	}
	defer rows.Close()

	// 遍历结果集
	for rows.Next() {
		var name string
		var age int

		err := rows.Scan(&name, &age)
		if err != nil {
			return err
		}
		fmt.Printf("%s is %d years old\n", name, age)
	}

	return nil
}

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
		err := rows.Scan(&event.Address, &event.ContractName, &event.ChainID, &event.Data,
			&event.BlockHash, &event.BlockNumber, &event.Name, &event.TxHash, &event.TxIndex)
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
		err := rows.Scan(&event.Address, &event.ContractName, &event.ChainID, &event.Data,
			&event.BlockHash, &event.BlockNumber, &event.Name, &event.TxHash, &event.TxIndex)
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
		err := rows.Scan(&event.Address, &event.ContractName, &event.ChainID, &event.Data,
			&event.BlockHash, &event.BlockNumber, &event.Name, &event.TxHash, &event.TxIndex)
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
		err := rows.Scan(&event.Address, &event.ContractName, &event.ChainID, &event.Data,
			&event.BlockHash, &event.BlockNumber, &event.Name, &event.TxHash, &event.TxIndex)
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
		err := rows.Scan(&event.Address, &event.ContractName, &event.ChainID, &event.Data,
			&event.BlockHash, &event.BlockNumber, &event.Name, &event.TxHash, &event.TxIndex)
		if err != nil {
			log.Fatal(err)
		}
		events = append(events, event)
	}

	return events, nil
}

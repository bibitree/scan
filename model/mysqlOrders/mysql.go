package mysqlOrders

import (
	"context"
	"encoding/json"
	"ethgo/model"
	"ethgo/sniffer"
	"time"

	"github.com/ethereum/go-ethereum/log"
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

func InsertContractData(event sniffer.ContractData) error {
	// 查询是否存在相同的txHash
	var count int
	err := model.MysqlPool.QueryRow("SELECT COUNT(*) FROM ercevent WHERE txHash=?", event.TxHash).Scan(&count)
	if err != nil {
		log.Error("查询是否存在相同的txHash时出错: ", err)
		return nil
	}

	if count == 0 { // 如果不存在相同的txHash，直接插入新数据
		sqlStr := `INSERT INTO event(contractName,data,  name, txHash, toAddress) VALUES (?, ?, ?, ?, ?)`
		// 使用ExecContext来执行sql语句，并且在执行时使用超时参数
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var data []byte
		data, _ = json.Marshal(event.Data)
		// 使用ExecContext执行sql语句，如果执行成功则返回nil
		_, err = model.MysqlPool.ExecContext(ctx, sqlStr, event.ContractName,
			string(data), event.Name, event.TxHash.Hex(),
			event.Contrac.Hex())
		if err != nil {
			log.Error("插入数据时出错: ", err)
			return err
		}
	}
	return nil
}

func InsertBlock(block sniffer.BlockData) error {
	// 查询是否存在相同的blockHash
	var count int
	err := model.MysqlPool.QueryRow("SELECT COUNT(*) FROM block WHERE blockHash=?", block.BlockHash).Scan(&count)
	if err != nil {
		log.Error("查询是否存在相同的blockHash时出错: ", err)
		return err
	}
	if count == 0 { // 如果不存在相同的blockHash，直接插入新数据
		sqlStr := `INSERT INTO block(blockHash, blockNumber, blockReward, minerAddress, size, timestamp) VALUES (?, ?, ?, ?, ?, ?)`
		// 使用ExecContext来执行sql语句，并且在执行时使用超时参数
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// 使用ExecContext执行sql语句，如果执行成功则返回nil
		_, err = model.MysqlPool.ExecContext(ctx, sqlStr, block.BlockHash.Hex(), block.BlockNumber, block.BlockReward,
			block.MinerAddress, block.Size, block.Timestamp)
		if err != nil {
			log.Info("插入数据时出错: ", err)
			return err
		}
	}
	return nil
}

func InsertEvent(event sniffer.EventData) error {
	// 查询是否存在相同的txHash
	var count int
	err := model.MysqlPool.QueryRow("SELECT COUNT(*) FROM event WHERE txHash=?", event.TxHash).Scan(&count)
	if err != nil {
		log.Error("查询是否存在相同的txHash时出错: ", err) //添加内容
		return nil
	}
	if count == 0 { // 如果不存在相同的txHash，直接插入新数据
		sqlStr := `INSERT INTO event(address, chainID, blockHash, blockNumber, txHash, txIndex, gas, gasPrice, gasTipCap, gasFeeCap, value, nonce, toAddress, status, timestamp, newAddress, newToAddress) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		// 使用ExecContext来执行sql语句，并且在执行时使用超时参数
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// 使用ExecContext执行sql语句，如果执行成功则返回nil
		_, err = model.MysqlPool.ExecContext(ctx, sqlStr, event.Address.Hex(), event.ChainID.String(),
			event.BlockHash.Hex(), event.BlockNumber, event.TxHash.Hex(), event.TxIndex,
			event.Gas, event.GasPrice.String(), event.GasTipCap.String(),
			event.GasFeeCap.String(), event.Value, event.Nonce,
			event.To.Hex(), event.Status, event.Timestamp,
			event.NewAddress, event.NewToAddress) // 添加新字段
		if err != nil {
			log.Error("插入数据时出错: ", err)
			return err
		}
	}
	return nil
}

func InsertErcTop(event sniffer.ErcTop) error {
	// 查询是否存在相同的contracaddress
	var count int
	err := model.MysqlPool.QueryRow("SELECT COUNT(*) FROM ercTop WHERE contracaddress=?", event.ContractAddress).Scan(&count)
	if err != nil {
		log.Error("查询是否存在相同的contracaddress时出错: ", err)
		return err
	}

	if count == 0 { // 如果不存在相同的contracaddress，直接插入新数据
		sqlStr := `INSERT INTO ercTop(contracaddress, name, value, newContracaddress, contractTxCount) VALUES (?, ?, ?, ?, ?)`
		// 使用ExecContext来执行sql语句，并且在执行时使用超时参数
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// 使用ExecContext执行sql语句，如果执行成功则返回nil
		_, err := model.MysqlPool.ExecContext(ctx, sqlStr, event.ContractAddress, event.ContractName, event.Value, event.NewContractAddress, event.ContractTxCount)
		if err != nil {
			log.Error("插入数据时出错: ", err)
			return err
		}
	} else { // 如果存在相同的contracaddress，更新数据
		sqlStr := `UPDATE ercTop SET name=?, value=?, newContracaddress=?, contractTxCount=? WHERE contracaddress=?`
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := model.MysqlPool.ExecContext(ctx, sqlStr, event.ContractName, event.Value, event.NewContractAddress, event.ContractTxCount, event.ContractAddress)
		if err != nil {
			log.Error("更新数据时出错: ", err)
			return err
		}
	}
	return nil
}

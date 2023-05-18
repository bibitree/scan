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

func InsertEvent(event sniffer.Event) error {
	// 查询是否存在相同的txHash
	var count int
	err := model.MysqlPool.QueryRow("SELECT COUNT(*) FROM event WHERE txHash=?", event.TxHash).Scan(&count)
	if err != nil {
		log.Error("查询是否存在相同的txHash时出错: ", err) //添加内容
		return nil
	}

	if count == 0 { // 如果不存在相同的txHash，直接插入新数据
		sqlStr := `INSERT INTO event2(address, contractName, chainID, data, blockHash, blockNumber, name, txHash, txIndex, gas, gasPrice, gasTipCap, gasFeeCap, value, nonce, toAddress, status, timestamp, minerAddress, size, blockReward, averageGasTipCap)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)` // 修改表名和新增字段
		// 使用ExecContext来执行sql语句，并且在执行时使用超时参数
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var data []byte
		data, _ = json.Marshal(event.Data)
		// 使用ExecContext执行sql语句，如果执行成功则返回nil
		_, err = model.MysqlPool.ExecContext(ctx, sqlStr, event.Address.Hex(), event.ContractName, event.ChainID.String(),
			string(data), event.BlockHash.Hex(),
			event.BlockNumber, event.Name, event.TxHash.Hex(), event.TxIndex,
			event.Gas, event.GasPrice.String(), event.GasTipCap.String(),
			event.GasFeeCap.String(), event.Value, event.Nonce,
			event.To.Hex(), event.Status, event.Timestamp,
			event.MinerAddress, event.Size,
			event.BlockReward, event.AverageGasTipCap) // 添加新字段
		if err != nil {
			log.Error("插入数据时出错: ", err)
			return err
		}

	}
	return nil
}

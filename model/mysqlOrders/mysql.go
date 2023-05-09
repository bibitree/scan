package mysqlOrders

import (
	"context"
	"database/sql"
	"encoding/json"
	"ethgo/model"
	"ethgo/sniffer"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
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

// 插入数据到MySQL
func InsertData() error {
	// 声明SQL语句
	sqlStr := `INSERT INTO user(name, age) VALUES (?, ?)`
	// 准备SQL语句
	stmt, err := model.MysqlPool.Prepare(sqlStr)
	if err != nil {
		return err
	}
	// 开始事务
	tx, err := model.MysqlPool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 执行SQL语句并传入参数
	_, err = tx.Stmt(stmt).Exec("Bob", 20)
	if err != nil {
		return err
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func InsertEvent(event sniffer.Event) error {
	// 查询是否存在相同的txHash
	var count int
	err := model.MysqlPool.QueryRow("SELECT COUNT(*) FROM event WHERE txHash=?", event.TxHash).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 { // 如果不存在相同的txHash，直接插入新数据
		sqlStr := `INSERT INTO event(address, contractName, chainID, data, blockHash, blockNumber, name, txHash, txIndex, gas, gasPrice, gasTipCap, gasFeeCap, value, nonce, to_address) 
                     VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		// 使用ExecContext来执行sql语句，并且在执行时使用超时参数
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var data []byte
		data, _ = sniffer.SerializeMap(event.Data)
		// 使用ExecContext执行sql语句，如果执行成功则返回nil
		_, err = model.MysqlPool.ExecContext(ctx, sqlStr, event.Address.Hex(), event.ContractName, event.ChainID.String(),
			string(data), event.BlockHash.Hex(),
			event.BlockNumber, event.Name, event.TxHash.Hex(), event.TxIndex,
			event.Gas, event.GasPrice.String(), event.GasTipCap.String(),
			event.GasFeeCap.String(), event.Value.String(), event.Nonce,
			event.To.Hex())
		if err != nil {
			return err
		}
	} else { // 如果存在相同的txHash，则查询是否有未填充数据，并将字段填充到空数据中
		// 查询未填充数据
		var emptyData sniffer.Event
		err = model.MysqlPool.QueryRow("SELECT * FROM event WHERE address='' OR contractName='' OR chainID='' OR data='' OR blockHash='' OR blockNumber=0 OR name='' OR txIndex=0 LIMIT 1").Scan(
			&emptyData.Address, &emptyData.ContractName, &emptyData.ChainID, &emptyData.Data, &emptyData.BlockHash,
			&emptyData.BlockNumber, &emptyData.Name, &emptyData.TxHash, &emptyData.TxIndex,
			&emptyData.Gas, &emptyData.GasPrice, &emptyData.GasTipCap, &emptyData.GasFeeCap,
			&emptyData.Value, &emptyData.Nonce, &emptyData.To,
		)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if emptyData.IsEmpty() { // 如果未找到未填充数据，则直接返回
			return nil
		} else { // 如果找到未填充数据，则检查是否有可填充的数据，进行补充
			updated := false

			if emptyData.Address == (common.Address{}) && event.Address != (common.Address{}) {
				emptyData.Address = event.Address
				updated = true
			}

			if emptyData.ContractName == "" && event.ContractName != "" {
				emptyData.ContractName = event.ContractName
				updated = true
			}

			if emptyData.ChainID == nil && event.ChainID != nil {
				emptyData.ChainID = event.ChainID
				updated = true
			}

			if emptyData.Data == nil && event.Data != nil && len(event.Data) > 0 {
				emptyData.Data = event.Data
				updated = true
			}

			if emptyData.BlockHash == (common.Hash{}) && event.BlockHash != (common.Hash{}) {
				emptyData.BlockHash = event.BlockHash
				updated = true
			}

			if emptyData.BlockNumber == "" && event.BlockNumber != "" {
				emptyData.BlockNumber = event.BlockNumber
				updated = true
			}

			if emptyData.Name == "" && event.Name != "" {
				emptyData.Name = event.Name
				updated = true
			}

			if emptyData.TxIndex == "" && event.TxIndex != "" {
				emptyData.TxIndex = event.TxIndex
				updated = true
			}

			if emptyData.Gas == 0 && event.Gas != 0 {
				emptyData.Gas = event.Gas
				updated = true
			}

			if emptyData.GasPrice == nil && event.GasPrice != nil {
				emptyData.GasPrice = event.GasPrice
				updated = true
			}

			if emptyData.GasTipCap == nil && event.GasTipCap != nil {
				emptyData.GasTipCap = event.GasTipCap
				updated = true
			}

			if emptyData.GasFeeCap == nil && event.GasFeeCap != nil {
				emptyData.GasFeeCap = event.GasFeeCap
				updated = true
			}

			if emptyData.Value == nil && event.Value != nil {
				emptyData.Value = event.Value
				updated = true
			}

			if emptyData.Nonce == 0 && event.Nonce != 0 {
				emptyData.Nonce = event.Nonce
				updated = true
			}

			if emptyData.To == nil && event.To != nil {
				emptyData.To = event.To
				updated = true
			}

			if updated { // 如果有可填充的数据，则更新数据表中的数据
				sqlStr := `UPDATE event SET address=?, contractName=?, chainID=?, data=?, blockHash=?, blockNumber=?, name=?, txIndex=?, 
                            gas=?, gasPrice=?, gasTipCap=?, gasFeeCap=?, value=?, nonce=?, to_address=? WHERE txHash=?`
				//使用ExecContext来执行sql语句，并且在执行时使用超时参数
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				var data []byte
				data, _ = sniffer.SerializeMap(event.Data)
				// 使用ExecContext执行sql语句，如果执行成功则返回nil
				_, err = model.MysqlPool.ExecContext(ctx, sqlStr, emptyData.Address.Hex(), emptyData.ContractName,
					emptyData.ChainID.String(), string(data),
					emptyData.BlockHash.Hex(), emptyData.BlockNumber, emptyData.Name,
					emptyData.TxIndex, emptyData.Gas, emptyData.GasPrice.String(),
					emptyData.GasTipCap.String(), emptyData.GasFeeCap.String(),
					emptyData.Value.String(), emptyData.Nonce, emptyData.To.Hex(),
					emptyData.TxHash.Hex())
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func InsertEventData(event sniffer.Event) error {
	// Check if event already exists in database
	var txCount int
	err := model.MysqlPool.QueryRow("SELECT count(*) FROM event WHERE txHash = ?", event.TxHash.Hex()).Scan(&txCount)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if txCount > 0 {
		// Event already exists, update empty fields if any
		var isUpdateRequired bool
		var updateFields []string
		query := "UPDATE event SET"
		if event.Address != (common.Address{}) { // check if address is not empty
			query += " address = ?,"
			updateFields = append(updateFields, "address")
			isUpdateRequired = true
		}
		if event.ContractName != "" { // check if contract name is not empty
			query += " contractName = ?,"
			updateFields = append(updateFields, "contractName")
			isUpdateRequired = true
		}
		if event.ChainID != nil { // check if chain ID is not nil
			query += " chainID = ?,"
			updateFields = append(updateFields, "chainID")
			isUpdateRequired = true
		}
		if len(event.Data) > 0 { // check if data is not empty
			_, err := json.Marshal(event.Data)
			if err != nil {
				return err
			}
			query += " data = ?,"
			updateFields = append(updateFields, "data")
			isUpdateRequired = true
		}
		if event.BlockHash != (common.Hash{}) { // check if block hash is not empty
			query += " blockHash = ?,"
			updateFields = append(updateFields, "blockHash")
			isUpdateRequired = true
		}
		if event.BlockNumber != "" { // check if block number is not empty
			query += " blockNumber = ?,"
			updateFields = append(updateFields, "blockNumber")
			isUpdateRequired = true
		}
		if event.Name != "" { // check if name is not empty
			query += " name = ?,"
			updateFields = append(updateFields, "name")
			isUpdateRequired = true
		}
		if event.TxIndex != "" { // check if transaction index is not empty
			query += " txIndex = ?,"
			updateFields = append(updateFields, "txIndex")
			isUpdateRequired = true
		}
		if event.Gas != 0 { // check if gas is not zero
			query += " gas = ?,"
			updateFields = append(updateFields, "gas")
			isUpdateRequired = true
		}
		if event.GasPrice != nil { // check if gas price is not nil
			query += " gasPrice = ?,"
			updateFields = append(updateFields, "gasPrice")
			isUpdateRequired = true
		}
		if event.GasTipCap != nil { // check if gas tip cap is not nil
			query += " gasTipCap = ?,"
			updateFields = append(updateFields, "gasTipCap")
			isUpdateRequired = true
		}
		if event.GasFeeCap != nil { // check if gas fee cap is not nil
			query += " gasFeeCap = ?,"
			updateFields = append(updateFields, "gasFeeCap")
			isUpdateRequired = true
		}
		if event.Value != nil { // check if value is not nil
			query += " value = ?,"
			updateFields = append(updateFields, "value")
			isUpdateRequired = true
		}
		if event.Nonce != 0 { // check if nonce is not zero
			query += " nonce = ?,"
			updateFields = append(updateFields, "nonce")
			isUpdateRequired = true
		}
		if event.To != nil { // check if recipient address is not nil
			query += " recipient = ?,"
			updateFields = append(updateFields, "recipient")
			isUpdateRequired = true
		}
		if isUpdateRequired {
			data, _ := json.Marshal(event.Data)
			query = strings.TrimSuffix(query, ",")
			query += " WHERE txHash = ?"
			updateFields = append(updateFields, "txHash")
			updateValues := make([]interface{}, len(updateFields))
			updateValues[0] = event.TxHash
			for i, field := range updateFields[1:] {
				switch field { // populate updateValues with values to insert in updated row
				case "address":
					updateValues[i+1] = event.Address
				case "contractName":
					updateValues[i+1] = event.ContractName
				case "chainID":
					updateValues[i+1] = event.ChainID.Int64()
				case "data":
					updateValues[i+1] = string(data)
				case "blockHash":
					updateValues[i+1] = event.BlockHash
				case "blockNumber":
					updateValues[i+1] = event.BlockNumber
				case "name":
					updateValues[i+1] = event.Name
				case "txIndex":
					updateValues[i+1] = event.TxIndex
				case "gas":
					updateValues[i+1] = event.Gas
				case "gasPrice":
					updateValues[i+1] = event.GasPrice.Int64()
				case "gasTipCap":
					updateValues[i+1] = event.GasTipCap.Int64()
				case "gasFeeCap":
					updateValues[i+1] = event.GasFeeCap.Int64()
				case "value":
					updateValues[i+1] = event.Value.Int64()
				case "nonce":
					updateValues[i+1] = event.Nonce
				case "recipient":
					updateValues[i+1] = event.To.Hex()
				}
			}
			_, err = model.MysqlPool.Exec(query, updateValues...)
			if err != nil {
				return err
			}
		}
	} else {
		// Event does not exist, insert new row
		sqlStr := `INSERT INTO event(address, contractName, chainID, data, blockHash, blockNumber, name, txHash, txIndex, gas, gasPrice, gasTipCap, gasFeeCap, value, nonce, recipient) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		stmt, err := model.MysqlPool.Prepare(sqlStr)
		if err != nil {
			return err
		}
		defer stmt.Close()
		serializedData, err := json.Marshal(event.Data)
		if err != nil {
			return err
		}
		var toStr string
		if event.To != nil {
			toStr = event.To.Hex()
		}
		_, err = stmt.Exec(event.Address, event.ContractName, event.ChainID.Int64(), string(serializedData), event.BlockHash.String(), event.BlockNumber, event.Name, event.TxHash.Hex(), event.TxIndex, event.Gas, event.GasPrice.Int64(), event.GasTipCap.Int64(), event.GasFeeCap.Int64(), event.Value.Int64(), event.Nonce, toStr)
		if err != nil {
			return err
		}
	}
	// Delete oldest record if more than 100 records in the table
	var rowCount int
	err = model.MysqlPool.QueryRow("SELECT count(*) FROM event").Scan(&rowCount)
	if err != nil {
		return err
	}
	if rowCount > 100 {
		_, err = model.MysqlPool.Exec("DELETE FROM event WHERE blockNumber = (SELECT MIN(blockNumber) FROM event)")
		if err != nil {
			return err
		}
	}
	return nil
}

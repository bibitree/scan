package mysqlOrders

import (
	"ethgo/model"
	"ethgo/sniffer"
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
	// 声明SQL语句
	sqlStr := `INSERT INTO event(address, contractName, chainID, data, blockHash, blockNumber, name, txHash, txIndex)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
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
	_, err = tx.Stmt(stmt).Exec(event.Address, event.ContractName, event.ChainID, event.Data, event.BlockHash, event.BlockNumber, event.Name, event.TxHash, event.TxIndex)
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

func InsertEventData(event sniffer.Event) error {
	// Declare SQL statement
	sqlStr := `INSERT INTO event(address, contractName, chainID, data, blockHash, blockNumber, name, txHash, txIndex)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Prepare SQL statement
	stmt, err := model.MysqlPool.Prepare(sqlStr)
	if err != nil {
		return err
	}

	// Start transaction
	tx, err := model.MysqlPool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute SQL statement and pass parameters
	_, err = tx.Stmt(stmt).Exec(event.Address, event.ContractName, event.ChainID, event.Data, event.BlockHash, event.BlockNumber, event.Name, event.TxHash, event.TxIndex)
	if err != nil {
		return err
	}

	// Fetch number of rows in table
	var count int
	rows, err := model.MysqlPool.Query("SELECT count(*) from event")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return err
		}
	}
	if count > 100 {
		// Remove row with minimum blockNumber
		_, err = tx.Exec("DELETE FROM event WHERE blockNumber = (SELECT MIN(blockNumber) FROM event)")
		if err != nil {
			return err
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

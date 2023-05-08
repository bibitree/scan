package sniffer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type AccessTuple struct {
	Address     common.Address `json:"address"        gencodec:"required"`
	StorageKeys []common.Hash  `json:"storageKeys"    gencodec:"required"`
}

type AccessList []AccessTuple

type TxData interface {
	txType() byte // returns the type ID
	copy() TxData // creates a deep copy and initializes all fields

	chainID() *big.Int
	accessList() AccessList
	data() []byte
	gas() uint64
	gasPrice() *big.Int
	gasTipCap() *big.Int
	gasFeeCap() *big.Int
	value() *big.Int
	nonce() uint64
	to() *common.Address

	rawSignatureValues() (v, r, s *big.Int)
	setSignatureValues(chainID, v, r, s *big.Int)
}

// 定义一个名为`MyTxData`的结构体，实现`TxData`接口中的所有方法
type MyTxData struct {
	TypeID       byte
	TxAccessList AccessList
	TxData       []byte
	TxGas        uint64
	TxGasPrice   *big.Int
	TxGasTipCap  *big.Int
	TxGasFeeCap  *big.Int
	TxValue      *big.Int
	TxNonce      uint64
	TxTo         *common.Address
	TxV          *big.Int
	TxR          *big.Int
	TxS          *big.Int
}

// 实现 `TxData`接口中的txType方法，返回结构体中的TypeID字段
func (m MyTxData) txType() byte {
	return m.TypeID
}

// 实现 `TxData`接口中的copy方法，返回当前结构体实例的一个深拷贝
func (m MyTxData) copy() MyTxData {
	copied := MyTxData{
		TypeID:       m.TypeID,
		TxAccessList: m.TxAccessList,
		TxData:       make([]byte, len(m.TxData)),
		TxGas:        m.TxGas,
		TxGasPrice:   m.TxGasPrice,
		TxGasTipCap:  m.TxGasTipCap,
		TxGasFeeCap:  m.TxGasFeeCap,
		TxValue:      m.TxValue,
		TxNonce:      m.TxNonce,
		TxTo:         m.TxTo,
		TxV:          m.TxV,
		TxR:          m.TxR,
		TxS:          m.TxS,
	}
	copy(copied.TxData, m.TxData)
	return copied
}

// 实现 `TxData`接口中的chainID方法，返回nil
func (m MyTxData) chainID() *big.Int {
	return nil
}

// 实现 `TxData`接口中的accessList方法，返回结构体中的TxAccessList字段
func (m MyTxData) accessList() AccessList {
	return m.TxAccessList
}

// 实现 `TxData`接口中的data方法，返回结构体中的TxData字段
func (m MyTxData) data() []byte {
	return m.TxData
}

// 实现 `TxData`接口中的gas方法，返回结构体中的TxGas字段
func (m MyTxData) gas() uint64 {
	return m.TxGas
}

// 实现 `TxData`接口中的gasPrice方法，返回结构体中的TxGasPrice字段
func (m MyTxData) gasPrice() *big.Int {
	return m.TxGasPrice
}

// 实现 `TxData`接口中的gasTipCap方法，返回结构体中的TxGasTipCap字段
func (m MyTxData) gasTipCap() *big.Int {
	return m.TxGasTipCap
}

// 实现 `TxData`接口中的gasFeeCap方法，返回结构体中的TxGasFeeCap字段
func (m MyTxData) gasFeeCap() *big.Int {
	return m.TxGasFeeCap
}

// 实现 `TxData`接口中的value方法，返回结构体中的TxValue字段
func (m MyTxData) value() *big.Int {
	return m.TxValue
}

// 实现 `TxData`接口中的nonce方法，返回结构体中的TxNonce字段
func (m MyTxData) nonce() uint64 {
	return m.TxNonce
}

// 实现 `TxData`接口中的to方法，返回结构体中的TxTo字段
func (m MyTxData) to() *common.Address {
	return m.TxTo
}

// 实现 `TxData`接口中的rawSignatureValues方法，
// 返回结构体中的TxV、TxR、TxS字段作为v、r、s的值
func (m MyTxData) rawSignatureValues() (v, r, s *big.Int) {
	return m.TxV, m.TxR, m.TxS
}

// 实现 `TxData`接口中的setSignatureValues方法，
// 将传入的chainID、v、r、s参数分别赋值给结构体中的TxV、TxR、TxS字段
func (m *MyTxData) setSignatureValues(chainID, v, r, s *big.Int) {
	m.TxV = v
	m.TxR = r
	m.TxS = s
}

package app

import (
	"context"
	"ethgo/model/orders"
	"ethgo/proto"
	"ethgo/tyche/types"
	"ethgo/util/ginx"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// BalanceAt
// @Description 获得钱包余额
// @Description
// @Tags 钱包
// @Accept application/json
// @Produce application/json
// @Param object body proto.Balance{} true "请求参数"
// @Success 200 {object}  proto.Response{data=proto.BalanceResponse{}}
// @Router /tyche/api/wallet/balance_at [post]
func (app *App) BalanceAt(c *ginx.Context) {
	var request = new(proto.Balance)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}

	if request.Address == "" {
		request.Address = app.conf.Tyche.Account
	}
	if !common.IsHexAddress(request.Address) {
		c.Failure(http.StatusBadRequest, "无效的参数: Address", nil)
		return
	}
	wei, err := app.backend.BalanceAt(context.Background(), common.HexToAddress(request.Address), nil)
	if err != nil {
		c.Failure(http.StatusInternalServerError, err.Error(), nil)
		return
	}
	c.Success(http.StatusOK, "succ", proto.BalanceResponse{
		Address: request.Address,
		Wei:     wei.String(),
	})
}

// Create
// @Description 创建一个钱包
// @Description
// @Tags 钱包
// @Accept application/json
// @Produce application/json
// @Param object body proto.Create{} true "请求参数"
// @Success 200 {object}  proto.Response{data=proto.CreateResponse{}}
// @Router /tyche/api/wallet/create [post]
func (app *App) Create(c *ginx.Context) {
	var request = new(proto.Create)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}

	key, err := crypto.GenerateKey()
	if err != nil {
		c.Failure(http.StatusInternalServerError, err.Error(), nil)
		return
	}

	c.Success(http.StatusOK, "succ", proto.CreateResponse{
		Address: crypto.PubkeyToAddress(key.PublicKey).Hex(),
		Key:     hexutil.Encode(crypto.FromECDSA(key))[2:],
	})
}

// Minter
// @Description 获取矿工信息
// @Description
// @Tags 钱包
// @Accept application/json
// @Produce application/json
// @Param object body proto.Minter{} true "请求参数"
// @Success 200 {object}  proto.Response{data=proto.MinterResponse{}}
// @Router /tyche/api/wallet/minter [post]
func (app *App) Minter(c *ginx.Context) {
	var request = new(proto.Minter)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}

	var address = app.conf.Tyche.Account
	var resp proto.MinterResponse
	resp.Address = address

	if request.Balance {
		balance, err := app.backend.BalanceAt(context.Background(), common.HexToAddress(address), nil)
		if err != nil {
			c.Failure(http.StatusInternalServerError, err.Error(), nil)
			return
		}
		resp.Balance = balance.String()
	}

	if request.ChainID {
		chainID, err := app.backend.ChainID(context.Background())
		if err != nil {
			c.Failure(http.StatusInternalServerError, err.Error(), nil)
			return
		}
		resp.ChainID = chainID.String()
	}

	if request.NonceAt {
		latestNonceAt, err := app.backend.NonceAt(context.Background(), common.HexToAddress(address), nil)
		if err != nil {
			c.Failure(http.StatusInternalServerError, err.Error(), nil)
			return
		}

		pendingNonceAt, err := app.backend.PendingNonceAt(context.Background(), common.HexToAddress(address))
		if err != nil {
			c.Failure(http.StatusInternalServerError, err.Error(), nil)
			return
		}

		localNonceAt, err := orders.NonceAt()
		if err != nil {
			c.Failure(http.StatusInternalServerError, err.Error(), nil)
			return
		}

		resp.PendingNonceAt = strconv.FormatUint(pendingNonceAt, 10)
		resp.LatestNonceAt = strconv.FormatUint(latestNonceAt, 10)
		resp.LocalNonceAt = strconv.FormatUint(localNonceAt, 10)
	}
	c.Success(http.StatusOK, "succ", resp)
}

// Sign
// @Description 对数据签名
// @Description
// @Tags 钱包
// @Accept application/json
// @Produce application/json
// @Param object body proto.Sign{} true "请求参数"
// @Success 200 {object}  proto.Response{data=proto.SignResponse{}}
// @Router /tyche/api/wallet/sign [post]
func (app *App) Sign(c *ginx.Context) {

	/*
		{
			"key": "164e15b4d90ee0b2fc2419308ba682eec15971e7600ae79cfcdb29854ae41d2a",  // 用于签名的私钥，默认使用 Minter 私钥
			"types": [                                                                  // 类型数组, 与 values 数组逐一匹配
				"address",
				"address",
				"uint256[]",
				"uint64",
				"uint40",
				"(uint256,string,(address,uint256))"
			],
			"values":[                                                                  			// 值数组
				"0x54987E5F03b503BFD7Df2c84f1981e2a7d3bC505",                           			// From
				"0xeD24FC36d5Ee211Ea25A80239Fb8C4Cfd80f12Ee",                           			// To
				[1, 5],                                                                 			// Token ID 列表
				80001,                                                                  			// 链ID
				1653966243,                                                             			// 时间戳
				[123456, "tupleString",[ "0x54987E5F03b503BFD7Df2c84f1981e2a7d3bC505", 654321]]		// 嵌套的Tuple
			]
		}
	*/
	var request = new(proto.Sign)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}

	if len(request.Types) == 0 {
		c.Failure(http.StatusBadRequest, "无效的参数: types", nil)
		return
	}

	key := request.Key
	if len(key) == 0 {
		key = app.conf.Tyche.PrivateKey
	}

	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		c.Failure(http.StatusInternalServerError, err.Error(), nil)
		return
	}

	messageData, err := types.Encode(request.Types, request.Values)
	if err != nil {
		c.Failure(http.StatusInternalServerError, err.Error(), nil)
		return
	}

	rawMessageHash := crypto.Keccak256Hash(messageData)

	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(rawMessageHash))
	messageHash := crypto.Keccak256Hash([]byte(prefixedMessage), rawMessageHash.Bytes())

	signedData, err := crypto.Sign(messageHash.Bytes(), privateKey)
	if err != nil {
		c.Failure(http.StatusInternalServerError, err.Error(), nil)
		return
	}
	signedData[64] += 27

	c.Success(http.StatusOK, "succ", proto.SignResponse{
		Hash: messageHash.String(),
		Sign: hexutil.Encode(signedData),
	})
}

func (app *App) UserSignHash(c *ginx.Context) {

	var request = new(proto.Sign)
	if err := c.BindJSONEx(request); err != nil {
		c.Failure(http.StatusBadRequest, err.Error(), nil)
		return
	}

	if len(request.Types) == 0 {
		c.Failure(http.StatusBadRequest, "无效的参数: types", nil)
		return
	}

	messageData, err := types.Encode(request.Types, request.Values)
	if err != nil {
		c.Failure(http.StatusInternalServerError, err.Error(), nil)
		return
	}

	rawMessageHash := crypto.Keccak256Hash(messageData)
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(rawMessageHash))
	messageHash := crypto.Keccak256Hash([]byte(prefixedMessage), rawMessageHash.Bytes())
	c.Success(http.StatusOK, "succ", proto.SignResponse{
		Hash: messageHash.String(),
	})
}

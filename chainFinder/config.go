package chainFinder

import (
	"errors"
)

type Config struct {
	Listen                       string `toml:"listen" json:"listen"`
	ReadTimeout                  int    `toml:"readTimeout" json:"readTimeout"`
	WriteTimeout                 int    `toml:"writeTimeout" json:"writeTimeout"`
	MaxHeaderBytes               int    `toml:"maxHeaderBytes" json:"maxHeaderBytes"`
	Account                      string `toml:"account" json:"account"`
	PrivateKey                   string `toml:"privateKey" json:"privateKey"`
	EstimatorJS                  string `toml:"estimatorJS" json:"estimatorJS"`
	ErrorURI                     string `toml:"errorURI" json:"errorURI"`
	FailedURI                    string `toml:"failedURI" json:"failedURI"`
	SucceedURI                   string `toml:"succeedURI" json:"succeedURI"`
	Callback                     string `toml:"callback" json:"callback"`
	ContractTxCount              string `toml:"contractTxCount" json:"contractTxCount"`
	ContractCreationTime         string `toml:"contractCreationTime" json:"contractCreationTime"`
	GetGasPrice                  string `toml:"getGasPrice" json:"getGasPrice"`
	BalanceAt                    string `toml:"balanceAt" json:"balanceAt"`
	CompareBytecodeAndSourceCode string `toml:"compareBytecodeAndSourceCode" json:"compareBytecodeAndSourceCode"`
	GasPriceUpdateInterval       int64  `toml:"gasPriceUpdateInterval" json:"gasPriceUpdateInterval"`
	PrefixChain                  string `toml:"prefixChain" json:"prefixChain"`
	MaxBumpingGasTimes           int64  `toml:"maxBumpingGasTimes" json:"maxBumpingGasTimes"`
	ErrorNumberOfConcurrent      int64  `toml:"errorNumberOfConcurrent" json:"errorNumberOfConcurrent"`
	FailedNumberOfConcurrent     int64  `toml:"failedNumberOfConcurrent" json:"failedNumberOfConcurrent"`
	PendingNumberOfConcurrent    int64  `toml:"pendingNumberOfConcurrent" json:"pendingNumberOfConcurrent"`
	SentNumberOfConcurrent       int64  `toml:"sentNumberOfConcurrent" json:"sentNumberOfConcurrent"`
	SucceedNumberOfConcurrent    int64  `toml:"succeedNumberOfConcurrent" json:"succeedNumberOfConcurrent"`
	PendingRetryInterval         int64  `toml:"pendingRetryInterval" json:"pendingRetryInterval"`
	SentRetryInterval            int64  `toml:"sentRetryInterval" json:"sentRetryInterval"`
	RedisRetryInterval           int64  `toml:"redisRetryInterval" json:"redisRetryInterval"`
	NetworkRetryInterval         int64  `toml:"networkRetryInterval" json:"networkRetryInterval"`
	WaitMinedRetryInterval       int64  `toml:"waitMinedRetryInterval" json:"waitMinedRetryInterval"`
	CallbackRetryInterval        int64  `toml:"callbackRetryInterval" json:"callbackRetryInterval"`
	EnableTLS                    bool   `toml:"enableTLS" json:"enableTLS"`
	CertFile                     string `toml:"certFile" json:"certFile"`
	KeyFile                      string `toml:"keyFile" json:"keyFile"`
}

type ContractConfig struct {
	Addr string `toml:"addr" json:"addr"`
	ABI  string `toml:"abi" json:"abi"`
}

func (c *Config) Init() error {
	if c.ErrorURI == "" {
		return errors.New("errorURI cannot be set to empty")
	}
	if c.EstimatorJS == "" {
		return errors.New("estimatorJS cannot be set to empty")
	}
	if c.FailedURI == "" {
		return errors.New("failedURI cannot be set to empty")
	}
	if c.SucceedURI == "" {
		return errors.New("succeedURI cannot be set to empty")
	}
	if c.Callback == "" {
		return errors.New("callback cannot be set to empty")
	}
	if c.PrefixChain == "" {
		return errors.New("prefixChain cannot be set to empty")
	}
	if c.ContractTxCount == "" {
		return errors.New("contractTxCount cannot be set to empty")
	}
	if c.ContractCreationTime == "" {
		return errors.New("contractCreationTime cannot be set to empty")
	}
	if c.GetGasPrice == "" {
		return errors.New("getGasPrice cannot be set to empty")
	}
	if c.BalanceAt == "" {
		return errors.New("balanceAt cannot be set to empty")
	}
	if c.CompareBytecodeAndSourceCode == "" {
		return errors.New("compareBytecodeAndSourceCode cannot be set to empty")
	}

	if c.GasPriceUpdateInterval < 1 {
		c.GasPriceUpdateInterval = 1
	}
	if c.MaxBumpingGasTimes < 0 {
		c.MaxBumpingGasTimes = 0
	}
	if c.ErrorNumberOfConcurrent < 1 {
		c.ErrorNumberOfConcurrent = 1
	}
	if c.FailedNumberOfConcurrent < 1 {
		c.FailedNumberOfConcurrent = 1
	}
	if c.PendingNumberOfConcurrent < 1 {
		c.PendingNumberOfConcurrent = 1
	}
	if c.SentNumberOfConcurrent < c.MaxBumpingGasTimes*2 {
		c.SentNumberOfConcurrent = c.MaxBumpingGasTimes * 2
	}
	if c.SucceedNumberOfConcurrent < 1 {
		c.SucceedNumberOfConcurrent = 1
	}
	if c.PendingRetryInterval < 1 {
		c.PendingRetryInterval = 1
	}
	if c.SentRetryInterval < 1 {
		c.SentRetryInterval = 1
	}
	if c.RedisRetryInterval < 1 {
		c.RedisRetryInterval = 1
	}
	if c.NetworkRetryInterval < 1 {
		c.NetworkRetryInterval = 1
	}
	if c.WaitMinedRetryInterval < 1 {
		c.WaitMinedRetryInterval = 1
	}
	if c.CallbackRetryInterval < 1 {
		c.CallbackRetryInterval = 1
	}
	if c.EnableTLS {
		switch {
		case c.CertFile == "":
			return errors.New("cretFile cannot be set to empty")
		case c.KeyFile == "":
			return errors.New("keyFile cannot be set to empty")
		}
	}
	return nil
}

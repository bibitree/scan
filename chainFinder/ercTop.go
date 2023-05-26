package chainFinder

type ErcTop struct {
	ContractAddress    string `json:"contractAddress"`
	ContractName       string `json:"contractName"`
	Value              string `json:"value"`
	NewContractAddress string `json:"nonce"`
	ContractTxCount    string `json:"contractTxCount"`
}

type ErcTopHandler func(*ErcTop) error

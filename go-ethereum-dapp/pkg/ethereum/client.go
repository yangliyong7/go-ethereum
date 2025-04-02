package ethereum

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client 以太坊客户端
type Client struct {
	ethClient *ethclient.Client // 以太坊客户端
	ctx       context.Context   // 上下文
}

// NewClient 创建新的以太坊客户端
func NewClient(networkURL string) (*Client, error) {
	// 连接到以太坊网络
	client, err := ethclient.Dial(networkURL)
	if err != nil {
		return nil, err
	}
	
	// 创建客户端
	return &Client{
		ethClient: client,
		ctx:       context.Background(),
	}, nil
}

// EthClient 获取原始以太坊客户端
func (c *Client) EthClient() *ethclient.Client {
	return c.ethClient
}

// GetLatestBlockNumber 获取最新区块号
func (c *Client) GetLatestBlockNumber() (uint64, error) {
	// 获取最新区块
	block, err := c.ethClient.BlockByNumber(c.ctx, nil)
	if err != nil {
		return 0, err
	}
	
	return block.NumberU64(), nil
}

// GetBalance 获取账户余额
func (c *Client) GetBalance(address common.Address) (*big.Int, error) {
	// 获取账户余额
	balance, err := c.ethClient.BalanceAt(c.ctx, address, nil)
	if err != nil {
		return nil, err
	}
	
	return balance, nil
}

// GetTransactionCount 获取账户nonce
func (c *Client) GetTransactionCount(address common.Address) (uint64, error) {
	// 获取账户nonce
	nonce, err := c.ethClient.PendingNonceAt(c.ctx, address)
	if err != nil {
		return 0, err
	}
	
	return nonce, nil
}

// GetChainID 获取链ID
func (c *Client) GetChainID() (*big.Int, error) {
	// 获取链ID
	chainID, err := c.ethClient.ChainID(c.ctx)
	if err != nil {
		return nil, err
	}
	
	return chainID, nil
}

// GetSuggestedGasPrice 获取建议的燃料价格
func (c *Client) GetSuggestedGasPrice() (*big.Int, error) {
	// 获取建议的燃料价格
	gasPrice, err := c.ethClient.SuggestGasPrice(c.ctx)
	if err != nil {
		return nil, err
	}
	
	return gasPrice, nil
}

// CreateTransactionOpts 创建交易选项
func (c *Client) CreateTransactionOpts(privateKey *ecdsa.PrivateKey, gasLimit uint64, gasPrice *big.Int) (*bind.TransactOpts, error) {
	// 获取链ID
	chainID, err := c.GetChainID()
	if err != nil {
		return nil, err
	}
	
	// 创建交易选项
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, err
	}
	
	// 如果未指定燃料价格，获取建议的燃料价格
	if gasPrice == nil {
		gasPrice, err = c.GetSuggestedGasPrice()
		if err != nil {
			return nil, err
		}
	}
	
	// 设置燃料限制和燃料价格
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice
	
	// 获取账户nonce
	nonce, err := c.GetTransactionCount(auth.From)
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	
	return auth, nil
}

// WaitForTransaction 等待交易确认
func (c *Client) WaitForTransaction(tx *types.Transaction) (*types.Receipt, error) {
	// 等待交易被打包到区块中
	receipt, err := bind.WaitMined(c.ctx, c.ethClient, tx)
	if err != nil {
		return nil, err
	}
	
	return receipt, nil
}

// Close 关闭客户端连接
func (c *Client) Close() {
	c.ethClient.Close()
}

// PrivateKeyFromHex 从十六进制字符串创建私钥
func PrivateKeyFromHex(hexKey string) (*ecdsa.PrivateKey, error) {
	// 解析私钥
	return crypto.HexToECDSA(hexKey)
}

// AddressFromPrivateKey 从私钥获取地址
func AddressFromPrivateKey(privateKey *ecdsa.PrivateKey) common.Address {
	// 获取公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("无法获取公钥")
	}
	
	// 从公钥获取地址
	return crypto.PubkeyToAddress(*publicKeyECDSA)
} 
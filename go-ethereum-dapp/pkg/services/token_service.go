package services

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/user/go-ethereum-dapp/pkg/config"
	"github.com/user/go-ethereum-dapp/pkg/contracts"
	"github.com/user/go-ethereum-dapp/pkg/ethereum"
)

// TokenService 代币服务
type TokenService struct {
	client      *ethereum.Client     // 以太坊客户端
	config      *config.Config       // 配置
	tokenContract *contracts.SimpleToken // 代币合约
}

// NewTokenService 创建新的代币服务
func NewTokenService(client *ethereum.Client, cfg *config.Config) (*TokenService, error) {
	service := &TokenService{
		client: client,
		config: cfg,
	}
	
	// 如果配置了代币合约地址，则加载合约
	if cfg.TokenContractAddress != "" {
		err := service.LoadContract(common.HexToAddress(cfg.TokenContractAddress))
		if err != nil {
			return nil, fmt.Errorf("加载代币合约失败: %w", err)
		}
	}
	
	return service, nil
}

// LoadContract 加载代币合约
func (s *TokenService) LoadContract(contractAddress common.Address) error {
	// 创建合约实例
	tokenContract, err := contracts.NewSimpleToken(contractAddress, s.client.EthClient())
	if err != nil {
		return err
	}
	
	s.tokenContract = tokenContract
	return nil
}

// DeployContract 部署代币合约
func (s *TokenService) DeployContract(privateKey string, name, symbol string, decimals uint8, initialSupply *big.Int) (common.Address, *types.Transaction, error) {
	// 解析私钥
	pk, err := ethereum.PrivateKeyFromHex(privateKey)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("解析私钥失败: %w", err)
	}
	
	// 创建交易选项
	auth, err := s.client.CreateTransactionOpts(pk, s.config.GasLimit, s.config.GasPrice)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("创建交易选项失败: %w", err)
	}
	
	// 部署合约
	address, tx, instance, err := contracts.DeploySimpleToken(auth, s.client.EthClient(), name, symbol, decimals, initialSupply)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("部署合约失败: %w", err)
	}
	
	// 保存合约实例
	s.tokenContract = instance
	
	return address, tx, nil
}

// GetContractAddress 获取代币合约地址
func (s *TokenService) GetContractAddress() (common.Address, error) {
	if s.tokenContract == nil {
		return common.Address{}, fmt.Errorf("合约未加载")
	}
	
	return s.tokenContract.address, nil
}

// GetTokenInfo 获取代币信息
func (s *TokenService) GetTokenInfo() (name string, symbol string, decimals uint8, totalSupply *big.Int, err error) {
	if s.tokenContract == nil {
		return "", "", 0, nil, fmt.Errorf("合约未加载")
	}
	
	// 创建调用选项
	opts := &bind.CallOpts{}
	
	// 获取代币名称
	name, err = s.tokenContract.Name(opts)
	if err != nil {
		return "", "", 0, nil, fmt.Errorf("获取代币名称失败: %w", err)
	}
	
	// 获取代币符号
	symbol, err = s.tokenContract.Symbol(opts)
	if err != nil {
		return "", "", 0, nil, fmt.Errorf("获取代币符号失败: %w", err)
	}
	
	// 获取代币精度
	decimals, err = s.tokenContract.Decimals(opts)
	if err != nil {
		return "", "", 0, nil, fmt.Errorf("获取代币精度失败: %w", err)
	}
	
	// 获取代币总供应量
	totalSupply, err = s.tokenContract.TotalSupply(opts)
	if err != nil {
		return "", "", 0, nil, fmt.Errorf("获取代币总供应量失败: %w", err)
	}
	
	return name, symbol, decimals, totalSupply, nil
}

// GetBalance 获取代币余额
func (s *TokenService) GetBalance(address common.Address) (*big.Int, error) {
	if s.tokenContract == nil {
		return nil, fmt.Errorf("合约未加载")
	}
	
	// 创建调用选项
	opts := &bind.CallOpts{}
	
	// 获取代币余额
	balance, err := s.tokenContract.BalanceOf(opts, address)
	if err != nil {
		return nil, fmt.Errorf("获取代币余额失败: %w", err)
	}
	
	return balance, nil
}

// Transfer 转账代币
func (s *TokenService) Transfer(privateKey string, to common.Address, amount *big.Int) (*types.Transaction, error) {
	if s.tokenContract == nil {
		return nil, fmt.Errorf("合约未加载")
	}
	
	// 解析私钥
	pk, err := ethereum.PrivateKeyFromHex(privateKey)
	if err != nil {
		return nil, fmt.Errorf("解析私钥失败: %w", err)
	}
	
	// 创建交易选项
	auth, err := s.client.CreateTransactionOpts(pk, s.config.GasLimit, s.config.GasPrice)
	if err != nil {
		return nil, fmt.Errorf("创建交易选项失败: %w", err)
	}
	
	// 转账代币
	tx, err := s.tokenContract.Transfer(auth, to, amount)
	if err != nil {
		return nil, fmt.Errorf("转账代币失败: %w", err)
	}
	
	return tx, nil
}

// Approve 授权代币
func (s *TokenService) Approve(privateKey string, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	if s.tokenContract == nil {
		return nil, fmt.Errorf("合约未加载")
	}
	
	// 解析私钥
	pk, err := ethereum.PrivateKeyFromHex(privateKey)
	if err != nil {
		return nil, fmt.Errorf("解析私钥失败: %w", err)
	}
	
	// 创建交易选项
	auth, err := s.client.CreateTransactionOpts(pk, s.config.GasLimit, s.config.GasPrice)
	if err != nil {
		return nil, fmt.Errorf("创建交易选项失败: %w", err)
	}
	
	// 授权代币
	tx, err := s.tokenContract.Approve(auth, spender, amount)
	if err != nil {
		return nil, fmt.Errorf("授权代币失败: %w", err)
	}
	
	return tx, nil
}

// WaitForTransaction 等待交易确认
func (s *TokenService) WaitForTransaction(tx *types.Transaction) (*types.Receipt, error) {
	return s.client.WaitForTransaction(tx)
} 
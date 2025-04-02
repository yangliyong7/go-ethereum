package config

import (
	"math/big"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

// Config 应用程序配置
type Config struct {
	// 以太坊网络配置
	EthereumNetworkURL string // 以太坊节点URL
	ChainID            *big.Int // 链ID
	GasLimit           uint64   // 燃料限制
	GasPrice           *big.Int // 燃料价格
	
	// 智能合约地址
	TokenContractAddress   string // 代币合约地址
	VotingContractAddress  string // 投票合约地址
	
	// 账户配置
	PrivateKey string // 私钥（用于签名交易）
}

// LoadConfig 加载配置信息
func LoadConfig() (*Config, error) {
	// 获取项目根目录
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "../..")
	
	// 加载.env文件
	err := godotenv.Load(filepath.Join(basePath, ".env"))
	if err != nil {
		// 如果.env文件不存在，继续执行，使用环境变量
	}
	
	// 默认使用本地开发网络
	ethereumURL := getEnv("ETHEREUM_URL", "http://localhost:8545")
	chainIDStr := getEnv("CHAIN_ID", "1337") // 默认使用Ganache的链ID
	
	// 解析链ID
	chainID := new(big.Int)
	chainID.SetString(chainIDStr, 10)
	
	// 创建配置
	cfg := &Config{
		EthereumNetworkURL:    ethereumURL,
		ChainID:               chainID,
		GasLimit:              uint64(3000000),                  // 默认燃料限制
		GasPrice:              big.NewInt(20000000000),          // 默认20 Gwei
		TokenContractAddress:  getEnv("TOKEN_CONTRACT_ADDRESS", ""),
		VotingContractAddress: getEnv("VOTING_CONTRACT_ADDRESS", ""),
		PrivateKey:            getEnv("PRIVATE_KEY", ""),
	}
	
	return cfg, nil
}

// 从环境变量获取值，如果不存在则使用默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
} 
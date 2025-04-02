package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/user/go-ethereum-dapp/pkg/config"
	"github.com/user/go-ethereum-dapp/pkg/ethereum"
	"github.com/user/go-ethereum-dapp/pkg/services"
)

func main() {
	// 捕获退出信号
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// 优雅退出
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("接收到退出信号，正在关闭...")
		cancel()
	}()
	
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	
	// 创建以太坊客户端
	client, err := ethereum.NewClient(cfg.EthereumNetworkURL)
	if err != nil {
		log.Fatalf("连接以太坊网络失败: %v", err)
	}
	defer client.Close()
	
	// 运行演示
	if err := runDemo(ctx, client, cfg); err != nil {
		log.Fatalf("运行演示失败: %v", err)
	}
}

// runDemo 运行完整演示
func runDemo(ctx context.Context, client *ethereum.Client, cfg *config.Config) error {
	log.Println("开始以太坊DApp演示")
	
	// 检查私钥是否配置
	if cfg.PrivateKey == "" {
		return fmt.Errorf("请在配置中设置私钥")
	}
	
	// 解析私钥
	privateKey, err := ethereum.PrivateKeyFromHex(cfg.PrivateKey)
	if err != nil {
		return fmt.Errorf("解析私钥失败: %v", err)
	}
	
	// 获取账户地址
	address := ethereum.AddressFromPrivateKey(privateKey)
	log.Printf("使用账户地址: %s", address.Hex())
	
	// 获取账户ETH余额
	balance, err := client.GetBalance(address)
	if err != nil {
		return fmt.Errorf("获取ETH余额失败: %v", err)
	}
	log.Printf("ETH余额: %s Wei", balance.String())
	
	// 创建服务
	tokenService, err := services.NewTokenService(client, cfg)
	if err != nil {
		return fmt.Errorf("创建代币服务失败: %v", err)
	}
	
	votingService, err := services.NewVotingService(client, cfg)
	if err != nil {
		return fmt.Errorf("创建投票服务失败: %v", err)
	}
	
	// 部署代币合约
	var tokenAddress common.Address
	var tokenTx *types.Transaction
	
	if cfg.TokenContractAddress == "" {
		log.Println("正在部署代币合约...")
		// 部署代币合约
		tokenAddress, tokenTx, err = tokenService.DeployContract(
			cfg.PrivateKey,
			"示例代币",     // 名称
			"EXT",      // 符号
			18,         // 精度
			big.NewInt(1000000), // 初始供应量（实际供应量 = 1000000 * 10^18）
		)
		if err != nil {
			return fmt.Errorf("部署代币合约失败: %v", err)
		}
		
		log.Printf("代币合约部署交易已提交，交易哈希: %s", tokenTx.Hash().Hex())
		log.Println("等待交易确认...")
		
		// 等待交易确认
		receipt, err := tokenService.WaitForTransaction(tokenTx)
		if err != nil {
			return fmt.Errorf("等待交易确认失败: %v", err)
		}
		
		if receipt.Status == 0 {
			return fmt.Errorf("部署代币合约失败，交易被回滚")
		}
		
		log.Printf("代币合约已部署，地址: %s", tokenAddress.Hex())
		
		// 更新配置
		cfg.TokenContractAddress = tokenAddress.Hex()
	} else {
		tokenAddress = common.HexToAddress(cfg.TokenContractAddress)
		log.Printf("使用已部署的代币合约，地址: %s", tokenAddress.Hex())
	}
	
	// 获取代币信息
	name, symbol, decimals, totalSupply, err := tokenService.GetTokenInfo()
	if err != nil {
		return fmt.Errorf("获取代币信息失败: %v", err)
	}
	
	log.Printf("代币名称: %s", name)
	log.Printf("代币符号: %s", symbol)
	log.Printf("代币精度: %d", decimals)
	log.Printf("代币总供应量: %s", totalSupply.String())
	
	// 获取代币余额
	tokenBalance, err := tokenService.GetBalance(address)
	if err != nil {
		return fmt.Errorf("获取代币余额失败: %v", err)
	}
	
	log.Printf("%s余额: %s", symbol, tokenBalance.String())
	
	// 创建其他测试账户的地址
	testAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")
	
	// 转账代币给测试账户
	log.Printf("转账100代币到地址: %s", testAddress.Hex())
	
	transferAmount := new(big.Int).Mul(big.NewInt(100), new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))
	transferTx, err := tokenService.Transfer(cfg.PrivateKey, testAddress, transferAmount)
	if err != nil {
		return fmt.Errorf("转账代币失败: %v", err)
	}
	
	log.Printf("代币转账交易已提交，交易哈希: %s", transferTx.Hash().Hex())
	log.Println("等待交易确认...")
	
	// 等待交易确认
	transferReceipt, err := tokenService.WaitForTransaction(transferTx)
	if err != nil {
		return fmt.Errorf("等待交易确认失败: %v", err)
	}
	
	if transferReceipt.Status == 0 {
		return fmt.Errorf("转账代币失败，交易被回滚")
	}
	
	log.Println("代币转账成功")
	
	// 检查余额变化
	newTokenBalance, err := tokenService.GetBalance(address)
	if err != nil {
		return fmt.Errorf("获取代币余额失败: %v", err)
	}
	
	log.Printf("%s新余额: %s", symbol, newTokenBalance.String())
	
	testTokenBalance, err := tokenService.GetBalance(testAddress)
	if err != nil {
		return fmt.Errorf("获取测试账户代币余额失败: %v", err)
	}
	
	log.Printf("测试账户%s余额: %s", symbol, testTokenBalance.String())
	
	// 部署投票合约
	var votingAddress common.Address
	var votingTx *types.Transaction
	
	if cfg.VotingContractAddress == "" {
		log.Println("正在部署投票合约...")
		
		// 投票选项
		voteOptions := []string{"选项A", "选项B", "选项C"}
		
		// 部署投票合约（投票持续10分钟）
		votingAddress, votingTx, err = votingService.DeployContract(
			cfg.PrivateKey,
			"示例投票",                 // 投票主题
			tokenAddress,            // 代币合约地址
			big.NewInt(10),         // 持续时间（分钟）
			voteOptions,            // 投票选项
		)
		if err != nil {
			return fmt.Errorf("部署投票合约失败: %v", err)
		}
		
		log.Printf("投票合约部署交易已提交，交易哈希: %s", votingTx.Hash().Hex())
		log.Println("等待交易确认...")
		
		// 等待交易确认
		votingReceipt, err := votingService.WaitForTransaction(votingTx)
		if err != nil {
			return fmt.Errorf("等待交易确认失败: %v", err)
		}
		
		if votingReceipt.Status == 0 {
			return fmt.Errorf("部署投票合约失败，交易被回滚")
		}
		
		log.Printf("投票合约已部署，地址: %s", votingAddress.Hex())
		
		// 更新配置
		cfg.VotingContractAddress = votingAddress.Hex()
	} else {
		votingAddress = common.HexToAddress(cfg.VotingContractAddress)
		log.Printf("使用已部署的投票合约，地址: %s", votingAddress.Hex())
	}
	
	// 获取投票信息
	topic, endTime, optionsCount, tokenContractAddress, err := votingService.GetVotingInfo()
	if err != nil {
		return fmt.Errorf("获取投票信息失败: %v", err)
	}
	
	log.Printf("投票主题: %s", topic)
	log.Printf("结束时间: %s", time.Unix(endTime.Int64(), 0).Format("2006-01-02 15:04:05"))
	log.Printf("选项数量: %d", optionsCount.Int64())
	log.Printf("代币合约地址: %s", tokenContractAddress.Hex())
	
	// 获取投票选项
	options, err := votingService.GetOptions()
	if err != nil {
		return fmt.Errorf("获取投票选项失败: %v", err)
	}
	
	log.Println("投票选项:")
	for _, option := range options {
		log.Printf("  ID: %d, 名称: %s, 得票数: %s", option.ID.Int64(), option.Name, option.VoteCount.String())
	}
	
	// 进行投票
	log.Println("进行投票...")
	voteTx, err := votingService.Vote(cfg.PrivateKey, big.NewInt(0)) // 投票给选项A
	if err != nil {
		return fmt.Errorf("投票失败: %v", err)
	}
	
	log.Printf("投票交易已提交，交易哈希: %s", voteTx.Hash().Hex())
	log.Println("等待交易确认...")
	
	// 等待交易确认
	voteReceipt, err := votingService.WaitForTransaction(voteTx)
	if err != nil {
		return fmt.Errorf("等待交易确认失败: %v", err)
	}
	
	if voteReceipt.Status == 0 {
		return fmt.Errorf("投票失败，交易被回滚")
	}
	
	log.Println("投票成功")
	
	// 检查是否已投票
	hasVoted, err := votingService.HasVoted(address)
	if err != nil {
		return fmt.Errorf("检查是否已投票失败: %v", err)
	}
	
	log.Printf("账户是否已投票: %v", hasVoted)
	
	if hasVoted {
		votedOption, err := votingService.GetVotedOption(address)
		if err != nil {
			return fmt.Errorf("获取投票选项失败: %v", err)
		}
		
		log.Printf("账户投票的选项ID: %d", votedOption.Int64())
	}
	
	// 获取实时投票结果
	log.Println("投票结果:")
	updatedOptions, err := votingService.GetOptions()
	if err != nil {
		return fmt.Errorf("获取投票选项失败: %v", err)
	}
	
	for _, option := range updatedOptions {
		log.Printf("  ID: %d, 名称: %s, 得票数: %s", option.ID.Int64(), option.Name, option.VoteCount.String())
	}
	
	// 获取获胜选项
	winningID, winningName, winningVotes, err := votingService.GetWinningOption()
	if err != nil {
		return fmt.Errorf("获取获胜选项失败: %v", err)
	}
	
	log.Printf("当前获胜选项: ID %d, 名称: %s, 得票数: %s", winningID.Int64(), winningName, winningVotes.String())
	
	// 检查投票是否已结束
	isEnded, err := votingService.IsVotingEnded()
	if err != nil {
		return fmt.Errorf("检查投票是否已结束失败: %v", err)
	}
	
	log.Printf("投票是否已结束: %v", isEnded)
	
	if !isEnded {
		remainingTime, err := votingService.RemainingTime()
		if err != nil {
			return fmt.Errorf("获取剩余时间失败: %v", err)
		}
		
		log.Printf("剩余时间: %d 秒", remainingTime.Int64())
	}
	
	log.Println("以太坊DApp演示完成")
	
	return nil
} 
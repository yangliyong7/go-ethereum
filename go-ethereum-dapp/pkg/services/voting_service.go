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

// VotingService 投票服务
type VotingService struct {
	client       *ethereum.Client // 以太坊客户端
	config       *config.Config   // 配置
	votingContract *contracts.Voting // 投票合约
}

// NewVotingService 创建新的投票服务
func NewVotingService(client *ethereum.Client, cfg *config.Config) (*VotingService, error) {
	service := &VotingService{
		client: client,
		config: cfg,
	}
	
	// 如果配置了投票合约地址，则加载合约
	if cfg.VotingContractAddress != "" {
		err := service.LoadContract(common.HexToAddress(cfg.VotingContractAddress))
		if err != nil {
			return nil, fmt.Errorf("加载投票合约失败: %w", err)
		}
	}
	
	return service, nil
}

// LoadContract 加载投票合约
func (s *VotingService) LoadContract(contractAddress common.Address) error {
	// 创建合约实例
	votingContract, err := contracts.NewVoting(contractAddress, s.client.EthClient())
	if err != nil {
		return err
	}
	
	s.votingContract = votingContract
	return nil
}

// DeployContract 部署投票合约
func (s *VotingService) DeployContract(
	privateKey string,
	votingTopic string,
	tokenAddress common.Address,
	durationInMinutes *big.Int,
	options []string,
) (common.Address, *types.Transaction, error) {
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
	address, tx, instance, err := contracts.DeployVoting(
		auth,
		s.client.EthClient(),
		votingTopic,
		tokenAddress,
		durationInMinutes,
		options,
	)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("部署合约失败: %w", err)
	}
	
	// 保存合约实例
	s.votingContract = instance
	
	return address, tx, nil
}

// GetContractAddress 获取投票合约地址
func (s *VotingService) GetContractAddress() (common.Address, error) {
	if s.votingContract == nil {
		return common.Address{}, fmt.Errorf("合约未加载")
	}
	
	return s.votingContract.address, nil
}

// GetVotingInfo 获取投票信息
func (s *VotingService) GetVotingInfo() (topic string, endTime *big.Int, optionsCount *big.Int, tokenAddress common.Address, err error) {
	if s.votingContract == nil {
		return "", nil, nil, common.Address{}, fmt.Errorf("合约未加载")
	}
	
	// 创建调用选项
	opts := &bind.CallOpts{}
	
	// 获取投票主题
	topic, err = s.votingContract.VotingTopic(opts)
	if err != nil {
		return "", nil, nil, common.Address{}, fmt.Errorf("获取投票主题失败: %w", err)
	}
	
	// 获取结束时间
	endTime, err = s.votingContract.EndTime(opts)
	if err != nil {
		return "", nil, nil, common.Address{}, fmt.Errorf("获取结束时间失败: %w", err)
	}
	
	// 获取选项数量
	optionsCount, err = s.votingContract.OptionsCount(opts)
	if err != nil {
		return "", nil, nil, common.Address{}, fmt.Errorf("获取选项数量失败: %w", err)
	}
	
	// 获取代币合约地址
	tokenAddress, err = s.votingContract.TokenContract(opts)
	if err != nil {
		return "", nil, nil, common.Address{}, fmt.Errorf("获取代币合约地址失败: %w", err)
	}
	
	return topic, endTime, optionsCount, tokenAddress, nil
}

// GetOptions 获取所有投票选项
func (s *VotingService) GetOptions() ([]struct {
	ID        *big.Int
	Name      string
	VoteCount *big.Int
}, error) {
	if s.votingContract == nil {
		return nil, fmt.Errorf("合约未加载")
	}
	
	// 创建调用选项
	opts := &bind.CallOpts{}
	
	// 获取选项数量
	optionsCount, err := s.votingContract.OptionsCount(opts)
	if err != nil {
		return nil, fmt.Errorf("获取选项数量失败: %w", err)
	}
	
	// 获取所有选项
	var options []struct {
		ID        *big.Int
		Name      string
		VoteCount *big.Int
	}
	
	for i := int64(0); i < optionsCount.Int64(); i++ {
		optionID := big.NewInt(i)
		option, err := s.votingContract.GetOption(opts, optionID)
		if err != nil {
			return nil, fmt.Errorf("获取选项信息失败: %w", err)
		}
		
		options = append(options, struct {
			ID        *big.Int
			Name      string
			VoteCount *big.Int
		}{
			ID:        optionID,
			Name:      option.Name,
			VoteCount: option.VoteCount,
		})
	}
	
	return options, nil
}

// IsVotingEnded 检查投票是否已结束
func (s *VotingService) IsVotingEnded() (bool, error) {
	if s.votingContract == nil {
		return false, fmt.Errorf("合约未加载")
	}
	
	// 创建调用选项
	opts := &bind.CallOpts{}
	
	// 检查投票是否已结束
	return s.votingContract.IsVotingEnded(opts)
}

// RemainingTime 获取剩余投票时间
func (s *VotingService) RemainingTime() (*big.Int, error) {
	if s.votingContract == nil {
		return nil, fmt.Errorf("合约未加载")
	}
	
	// 创建调用选项
	opts := &bind.CallOpts{}
	
	// 获取剩余投票时间
	return s.votingContract.RemainingTime(opts)
}

// HasVoted 检查地址是否已投票
func (s *VotingService) HasVoted(address common.Address) (bool, error) {
	if s.votingContract == nil {
		return false, fmt.Errorf("合约未加载")
	}
	
	// 创建调用选项
	opts := &bind.CallOpts{}
	
	// 检查地址是否已投票
	return s.votingContract.HasVoted(opts, address)
}

// GetVotedOption 获取地址投票的选项
func (s *VotingService) GetVotedOption(address common.Address) (*big.Int, error) {
	if s.votingContract == nil {
		return nil, fmt.Errorf("合约未加载")
	}
	
	// 创建调用选项
	opts := &bind.CallOpts{}
	
	// 获取地址投票的选项
	return s.votingContract.VotedOption(opts, address)
}

// Vote 进行投票
func (s *VotingService) Vote(privateKey string, optionID *big.Int) (*types.Transaction, error) {
	if s.votingContract == nil {
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
	
	// 投票
	tx, err := s.votingContract.Vote(auth, optionID)
	if err != nil {
		return nil, fmt.Errorf("投票失败: %w", err)
	}
	
	return tx, nil
}

// GetWinningOption 获取获胜选项
func (s *VotingService) GetWinningOption() (id *big.Int, name string, voteCount *big.Int, err error) {
	if s.votingContract == nil {
		return nil, "", nil, fmt.Errorf("合约未加载")
	}
	
	// 创建调用选项
	opts := &bind.CallOpts{}
	
	// 获取获胜选项ID
	id, err = s.votingContract.WinningOptionId(opts)
	if err != nil {
		return nil, "", nil, fmt.Errorf("获取获胜选项ID失败: %w", err)
	}
	
	// 获取获胜选项名称
	name, err = s.votingContract.WinningOptionName(opts)
	if err != nil {
		return nil, "", nil, fmt.Errorf("获取获胜选项名称失败: %w", err)
	}
	
	// 获取获胜选项得票数
	option, err := s.votingContract.GetOption(opts, id)
	if err != nil {
		return nil, "", nil, fmt.Errorf("获取获胜选项信息失败: %w", err)
	}
	
	voteCount = option.VoteCount
	
	return id, name, voteCount, nil
}

// WaitForTransaction 等待交易确认
func (s *VotingService) WaitForTransaction(tx *types.Transaction) (*types.Receipt, error) {
	return s.client.WaitForTransaction(tx)
} 
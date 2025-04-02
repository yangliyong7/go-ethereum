// Code generated - DO NOT EDIT.
// This file is a manually written Go binding of the Voting contract.

package contracts

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// 投票合约ABI
var VotingABI = `[
	{
		"inputs": [
			{"internalType":"string","name":"_votingTopic","type":"string"},
			{"internalType":"address","name":"_tokenAddress","type":"address"},
			{"internalType":"uint256","name":"_durationInMinutes","type":"uint256"},
			{"internalType":"string[]","name":"_options","type":"string[]"}
		],
		"stateMutability":"nonpayable",
		"type":"constructor"
	},
	{
		"anonymous":false,
		"inputs":[
			{"indexed":true,"internalType":"address","name":"voter","type":"address"},
			{"indexed":false,"internalType":"uint256","name":"optionId","type":"uint256"}
		],
		"name":"Voted",
		"type":"event"
	},
	{
		"inputs":[],
		"name":"endTime",
		"outputs":[{"internalType":"uint256","name":"","type":"uint256"}],
		"stateMutability":"view",
		"type":"function"
	},
	{
		"inputs":[{"internalType":"address","name":"","type":"address"}],
		"name":"hasVoted",
		"outputs":[{"internalType":"bool","name":"","type":"bool"}],
		"stateMutability":"view",
		"type":"function"
	},
	{
		"inputs":[],
		"name":"isVotingEnded",
		"outputs":[{"internalType":"bool","name":"","type":"bool"}],
		"stateMutability":"view",
		"type":"function"
	},
	{
		"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],
		"name":"options",
		"outputs":[
			{"internalType":"string","name":"name","type":"string"},
			{"internalType":"uint256","name":"voteCount","type":"uint256"}
		],
		"stateMutability":"view",
		"type":"function"
	},
	{
		"inputs":[],
		"name":"optionsCount",
		"outputs":[{"internalType":"uint256","name":"","type":"uint256"}],
		"stateMutability":"view",
		"type":"function"
	},
	{
		"inputs":[],
		"name":"remainingTime",
		"outputs":[{"internalType":"uint256","name":"","type":"uint256"}],
		"stateMutability":"view",
		"type":"function"
	},
	{
		"inputs":[],
		"name":"tokenContract",
		"outputs":[{"internalType":"contract SimpleToken","name":"","type":"address"}],
		"stateMutability":"view",
		"type":"function"
	},
	{
		"inputs":[{"internalType":"uint256","name":"_optionId","type":"uint256"}],
		"name":"vote",
		"outputs":[],
		"stateMutability":"nonpayable",
		"type":"function"
	},
	{
		"inputs":[{"internalType":"address","name":"","type":"address"}],
		"name":"votedOption",
		"outputs":[{"internalType":"uint256","name":"","type":"uint256"}],
		"stateMutability":"view",
		"type":"function"
	},
	{
		"inputs":[],
		"name":"votingTopic",
		"outputs":[{"internalType":"string","name":"","type":"string"}],
		"stateMutability":"view",
		"type":"function"
	},
	{
		"inputs":[],
		"name":"winningOptionId",
		"outputs":[{"internalType":"uint256","name":"","type":"uint256"}],
		"stateMutability":"view",
		"type":"function"
	},
	{
		"inputs":[],
		"name":"winningOptionName",
		"outputs":[{"internalType":"string","name":"","type":"string"}],
		"stateMutability":"view",
		"type":"function"
	}
]`

// Voting Go绑定结构体
type Voting struct {
	address common.Address      // 合约地址
	backend bind.ContractBackend // 合约后端
	abi     abi.ABI            // 合约ABI
}

// NewVoting 创建新的Voting实例
func NewVoting(address common.Address, backend bind.ContractBackend) (*Voting, error) {
	// 解析合约ABI
	votingAbi, err := abi.JSON(strings.NewReader(VotingABI))
	if err != nil {
		return nil, err
	}
	
	// 创建合约实例
	voting := &Voting{
		address: address,
		backend: backend,
		abi:     votingAbi,
	}
	
	return voting, nil
}

// 部署新的投票合约
func DeployVoting(
	auth *bind.TransactOpts,
	backend bind.ContractBackend,
	votingTopic string,
	tokenAddress common.Address,
	durationInMinutes *big.Int,
	options []string,
) (common.Address, *types.Transaction, *Voting, error) {
	// 解析ABI
	parsed, err := abi.JSON(strings.NewReader(VotingABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	
	// 编码部署参数
	input, err := parsed.Pack("", votingTopic, tokenAddress, durationInMinutes, options)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	
	// 部署合约
	address, tx, err := bind.DeployContract(auth, parsed, input, backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	
	// 创建合约实例
	contract := &Voting{
		address: address,
		backend: backend,
		abi:     parsed,
	}
	
	return address, tx, contract, nil
}

// VotingTopic 获取投票主题
func (v *Voting) VotingTopic(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := v.call(opts, &out, "votingTopic")
	if err != nil {
		return "", err
	}
	return out[0].(string), nil
}

// EndTime 获取投票结束时间
func (v *Voting) EndTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := v.call(opts, &out, "endTime")
	if err != nil {
		return nil, err
	}
	return out[0].(*big.Int), nil
}

// OptionsCount 获取选项数量
func (v *Voting) OptionsCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := v.call(opts, &out, "optionsCount")
	if err != nil {
		return nil, err
	}
	return out[0].(*big.Int), nil
}

// GetOption 获取选项信息
func (v *Voting) GetOption(opts *bind.CallOpts, optionId *big.Int) (struct {
	Name      string
	VoteCount *big.Int
}, error) {
	var out []interface{}
	err := v.call(opts, &out, "options", optionId)
	if err != nil {
		return struct {
			Name      string
			VoteCount *big.Int
		}{}, err
	}
	
	result := struct {
		Name      string
		VoteCount *big.Int
	}{
		Name:      out[0].(string),
		VoteCount: out[1].(*big.Int),
	}
	
	return result, nil
}

// HasVoted 查询地址是否已投票
func (v *Voting) HasVoted(opts *bind.CallOpts, voter common.Address) (bool, error) {
	var out []interface{}
	err := v.call(opts, &out, "hasVoted", voter)
	if err != nil {
		return false, err
	}
	return out[0].(bool), nil
}

// VotedOption 获取地址投票的选项ID
func (v *Voting) VotedOption(opts *bind.CallOpts, voter common.Address) (*big.Int, error) {
	var out []interface{}
	err := v.call(opts, &out, "votedOption", voter)
	if err != nil {
		return nil, err
	}
	return out[0].(*big.Int), nil
}

// TokenContract 获取代币合约地址
func (v *Voting) TokenContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := v.call(opts, &out, "tokenContract")
	if err != nil {
		return common.Address{}, err
	}
	return out[0].(common.Address), nil
}

// IsVotingEnded 查询投票是否已结束
func (v *Voting) IsVotingEnded(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := v.call(opts, &out, "isVotingEnded")
	if err != nil {
		return false, err
	}
	return out[0].(bool), nil
}

// RemainingTime 获取剩余投票时间
func (v *Voting) RemainingTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := v.call(opts, &out, "remainingTime")
	if err != nil {
		return nil, err
	}
	return out[0].(*big.Int), nil
}

// WinningOptionId 获取获胜选项ID
func (v *Voting) WinningOptionId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := v.call(opts, &out, "winningOptionId")
	if err != nil {
		return nil, err
	}
	return out[0].(*big.Int), nil
}

// WinningOptionName 获取获胜选项名称
func (v *Voting) WinningOptionName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := v.call(opts, &out, "winningOptionName")
	if err != nil {
		return "", err
	}
	return out[0].(string), nil
}

// Vote 进行投票
func (v *Voting) Vote(opts *bind.TransactOpts, optionId *big.Int) (*types.Transaction, error) {
	return v.transact(opts, "vote", optionId)
}

// 内部调用方法
func (v *Voting) call(opts *bind.CallOpts, result *[]interface{}, method string, args ...interface{}) error {
	if opts == nil {
		opts = &bind.CallOpts{}
	}
	
	// 打包方法调用数据
	input, err := v.abi.Pack(method, args...)
	if err != nil {
		return err
	}
	
	// 创建调用消息
	msg := ethereum.CallMsg{
		From:     opts.From,
		To:       &v.address,
		Data:     input,
		GasPrice: opts.GasPrice,
	}
	
	// 执行调用
	output, err := v.backend.CallContract(opts.Context, msg, opts.BlockNumber)
	if err != nil {
		return err
	}
	
	// 解包结果
	return v.abi.UnpackIntoInterface(result, method, output)
}

// 内部交易方法
func (v *Voting) transact(opts *bind.TransactOpts, method string, args ...interface{}) (*types.Transaction, error) {
	// 打包方法调用数据
	input, err := v.abi.Pack(method, args...)
	if err != nil {
		return nil, err
	}
	
	// 发送交易
	return bind.DoTransaction(opts, v.backend, &v.address, input)
}

// FilterVoted 过滤Voted事件
func (v *Voting) FilterVoted(opts *bind.FilterOpts, voter []common.Address) (*VotingVotedIterator, error) {
	// 创建过滤器
	var voterRule []interface{}
	for _, addr := range voter {
		voterRule = append(voterRule, addr)
	}

	logs, sub, err := v.contract().FilterLogs(opts.Context, "Voted", voterRule)
	if err != nil {
		return nil, err
	}
	return &VotingVotedIterator{logs: logs, sub: sub}, nil
}

// WatchVoted 监听Voted事件
func (v *Voting) WatchVoted(opts *bind.WatchOpts, sink chan<- *VotingVoted, voter []common.Address) (event.Subscription, error) {
	// 创建过滤器
	var voterRule []interface{}
	for _, addr := range voter {
		voterRule = append(voterRule, addr)
	}

	logs, sub, err := v.contract().WatchLogs(opts.Context, "Voted", sink, voterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// 解析事件
				event := new(VotingVoted)
				if err := v.abi.UnpackIntoInterface(event, "Voted", log.Data); err != nil {
					return err
				}
				event.Voter = common.BytesToAddress(log.Topics[1].Bytes())
				event.Raw = log

				select {
				case sink <- event:
				case <-quit:
					return nil
				}
			case <-quit:
				return nil
			}
		}
	}), nil
}

// VotingVoted 投票事件
type VotingVoted struct {
	Voter    common.Address
	OptionId *big.Int
	Raw      types.Log
}

// VotingVotedIterator 投票事件迭代器
type VotingVotedIterator struct {
	logs []types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VotingVotedIterator) Next() bool {
	// 检查是否已完成或出错
	if it.done || it.fail != nil {
		return false
	}
	// 检查是否还有日志
	if len(it.logs) == 0 {
		it.done = true
		return false
	}
	return true
}

func (it *VotingVotedIterator) Error() error {
	return it.fail
}

func (it *VotingVotedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

func (it *VotingVotedIterator) Event() *VotingVoted {
	if len(it.logs) == 0 {
		return nil
	}
	
	// 获取并移除第一条日志
	log := it.logs[0]
	it.logs = it.logs[1:]
	
	// 解析事件
	event := new(VotingVoted)
	if err := abi.UnpackIntoInterface(event, "Voted", log.Data); err != nil {
		it.fail = err
		return nil
	}
	event.Voter = common.BytesToAddress(log.Topics[1].Bytes())
	event.Raw = log
	
	return event
}

// contract 返回合约对象
func (v *Voting) contract() bind.BoundContract {
	return bind.NewBoundContract(v.address, v.abi, v.backend, v.backend, v.backend)
} 
// Code generated - DO NOT EDIT.
// This file is a manually written Go binding of the SimpleToken contract.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// 代币合约ABI
var TokenABI = `[
	{
		"inputs": [
			{"internalType":"string","name":"_name","type":"string"},
			{"internalType":"string","name":"_symbol","type":"string"},
			{"internalType":"uint8","name":"_decimals","type":"uint8"},
			{"internalType":"uint256","name":"_initialSupply","type":"uint256"}
		],
		"stateMutability":"nonpayable",
		"type":"constructor"
	},
	{
		"anonymous":false,
		"inputs":[
			{"indexed":true,"internalType":"address","name":"owner","type":"address"},
			{"indexed":true,"internalType":"address","name":"spender","type":"address"},
			{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}
		],
		"name":"Approval",
		"type":"event"
	},
	{
		"anonymous":false,
		"inputs":[
			{"indexed":true,"internalType":"address","name":"from","type":"address"},
			{"indexed":true,"internalType":"address","name":"to","type":"address"},
			{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}
		],
		"name":"Transfer",
		"type":"event"
	},
	{
		"inputs":[
			{"internalType":"address","name":"","type":"address"},
			{"internalType":"address","name":"","type":"address"}
		],
		"name":"allowance",
		"outputs":[{"internalType":"uint256","name":"","type":"uint256"}],
		"stateMutability":"view",
		"type":"function"
	},
	{
		"inputs":[
			{"internalType":"address","name":"_spender","type":"address"},
			{"internalType":"uint256","name":"_value","type":"uint256"}
		],
		"name":"approve",
		"outputs":[{"internalType":"bool","name":"","type":"bool"}],
		"stateMutability":"nonpayable",
		"type":"function"
	},
	{
		"inputs":[{"internalType":"address","name":"","type":"address"}],
		"name":"balanceOf",
		"outputs":[{"internalType":"uint256","name":"","type":"uint256"}],
		"stateMutability":"view",
		"type":"function"
	},
	{
		"inputs":[],
		"name":"decimals",
		"outputs":[{"internalType":"uint8","name":"","type":"uint8"}],
		"stateMutability":"view",
		"type":"function"
	},
	{
		"inputs":[],
		"name":"name",
		"outputs":[{"internalType":"string","name":"","type":"string"}],
		"stateMutability":"view",
		"type":"function"
	},
	{
		"inputs":[],
		"name":"symbol",
		"outputs":[{"internalType":"string","name":"","type":"string"}],
		"stateMutability":"view",
		"type":"function"
	},
	{
		"inputs":[],
		"name":"totalSupply",
		"outputs":[{"internalType":"uint256","name":"","type":"uint256"}],
		"stateMutability":"view",
		"type":"function"
	},
	{
		"inputs":[
			{"internalType":"address","name":"_to","type":"address"},
			{"internalType":"uint256","name":"_value","type":"uint256"}
		],
		"name":"transfer",
		"outputs":[{"internalType":"bool","name":"","type":"bool"}],
		"stateMutability":"nonpayable",
		"type":"function"
	},
	{
		"inputs":[
			{"internalType":"address","name":"_from","type":"address"},
			{"internalType":"address","name":"_to","type":"address"},
			{"internalType":"uint256","name":"_value","type":"uint256"}
		],
		"name":"transferFrom",
		"outputs":[{"internalType":"bool","name":"","type":"bool"}],
		"stateMutability":"nonpayable",
		"type":"function"
	}
]`

// SimpleToken Go绑定结构体
type SimpleToken struct {
	address common.Address     // 合约地址
	backend bind.ContractBackend // 合约后端
	abi     abi.ABI           // 合约ABI
}

// NewSimpleToken 创建新的SimpleToken实例
func NewSimpleToken(address common.Address, backend bind.ContractBackend) (*SimpleToken, error) {
	// 解析合约ABI
	tokenAbi, err := abi.JSON(strings.NewReader(TokenABI))
	if err != nil {
		return nil, err
	}
	
	// 创建合约实例
	token := &SimpleToken{
		address: address,
		backend: backend,
		abi:     tokenAbi,
	}
	
	return token, nil
}

// 部署新的代币合约
func DeploySimpleToken(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string, decimals uint8, initialSupply *big.Int) (common.Address, *types.Transaction, *SimpleToken, error) {
	// 解析ABI
	parsed, err := abi.JSON(strings.NewReader(TokenABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	
	// 编码部署参数
	input, err := parsed.Pack("", name, symbol, decimals, initialSupply)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	
	// 部署合约
	address, tx, err := bind.DeployContract(auth, parsed, input, backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	
	// 创建合约实例
	contract := &SimpleToken{
		address: address,
		backend: backend,
		abi:     parsed,
	}
	
	return address, tx, contract, nil
}

// Name 获取代币名称
func (t *SimpleToken) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := t.call(opts, &out, "name")
	if err != nil {
		return "", err
	}
	return out[0].(string), nil
}

// Symbol 获取代币符号
func (t *SimpleToken) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := t.call(opts, &out, "symbol")
	if err != nil {
		return "", err
	}
	return out[0].(string), nil
}

// Decimals 获取代币精度
func (t *SimpleToken) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := t.call(opts, &out, "decimals")
	if err != nil {
		return 0, err
	}
	return out[0].(uint8), nil
}

// TotalSupply 获取代币总供应量
func (t *SimpleToken) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := t.call(opts, &out, "totalSupply")
	if err != nil {
		return nil, err
	}
	return out[0].(*big.Int), nil
}

// BalanceOf 获取指定地址的代币余额
func (t *SimpleToken) BalanceOf(opts *bind.CallOpts, who common.Address) (*big.Int, error) {
	var out []interface{}
	err := t.call(opts, &out, "balanceOf", who)
	if err != nil {
		return nil, err
	}
	return out[0].(*big.Int), nil
}

// Allowance 获取授权额度
func (t *SimpleToken) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := t.call(opts, &out, "allowance", owner, spender)
	if err != nil {
		return nil, err
	}
	return out[0].(*big.Int), nil
}

// Transfer 转账代币
func (t *SimpleToken) Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return t.transact(opts, "transfer", to, amount)
}

// Approve 授权代币
func (t *SimpleToken) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return t.transact(opts, "approve", spender, amount)
}

// TransferFrom 从授权地址转账
func (t *SimpleToken) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return t.transact(opts, "transferFrom", from, to, amount)
}

// 内部调用方法
func (t *SimpleToken) call(opts *bind.CallOpts, result *[]interface{}, method string, args ...interface{}) error {
	if opts == nil {
		opts = &bind.CallOpts{}
	}
	
	// 打包方法调用数据
	input, err := t.abi.Pack(method, args...)
	if err != nil {
		return err
	}
	
	// 创建调用消息
	msg := ethereum.CallMsg{
		From:     opts.From,
		To:       &t.address,
		Data:     input,
		GasPrice: opts.GasPrice,
	}
	
	// 执行调用
	output, err := t.backend.CallContract(opts.Context, msg, opts.BlockNumber)
	if err != nil {
		return err
	}
	
	// 解包结果
	return t.abi.UnpackIntoInterface(result, method, output)
}

// 内部交易方法
func (t *SimpleToken) transact(opts *bind.TransactOpts, method string, args ...interface{}) (*types.Transaction, error) {
	// 打包方法调用数据
	input, err := t.abi.Pack(method, args...)
	if err != nil {
		return nil, err
	}
	
	// 发送交易
	return bind.DoTransaction(opts, t.backend, &t.address, input)
}

// FilterTransfer 过滤Transfer事件
func (t *SimpleToken) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SimpleTokenTransferIterator, error) {
	// 创建过滤器
	var fromRule []interface{}
	for _, addr := range from {
		fromRule = append(fromRule, addr)
	}
	var toRule []interface{}
	for _, addr := range to {
		toRule = append(toRule, addr)
	}

	logs, sub, err := t.contract().FilterLogs(opts.Context, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &SimpleTokenTransferIterator{logs: logs, sub: sub}, nil
}

// WatchTransfer 监听Transfer事件
func (t *SimpleToken) WatchTransfer(opts *bind.WatchOpts, sink chan<- *SimpleTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {
	// 创建过滤器
	var fromRule []interface{}
	for _, addr := range from {
		fromRule = append(fromRule, addr)
	}
	var toRule []interface{}
	for _, addr := range to {
		toRule = append(toRule, addr)
	}

	logs, sub, err := t.contract().WatchLogs(opts.Context, "Transfer", sink, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// 解析事件
				event := new(SimpleTokenTransfer)
				if err := t.abi.UnpackIntoInterface(event, "Transfer", log.Data); err != nil {
					return err
				}
				event.From = common.BytesToAddress(log.Topics[1].Bytes())
				event.To = common.BytesToAddress(log.Topics[2].Bytes())
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

// SimpleTokenTransfer 转账事件
type SimpleTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log
}

// SimpleTokenTransferIterator 转账事件迭代器
type SimpleTokenTransferIterator struct {
	logs []types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *SimpleTokenTransferIterator) Next() bool {
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

func (it *SimpleTokenTransferIterator) Error() error {
	return it.fail
}

func (it *SimpleTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

func (it *SimpleTokenTransferIterator) Event() *SimpleTokenTransfer {
	if len(it.logs) == 0 {
		return nil
	}
	
	// 获取并移除第一条日志
	log := it.logs[0]
	it.logs = it.logs[1:]
	
	// 解析事件
	event := new(SimpleTokenTransfer)
	if err := abi.UnpackIntoInterface(event, "Transfer", log.Data); err != nil {
		it.fail = err
		return nil
	}
	event.From = common.BytesToAddress(log.Topics[1].Bytes())
	event.To = common.BytesToAddress(log.Topics[2].Bytes())
	event.Raw = log
	
	return event
}

// contract 返回合约对象
func (t *SimpleToken) contract() bind.BoundContract {
	return bind.NewBoundContract(t.address, t.abi, t.backend, t.backend, t.backend)
} 
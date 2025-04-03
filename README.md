# 以太坊Golang DApp示例

这是一个使用Golang和Solidity开发的以太坊DApp示例项目，展示了如何：

1. 部署和与ERC20代币合约交互
2. 实现基于代币的投票系统
3. 使用Go语言调用以太坊智能合约

## 项目结构

```
.
├── contracts/                 # Solidity智能合约
│   ├── Token.sol              # ERC20代币合约
│   └── Voting.sol             # 投票合约
│
├── go-ethereum-dapp/          # Go后端
│   ├── cmd/                   # 命令行程序
│   │   └── main.go            # 主程序
│   │
│   └── pkg/                   # 包
│       ├── config/            # 配置
│       ├── contracts/         # 合约绑定
│       ├── ethereum/          # 以太坊客户端
│       └── services/          # 服务层
│
├── .env.example               # 环境变量示例
└── README.md                  # 说明文档
```

## 技术栈

- **智能合约**: Solidity 0.8.0+
- **区块链**: 以太坊
- **后端**: Golang
- **以太坊交互**: go-ethereum (geth)

## 前置条件

- Golang 1.16+
- 以太坊开发网络 (如Ganache、Hardhat)
- 基本的以太坊和智能合约知识

## 快速开始

### 1. 克隆项目并安装依赖

```bash
git clone https://github.com/yangliyong7/go-ethereum.git
cd go-ethereum
go mod download
```

### 2. 配置环境变量

复制`.env.example`文件为`.env`并根据需要修改：

```bash
cp .env.example .env
```

确保修改以下内容：
- `ETHEREUM_URL`: 你的以太坊节点URL
- `CHAIN_ID`: 对应网络的链ID
- `PRIVATE_KEY`: 你的以太坊账户私钥

### 3. 启动以太坊开发网络

你可以使用Ganache或Hardhat运行本地开发网络：

```bash
# 使用Ganache（如已安装）
ganache-cli

# 或使用Hardhat（如已安装）
npx hardhat node
```

### 4. 编译和运行

```bash
cd go-ethereum-dapp
go build -o eth-app ./cmd
./eth-app
```

程序将执行以下步骤：
1. 部署ERC20代币合约
2. 转账一些代币到测试账户
3. 部署投票合约
4. 进行投票并展示结果

## 智能合约

### SimpleToken (ERC20)

`SimpleToken`是一个简单的ERC20代币合约，实现了以下功能：
- 代币的基本属性（名称、符号、精度）
- 代币的发行和转账
- 代币的授权和从授权地址转账

### Voting

`Voting`是一个基于代币的投票合约，实现了以下功能：
- 创建投票主题和选项
- 持有代币可以进行投票（一个地址只能投票一次）
- 投票权重与持有代币数量成正比
- 查询投票结果和获胜选项

## Golang后端

### 主要组件

- **ethereum.Client**: 封装以太坊网络交互
- **contracts**: 智能合约的Go绑定
- **services**: 业务逻辑层
  - **TokenService**: 管理代币合约操作
  - **VotingService**: 管理投票合约操作

## 重要概念

1. **ABI绑定**: 将Solidity合约转换为Go代码
2. **交易签名**: 使用私钥签名交易
3. **合约部署**: 在区块链上部署智能合约
4. **合约调用**: 调用智能合约的方法

## 进阶主题

- 使用事件过滤器监听合约事件
- 优化gas使用
- 错误处理和交易重试
- 使用钱包或密钥管理系统代替私钥

## 生产环境注意事项

- **不要在生产代码中硬编码私钥**
- 实现适当的错误处理和日志记录
- 考虑使用连接池管理ethclient连接
- 实现交易确认和重试机制
- 添加适当的测试用例

## 贡献

欢迎贡献代码或提出建议！请提交Pull Request或创建Issue。

## 许可证

MIT 
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./Token.sol";

/**
 * @title 投票合约
 * @dev 一个使用代币进行投票的合约示例
 */
contract Voting {
    // 投票的主题
    string public votingTopic;
    // 投票的截止时间
    uint256 public endTime;
    // 投票选项的数量
    uint256 public optionsCount;
    // 代币合约的地址
    SimpleToken public tokenContract;
    
    // 投票选项结构
    struct Option {
        string name;       // 选项名称
        uint256 voteCount; // 得票数
    }
    
    // 保存所有投票选项
    mapping(uint256 => Option) public options;
    // 记录每个地址是否已投票
    mapping(address => bool) public hasVoted;
    // 记录每个地址投票的选项
    mapping(address => uint256) public votedOption;
    
    // 投票事件
    event Voted(address indexed voter, uint256 optionId);
    
    /**
     * @dev 构造函数
     * @param _votingTopic 投票主题
     * @param _tokenAddress 代币合约地址
     * @param _durationInMinutes 投票持续时间（分钟）
     * @param _options 投票选项名称数组
     */
    constructor(
        string memory _votingTopic,
        address _tokenAddress,
        uint256 _durationInMinutes,
        string[] memory _options
    ) {
        require(_options.length >= 2, "至少需要两个投票选项");
        
        votingTopic = _votingTopic;
        tokenContract = SimpleToken(_tokenAddress);
        endTime = block.timestamp + _durationInMinutes * 1 minutes;
        
        // 初始化投票选项
        for (uint256 i = 0; i < _options.length; i++) {
            options[i] = Option({
                name: _options[i],
                voteCount: 0
            });
        }
        
        optionsCount = _options.length;
    }
    
    /**
     * @dev 投票函数
     * @param _optionId 选择的投票选项ID
     */
    function vote(uint256 _optionId) public {
        require(block.timestamp < endTime, "投票已结束");
        require(_optionId < optionsCount, "无效的选项ID");
        require(!hasVoted[msg.sender], "已经投过票了");
        require(tokenContract.balanceOf(msg.sender) > 0, "需要持有代币才能投票");
        
        // 记录投票信息
        hasVoted[msg.sender] = true;
        votedOption[msg.sender] = _optionId;
        
        // 增加选项票数（票数权重与持有代币数量成正比）
        uint256 voteWeight = tokenContract.balanceOf(msg.sender);
        options[_optionId].voteCount += voteWeight;
        
        // 触发投票事件
        emit Voted(msg.sender, _optionId);
    }
    
    /**
     * @dev 获取获胜选项ID
     * @return 获胜选项的ID
     */
    function winningOptionId() public view returns (uint256) {
        uint256 winningVoteCount = 0;
        uint256 winningId = 0;
        
        for (uint256 i = 0; i < optionsCount; i++) {
            if (options[i].voteCount > winningVoteCount) {
                winningVoteCount = options[i].voteCount;
                winningId = i;
            }
        }
        
        return winningId;
    }
    
    /**
     * @dev 获取获胜选项名称
     * @return 获胜选项的名称
     */
    function winningOptionName() public view returns (string memory) {
        return options[winningOptionId()].name;
    }
    
    /**
     * @dev 查询投票是否已结束
     * @return 投票是否已结束
     */
    function isVotingEnded() public view returns (bool) {
        return block.timestamp >= endTime;
    }
    
    /**
     * @dev 获取剩余投票时间（秒）
     * @return 剩余时间（秒）
     */
    function remainingTime() public view returns (uint256) {
        if (block.timestamp >= endTime) {
            return 0;
        }
        return endTime - block.timestamp;
    }
} 
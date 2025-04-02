// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title 简单的ERC20代币合约
 * @dev 实现ERC20标准的代币合约
 */
contract SimpleToken {
    // 代币名称
    string public name;
    // 代币符号
    string public symbol;
    // 代币精度（小数位数）
    uint8 public decimals;
    // 代币总供应量
    uint256 public totalSupply;

    // 记录每个地址的代币余额
    mapping(address => uint256) public balanceOf;
    // 记录每个地址授权给其他地址的代币数量
    mapping(address => mapping(address => uint256)) public allowance;
    
    // 转账事件
    event Transfer(address indexed from, address indexed to, uint256 value);
    // 授权事件
    event Approval(address indexed owner, address indexed spender, uint256 value);
    
    /**
     * @dev 构造函数，初始化代币信息
     * @param _name 代币名称
     * @param _symbol 代币符号
     * @param _decimals 代币精度
     * @param _initialSupply 初始代币供应量
     */
    constructor(
        string memory _name,
        string memory _symbol,
        uint8 _decimals,
        uint256 _initialSupply
    ) {
        name = _name;
        symbol = _symbol;
        decimals = _decimals;
        
        // 计算实际供应量（考虑精度）
        totalSupply = _initialSupply * (10 ** uint256(decimals));
        // 初始时将所有代币分配给合约部署者
        balanceOf[msg.sender] = totalSupply;
        
        // 触发转账事件
        emit Transfer(address(0), msg.sender, totalSupply);
    }
    
    /**
     * @dev 转账代币
     * @param _to 接收者地址
     * @param _value 转账金额
     * @return 是否转账成功
     */
    function transfer(address _to, uint256 _value) public returns (bool) {
        require(_to != address(0), "不能转账到零地址");
        require(balanceOf[msg.sender] >= _value, "余额不足");
        
        // 更新余额
        balanceOf[msg.sender] -= _value;
        balanceOf[_to] += _value;
        
        // 触发转账事件
        emit Transfer(msg.sender, _to, _value);
        return true;
    }
    
    /**
     * @dev 从授权地址转账代币
     * @param _from 转出地址
     * @param _to 接收地址
     * @param _value 转账金额
     * @return 是否转账成功
     */
    function transferFrom(address _from, address _to, uint256 _value) public returns (bool) {
        require(_from != address(0), "不能从零地址转账");
        require(_to != address(0), "不能转账到零地址");
        require(balanceOf[_from] >= _value, "余额不足");
        require(allowance[_from][msg.sender] >= _value, "授权额度不足");
        
        // 更新余额和授权额度
        balanceOf[_from] -= _value;
        balanceOf[_to] += _value;
        allowance[_from][msg.sender] -= _value;
        
        // 触发转账事件
        emit Transfer(_from, _to, _value);
        return true;
    }
    
    /**
     * @dev 授权其他地址使用代币
     * @param _spender 被授权地址
     * @param _value 授权金额
     * @return 是否授权成功
     */
    function approve(address _spender, uint256 _value) public returns (bool) {
        require(_spender != address(0), "不能授权给零地址");
        
        // 设置授权额度
        allowance[msg.sender][_spender] = _value;
        
        // 触发授权事件
        emit Approval(msg.sender, _spender, _value);
        return true;
    }
} 
# Aztecs 区块链实现

## 区块链原理与技术实现
### 核心概念
- **工作量证明(PoW)**：通过 `consensus/pow.go` 实现的挖矿算法，调整难度值控制区块生成速度
- **UTXO模型**：在 `core/utxo.go` 中实现未花费交易输出模型，确保交易可验证且防双花
- **椭圆曲线加密(ECDSA)**：`crypto/wallet.go` 使用 secp256k1 曲线生成密钥对，保障资产安全
- **默克尔树**：`core/block.go` 中实现交易哈希树，快速验证区块完整性

### 数据结构
```go
// 区块结构 (core/block.go)
type Block struct {
	Timestamp    int64
	Transactions []*Transaction
	PrevHash     []byte
	Hash         []byte
	Nonce        int
	Height       int
}

// 交易结构 (core/transaction.go)
type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}
```

## 安装与运行
### 环境准备
1. 安装 Go 1.18+ [下载地址](https://go.dev/dl/)
2. 安装 Node.js 16+ [下载地址](https://nodejs.org/)
3. 安装 LevelDB 依赖: `brew install leveldb` (macOS)

### 后端启动
```bash
# 1. 克隆仓库
git clone https://github.com/your-repo/aztecs.git
cd aztecs

# 2. 创建新钱包（获取矿工地址）
go run main.go createwallet
# 输出示例：新钱包地址: 1ABC... (复制这个地址)

# 3. 启动节点
go run main.go startnode \
  --port=8080 \
  "--miner-address=1ABC..." # 粘贴上一步复制的地址
```

### 前端启动
```bash
cd frontend
npm install  # 安装依赖
REACT_APP_API_URL=http://localhost:8080 npm start
```

### 完整执行示例
```bash
# 终端1：启动后端
$ cd aztecs
$ go run main.go createwallet
新钱包地址: 1F1tAaz5x1HUXrCNLbtMDqcw6o5GNn4xqX

$ go run main.go startnode \
    --port=8080 \
    "--miner-address=1F1tAaz5x1HUXrCNLbtMDqcw6o5GNn4xqX"

# 终端2：启动前端
$ cd aztecs/frontend
$ REACT_APP_API_URL=http://localhost:8080 npm start
```

## 操作指南
### 1. 创建钱包
通过前端界面或API创建：
```bash
# API方式
curl -X POST http://localhost:8080/wallet/new
# 返回：{"address":"1ABC...","private_key":"..."}
```

### 2. 发送交易
1. 在前端 "Transaction Sender" 界面
2. 输入接收方地址和金额
3. 使用私钥签名交易
4. 广播到网络

### 3. 查看区块链
- 区块浏览器：`http://localhost:3000/blocks`
- API端点：`GET /blocks`

### 4. 挖矿操作
节点自动执行PoW挖矿，可通过API手动触发：
```bash
curl -X POST http://localhost:8080/mine
```

## API参考
| 端点 | 方法 | 功能 |
|------|------|------|
| `/wallet/new` | POST | 创建新钱包 |
| `/transaction` | POST | 创建交易 |
| `/blocks` | GET | 获取区块链 |
| `/utxo` | GET | 查询UTXO |
| `/balance/:address` | GET | 查询余额 |
| `/mine` | POST | 手动挖矿 |

## 项目结构
```
aztecs/
├── api/          # REST API
├── consensus/    # PoW共识算法
├── core/         # 区块链核心
├── crypto/       # 加密模块
├── frontend/     # React前端
└── storage/      # LevelDB存储
```

## 贡献
欢迎提交PR，请确保：
1. 包含单元测试（Go测试或Jest）
2. 遵循Go代码规范 `gofmt`
3. 更新相关文档

许可证: [MIT](LICENSE)

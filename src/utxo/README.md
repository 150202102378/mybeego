# utxo
可以理解为一个比特币区块链的监听器，通过配置需要监听的地址，记录该地址集相关的比特币资产转出转入详情。

# Install
git clone [待定]

# Introduce
## tables
- blocks : 记录区块高度信息
- txs : 记录监听地址相关的交易订单
- vout : 记录监听地址转入资产的vouts信息
- vins : 记录系统地址转出资产的信息
- addr : 记录监听地址的资产汇总情况
- xcurrency_address : 配置监听地址
## logic
start:

check need Revert:
if 当前记录区块 > 当前区块链的最新区块 then Revert

execute Revert:
revertNum = 当前记录区块 - 当前区块链的最新区块
foreach need revertBlocks: delete vins, vouts, txs, blocks and change addrs.

execute Tally Blocks:
更新vins, vouts 的confirmatons + 1
get txs by 当前记录区块
过滤出监听地址相关的交易订单
记录vins, vouts, txs, blocks

# Usage
insert Addresses to db_table : xcurrency_address 
config start block_height
go run main.go release

# transfer can not be here

package virtual_machine

import "math/big"

/*
CoinList
@dbpath db.coin
@desc Store all coins' name deployed by OpDeploy.
*/
type CoinList struct {
	names []string
}

/*
CoinMeta
@db_path db.coin_[coin-name]
@desc Store a coin's metadata deployed.
@writing_op OpDeploy, OpMint
*/
type CoinMeta struct {
	supply          big.Float
	minted          big.Float
	limitPerMint    big.Float
	limitPerAccount big.Float
	deployedTx      string
}

/*
CoinBalance
@db_path db.coin_[coin-name]_[address]
@desc Store a coin's balance at a single account.
@writing_op OpMint, OpTransfer, OpAddLiquidity, OpRemoveLiquidity OpSwap
*/
type CoinBalance struct {
	value big.Float
}

/*
Transactions
@db_path db.tx_[coin-name]_[instruction-id]
@desc Store a coin's transactions
@writing_op OpDeploy, OpMint, OpTransfer, OpAddLiquidity, OpRemoveLiquidity OpSwap
*/
type Transactions struct {
	txId   string
	op     string
	params map[string]interface{}
}

/*
LiquidityPairList
@db_path db.lp_[coin-name]
@desc Store all liquidity providers associated with a coin
@writing_op OpAddLiquidity
*/
type LiquidityPairList struct {
	pairs []string
}

/*
LiquidityPair
@db_path db.lp_[left-coin-name|right-coin-name]
@desc Store liquidity provider metadata, totalValue, price
@writing_op OpAddLiquidity, OpRemoveLiquidity, OpSwap
*/
type LiquidityPair struct {
	total big.Float
	price big.Float
}

/*
LiquidityPairBalance
@db_path db.lp_[left-coin-name|right-coin-name]_[address]
@desc Store liquidity provider balance of a single value
@writing_op OpAddLiquidity, OpRemoveLiquidity, OpSwap
*/
type LiquidityPairBalance struct {
	value big.Float
}

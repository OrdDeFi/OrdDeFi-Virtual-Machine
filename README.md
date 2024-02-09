# Introduction

In the financial sector, there are three crucial applications: asset issuance, trading, and lending. BRC-20 has made significant progress in asset issuance. OrdDeFi expands the BRC-20 instruction set and provides a native Automated Market Maker (AMM) protocol on Bitcoin L1, achieving the second application: trading.

### AMM Operations in OrdDeFi

In addition to the deploy, mint, and transfer instructions in BRC-20, OrdDeFi introduces additional instructions: `addlp` (add liquidity provider), `swap`, and `rmlp` (remove liquidity provider). These instructions enable users to utilize AMM algorithms for automated market-making and trading of OrdDeFi protocol assets on L1.

### Fair Mint Made Even Fairer

In the current BRC-20 fair mint process, batch minting and repeated minting allow for the engraving of a large number of tokens by splitting the outputs or splitting the transaction. This leads to a more concentrated distribution of tokens. As a result, transaction fees on the Bitcoin network are driven up, and participating users incur higher costs without receiving a proportional increase in tokens.

OrdDeFi has addressed this issue with optimizations. On one hand, by imposing restrictions on `TxIn[0]` and the output index from the previous transaction, only the first mint instruction becomes effective when generating multiple mint instructions through repeat mint. On the other hand, OrdDeFi introduces the address limit setting in the `deploy` operation, allowing the deployer to restrict the maximum number of mints per individual address. These measures promote a more decentralized token distribution within the OrdDeFi protocol and enhance the fairness of the fair mint process.

### Extended Deploy and Transfer
OrdDeFi has made the following expansions to the `deploy` and `transfer` instructions:

* The deploy instruction in OrdDeFi now includes the `desc` and `icon` parameters, allowing the deployer to provide a short description and an icon image. Additionally, the deploy instruction introduces the `alim` parameter (address limit) to impose restrictions on the maximum quantity of tokens that can be minted by a single address.
* The transfer instruction in OrdDeFi adds the "to" parameter, enabling users to transfer assets in a single transaction, eliminating the need for separate "inscribe transfer" and "send transfer UTXO" steps. 

See [docs/1.Intro.md](https://github.com/OrdDefi/OrdDefi-Virtual-Machine/blob/main/docs/1.Intro.md) for details.

# Ecosystem Design

In OrdDeFi, there are two enhanced assets: ODFI and ODGV.  

ODFI is the primary coin of the OrdDeFi protocol, offering the following features:

* Fee collection by adding liquidity provider
* Transaction fee discounts at swap

ODGV is the governance coin of the OrdDeFi protocol.

See [docs/2.Ecosystem.md](https://github.com/OrdDefi/OrdDefi-Virtual-Machine/blob/main/docs/2.Ecosystem.md) for details.

# Build, Index and Query assets

To start indexer:  

```
go build
./OrdDeFi-Virtual-Machine
```

See [docs/3.BuildAndUseIndexer.md](https://github.com/OrdDefi/OrdDefi-Virtual-Machine/blob/main/docs/3.BuildAndUseIndexer.md) for details.

# Operations

## Operation Compiler

See [docs/4.1.OperationCompiler.md](https://github.com/OrdDefi/OrdDefi-Virtual-Machine/blob/main/docs/4.1.OperationCompiler.md) for details.

## Operations

A typical `deploy` instruction:

```
{
  "p":"orddefi",
  "op":"deploy",
  "tick":"half",
  "max":"210000000",
  "lim":"1000",
  "alim":"1000",
  "desc":"meme coin of Hal Finney",
  "icon":"Base64 string of icon"
}
```

A typical `mint` instruction:

```
{
  "p":"orddefi",
  "op":"mint",
  "tick":"odfi",
  "amt":"1000"
}
```

A typical UTXO transfer instruction:

```
{
  "p":"orddefi",
  "op":"transfer",
  "tick":"odfi",
  "amt":"1000"
}
```

A typical `addlp` instruction:

```
{
  "p":"orddefi",
  "op":"addlp",
  "ltick":"odfi",
  "rtick":"odgv",
  "lamt":"1000",
  "ramt":"1000"
}
```

A typical `swap` instruction:

```
{
  "p":"orddefi",
  "op":"swap",
  "ltick":"odfi",
  "rtick":"odgv",
  "spend":"odfi",
  "amt":"1000"
}
```

A `swap` instruction with `threshold` parameter:

```
{
  "p":"orddefi",
  "op":"swap",
  "ltick":"odfi",
  "rtick":"odgv",
  "spend":"odfi",
  "amt":"1000",
  "threshold":"0.005"
}
```

A typical `rmlp` instruction:

```
{
  "p":"orddefi",
  "op":"rmlp",
  "ltick":"odfi",
  "rtick":"odgv",
  "amt":"1000"
}
```

See [docs/4.2.Operations.md](https://github.com/OrdDefi/OrdDefi-Virtual-Machine/blob/main/docs/4.2.Operations.md) for details.

# Inscriber

For enhancing security during the inscription of operations such as mint, addlp, rmlp, swap, and direct-transfer (specified as transfer with a to parameter), it is crucial to adhere to specific transaction (tx) guidelines. These guidelines dictate that the first transaction input (TxIn[0]) in the commit transaction must match the first transaction output (TxOut[0]) in the reveal transaction. Additionally, the final transaction output (TxOut[-1]) in the commit transaction must be an OpReturn containing the data orddefi:auth. Any instruction failing to meet these criteria will be terminated for safety reasons.

To successfully execute the operations mint, addlp, rmlp, swap, and transfer with a specified to parameter, utilize the [OrdDeFi-Inscribe](https://github.com/OrdDeFi/OrdDeFi-Inscribe) tool. For detailed guidance and to ensure compliance with the required transaction structures, thoroughly review the [readme file](https://github.com/OrdDeFi/OrdDeFi-Inscribe/blob/main/README.md) of the [OrdDeFi-Inscribe](https://github.com/OrdDeFi/OrdDeFi-Inscribe) repository.

# LICENSE

[GNU GENERAL PUBLIC LICENSE](https://github.com/OrdDefi/OrdDefi-Virtual-Machine/blob/main/LICENSE)

# Get the Latest News

Follow the [Official OrdDeFi Twitter](https://twitter.com/OrdDeFi)
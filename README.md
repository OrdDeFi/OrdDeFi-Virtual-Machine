# Introduction

In the financial sector, there are three crucial applications: asset issuance, trading, and lending. BRC-20 has made significant progress in asset issuance. OrdDeFi expands the BRC-20 instruction set and provides a native Automated Market Maker (AMM) protocol on Bitcoin L1, achieving the second application: trading.

### AMM Operations in OrdDeFi

In addition to the deploy, mint, and transfer instructions in BRC-20, OrdDeFi introduces additional instructions: `addlp` (add liquidity provider), `swap`, and `rmlp` (remove liquidity provider). These instructions enable users to utilize AMM algorithms for automated market-making and trading of OrdDeFi protocol assets on L1.

When using AMM-related instructions, the sender and recipient of the instruction must be the same address to prevent other users from controlling the assets.

### Fair Mint Made Even Fairer

Currently, in the fair mint process of BRC-20, some third-party services offer a "repeat mint" feature. This feature allows a single UTXO to be split into 1,000 small-value UTXOs, which are then individually inscribed, enabling the creation of 1,000 inscriptions in a single operation. With the use of the repeat mint feature, which originally required 21,000,000 individual operations, can now be completed by just 21,000 users, each performing a single repeat mint. While this improves the efficiency of inscription, it leads to a more concentrated token distribution. As a result, transaction fees on the Bitcoin network are driven up, and participating users incur higher costs without receiving a proportional increase in tokens.

OrdDeFi has addressed this issue with optimizations. On one hand, by imposing restrictions on `TxIn[0]` and the output index from the previous transaction, only the first mint instruction becomes effective when generating multiple mint instructions through repeat mint. On the other hand, OrdDeFi introduces the address limit setting in the `deploy` operation, allowing the deployer to restrict the maximum number of mints per individual address. These measures promote a more decentralized token distribution within the OrdDeFi protocol and enhance the fairness of the fair mint process.

### Suppress the Inflation of the UTXO Dataset

Bitcoin maximalists have criticized ordinals and BRC-20 for their inscription mechanism, as they believe that continuously creating new inscriptions in ordinals and BRC-20 leads to rapid expansion of the UTXO dataset, resulting in wastage of Bitcoin node resources.

In OrdDeFi, it is legitimate to inscribe witness scripts for the same sat in different transactions. Unlike in the original version of the ordinal client, which does not allow different inscriptions for the same sat. Although the format of the witness script remains consistent with ordinals, OrdDeFi eliminates the concept of inscriptions and replaces it with instructions.

OrdDeFi cannot enforce users to refrain from generating new inscriptions. However, if instructions that do not generate new UTXOs are widely adopted, the issue of UTXO dataset inflation can be resolved.

To solve this concern, OrdDeFi has independently implemented a witness script indexer that is compatible with the ordinals protocol but does not rely on it.

### Extended Deploy and Transfer
OrdDeFi has made the following expansions to the `deploy` and `transfer` instructions:

* The deploy instruction in OrdDeFi now includes the `desc` and `icon` parameters, allowing the deployer to provide a short description and an icon image. Additionally, the deploy instruction introduces the `alim` parameter (address limit) to impose restrictions on the maximum quantity of tokens that can be minted by a single address.
* The transfer instruction in OrdDeFi adds the "to" parameter, enabling users to transfer assets in a single transaction, eliminating the need for separate "inscribe transfer" and "transfer" steps. When using the "to" parameter, the sender and recipient of the instruction must be the same address to prevent asset theft by other users.

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

## Deploy, Mint and Transfer

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

Although the `deploy`, `mint`, and `transfer` instructions may appear similar to BRC-20, they have been expanded compared to BRC-20. See [docs/4.2.BasicOperations.md](https://github.com/OrdDefi/OrdDefi-Virtual-Machine/blob/main/docs/4.2.BasicOperations.md) for details.

## AMM operations

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

See [docs/4.3.AMMOperations.md](https://github.com/OrdDefi/OrdDefi-Virtual-Machine/blob/main/docs/4.3.AMMOperations.md) for details.

# LICENSE

[GNU GENERAL PUBLIC LICENSE](https://github.com/OrdDefi/OrdDefi-Virtual-Machine/blob/main/LICENSE)
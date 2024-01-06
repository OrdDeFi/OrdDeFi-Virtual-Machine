# Introduction

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

See [docs/4.3.AMMOperations.md](https://github.com/OrdDefi/OrdDefi-Virtual-Machine/blob/main/docs/4.3.AMMOperations.md) for details.

# LICENSE

[GNU GENERAL PUBLIC LICENSE](https://github.com/OrdDefi/OrdDefi-Virtual-Machine/blob/main/LICENSE)
# Bitcoin wallet API

This project is implemented in Go, with APIs supporting the following operations:

1. Generate a random mnemonic words following BIP39 standard
2. Generate a Hierarchical Deterministic (HD) Segregated Witness (SegWit) bitcoin address from a given seed and path

## Used Libraries

1. [go-bip39](https://github.com/tyler-smith/go-bip39)
2. [go-bip32](https://github.com/tyler-smith/go-bip32)
3. [btcd](https://github.com/btcsuite/btcd)
4. [btcutil](https://github.com/btcsuite/btcutil)

### run on debug console

```sh
GO111MODULE=on go run main.go
```

### run on docker

```sh
make build # build docker image
docker run -d -p 20200:20200 --restart=always --name=bitcoin-wallet github.com/clcng/bitcoin-wallet:0.0.1
```

# Bitcoin wallet API

This project is implemented in Go, with APIs supporting the following operations:

1. Generate a random mnemonic words following BIP39 standard

## Used Libraries

1. [go-bip39](https://github.com/tyler-smith/go-bip39)

### run on debug console

```sh
GO111MODULE=on go run main.go
```

### run on docker

```sh
make build # build docker image
docker run -d -p 20200:20200 --restart=always --name=bitcoin-wallet github.com/clcng/bitcoin-wallet:0.0.1
```

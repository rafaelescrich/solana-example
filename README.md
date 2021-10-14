# solana-example

Rest API to receive solana tokens in testnet just like a faucet

## Running

```bash
go run main.go
```

## Test

Request airdrop

```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"address":"my_address","quantity":"value"}' \
  http://localhost:1337/airdrop
```

Get Balance

```bash
curl -X GET "http://localhost:1337/address/my_address"
```
# EYE - Decentralized Oracle Service

___Description___:

A Proof of Concept of a decentralized Oracle using the DERO service for the dARCH event 0.5. Everyone will use the dummy wallet like a pool of inputs and use the private Operator wallet to send result to the wallet that ask for requests.


___2021 dARCH event 0.5___:

https://forum.dero.io/t/darch-2021-event-0-5-services-only/1330

#
___How to run___:

### 0. Launch the Dero Daemon:
```
./derod-linux-amd64 --testnet
```
### 1. Recover / Launch the dummy wallet: (wallet shared by everyone)
``` 
./dero-wallet-cli-linux-amd64 --testnet --wallet-file=dum.db --rpc-server --rpc-bind 127.0.0.1:40404
```
!!! DUMMY WALLET (25 seeds)!!!: 
```
trash gables adjust tufts affair alerts inorganic banjo directed fossil somewhere sipped gleeful aisle wagtail bias pockets tudor voice sixteen duke incur motherly algebra pockets
```
### 2. Create / Launch the operator wallet (Your private wallet)
```
./dero-wallet-cli-linux-amd64 --testnet --wallet-file=ope.db --rpc-server --rpc-bind 127.0.0.1:40405 
```
### 3. Create / Launch the request wallet: (test wallet that send input to the dummy)
```
./dero-wallet-cli-linux-amd64 --testnet --wallet-file=req.db --rpc-server --rpc-bind 127.0.0.1:40403 
```

### 4. Launch the oracle:
```
./eye-linux-amd64
```
or
``` 
go run main.go
```
### 5a. Use the command below: (bitcoin to usd price)
``` 
curl http://127.0.0.1:40403/json_rpc -d ' { "jsonrpc": "2.0", "id": "1", "method": "transfer", "params": { "transfers": [ { "amount": 1, "destination": "detoi1qxszv4ell3de4ur8lsrfys9k4dhgzw9kgu3cp25z4yu5382h6wtzc29pvfz92x774klw7y352euqwdxwuw", "payload_rpc": [ { "name":"D", "datatype":"U", "value":16045690981402826360 }, { "name":"NAME", "datatype":"S", "value":"bitcoin" }, { "name":"CURR", "datatype":"S", "value":"usd" } ] } ] } }' -H 'Content-Type: application/json' 
```

### 5b. Or launch the script:
``` 
bash request_dero_price.sh 
```
#
___Donation___:
### Dero:
```
dERoNFYEXufYzcs64qnfocZ48nvmbxbKR1x2MZBqbrHn5dULSfRRdN3d4EsbwKLGeHE5k3Vrh77BWFufe2gBcrDF57PqDCaJoc
```

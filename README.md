Commit-Blockchain
===

## Installation


**Note:** Requires [Go 1.12+](https://golang.org/dl/)

Golang : installation in  ubuntu

### Download the go 

Add the Golang PPA repository to get the latest version of Golang.

`sudo add-apt-repository ppa:longsleep/golang-backports`

After adding the PPA, update packages list using the below command.   

`sudo apt-get update`

Install the latest version of Golang and other required packages

`sudo apt-get install -y git golang-go make`

Setup Environment Variables

`export GOROOT=/usr/lib/go`

`export GOPATH=$HOME/go`

`export GOBIN=$GOPATH/bin`

`export PATH=$PATH:$GOROOT/bin:$GOBIN`

You can also append the above lines to $HOME/.bashrc file and run the following command to reflect in current Terminal session

`source $HOME/.bashrc`

`go env`

```
git clone https://github.com/commitHub/commitBlockchain
cd commitBlockchain
```
```
make all
```

## Validator Setup 
### Setup a New Node

**Note** : You need to install **commit-blockchain** before you go further


> Warning : Monikers can contain only ASCII characters. Using Unicode characters will render your node unreachable.

```
maind init --name [moniker name] 
```
You can edit this ```moniker``` later, in the ```~/.maind/config/config.toml``` file:
```
# A custom human readable name for this node
moniker = "<your_custom_moniker>"
```



### Adding seeds
Your node needs to know how to find peers. You'll need to add healthy seed nodes to ```$HOME/.maind/config/config.toml```. The launch repo contains links to some seed nodes. 

Seeds of local network 
```
da71eeaf2911ae76ed3811c29f10eb3924df620c@192.168.2.98:36656,
ab940fb60d51ad8aa2b06c03f9e21a0fe5e47f40@192.168.2.98:26656
```

### Gas and Fees

> WARNING
On commit-blockchain mainnet, the accepted denom is `commit`

```
fees = ceil(gas * gasPrices)
```
The```gas``` is dependents on transaction. but we can estimate the amount of gas required for the transaction using flag ```--gas-adjustment```(default 1.0)

The ```gasPrice``` is the price of each unit of gas. 

The higher the ```gasPrice/fees```, the higher the chance that your transaction will get included in a block.



### Upgrade to Validator Node
You now have an active full node. What's the next step? You can upgrade your full node to become a commit-blockchain Validator.

## Run Validator on Commit-blockchain


To know validator pubkey 
```
maind tendermint show-validator 
```

To add Validator

```
maincli stake create-validator \
    --amount=1commit  \
    --from="<key_name>" \ 
    --pubkey=$(maind tendermint show-validator) \
    --moniker="choose a moniker" \
    --chain-id="<chain_id>"\
    --identity="identity of validator"\ 
    --website="commit.kgf"\ 
    --details="validator details"\
    --gas-adjuestment=1.0

```
Minimun amount is 1commit


>To avoid the ddos attack on the validator, it recommanded to follow youre own implementation of securing validator or use [Sentry node architecture](https://forum.cosmos.network/t/sentry-node-architecture-overview/454)



### Delegator : Delegate tokens to Validator

- download the latest build of maincli binary file

To get the list of validators
```
maincli stake validators

```
Delegate Tokens

```
./maincli stake delegate \
    --validator=<validator_address> \
    --amount=1commit \
    --chain-id=test-chain-fkp8YY \
    --gas-adjustment=1.0 \
    --from=<key_name> 

```

To query delegations 
```
./maincli stake delegation [delegator address]
```



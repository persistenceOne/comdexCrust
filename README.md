Commit-Blockchain
===

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

### Copy the genesis file
[Link of genesis file]() //TODO


### Adding seeds
Your node needs to know how to find peers. You'll need to add healthy seed nodes to ```$HOME/.maind/config/config.toml```. The launch repo contains links to some seed nodes. 

Seeds of local network 
```
da71eeaf2911ae76ed3811c29f10eb3924df620c@192.168.2.98:36656,
ab940fb60d51ad8aa2b06c03f9e21a0fe5e47f40@192.168.2.98:26656
```

### Gas and Fees

> WARNING
On commit-blockchain mainnet, the accepted denom is steak

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
    --amount=1steak  \
    --from="<key_name>" \ 
    --pubkey=$(maind tendermint show-validator) \
    --moniker="choose a moniker" \
    --chain-id="<chain_id>"\
    --identity="identity of validator"\ 
    --website="commit.kgf"\ 
    --details="validator details"\
    --gas-adjuestment=1.0

```
> Minimun amount is 1steak


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
    --amount=1steak \
    --chain-id=test-chain-fkp8YY \
    --gas-adjustment=1.0 \
    --from=<key_name> 

```

To query delegations 
```
./maincli stake delegation [delegator address]
```



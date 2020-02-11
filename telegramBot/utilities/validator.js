const HttpUtils = require('./httpRequest');
const errors = require('./errors');
const httpUtils = new HttpUtils();
const config = require('../config.json');
const dataUtils = require('./data');
const jsonUtils = require('./json');
const botUtils = require('./bot');
const WebSocket = require('ws');
const wsConstants = require('../constants/websocket');

const bech32 = require('bech32');
const hash = require('tendermint/lib/hash.js');
const tmhash = hash.tmhash;

const addressOperations = {
    pubkeyToBech32(pubkey, prefix) {
        let pubkeyAminoPrefix = Buffer.from('1624DE6420', 'hex');
        let buffer = Buffer.alloc(37);
        pubkeyAminoPrefix.copy(buffer, 0);
        Buffer.from(pubkey, 'base64').copy(buffer, pubkeyAminoPrefix.length);
        return bech32.encode(prefix, bech32.toWords(buffer));
    },
    bech32ToPubkey(pubkey) {
        let pubkeyAminoPrefix = Buffer.from('1624DE6420', 'hex')
        let buffer = Buffer.from(bech32.fromWords(bech32.decode(pubkey).words));
        return buffer.slice(pubkeyAminoPrefix.length).toString('base64');
    },
    getHexAddress(pubkeyValue) {
        let bytes = Buffer.from(pubkeyValue, 'base64');
        return tmhash(bytes).slice(0, 20).toString('hex').toUpperCase();
    },
    toPubKey(address) {
        return bech32.decode(config.prefix, address);
    },
    createAddress(publicKey) {
        const message = CryptoJS.enc.Hex.parse(publicKey.toString(`hex`));
        const hash = CryptoJS.RIPEMD160(CryptoJS.SHA256(message)).toString();
        const addr = Buffer.from(hash, `hex`);
        return bech32ify(addr, config.prefix);
    },
    getDelegatorAddrFromOperatorAddr(operatorAddr) {
        let address = bech32.decode(operatorAddr);
        return bech32.encode(config.prefix, address.words);
    },
    verifyValidatorOperatorAddress(validatorOperatorAddr) {
        const validatorOperatorAddrRegex = new RegExp('\^' + config.prefix + 'valoper' + '\[a-z0-9]{39}$');
        return validatorOperatorAddrRegex.test(validatorOperatorAddr);
    },
};

function bech32ify(address, prefix) {
    const words = bech32.toWords(address);
    return bech32.encode(prefix, words);
}

let wsTx;

const reinitWSTx = () => {
    if (wsTx === undefined) {
        wsTx = new WebSocket(wsConstants.url);
    } else {
        if (wsTx.url === wsConstants.url) {
            wsTx = new WebSocket(wsConstants.backupURL);
        }
        if (wsTx.url === wsConstants.backupURL) {
            wsTx = new WebSocket(wsConstants.url);
        }
    }
    try {
        wsTx.on('open', wsTxOpen);
        wsTx.on('close', wsTxClose);
        wsTx.on('message', wsTxIncoming);
        wsTx.on('error', wsTxError);
    } catch (e) {
        errors.Log(e, 'WS_TX_CONNECTION');
        wsTx.send(JSON.stringify(wsConstants.unsubscribeAllMsg));
        reinitWSTx();
    }
};

reinitWSTx();

function wsTxOpen() {
    wsTx.send(JSON.stringify(wsConstants.subscribeTxMsg));
}

function wsTxClose(code, reason) {
    let err = {statusCode: code, message: 'WS TX connection closed:    ' + reason};
    errors.Log(err, 'WS_TX_CONNECTION');
    reinitWSTx();
}

function wsTxError(err) {
    errors.Log(err, 'WS_TX_CONNECTION');
    wsTx.send(JSON.stringify(wsConstants.unsubscribeAllMsg));
    wsTx.close();
}

//If this doesn't work when there are more than one transactions in one block,
// use httpUtils.httpGet(botUtils.nodeURL, config.node.abciPort, `/tx_search?query="tx.height=${height}"&per_page=30`)
// to query and update.
function wsTxIncoming(data) {
    let json = JSON.parse(data);
    if (json === undefined) {
        errors.Log('Cannot parse data from ws connection.');
    }
    if (errors.isEmpty(json.result)) {
        console.log('ws Tx Connected!');
    } else {
        let txs = JSON.parse(json.result.data.value.TxResult.result.log);
        txs.forEach((tx) => {
            if (tx.success) {
                tx.events.forEach((event) => {
                    event.attributes.forEach((attribute) => {
                        if (attribute.value === 'unjail' || attribute.value === 'edit_validator') {
                            findAndUpdateValidator(tx.events)
                                .catch(err => {
                                    errors.Log(err, 'FIND_AND_UPDATE_VALIDATOR')
                                })
                        }
                    });
                });
            }
        });
    }
}

async function findAndUpdateValidator(events) {
    let messageEvent;
    let operatorAddress;
    for (let i = 0; i < events.length; i++) {
        if (events[i].type === 'message') {
            messageEvent = events[i];
            break;
        }
    }
    if (messageEvent) {
        for (let i = 0; i < messageEvent.attributes.length; i++) {
            if (messageEvent.attributes[i].key === 'sender') {
                operatorAddress = messageEvent.attributes[i].value;
                break;
            }
        }
    }
    if (operatorAddress) {
        console.log('Updating validator ' + operatorAddress + '...');
        updateValidatorDetails(operatorAddress);
    }
}


function updateValidatorDetails(operatorAddress) {
    httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/validators/${operatorAddress}`)
        .then(data => JSON.parse(data))
        .then(json => {
            let validator = json.result;       // with cosmos version upgrade, change here
            let hexAddress = addressOperations.getHexAddress(addressOperations.bech32ToPubkey(validator.consensus_pubkey));
            let selfDelegationAddress = addressOperations.getDelegatorAddrFromOperatorAddr(validator.operator_address);
            let validatorData = newValidatorObject(hexAddress, selfDelegationAddress, validator.operator_address,
                validator.consensus_pubkey, validator.jailed, validator.description);
            dataUtils.upsertOne(dataUtils.validatorCollection, {operatorAddress: validator.operator_address}, {$set: validatorData})
                .then(console.log(validator.operator_address + ' was updated.'))
                .catch(err => errors.Log(err, 'UPDATING_VALIDATORS'));
        })
        .catch(err => {
            errors.Log(err, 'UPDATING_VALIDATORS');
        });
}

function newValidatorObject(hexAddress, selfDelegateAddress, operatorAddress, consensusPublicKey, jailed, description) {
    return {
        hexAddress: hexAddress,
        selfDelegateAddress: selfDelegateAddress,
        operatorAddress: operatorAddress,
        consensusPublicKey: consensusPublicKey,
        jailed: jailed,
        description: description,
    };
}

module.exports = {addressOperations, updateValidatorDetails, newValidatorObject, wsTxError};
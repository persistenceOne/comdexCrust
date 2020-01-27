const errors = require('./errors');
const jsonUtils = require('./json');
const WebSocket = require('ws');
const wsConstants = require('../constants/websocket');

//When using labels, bot directly goes into subEvent when `bot.sendMessage(chatID, message, {ask: subEvent})`
// is called without waiting for the user to reply. This method makes it await somehow.
async function waitForUserReply(bot, chatID, message, subEvent, options) {
    if (!options) {
        await bot.eventList;
        bot.sendMessage(chatID, message, {ask: subEvent})
        .catch((err) => errors.Log(err))
    } else {
        await bot.eventList;
        bot.sendMessage(chatID, message, {ask: subEvent}, options)
        .catch((err) => errors.Log(err))
    }
}

function sendMessage(bot, chatID, message, options) {
    if (!options) {
        bot.sendMessage(chatID, message)
        .catch((err) => errors.Log(err))
    } else {
        bot.sendMessage(chatID, message, options)
        .catch((err) => errors.Log(err))
    }
}

let wsTx;

const reinitWSTx = () => {
    wsTx = new WebSocket(wsConstants.url);
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
// use httpUtil.httpGet(config.node.url, config.node.abciPort, `/tx_search?query="tx.height=${height}"&per_page=30`)
// to query and update.
function wsTxIncoming(data) {
    let json = jsonUtils.Parse(data, 'WS_INCOMING');
    if (json === undefined) {
        errors.Log('Error empty data from ws connection.');
    }
    if (errors.isEmpty(json.result)) {
        console.log('ws Tx Connected!');
    } else {
        let txs = JSON.parse(json.result.data.value.TxResult.result.log);
        txs.forEach((tx) => {
            if (tx.success) {
                tx.events.forEach((event) => {
                    if (event.type === 'edit_validator') {
                        findAndUpdateValidator(tx.events)
                            .catch(err => {
                                errors.Log(err, 'FIND_AND_UPDATE_VALIDATOR')
                            })
                    }
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
    httpUtil.httpGet(config.node.url, config.node.lcdPort, `/staking/validators/${operatorAddress}`)
        .then(data => JSON.parse(data))
        .then(json => {
            let validator = json;       // with cosmos version upgrade, change here   
            let hexAddress = validatorUtils.getHexAddress(validatorUtils.bech32ToPubkey(validator.consensus_pubkey));
            let selfDelegationAddress = validatorUtils.getDelegatorAddrFromOperatorAddr(validator.operator_address);
            let validatorData = newValidatorObject(hexAddress, selfDelegationAddress, validator.operator_address,
                validator.consensus_pubkey, validator.jailed, validator.description);
            dataUtil.upsertOne(dataUtil.validatorCollection, {operatorAddress: validator.operator_address}, {$set: validatorData})
                .then((res, err) => {
                    console.log(validator.operator_address + ' was updated.');
                })
                .catch(err => errors.Log(err, 'UPDATING_VALIDATORS'));
        })
        .catch(err => {
            errors.Log(err, 'UPDATING_VALIDATORS');
        });
}

module.exports = {waitForUserReply, sendMessage, updateValidatorDetails, wsTxError};
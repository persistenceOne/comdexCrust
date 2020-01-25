const TeleBot = require('telebot');
const WebSocket = require('ws');
const wsConstants = require('./constants/websocket');
const config = require('./config.json');
const dataUtils = require('./utilities/data');
const errors = require('./utilities/errors');
const jsonUtils = require('./utilities/json');
const Buttons = require('./constants/buttons');
const validatorUtils = require('./utilities/validator');
const chainUtils = require('./utilities/chain');
const botUtils = require('./utilities/bot');
const HttpUtils = require('./utilities/httpRequest');
const httpUtils = new HttpUtils();

dataUtils.InitializeDB();

const bot = new TeleBot({
    token: config.botToken,
    usePlugins: ['namedButtons', 'askUser'],
    pluginFolder: '../plugins/',
    pluginConfig: {
        namedButtons: {
            buttons: Buttons
        }
    }
});

bot.on(['/start', '/home'], msg => {
    let replyMarkup = bot.keyboard([
        [Buttons.chain.label, Buttons.hide.label],
    ], {resize: true});
    return botUtils.sendMessage(bot, msg.chat.id, `How can I help you?`, {replyMarkup});
});

bot.on('/hide', msg => {
    return botUtils.sendMessage(bot, msg.chat.id, 'Keyboard is now hidden. Type /start to re-enable.', {replyMarkup: 'hide'});
});

bot.on('/help', (msg) => {
    return botUtils.sendMessage(bot, msg.chat.id, `\`/start\` to start using the bot.`, {parseMode: 'Markdown'});
});

bot.on(/^\/say (.+)$/, (msg, props) => {
    const text = props.match[1];
    return botUtils.sendMessage(bot, msg.chat.id, text, {replyToMessage: msg.message_id});
});

bot.on('edit', (msg) => {
    return msg.reply.text('No editing of commands supported. Please re-enter the command.', {asReply: true});
});

bot.on(['keyboard', 'button', 'inlineKeyboard', 'inlineQueryKeyboard', 'inlineButton'], (msg) => {
    return msg.reply.text('No editing of commands supported. Please re-enter the command.', {asReply: true});
});

bot.on(['/chain', '/back'], msg => {
    let replyMarkup = bot.keyboard([
        [Buttons.nodeQuery.label, Buttons.chainQuery.label],
        [Buttons.alerts.label, Buttons.lcdQuery.label],
        [Buttons.home.label, Buttons.hide.label]
    ], {resize: true});
    return botUtils.sendMessage(bot, msg.chat.id, 'How can I help you?', {replyMarkup});
});

bot.on(['/node_queries'], msg => {
    let replyMarkup = bot.keyboard([
        [Buttons.nodeStatus.label, Buttons.lastBlock.label],
        [Buttons.peersCount.label, Buttons.peersList.label],
        [Buttons.back.label, Buttons.home.label, Buttons.hide.label]
    ], {resize: true});

    return botUtils.sendMessage(bot, msg.chat.id, 'What would you like to query?', {replyMarkup});
});

bot.on(['/chain_queries'], msg => {
    let replyMarkup = bot.keyboard([
        [Buttons.consensusState.label, Buttons.consensusParams.label],
        [Buttons.validatorsCount.label, Buttons.validatorsList.label],
        [Buttons.validatorInfo.label, Buttons.blockLookup.label],
        [Buttons.txLookup.label, Buttons.txByHeight.label],
        [Buttons.back.label, Buttons.home.label, Buttons.hide.label]
    ], {resize: true});
    return botUtils.sendMessage(bot, msg.chat.id, 'What would you like to query?', {replyMarkup});
});

bot.on(['/alerts'], msg => {
    let replyMarkup = bot.keyboard([
        [Buttons.subscribe.label, Buttons.unsubscribe.label],
        [Buttons.back.label, Buttons.home.label, Buttons.hide.label]
    ], {resize: true});
    return botUtils.sendMessage(bot, msg.chat.id, 'What would you like to query?', {replyMarkup});
});

bot.on(['/lcd_queries'], msg => {
    let replyMarkup = bot.keyboard([
        [Buttons.accountBalance.label, Buttons.delegatorRewards.label, Buttons.validatorRewards.label],
        [Buttons.stakingPool.label, Buttons.stakingParams.label, Buttons.mintingInflation.label],
        [Buttons.slashingParams.label, Buttons.mintingParams.label, Buttons.validatorSigning.label],
        [Buttons.back.label, Buttons.home.label, Buttons.hide.label]
    ], {resize: true});
    return botUtils.sendMessage(bot, msg.chat.id, 'What would you like to query?', {replyMarkup});
});

bot.on('/subscribe', async (msg) => {
    return botUtils.waitForUserReply(bot, msg.chat.id, `What\'s the validator\'s operator address?`, 'valAddr', {parseMode: 'Markdown'});
});

bot.on('ask.valAddr', msg => {
    const valAddr = msg.text;
    const chatID = msg.chat.id;

    if (!validatorUtils.verifyValidatorOperatorAddress(valAddr)) {
        return botUtils.sendMessage(bot, chatID, errors.INVALID_ADDRESS, {parseMode: 'Markdown'});
    }

    httpUtils.httpGet(config.node.url, config.node.lcdPort, `/staking/validators/${valAddr}`)
        .then(data => JSON.parse(data))
        .then(async json => {
            let validator = json;               // with cosmos version upgrade, change here
            if (validator.jailed) {
                return botUtils.sendMessage(bot, chatID, `Validator is jailed right now. Cannot subscribe to it.`, {parseMode: 'Markdown'});
            }
            let hexAddress = validatorUtils.getHexAddress(validatorUtils.bech32ToPubkey(validator.consensus_pubkey));
            let selfDelegationAddress = validatorUtils.getDelegatorAddrFromOperatorAddr(validator.operator_address);
            let validatorData = chainUtils.newValidatorObject(hexAddress, selfDelegationAddress, validator.operator_address,
                validator.consensus_pubkey, validator.jailed, validator.description);
            await dataUtils.upsertOne(dataUtils.validatorCollection, {operatorAddress: validator.operator_address}, {$set: validatorData})
                .then((res, err) => {
                    if (err) {
                        errors.Log(err, 'SUBSCRIBE_UPDATING_VALIDATORS');
                        return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                    }
                })
                .catch(err => {
                    errors.Log(err, 'SUBSCRIBE_UPDATING_VALIDATORS');
                    return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                });
            let query = {operatorAddress: valAddr};
            dataUtils.find(dataUtils.subscriberCollection, query)
                .then((result, err) => {
                    if (err) {
                        errors.Log(err, 'SUBSCRIBE_FIND');
                        return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                    }
                    let validatorSubscribers = result[0];
                    if (result.length === 0) {
                        let subscribers = [];
                        subscribers.push({chatID: chatID});
                        dataUtils.insertOne(dataUtils.subscriberCollection, {
                            operatorAddress: valAddr,
                            subscribers: subscribers
                        })
                            .then((res, err) => {
                                return botUtils.sendMessage(bot, chatID, `You are subscribed.`, {parseMode: 'Markdown'});
                            })
                            .catch(err => {
                                errors.Log(err, 'SUBSCRIBE_INSERT');
                                return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                            });
                    } else {
                        let subscribers = validatorSubscribers.subscribers;
                        let newSubscribers = [];
                        if (subscribers.length === 0) {
                            newSubscribers.push({chatID: chatID});
                        } else {
                            for (let i = 0; i < subscribers.length; i++) {
                                if (subscribers[i].chatID === chatID) {
                                    return botUtils.sendMessage(bot, chatID, `You are already subscribed to the validator: \`${valAddr}\`.`, {parseMode: 'Markdown'});
                                }
                            }
                            subscribers.push({chatID: chatID});
                            newSubscribers = subscribers;
                        }
                        dataUtils.updateOne(dataUtils.subscriberCollection, query, {
                            $set: {
                                operatorAddress: valAddr,
                                subscribers: newSubscribers
                            }
                        })
                            .then((res, err) => {
                                return botUtils.sendMessage(bot, chatID, `You are subscribed.`, {parseMode: 'Markdown'});
                            })
                            .catch(err => {
                                errors.Log(err, 'SUBSCRIBE_UPDATE');
                                return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                            });
                    }
                })
                .catch(err => {
                    errors.Log(err, 'SUBSCRIBE');
                    botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                })
        })
        .catch(e => {
            errors.Log(e, 'SUBSCRIBE');
            if (e.statusCode === 400) {
                botUtils.sendMessage(bot, chatID, errors.INVALID_ADDRESS, {parseMode: 'Markdown'});
            } else {
                botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
            }
        });
});

bot.on('/unsubscribe', (msg) => {
    return botUtils.waitForUserReply(bot, msg.chat.id, `What\'s the validator\'s operator address to unsubscribe?`, 'valAddrUnsub');
});

bot.on('ask.valAddrUnsub', msg => {
    const valAddr = msg.text;
    const chatID = msg.chat.id;

    if (!validatorUtils.verifyValidatorOperatorAddress(valAddr)) {
        return botUtils.sendMessage(bot, chatID, errors.INVALID_ADDRESS, {parseMode: 'Markdown'});
    }
    let query = {operatorAddress: valAddr};
    dataUtils.find(dataUtils.subscriberCollection, query)
        .then((result, err) => {
            if (err) {
                errors.Log(e, 'UNSUBSCRIBE_FIND');
                return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
            }
            if (result.length !== 1) {
                errors.Log(e, 'UNSUBSCRIBE_FIND');
                return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
            }

            let validatorSubscribers = result[0];
            if (result.length === 0 || validatorSubscribers.subscribers.length === 0) {
                return botUtils.sendMessage(bot, chatID, `You are not subscribed to validator.`, {parseMode: 'Markdown'});
            } else {
                let oldSubscribers = validatorSubscribers.subscribers;

                let removeByAttribute = jsonUtils.RemoveByAttribute(oldSubscribers, 'chatID', chatID);
                if (!removeByAttribute.removed) {
                    return botUtils.sendMessage(bot, chatID, `You are not subscribed to validator.`, {parseMode: 'Markdown'});
                } else {
                    dataUtils.updateOne(dataUtils.subscriberCollection, query, {
                        $set: {
                            operatorAddress: valAddr,
                            subscribers: removeByAttribute.newList
                        }
                    })
                        .then((res, err) => {
                            return botUtils.sendMessage(bot, chatID, `You are now unsubscribed to the validator: \`${valAddr}\`.`, {parseMode: 'Markdown'});
                        })
                        .catch(err => {
                            errors.Log(e, 'UNSUBSCRIBE_UPDATE');
                            return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                        });
                }
            }
        })
        .catch(err => {
            errors.Log(err, 'UNSUBSCRIBE');
            return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
        });
});

// node info
bot.on('/node_status', (msg) => {
    chainUtils.queries.sendNodeInfo(bot, msg.chat.id);
});

// last block
bot.on('/last_block', (msg) => {
    chainUtils.queries.sendLastBlock(bot, msg.chat.id);
});


// peers count
bot.on('/peers_count', (msg) => {
    chainUtils.queries.sendPeersCount(bot, msg.chat.id);
});

// peers list
bot.on('/peers_list', (msg) => {
    chainUtils.queries.sendPeersList(bot, msg.chat.id);
});

// consensus state
bot.on('/consensus_state', (msg) => {
    chainUtils.queries.sendConsensusState(bot, msg.chat.id);
});

// consensus params
bot.on('/consensus_params', (msg) => {
    chainUtils.queries.sendConsensusParams(bot, msg.chat.id);
});

// validators count
bot.on('/validators_count', (msg) => {
    chainUtils.queries.sendValidatorsCount(bot, msg.chat.id);
});

// validators list
bot.on('/validators_list', (msg) => {
    chainUtils.queries.sendValidators(bot, msg.chat.id);
});

//validator info
bot.on('/validator_info', async (msg) => {
    return botUtils.waitForUserReply(bot, msg.chat.id, `Please provide an operator address.`, 'validatorInfo', {parseMode: 'Markdown'});
});

bot.on(['ask.validatorInfo'], async msg => {
    const addr = msg.text;
    const chatID = msg.chat.id;
    if (!validatorUtils.verifyValidatorOperatorAddress(addr)) {
        return botUtils.sendMessage(bot, chatID, 'Address is invalid!');
    }
    chainUtils.queries.sendValidatorInfo(bot, chatID, addr)
});

// block lookup
bot.on('/block_lookup', async (msg) => {
    return botUtils.waitForUserReply(bot, msg.chat.id, `Please provide a block height.`, 'blockHeight', {parseMode: 'Markdown'});
});

bot.on(['ask.blockHeight'], async msg => {
    chainUtils.queries.sendBlockInfo(bot, msg.chat.id, msg.text);
});

// tx by hash
bot.on('/tx_lookup', async (msg) => {
    return botUtils.waitForUserReply(bot, msg.chat.id, `Please provide a tx hash.`, 'txByHash', {parseMode: 'Markdown'});
});

bot.on(['ask.txByHash'], async msg => {
    chainUtils.queries.sendTxByHash(bot, msg.chat.id, msg.text);
});

// tx by height
bot.on('/tx_by_height', async (msg) => {
    return botUtils.waitForUserReply(bot, msg.chat.id, `Please provide a block height.`, 'txByHeight', {parseMode: 'Markdown'});
});

bot.on(['ask.txByHeight'], async msg => {
    chainUtils.queries.sendTxByHeight(bot, msg.chat.id, msg.text);
});

// account balance
bot.on('/account_balance', async (msg) => {
    return botUtils.waitForUserReply(bot, msg.chat.id, `Please provide an address.`, 'accountBalance', {parseMode: 'Markdown'});
});

bot.on(['ask.accountBalance'], async msg => {
    const addr = msg.text;
    const chatID = msg.chat.id;
    if (addr.length !== 45) {
        return botUtils.sendMessage(bot, chatID, 'Address is invalid!');
    } else {
        chainUtils.queries.sendBalance(bot, chatID, addr)
    }
});

// delegator rewards
bot.on('/delegator_rewards', async (msg) => {
    return botUtils.waitForUserReply(bot, msg.chat.id, `Please provide a delegator address.`, 'delegatorRewards', {parseMode: 'Markdown'});
});

bot.on(['ask.delegatorRewards'], async msg => {
    const addr = msg.text;
    const id = msg.chat.id;
    if (addr.length !== 45) {
        return botUtils.sendMessage(bot, id, 'Address is invalid!');
    } else {
        chainUtils.queries.sendDelRewards(bot, msg.chat.id, addr);
    }
});

// validator rewards
bot.on('/validator_rewards', async (msg) => {
    return botUtils.waitForUserReply(bot, msg.chat.id, `Please provide a validator address.`, 'validatorRewards', {parseMode: 'Markdown'});
});

bot.on(['ask.validatorRewards'], async msg => {
    const addr = msg.text;
    const id = msg.chat.id;
    if (addr.length !== 52) {
        return botUtils.sendMessage(bot, id, 'Address is invalid!');
    } else {
        chainUtils.queries.sendValRewards(bot, msg.chat.id, addr);
    }
});

// staking pool
bot.on('/staking_pool', async (msg) => {
    chainUtils.queries.sendStakingPool(bot, msg.chat.id);
});

// staking params
bot.on('/staking_params', async (msg) => {
    chainUtils.queries.sendStakingParams(bot, msg.chat.id)
});

// minting inflation
bot.on('/minting_inflation', async (msg) => {
    chainUtils.queries.sendMintingInflation(bot, msg.chat.id);
});

// slashing params
bot.on('/slashing_params', async (msg) => {
    chainUtils.queries.sendSlashingParams(bot, msg.chat.id);
});

// minting params
bot.on('/minting_params', async (msg) => {
    chainUtils.queries.sendMintingParams(bot, msg.chat.id);
});

// validator signing-info
bot.on('/validator_signing', async (msg) => {
    return botUtils.waitForUserReply(bot, msg.chat.id, `Please provide a validator public key address.`, 'validatorSigning', {parseMode: 'Markdown'});
});

bot.on(['ask.validatorSigning'], async msg => {
    const addr = msg.text;
    const chatID = msg.chat.id;
    if (addr.length !== 83) {
        return botUtils.sendMessage(bot, chatID, 'Address is invalid!');
    } else {
        chainUtils.queries.sendValSigningInfo(bot, msg.chat.id, addr);
    }
});

bot.connect();

let ws;

const reinitWS = () => {
    ws = new WebSocket(wsConstants.url);
    try {
        ws.on('open', wsOpen);
        ws.on('close', wsClose);
        ws.on('message', wsIncoming);
        ws.on('error', wsError);
    } catch (e) {
        errors.Log(e, 'WS_CONNECTION');
        ws.send(JSON.stringify(unsubscribeAllMsg));
        reinitWS();
    }
};

reinitWS();

function wsOpen() {
    ws.send(JSON.stringify(wsConstants.subscribeNewBlockMsg));
}

function wsClose(code, reason) {
    let err = {statusCode: code, message: 'WS connection closed:    ' + reason};
    errors.Log(err, 'WS_CONNECTION');
    reinitWS();
}

function wsError(err) {
    errors.Log(err, 'WS_CONNECTION');
    ws.send(JSON.stringify(wsConstants.unsubscribeAllMsg));
    ws.close();
}

function wsIncoming(data) {
    let json = jsonUtils.Parse(data, 'WS_INCOMING');
    if (json === undefined) {
        errors.Log('Error empty data from ws connection.');
    }
    if (errors.isEmpty(json.result)) {
        console.log('ws Connected!');
    } else {
        let latestBlockHeight = json.result.data.value.block.header.height;
        console.log(latestBlockHeight);
        checkAndSendMsgOnValidatorsAbsence(json, latestBlockHeight);
    }
}

async function checkAndSendMsgOnValidatorsAbsence(json, latestBlockHeight) {
    dataUtils.find(dataUtils.subscriberCollection, {})
        .then((result, err) => {
            if (err) {
                errors.Log(err, 'SEND_ALERT');
                return;
            }
            result.forEach((validatorSubscribers) => {
                let found = false;
                let i = 0;
                dataUtils.find(dataUtils.validatorCollection, {operatorAddress: validatorSubscribers.operatorAddress})
                    .then((result, err) => {
                        if (err) {
                            errors.Log(err);
                        }
                        if (result.length === 1) {
                            let validatorDetails = result[0];
                            let moniker = validatorDetails.description.moniker;
                            do {
                                if (!errors.isEmpty(json.result.data.value.block.last_commit.precommits[i])) {
                                    let hexAddress = json.result.data.value.block.last_commit.precommits[i].validator_address;
                                    if (validatorDetails.hexAddress === hexAddress) {
                                        found = true;
                                    }
                                }
                                i += 1;
                            } while (!found && i < json.result.data.value.block.last_commit.precommits.length);

                            if (!found) {
                                sendMsgToSubscribers(moniker, validatorSubscribers.subscribers, latestBlockHeight);
                            }
                        } else {
                            if (result.length === 0) {
                                botUtils.updateValidatorDetails(validatorSubscribers.operatorAddress);
                            } else {
                                errors.Log('Incorrect database');
                            }
                        }
                    })
                    .catch(err => {
                        errors.Log(err, 'SENDING_MESSAGE');
                    });
            });
        })
        .catch(err => {
            errors.Log(err, 'SENDING_MESSAGE');
        })
}

async function sendMsgToSubscribers(moniker, subscribersList, latestBlockHeight) {
    subscribersList.forEach((subscriber) => {
        botUtils.sendMessage(bot, subscriber.chatID, `Alert: \`${moniker} is absent at height \`${latestBlockHeight}`, {
            parseMode: 'Markdown',
            notification: true
        });
    });
}

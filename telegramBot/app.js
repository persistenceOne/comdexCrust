const TeleBot = require('telebot');
const WebSocket = require('ws');
const wsConstants = require('./constants/websocket');
const config = require('./config.json');
const dataUtils = require('./utilities/data');
const errors = require('./utilities/errors');
const jsonUtils = require('./utilities/json');
const Buttons = require('./constants/buttons');
const validatorUtils = require('./utilities/validator');
const subscriberUtils = require('./utilities/subscriber');
const chainUtils = require('./utilities/chain');
const botUtils = require('./utilities/bot');
const blockchainConstants = require('./constants/blockchain');
const blockchainUtils = require('./utilities/blockchain');
const HttpUtils = require('./utilities/httpRequest');
const httpUtils = new HttpUtils();

dataUtils.SetupDB(function () {
    validatorUtils.initializeValidatorDB();
    subscriberUtils.initializeSubscriberDB();
    blockchainUtils.initializeBlockchainDB();
});

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
        [Buttons.chain.command, Buttons.hide.command],
    ], {resize: true});
    return botUtils.sendMessage(bot, msg.chat.id, `How can I help you?`, {replyMarkup});
});

bot.on('/hide_keyboard', msg => {
    return botUtils.sendMessage(bot, msg.chat.id, 'Keyboard is now hidden. Type /start to re-enable.', {replyMarkup: 'hide'});
});

bot.on('/help', (msg) => {
    return botUtils.sendMessage(bot, msg.chat.id, `\`/start\` to start using the bot.`, {parseMode: 'Markdown'});
});

bot.on(/^\/say (.+)$/, (msg, props) => {
    const text = props.match[1];
    return botUtils.sendMessage(bot, msg.chat.id, text, {replyToMessage: msg.message_id});
});


bot.on(['/chain', '/back'], msg => {
    let replyMarkup = bot.keyboard([
        [Buttons.validatorQuery.command, Buttons.chainQuery.command],
        [Buttons.subscribe.command, Buttons.analyticsQuery.command],
        [Buttons.home.command, Buttons.hide.command]
    ], {resize: true});
    return botUtils.sendMessage(bot, msg.chat.id, 'How can I help you?', {replyMarkup});
});

bot.on(['/chain_queries'], msg => {
    let replyMarkup = bot.keyboard([
        [Buttons.accountBalance.command, Buttons.delegatorRewards.command],
        [Buttons.lastBlock.command, Buttons.blockLookup.command],
        [Buttons.txLookup.command, Buttons.txByHeight.command],
        [Buttons.back.command, Buttons.home.command, Buttons.hide.command]
    ], {resize: true});
    return botUtils.sendMessage(bot, msg.chat.id, 'What would you like to query?', {replyMarkup});
});

bot.on(['/validator_queries'], msg => {
    let replyMarkup = bot.keyboard([
        [Buttons.validatorsCount.command, Buttons.validatorInfo.command],
        [Buttons.validatorRewards.command, Buttons.lastMissedBlock.command],
        [Buttons.validatorsList.command, Buttons.validatorReport.command],
        [Buttons.back.command, Buttons.home.command, Buttons.hide.command]
    ], {resize: true});
    return botUtils.sendMessage(bot, msg.chat.id, 'What would you like to query?', {replyMarkup});
});

bot.on(['/subscribe'], msg => {
    let replyMarkup = bot.keyboard([
        [Buttons.validator.command, Buttons.unsubValidator.command],
        [Buttons.allValidators.command, Buttons.unsubAllValidators.command],
        [Buttons.back.command, Buttons.home.command, Buttons.hide.command]
    ], {resize: true});
    return botUtils.sendMessage(bot, msg.chat.id, 'What would you like to query?', {replyMarkup});
});

bot.on('/analytics_queries', msg => {
    let replyMarkup = bot.keyboard([
        [Buttons.votingPower.command, Buttons.commission.command],
        [Buttons.uptime.command, Buttons.topValidator.command],
        [Buttons.back.command, Buttons.home.command, Buttons.hide.command]
    ], {resize: true});
    return botUtils.sendMessage(bot, msg.chat.id, 'What would you like to query?', {replyMarkup});
});

bot.on('/top_validators_wrt_voting_power', async (msg) => {
    const chatID = msg.chat.id;
    let blockHeight = latestBlockHeight;
    httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/pool`)
        .then(data => JSON.parse(data))
        .then(json => {
            const totalBondedToken = parseInt(json.result.bonded_tokens, 10);      // with cosmos version upgrade, change here
            httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/validators`)
                .then(data => JSON.parse(data))
                .then(async (json) => {
                    let activeValidators = json.result;       // with cosmos version upgrade, change here
                    activeValidators.sort((a, b) => parseInt(b.tokens, 10) - parseInt(a.tokens, 10));
                    let topValidatorsLength;
                    if (activeValidators.length > 10) {
                        topValidatorsLength = 10;
                    } else {
                        topValidatorsLength = activeValidators.length;
                    }
                    dataUtils.find(dataUtils.subscriberCollection, {})
                        .then(async (validatorSubscribersList, err) => {
                            if (err) {
                                errors.Log(err, 'VOTING_POWER');
                                return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                            }
                            let validatorList = [];
                            for (let i = 0; i < topValidatorsLength; i++) {
                                let validatorSubscribe = validatorSubscribersList.filter(validatorSubscribers => (validatorSubscribers.operatorAddress === activeValidators[i].operator_address));
                                if (validatorSubscribe.length === 0) {
                                    subscriberUtils.initializeValidatorSubscriber(activeValidators[i].operator_address, blockHeight);
                                } else {
                                    validatorList.push(activeValidators[i]);
                                }
                                if (validatorList.length === topValidatorsLength) {
                                    break;
                                }
                            }
                            let message = `Top validators by voting power at block \`${blockHeight}\` are:\n\n`;
                            for (let i = 0; i < validatorList.length; i++) {
                                let valMessage = validatorUtils.getValidatorVotingPowerMessage(validatorList[i], totalBondedToken);
                                message = message + `(${i + 1})\n\n` + valMessage;
                            }
                            botUtils.sendMessage(bot, chatID, message, {parseMode: 'Markdown'});
                        })
                        .catch(err => {
                            errors.Log(err, 'VOTING_POWER');
                            botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                        })
                })
                .catch(err => {
                    errors.Log(err, 'VOTING_POWER');
                    botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                })
        })
        .catch(err => {
            errors.Log(err, 'VOTING_POWER');
            botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
        })
});

bot.on('/top_validators_wrt_commission', async (msg) => {
    const chatID = msg.chat.id;
    let blockHeight = latestBlockHeight;
    httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/validators`)
        .then(data => JSON.parse(data))
        .then(async json => {
            let activeValidators = json.result;       // with cosmos version upgrade, change here
            activeValidators.sort((a, b) => parseFloat(a.commission.commission_rates.rate) - parseFloat(b.commission.commission_rates.rate));
            let lowestCommissionRate = parseFloat(activeValidators[0].commission.commission_rates.rate);
            dataUtils.find(dataUtils.subscriberCollection, {})
                .then(async (validatorSubscribersList, err) => {
                    if (err) {
                        errors.Log(err, 'COMMISSION');
                        return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                    }
                    let message = `Validators by lowest commission rate \`${(lowestCommissionRate * 100.0).toFixed(2)}\` % at block \`${blockHeight}\` are:\n\n`;
                    for (let i = 0; i < activeValidators.length; i++) {
                        if (parseFloat(activeValidators[i].commission.commission_rates.rate) > lowestCommissionRate) {
                            break;
                        }
                        let validatorSubscribe = validatorSubscribersList.filter(validatorSubscribers => (validatorSubscribers.operatorAddress === activeValidators[i].operator_address));
                        if (validatorSubscribe.length === 0) {
                            subscriberUtils.initializeValidatorSubscriber(activeValidators[i].operator_address, blockHeight);
                        }
                        let valMessage = validatorUtils.getValidatorCommissionMessage(activeValidators[i]);
                        message = message + `(${i + 1})\n\n` + valMessage;
                        if ((i +1) % 20 === 0) {
                            await bot.sendMessage(chatID, message, {parseMode: 'Markdown'});
                            message = ``;
                        }
                    }
                    if (message !== ``) {
                        botUtils.sendMessage(bot, chatID, message, {parseMode: 'Markdown'});
                    }
                })
                .catch(err => {
                    errors.Log(err, 'COMMISSION');
                    botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                })
        })
        .catch(err => {
            errors.Log(err, 'COMMISSION');
            botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
        })
});

bot.on('/top_validators_wrt_commission_voting_power', async (msg) => {
    const chatID = msg.chat.id;
    let blockHeight = latestBlockHeight;
    httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/pool`)
        .then(data => JSON.parse(data))
        .then(json => {
            const totalBondedToken = parseInt(json.result.bonded_tokens, 10);      // with cosmos version upgrade, change here
            httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/validators`)
                .then(data => JSON.parse(data))
                .then(async json => {
                    let activeValidators = json.result;       // with cosmos version upgrade, change here
                    dataUtils.find(dataUtils.subscriberCollection, {})
                        .then(async (validatorSubscribersList, err) => {
                            if (err) {
                                errors.Log(err, 'VOTING_POWER_COMMISSION');
                                return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                            }
                            activeValidators.sort((a, b) => parseInt(b.tokens, 10) - parseInt(a.tokens, 10));
                            let topValidatorsLength;
                            if (activeValidators.length > 5) {
                                topValidatorsLength = 5;
                            } else {
                                topValidatorsLength = activeValidators.length;
                            }
                            let slicedValidator = activeValidators.slice(0, topValidatorsLength);
                            slicedValidator.sort((a, b) => parseFloat(a.commission.commission_rates.rate) - parseFloat(b.commission.commission_rates.rate));
                            let validatorSubscribe = validatorSubscribersList.filter(validatorSubscribers => (validatorSubscribers.operatorAddress === slicedValidator[0].operator_address));
                            if (validatorSubscribe.length === 0) {
                                subscriberUtils.initializeValidatorSubscriber(slicedValidator[0].operator_address, blockHeight);
                            }
                            let message = validatorUtils.getTopValidatorMessage(slicedValidator[0], totalBondedToken);
                            botUtils.sendMessage(bot, chatID, `Validator with voting power in top \`${topValidatorsLength}\` and lowest commission rate:\n\n` + message, {parseMode: 'Markdown'});
                        })
                        .catch(err => {
                            errors.Log(err, 'VOTING_POWER_COMMISSION');
                            botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                        })
                })
                .catch(err => {
                    errors.Log(err, 'VOTING_POWER_COMMISSION');
                    botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                })
        })
        .catch(err => {
            errors.Log(err, 'VOTING_POWER_COMMISSION');
            botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
        })
});

bot.on('/top_validators_wrt_uptime', async (msg) => {
    const chatID = msg.chat.id;
    let blockHeight = latestBlockHeight;
    httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/validators`)
        .then(data => JSON.parse(data))
        .then(json => {
            let activeValidators = json.result;       // with cosmos version upgrade, change here
            dataUtils.find(dataUtils.subscriberCollection, {})
                .then(async (validatorSubscribersList, err) => {
                    if (err) {
                        errors.Log(err, 'UPTIME_FIND');
                        return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                    }
                    validatorSubscribersList.sort((a, b) => (subscriberUtils.calculateUptime(b.blocksHistory) - subscriberUtils.calculateUptime(a.blocksHistory)));
                    let validatorList = [];
                    let highestUptime = 0.0;
                    for (let i = 0; i < activeValidators.length; i++) {
                        let validator = validatorSubscribersList.filter(validatorSubscribe => (activeValidators[i].operator_address === validatorSubscribe.operatorAddress));
                        if (validator.length !== 0) {
                            activeValidators[i].blocksHistory = validator[0].blocksHistory;
                            validatorList.push(activeValidators[i]);
                            let validatorUptime = subscriberUtils.calculateUptime(validatorList[i].blocksHistory);
                            if(i === 0) {
                                highestUptime = validatorUptime;
                            }
                            if (validatorUptime < highestUptime) {
                                break;
                            }
                        }
                    }
                    if (validatorList.length !== 0) {
                        let message = `Top validators with highest uptime \`${highestUptime}\`% (based on \`${config.subscriberBlockHistoryLimit}\` blocks if not shown) at height \`${blockHeight}\` are:\n\n`;
                        for (let i = 0; i < validatorList.length; i++) {
                            let validatorUptime = subscriberUtils.calculateUptime(validatorList[i].blocksHistory);
                            if (validatorUptime < highestUptime) {
                                break;
                            }
                            let valMessage = ``;
                            if (validatorList[i].blocksHistory.length !== config.subscriberBlockHistoryLimit) {
                                valMessage = `\`${validatorList[i].description.moniker}\` (based on \`${validatorList[i].blocksHistory.length}\` blocks) \n\n`;
                            } else {
                                valMessage = `\`${validatorList[i].description.moniker}\` \n\n`;
                            }
                            
                            message = message + valMessage;
                            if ((i +1) % 20 === 0) {
                                await bot.sendMessage(chatID, message, {parseMode: 'Markdown'});
                                message = ``;
                            }
                        }
                        if (message !== ``) {
                            botUtils.sendMessage(bot, chatID, message, {parseMode: 'Markdown'});
                        }
                    } else {
                        return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                    }
                })
                .catch(err => {
                    errors.Log(err, 'UPTIME');
                    botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                })
        })
        .catch(err => {
            errors.Log(err, 'UPTIME');
            botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
        })
});

bot.on('/one_validator', async (msg) => {
    return botUtils.waitForUserReply(bot, msg.chat.id, `What\'s the validator\'s operator address?`, 'validatorAddress', {parseMode: 'Markdown'});
});

bot.on('ask.validatorAddress', msg => {
    const valAddr = msg.text;
    const chatID = msg.chat.id;

    if (!validatorUtils.addressOperations.verifyValidatorOperatorAddress(valAddr)) {
        return botUtils.sendMessage(bot, chatID, errors.INVALID_ADDRESS, {parseMode: 'Markdown'});
    }

    httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/validators/${valAddr}`)
        .then(data => JSON.parse(data))
        .then(json => {
            let validator = json.result;               // with cosmos version upgrade, change here
            if (validator.jailed) {
                return botUtils.sendMessage(bot, chatID, `Validator is jailed right now. Cannot subscribe to it.`, {parseMode: 'Markdown'});
            }
            let query = {operatorAddress: valAddr};
            dataUtils.find(dataUtils.subscriberCollection, query)
                .then(async (result, err) => {
                    if (err) {
                        errors.Log(err, 'SUBSCRIBE_FIND');
                        return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                    }
                    if (result.length === 0) {
                        let subscribers = [];
                        subscribers.push({chatID: chatID});
                        let validatorSubscriber = subscriberUtils.newValidatorSubscribers(valAddr, latestBlockHeight, subscribers);
                        validatorUtils.updateValidatorDetails(valAddr);
                        dataUtils.insertOne(dataUtils.subscriberCollection, validatorSubscriber)
                            .then(botUtils.sendMessage(bot, chatID, `You are subscribed.`, {parseMode: 'Markdown'}))
                            .catch(err => {
                                errors.Log(err, 'SUBSCRIBE_INSERT');
                                return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                            });
                    } else {
                        let validatorSubscribers = result[0];
                        let subscribers = validatorSubscribers.subscribers;
                        let subscriberExists = false;
                        for (let i = 0; i < subscribers.length; i++) {
                            if (subscribers[i].chatID === chatID) {
                                subscriberExists = true;
                                break;
                            }
                        }
                        if (!subscriberExists || subscribers.length === 0) {
                            subscribers.push({chatID: chatID});
                            dataUtils.updateOne(dataUtils.subscriberCollection, query, {
                                $set: {
                                    subscribers: subscribers
                                }
                            })
                                .then(botUtils.sendMessage(bot, chatID, `You are subscribed.`, {parseMode: 'Markdown'}))
                                .catch(err => {
                                    errors.Log(err, 'SUBSCRIBE_UPDATE');
                                    return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                                });
                        } else {
                            return botUtils.sendMessage(bot, chatID, `You are already subscribed to the validator: \`${valAddr}\`.`, {parseMode: 'Markdown'});
                        }
                    }
                })
                .catch(err => {
                    errors.Log(err, 'SUBSCRIBE');
                    botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                })
        })
        .catch(e => {
            errors.Log(e, 'SUBSCRIBE');
            if (e.statusCode === 400 || e.statusCode === 404) {
                botUtils.sendMessage(bot, chatID, errors.INVALID_ADDRESS, {parseMode: 'Markdown'});
            } else {
                botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
            }
        });
});

bot.on('/all_validator', msg => {
    const chatID = msg.chat.id;
    let blockHeight = latestBlockHeight;
    httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/validators`)
        .then(data => JSON.parse(data))
        .then(async json => {
            let validators = json.result;               // with cosmos version upgrade, change here
            validators.forEach((validator) => {
                let valAddr = validator.operator_address;
                let query = {operatorAddress: valAddr};
                dataUtils.find(dataUtils.subscriberCollection, query)
                    .then((result, err) => {
                        if (err) {
                            errors.Log(err, 'SUBSCRIBE_FIND');
                            return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                        }
                        if (result.length === 0) {
                            let subscribers = [];
                            subscribers.push({chatID: chatID});
                            let validatorSubscriber = subscriberUtils.newValidatorSubscribers(valAddr, blockHeight, subscribers);
                            validatorUtils.updateValidatorDetails(valAddr);
                            dataUtils.insertOne(dataUtils.subscriberCollection, validatorSubscriber)
                                .catch(err => {
                                    errors.Log(err, 'SUBSCRIBE_INSERT');
                                    return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                                });
                        } else {
                            let validatorSubscribers = result[0];
                            let subscribers = validatorSubscribers.subscribers;
                            let subscriberExists = false;
                            for (let i = 0; i < subscribers.length; i++) {
                                if (subscribers[i].chatID === chatID) {
                                    subscriberExists = true;
                                    break;
                                }
                            }
                            if (!subscriberExists || subscribers.length === 0) {
                                subscribers.push({chatID: chatID});
                                dataUtils.updateOne(dataUtils.subscriberCollection, query, {
                                    $set: {
                                        subscribers: subscribers
                                    }
                                })
                                    .catch(err => {
                                        errors.Log(err, 'SUBSCRIBE_UPDATE');
                                        return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                                    });
                            }
                        }
                    })
                    .catch(err => {
                        errors.Log(err, 'SUBSCRIBE');
                        botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                    })
            });
        })
        .then(() => botUtils.sendMessage(bot, chatID, 'You have been subscribed to all validators.', {parseMode: 'Markdown'}))
        .catch(e => {
            errors.Log(e, 'SUBSCRIBE');
            if (e.statusCode === 400 || e.statusCode === 404) {
                botUtils.sendMessage(bot, chatID, errors.INVALID_ADDRESS, {parseMode: 'Markdown'});
            } else {
                botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
            }
        });
});

bot.on('/unsub_one_validator', (msg) => {
    return botUtils.waitForUserReply(bot, msg.chat.id, `What\'s the validator\'s operator address to unsubscribe?`, 'unsubValidatorAddress');
});

bot.on('ask.unsubValidatorAddress', msg => {
    const valAddr = msg.text;
    const chatID = msg.chat.id;

    if (!validatorUtils.addressOperations.verifyValidatorOperatorAddress(valAddr)) {
        return botUtils.sendMessage(bot, chatID, errors.INVALID_ADDRESS, {parseMode: 'Markdown'});
    }
    let query = {operatorAddress: valAddr};
    dataUtils.find(dataUtils.subscriberCollection, query)
        .then((result, err) => {
            if (err) {
                errors.Log(err, 'UNSUBSCRIBE_FIND');
                return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
            }
            if (result.length !== 1) {
                errors.Log('More than one validator object for same operator address.', 'UNSUBSCRIBE_FIND');
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
                            subscribers: removeByAttribute.newList
                        }
                    })
                        .then(botUtils.sendMessage(bot, chatID, `You are now unsubscribed to the validator: \`${valAddr}\`.`, {parseMode: 'Markdown'}))
                        .catch(err => {
                            errors.Log(err, 'UNSUBSCRIBE_UPDATE');
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

bot.on('/unsub_all_validator', msg => {
    const chatID = msg.chat.id;

    dataUtils.find(dataUtils.subscriberCollection, {})
        .then((result, err) => {
            if (err) {
                errors.Log(err, 'UNSUBSCRIBE_FIND');
                return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
            }
            result.forEach((validatorSubscribers) => {
                if (validatorSubscribers.subscribers.length !== 0) {
                    let oldSubscribers = validatorSubscribers.subscribers;
                    let removeByAttribute = jsonUtils.RemoveByAttribute(oldSubscribers, 'chatID', chatID);
                    let query = {operatorAddress: validatorSubscribers.operatorAddress};
                    dataUtils.updateOne(dataUtils.subscriberCollection, query, {
                        $set: {
                            subscribers: removeByAttribute.newList
                        }
                    })
                        .catch(err => {
                            errors.Log(err, 'UNSUBSCRIBE_UPDATE');
                            return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                        });
                }
            });
        })
        .then(() => botUtils.sendMessage(bot, chatID, `You have been unsubscribed to all validators you were subscribed to.`, {parseMode: 'Markdown'}))
        .catch(err => {
            errors.Log(err, 'UNSUBSCRIBE');
            return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
        });
});

// last block
bot.on('/last_block', (msg) => {
    chainUtils.queries.sendLastBlock(bot, msg.chat.id);
});
// validators count
bot.on('/validators_count', (msg) => {
    chainUtils.queries.sendValidatorsCount(bot, msg.chat.id);
});

// validators list
bot.on('/validators_list', (msg) => {
    chainUtils.queries.sendValidators(bot, msg.chat.id, latestBlockHeight);
});

//validator info
bot.on('/validator_info', async (msg) => {
    return botUtils.waitForUserReply(bot, msg.chat.id, `Please provide an operator address.`, 'validatorInfo', {parseMode: 'Markdown'});
});

bot.on(['ask.validatorInfo'], async msg => {
    const addr = msg.text;
    const chatID = msg.chat.id;
    if (!validatorUtils.addressOperations.verifyValidatorOperatorAddress(addr)) {
        return botUtils.sendMessage(bot, chatID, 'Address is invalid!');
    }
    chainUtils.queries.sendValidatorInfo(bot, chatID, addr, latestBlockHeight)
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

bot.on(['ask.txByHash'], async (msg) => {
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
        chainUtils.queries.sendBalance(bot, chatID, addr, latestBlockHeight)
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
        chainUtils.queries.sendDelRewards(bot, msg.chat.id, addr, latestBlockHeight);
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
        chainUtils.queries.sendValRewards(bot, msg.chat.id, addr, latestBlockHeight);
    }
});

// Last Missed Blocks
bot.on('/last_missed_blocks', async (msg) => {
    return botUtils.waitForUserReply(bot, msg.chat.id, `Please provide a validator address.`, 'lastMissedBlockValidatorAddress', {parseMode: 'Markdown'});
});

bot.on(['ask.lastMissedBlockValidatorAddress'], async msg => {
    const addr = msg.text;
    const chatID = msg.chat.id;
    let blockHeight = latestBlockHeight;
    if (!validatorUtils.addressOperations.verifyValidatorOperatorAddress(addr)) {
        return botUtils.sendMessage(bot, chatID, errors.INVALID_ADDRESS, {parseMode: 'Markdown'});
    } else {
        httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/validators/${addr}`)
            .then(data => JSON.parse(data))
            .then(json => {
                let validator = json.result;       // with cosmos version upgrade, change here
                dataUtils.find(dataUtils.subscriberCollection, {operatorAddress: addr})
                    .then((validatorSubscribers, err) => {
                        if (err) {
                            errors.Log(err, 'LAST_MISSED_BLOCK');
                            return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                        }
                        if (validatorSubscribers.length === 0) {
                            subscriberUtils.initializeValidatorSubscriber(validator.operatorAddress, blockHeight);
                            return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                        } else {
                            let moniker = validator.description.moniker;
                            let missedBlocks = [];
                            let missedBlocksEmoticonMessage = '';
                            for (let i = 0; i < validatorSubscribers[0].blocksHistory.length; i++) {
                                if (validatorSubscribers[0].blocksHistory[i].found) {
                                    missedBlocksEmoticonMessage = missedBlocksEmoticonMessage + blockchainConstants.emoticon.presentBlock + ' ';
                                } else {
                                    missedBlocks.unshift(validatorSubscribers[0].blocksHistory[i].block);
                                    missedBlocksEmoticonMessage = missedBlocksEmoticonMessage + blockchainConstants.emoticon.missedBlock + ' ';
                                }
                                if ((i + 1) % 10 === 0) {
                                    missedBlocksEmoticonMessage = missedBlocksEmoticonMessage + '\n';
                                }
                            }
                            missedBlocks = JSON.stringify(missedBlocks);
                            let message = '';
                            if (validatorSubscribers[0].lastMissedBlock !== 0) {
                                message = `\`${moniker}\` last missed \`${validatorSubscribers[0].lastMissedBlock}\` block.\n\n`;
                            } else {
                                message = `\`${moniker}\` has not missed any blocks since start of this Bot.\n\n`;
                            }
                            message = message + `Latest Missed Block List: \`${missedBlocks}\`\n\n` + missedBlocksEmoticonMessage;
                            botUtils.sendMessage(bot, chatID, message, {parseMode: 'Markdown'});
                        }
                    })
            })
            .catch(err => botUtils.handleErrors(bot, chatID, err, 'LAST_MISSED_BLOCK'));
    }
});

bot.on('/validator_report', async (msg) => {
    return botUtils.waitForUserReply(bot, msg.chat.id, `Please provide a validator address.`, 'validatorAddressReport', {parseMode: 'Markdown'});
});

bot.on(['ask.validatorAddressReport'], async msg => {
    const addr = msg.text;
    const chatID = msg.chat.id;
    let blockHeight = latestBlockHeight;
    let baseBlockHeight = Math.floor(blockHeight / 10000) * 10000;
    if (baseBlockHeight === 0) {
        baseBlockHeight = 1;
    }
    if (!validatorUtils.addressOperations.verifyValidatorOperatorAddress(addr)) {
        return botUtils.sendMessage(bot, chatID, errors.INVALID_ADDRESS, {parseMode: 'Markdown'});
    } else {
        httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/validators/${addr}?height=${blockHeight}`)
            .then(newData => JSON.parse(newData))
            .then(newJson => {
                let newValidator = newJson.result;       // with cosmos version upgrade, change here
                httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/validators/${addr}?height=${baseBlockHeight}`)
                    .then(oldData => JSON.parse(oldData))
                    .then(oldJson => {
                        let oldValidator = oldJson.result;       // with cosmos version upgrade, change here
                        dataUtils.find(dataUtils.subscriberCollection, {operatorAddress: addr})
                            .then((validatorSubscribers, err) => {
                                if (err) {
                                    errors.Log(err, 'SENDING_REPORT');
                                    return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                                }
                                if (validatorSubscribers.length === 0) {
                                    subscriberUtils.initializeValidatorSubscriber(validator.operatorAddress, blockHeight);
                                    return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                                } else {
                                    httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/pool?height=${baseBlockHeight}`)
                                        .then(data => JSON.parse(data))
                                        .then(json => {
                                            let oldTotalBondedTokens = json.result.bonded_tokens;       // with cosmos version upgrade, change here
                                            httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/pool?height=${blockHeight}`)
                                                .then(data => JSON.parse(data))
                                                .then(json => {
                                                    let newTotalBondedTokens = json.result.bonded_tokens;       // with cosmos version upgrade, change here
                                                    dataUtils.find(dataUtils.blockchainCollection, {})
                                                        .then((blockchainHistory, err) => {
                                                            if (err) {
                                                                errors.Log(err, 'SENDING_REPORT');
                                                                return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                                                            }
                                                            let message = validatorUtils.getValidatorReport(oldValidator, newValidator, oldTotalBondedTokens, newTotalBondedTokens, blockHeight, baseBlockHeight, validatorSubscribers[0].blocksHistory, blockchainHistory);
                                                            botUtils.sendMessage(bot, chatID, message, {
                                                                parseMode: 'Markdown',
                                                                notification: true
                                                            })
                                                        })
                                                        .catch(err => botUtils.handleErrors(bot, chatID, err, 'SENDING_REPORT'));
                                                })
                                                .catch(err => botUtils.handleErrors(bot, chatID, err, 'SENDING_REPORT'));
                                        })
                                        .catch(err => botUtils.handleErrors(bot, chatID, err, 'SENDING_REPORT'));
                                }
                            })
                            .catch(err => botUtils.handleErrors(bot, chatID, err, 'SENDING_REPORT'));
                    })
                    .catch(err => botUtils.handleErrors(bot, chatID, err, 'SENDING_REPORT'));
            })
            .catch(err => botUtils.handleErrors(bot, chatID, err, 'SENDING_REPORT'));
    }
});

bot.connect();

let ws;

const reinitWS = () => {
    if (ws === undefined) {
        ws = new WebSocket(wsConstants.url);
    } else {
        if (ws.url === wsConstants.url) {
            ws = new WebSocket(wsConstants.backupURL);
            botUtils.nodeURL = config.node.backupURL;
        }
        if (ws.url === wsConstants.backupURL) {
            ws = new WebSocket(wsConstants.url);
            botUtils.nodeURL = config.node.url;
        }
    }
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

let latestBlockHeight = 1;
let oldBlockHeight = 0;
let slashingWindow = config.slashingWindow;

function scheduler() {
    if (latestBlockHeight === oldBlockHeight) {
        wsError('WS Connection Freezed');
    } else {
        oldBlockHeight = latestBlockHeight;
    }
}

setInterval(scheduler, 120000);

function wsIncoming(data) {
    let json = JSON.parse(data);
    if (errors.isEmpty(json.result)) {
        console.log('ws Connected!');
    } else {
        let currentBlockHeight = parseInt(json.result.data.value.block.header.height, 10);
        latestBlockHeight = currentBlockHeight;
        console.log(currentBlockHeight);
        validatorUtils.checkTxs(currentBlockHeight)
            .catch(err => errors.Log(err, 'CHECKING_TXS'));
        checkAndSendMsgOnValidatorsAbsence(json)
            .catch(err => errors.Log(err, 'CHECK_UPDATE_COUNTER_SEND_MESSAGE'));
        blockchainUtils.updateBlock(currentBlockHeight)
            .catch(err => errors.Log(err, 'UPDATING_BLOCKCHAIN_DB'));
        if (currentBlockHeight % 10000 === 0) {
            sendDailyReports(currentBlockHeight)
                .catch(err => errors.Log(err, 'SENDING_REPORTS'));
            httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/slashing/parameters`)
                .then(data => JSON.parse(data))
                .then(async (json) => {
                    let signed_blocks_window = parseInt(json.result.signed_blocks_window, 10);
                    let min_signed_per_window = parseFloat(json.result.min_signed_per_window);
                    slashingWindow = signed_blocks_window * min_signed_per_window;
                })
                .catch(err => errors.Log(err, 'UPDATING_SLASHING_WINDOW'));
        }
    }
}

async function checkAndSendMsgOnValidatorsAbsence(json) {
    let blockHeight = parseInt(json.result.data.value.block.header.height, 10);
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
                        if (result.length === 0) {
                            validatorUtils.updateValidatorDetails(validatorSubscribers.operatorAddress);
                        }
                        if (result.length === 1) {
                            let validatorDetails = result[0];
                            do {
                                if (!errors.isEmpty(json.result.data.value.block.last_commit.precommits[i])) {
                                    let hexAddress = json.result.data.value.block.last_commit.precommits[i].validator_address;
                                    if (validatorDetails.hexAddress === hexAddress) {
                                        found = true;
                                    }
                                }
                                i += 1;
                            } while (!found && i < json.result.data.value.block.last_commit.precommits.length);
                            updateCounterAndSendMessage(validatorSubscribers, validatorDetails, blockHeight, found);
                        } else {
                            errors.Log('Incorrect database');
                        }
                    })
                    .catch(err => errors.Log(err, 'SENDING_MESSAGE'));
            });
        })
        .catch(err => errors.Log(err, 'SENDING_MESSAGE'));
}

function updateCounterAndSendMessage(validatorSubscribers, validatorDetails, blockHeight, found) {
    let query = {operatorAddress: validatorSubscribers.operatorAddress};
    if (!validatorDetails.jailed) {
        let blocksHistory = validatorSubscribers.blocksHistory;
        blocksHistory.push({block: blockHeight, found: found});
        if (blocksHistory.length > config.subscriberBlockHistoryLimit) {
            blocksHistory.shift();
        }
        if (!found) {
            let consecutiveCounter = validatorSubscribers.consecutiveCounter;
            let alertLevel = subscriberUtils.getAlertLevel(slashingWindow, consecutiveCounter);
            if (blockHeight - validatorSubscribers.lastMissedBlock === 1) {
                consecutiveCounter = validatorSubscribers.consecutiveCounter + 1;
            }
            dataUtils.updateOne(dataUtils.subscriberCollection, query, {
                $set: {
                    consecutiveCounter: consecutiveCounter,
                    alertLevel: alertLevel,
                    lastMissedBlock: blockHeight,
                    blocksHistory: blocksHistory
                }
            })
                .then((result) => {
                    let blocksLevel = subscriberUtils.getBlocksLevel(slashingWindow, alertLevel);
                    if (consecutiveCounter !== 0 && consecutiveCounter % blocksLevel === 0) {
                        sendMissedMsgToSubscribers(validatorDetails.description.moniker, validatorSubscribers.subscribers, consecutiveCounter, alertLevel)
                            .catch(err => errors.Log(err));
                        checkJailedStatusAndSendMessage(validatorDetails, validatorSubscribers)
                            .catch(err => errors.Log(err));
                    }
                })
                .catch(err => errors.Log(err, 'UPDATING_COUNTER_AND_SENDING_MESSAGE'));
        } else {
            dataUtils.updateOne(dataUtils.subscriberCollection, query, {
                $set: {
                    consecutiveCounter: 0,
                    alertLevel: 1,
                    blocksHistory: blocksHistory
                }
            })
                .catch(err => errors.Log(err, 'UPDATING_COUNTER_AND_SENDING_MESSAGE'));
        }
    }
}

async function sendMissedMsgToSubscribers(moniker, subscribersList, consecutiveCounter, alertLevel) {
    let emoji = '';
    for (let i = 0; i < alertLevel; i++) {
        emoji = emoji + ' ' + blockchainConstants.emoticon.alert;
        if (i === 4) {
            break;
        }
    }
    subscribersList.forEach((subscriber) => {
        botUtils.sendMessage(bot, subscriber.chatID, `Alert: \`${moniker}\` has consecutively missed \`${consecutiveCounter}\` blocks \`${emoji}\`.`, {
            parseMode: 'Markdown',
            notification: true
        });
    });
}

async function checkJailedStatusAndSendMessage(validatorDetails, validatorSubscribers) {
    let operatorAddress = validatorDetails.operatorAddress;
    httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/validators/${operatorAddress}`)
        .then(async data => {
            let json = JSON.parse(data);
            if (json.error) {
                errors.Log('Invalid Operator Address', 'UPDATE_COUNTER_QUERY_VALIDATOR')
            } else {
                let validator = json.result;       // with cosmos version upgrade, change here
                if (validator.jailed) {
                    dataUtils.updateOne(dataUtils.validatorCollection, {operatorAddress: operatorAddress}, {
                        $set: {
                            jailed: true,
                        }
                    })
                        .then(sendJailedMsgToSubscribers(validatorDetails.description.moniker, validatorSubscribers.subscribers))
                        .catch(err => errors.Log(err, 'UPDATING_COUNTER_UPDATE_VALIDATOR'));
                }
            }
        })
        .catch(err => errors.Log(err, 'UPDATE_COUNTER_QUERY_VALIDATOR'));
}

async function sendJailedMsgToSubscribers(moniker, subscribersList) {
    subscribersList.forEach((subscriber) => {
        botUtils.sendMessage(bot, subscriber.chatID, `Alert: \`${moniker}\` has been jailed.`, {
            parseMode: 'Markdown',
            notification: true
        });
    });
}

async function sendDailyReports(blockHeight) {
    dataUtils.find(dataUtils.subscriberCollection, {})
        .then((result, err) => {
            if (err) {
                errors.Log(err, 'SEND_REPORT');
                return;
            }
            let oldBlockHeight = blockHeight - 10000;
            if (oldBlockHeight <= 0) {
                oldBlockHeight = 1;
            }
            dataUtils.find(dataUtils.blockchainCollection, {})
                .then((blockchainHistory, err) => {
                    if (err) {
                        errors.Log(err, 'SENDING_REPORT');
                        return;
                    }
                    result.forEach((validatorSubscribers) => {
                        httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/validators/${validatorSubscribers.operatorAddress}?height=${oldBlockHeight}`)
                            .then(data => JSON.parse(data))
                            .then(async (json) => {
                                let oldValidatorDetails = json.result;       // with cosmos version upgrade, change here
                                httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/validators/${validatorSubscribers.operatorAddress}?height=${blockHeight}`)
                                    .then(data => JSON.parse(data))
                                    .then(json => {
                                        let latestValidatorDetails = json.result;       // with cosmos version upgrade, change here
                                        if (!latestValidatorDetails.jailed) {
                                            httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/pool?height=${oldBlockHeight}`)
                                                .then(data => JSON.parse(data))
                                                .then(json => {
                                                    let oldTotalBondedTokens = json.result.bonded_tokens;       // with cosmos version upgrade, change here
                                                    httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/pool?height=${blockHeight}`)
                                                        .then(data => JSON.parse(data))
                                                        .then(json => {
                                                            let newTotalBondedTokens = json.result.bonded_tokens;       // with cosmos version upgrade, change here
                                                            let message = validatorUtils.getValidatorReport(oldValidatorDetails, latestValidatorDetails, oldTotalBondedTokens, newTotalBondedTokens, blockHeight, oldBlockHeight, validatorSubscribers.blocksHistory, blockchainHistory);
                                                            validatorSubscribers.subscribers.forEach((subscriber) => {
                                                                botUtils.sendMessage(bot, subscriber.chatID, message, {
                                                                    parseMode: 'Markdown',
                                                                    notification: true
                                                                });
                                                            });
                                                        })
                                                        .catch(err => errors.Log(err, 'SENDING_REPORTS'));
                                                })
                                                .catch(err => errors.Log(err, 'SENDING_REPORTS'));
                                        }
                                    })
                                    .catch(err => errors.Log(err, 'SENDING_REPORTS'));
                            })
                            .catch(err => errors.Log(err, 'SENDING_REPORTS'));
                    });
                })
                .catch(err => errors.Log(err, 'SENDING_REPORTS'));
        })
        .catch(err => errors.Log(err, 'SENDING_REPORTS'));
}

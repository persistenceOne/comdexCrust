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
const HttpUtils = require('./utilities/httpRequest');
const httpUtils = new HttpUtils();

dataUtils.SetupDB(function () {
    validatorUtils.initializeValidatorDB();
    subscriberUtils.initializeSubscriberDB();
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


bot.on(['/chain', '/back'], msg => {
    let replyMarkup = bot.keyboard([
        [Buttons.validatorQuery.label, Buttons.chainQuery.label],
        [Buttons.alerter.label, Buttons.analyticsQuery.label],
        [Buttons.home.label, Buttons.hide.label]
    ], {resize: true});
    return botUtils.sendMessage(bot, msg.chat.id, 'How can I help you?', {replyMarkup});
});

bot.on(['/chain_queries'], msg => {
    let replyMarkup = bot.keyboard([
        [Buttons.accountBalance.label, Buttons.delegatorRewards.label],
        [Buttons.lastBlock.label, Buttons.blockLookup.label],
        [Buttons.txLookup.label, Buttons.txByHeight.label],
        [Buttons.back.label, Buttons.home.label, Buttons.hide.label]
    ], {resize: true});
    return botUtils.sendMessage(bot, msg.chat.id, 'What would you like to query?', {replyMarkup});
});

bot.on(['/validator_queries'], msg => {
    let replyMarkup = bot.keyboard([
        [Buttons.validatorsCount.label, Buttons.validatorInfo.label],
        [Buttons.validatorRewards.label, Buttons.lastMissedBlock.label],
        [Buttons.validatorsList.label],
        [Buttons.back.label, Buttons.home.label, Buttons.hide.label]
    ], {resize: true});
    return botUtils.sendMessage(bot, msg.chat.id, 'What would you like to query?', {replyMarkup});
});

bot.on(['/alerter'], msg => {
    let replyMarkup = bot.keyboard([
        [Buttons.validator.label, Buttons.unsubValidator.label],
        [Buttons.allValidators.label, Buttons.unsubAllValidators.label],
        [Buttons.back.label, Buttons.home.label, Buttons.hide.label]
    ], {resize: true});
    return botUtils.sendMessage(bot, msg.chat.id, 'What would you like to query?', {replyMarkup});
});

bot.on('/analytics_queries', msg => {
    let replyMarkup = bot.keyboard([
        [Buttons.votingPower.label, Buttons.commission.label],
        [Buttons.uptime.label, Buttons.topValidator.label],
        [Buttons.back.label, Buttons.home.label, Buttons.hide.label]
    ], {resize: true});
    return botUtils.sendMessage(bot, msg.chat.id, 'What would you like to query?', {replyMarkup});
});

bot.on('/voting_power', async (msg) => {
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
                                errors.Log(err, 'UPTIME_FIND');
                                return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                            }
                            let validatorList = [];
                            for (let i = 0; i < topValidatorsLength; i++) {
                                let validatorSubscribe = validatorSubscribersList.filter(validatorSubscribers => (validatorSubscribers.operatorAddress === activeValidators[i].operator_address));
                                if (validatorSubscribe.length === 0) {
                                    subscriberUtils.initializeValidatorSubscriber(activeValidators[i].operator_address, blockHeight);
                                } else {
                                    activeValidators[i].counter = validatorSubscribe[0].counter;
                                    activeValidators[i].initHeight = validatorSubscribe[0].initHeight;
                                    validatorList.push(activeValidators[i]);
                                }
                                if (validatorList.length === topValidatorsLength) {
                                    break;
                                }
                            }
                            if (validatorList.length !== 0) {
                                await bot.sendMessage(chatID, `Top validators by voting power are:`, {parseMode: 'Markdown'});
                                for (let i = 0; i < validatorList.length; i++) {
                                    let message = validatorUtils.getValidatorMessage(validatorList[i], totalBondedToken, validatorList[i].counter, validatorList[i].initHeight, blockHeight);
                                    await bot.sendMessage(chatID, `(${i + 1})\n\n` + message, {parseMode: 'Markdown'});
                                }
                            } else {
                                await bot.sendMessage(chatID, `Top \`${topValidatorsLength}\` validators by voting power at current height are:`, {parseMode: 'Markdown'});
                                for (let i = 0; i < topValidatorsLength; i++) {
                                    let message = validatorUtils.getValidatorMessage(activeValidators[i], totalBondedToken, 0, 0, blockHeight);
                                    await bot.sendMessage(chatID, `(${i + 1})\n\n` + message, {parseMode: 'Markdown'});
                                }
                            }
                        })
                        .catch(err => {
                            errors.Log(err, 'UPTIME');
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

bot.on('/commission', async (msg) => {
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
                    activeValidators.sort((a, b) => parseFloat(a.commission.commission_rates.rate) - parseFloat(b.commission.commission_rates.rate));
                    let lowestCommissionRate = parseFloat(activeValidators[0].commission.commission_rates.rate);
                    dataUtils.find(dataUtils.subscriberCollection, {})
                        .then(async (validatorSubscribersList, err) => {
                            if (err) {
                                errors.Log(err, 'COMMISSION');
                                return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                            }
                            await bot.sendMessage(chatID, `Validators by lowest commission rate \`${(lowestCommissionRate * 100.0).toFixed(2)}\` % at current height are:`, {parseMode: 'Markdown'});
                            for (let i = 0; i < activeValidators.length; i++) {
                                if (parseFloat(activeValidators[i].commission.commission_rates.rate) > lowestCommissionRate) {
                                    break;
                                }
                                let validatorSubscribe = validatorSubscribersList.filter(validatorSubscribers => (validatorSubscribers.operatorAddress === activeValidators[i].operator_address));
                                let counter = 0;
                                let initHeight = 0;
                                if (validatorSubscribe.length === 0) {
                                    subscriberUtils.initializeValidatorSubscriber(activeValidators[i].operator_address, blockHeight);
                                } else {
                                    counter = validatorSubscribe[0].counter;
                                    initHeight = validatorSubscribe[0].initHeight;
                                }
                                let message = validatorUtils.getValidatorMessage(activeValidators[i], totalBondedToken, counter, initHeight, blockHeight);
                                await bot.sendMessage(chatID, `(${i + 1})\n\n` + message, {parseMode: 'Markdown'});
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
        })
        .catch(err => {
            errors.Log(err, 'COMMISSION');
            botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
        })

});

bot.on('/topValidator', async (msg) => {
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
                                errors.Log(err, 'COMMISSION');
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
                            let counter = 0;
                            let initHeight = 0;
                            if (validatorSubscribe.length === 0) {
                                subscriberUtils.initializeValidatorSubscriber(slicedValidator[0].operator_address, blockHeight);
                            } else {
                                counter = validatorSubscribe[0].counter;
                                initHeight = validatorSubscribe[0].initHeight;
                            }
                            let message = validatorUtils.getValidatorMessage(slicedValidator[0], totalBondedToken, counter, initHeight, blockHeight);
                            await bot.sendMessage(chatID, `Validator with voting power in top \`${topValidatorsLength}\` and lowest commission rate:\n\n\n` + message, {parseMode: 'Markdown'});

                            activeValidators.sort((a, b) => parseFloat(a.commission.commission_rates.rate) - parseFloat(b.commission.commission_rates.rate));
                            let lowestCommissionRate = parseFloat(activeValidators[0].commission.commission_rates.rate);
                            for (let i = 0; i < activeValidators.length; i++) {
                                if (parseFloat(activeValidators[i].commission.commission_rates.rate) > lowestCommissionRate) {
                                    topValidatorsLength = i + 1;
                                    break;
                                }
                            }
                            slicedValidator = activeValidators.slice(0, topValidatorsLength);
                            slicedValidator.sort((a, b) => parseInt(b.tokens, 10) - parseInt(a.tokens, 10));
                            validatorSubscribe = validatorSubscribersList.filter(validatorSubscribers => (validatorSubscribers.operatorAddress === slicedValidator[0].operator_address));
                            counter = 0;
                            initHeight = 0;
                            if (validatorSubscribe.length === 0) {
                                subscriberUtils.initializeValidatorSubscriber(slicedValidator[0].operator_address, blockHeight);
                            } else {
                                counter = validatorSubscribe[0].counter;
                                initHeight = validatorSubscribe[0].initHeight;
                            }
                            message = validatorUtils.getValidatorMessage(slicedValidator[0], totalBondedToken, counter, initHeight, blockHeight);
                            botUtils.sendMessage(bot, chatID, `Validator among lowest commission rates \`${(lowestCommissionRate * 100.0).toFixed(2)}\`% and highest voting power:\n\n\n` + message, {parseMode: 'Markdown'});

                        })
                        .catch(err => {
                            errors.Log(err, 'COMMISSION');
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

bot.on('/uptime', async (msg) => {
    const chatID = msg.chat.id;
    let blockHeight = latestBlockHeight;
    httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/pool`)
        .then(data => JSON.parse(data))
        .then(json => {
            const totalBondedToken = parseInt(json.result.bonded_tokens, 10);      // with cosmos version upgrade, change here
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
                            validatorSubscribersList.sort((a, b) => (b.counter / b.counterHeight) - (a.counter / a.counterHeight));
                            let validatorList = [];
                            for (let i = 0; i < activeValidators.length; i++) {
                                let validator = validatorSubscribersList.filter(validatorSubscribe => (activeValidators[i].operator_address === validatorSubscribe.operatorAddress));
                                if (validator.length !== 0) {
                                    activeValidators[i].counter = validator[0].counter;
                                    activeValidators[i].initHeight = validator[0].initHeight;
                                    validatorList.push(activeValidators[i]);
                                }
                            }
                            if (validatorList.length !== 0) {
                                await bot.sendMessage(chatID, `Top active validators by uptime are:`, {parseMode: 'Markdown'});
                                for (let i = 0; i < validatorList.length; i++) {
                                    if (validatorList[i].counter / validatorList[i].initHeight > validatorList[0].counter / validatorList[0].initHeight) {
                                        break;
                                    }
                                    let message = validatorUtils.getValidatorMessage(validatorList[i], totalBondedToken, validatorList[i].counter, validatorList[i].initHeight, blockHeight);
                                    await bot.sendMessage(chatID, `(${i + 1})\n\n` + message, {parseMode: 'Markdown'});
                                }
                            } else {
                                return botUtils.sendMessage(bot, chatID, `Not enough data now. Please try after sometime.`, {parseMode: 'Markdown'})
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
        })
        .catch(err => {
            errors.Log(err, 'VOTING_POWER');
            botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
        })
});

bot.on('/validator', async (msg) => {
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

bot.on('/allValidators', msg => {
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

bot.on('/unsubValidator', (msg) => {
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

bot.on('/unsubAllValidators', msg => {
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

// Last Missed Blocks
bot.on('/lastMissedBlock', async (msg) => {
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
                            if (validatorSubscribers[0].lastMissedBlock) {
                                return botUtils.sendMessage(bot, chatID, `\`${moniker}\` last missed \`${validatorSubscribers[0].lastMissedBlock}\` block.`, {parseMode: 'Markdown'});
                            } else {
                                return botUtils.sendMessage(bot, chatID, `\`${moniker}\` has not missed any blocks since birth of this Bot.`, {parseMode: 'Markdown'});
                            }
                        }
                    })
            })
            .catch(err => botUtils.handleErrors(bot, chatID, err,'LAST_MISSED_BLOCK'));
    }
});

// Report
// bot.on('/report', async (msg) => {
//     return botUtils.waitForUserReply(bot, msg.chat.id, `Please provide a validator address.`, 'validatorAddressReport', {parseMode: 'Markdown'});
// });
//
// bot.on(['ask.validatorAddressReport'], async msg => {
//     const addr = msg.text;
//     const chatID = msg.chat.id;
//     let blockHeight = latestBlockHeight;
//     if (!validatorUtils.addressOperations.verifyValidatorOperatorAddress(addr)) {
//         return botUtils.sendMessage(bot, chatID, errors.INVALID_ADDRESS, {parseMode: 'Markdown'});
//     } else {
//         httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/validators/${addr}`)
//             .then(data => JSON.parse(data))
//             .then(json => {
//                 let validator = json.result;       // with cosmos version upgrade, change here
//                 dataUtils.find(dataUtils.subscriberCollection, {operatorAddress: addr})
//                     .then((validatorSubscribers, err) => {
//                         if (err) {
//                             errors.Log(err, 'LAST_MISSED_BLOCK');
//                             return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
//                         }
//                         if (validatorSubscribers.length === 0) {
//                             subscriberUtils.initializeValidatorSubscriber(validator.operatorAddress, blockHeight);
//                             return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
//                         } else {
//                             let moniker = validator.description.moniker;
//                             if (validatorSubscribers[0].lastMissedBlock) {
//                                 return botUtils.sendMessage(bot, chatID, `\`${moniker}\` last missed \`${validatorSubscribers[0].lastMissedBlock}\` block.`, {parseMode: 'Markdown'});
//                             } else {
//                                 return botUtils.sendMessage(bot, chatID, `\`${moniker}\` has not missed any blocks since birth of this Bot.`, {parseMode: 'Markdown'});
//                             }
//                         }
//                     })
//             })
//             .catch(err => {
//                 errors.Log(err, 'UPDATING_VALIDATORS');
//             });
//     }
// });

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
let initHeight = 0;

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
        latestBlockHeight = parseInt(json.result.data.value.block.header.height, 10);
        console.log(latestBlockHeight);
        if (initHeight === 0) {
            initHeight = latestBlockHeight;
        }
        checkAndSendMsgOnValidatorsAbsence(json);
        if (latestBlockHeight % 10000 === 0) {
            sendReports(latestBlockHeight)
                .catch(err => errors.Log(err, 'SENDING_REPORTS'));
        }
        validatorUtils.checkTxs(latestBlockHeight);
    }
}

function checkAndSendMsgOnValidatorsAbsence(json) {
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

                            if (!found) {
                                updateCounterAndSendMessage(validatorSubscribers, validatorDetails, blockHeight);
                            }
                        } else {
                            errors.Log('Incorrect database');
                        }
                    })
                    .catch(err => errors.Log(err, 'SENDING_MESSAGE'));
            });
        })
        .catch(err => errors.Log(err, 'SENDING_MESSAGE'))
}

function updateCounterAndSendMessage(validatorSubscribers, validatorDetails, blockHeight) {
    let query = {operatorAddress: validatorSubscribers.operatorAddress};
    if (!validatorDetails.jailed) {
        if ((validatorSubscribers.lastMissedBlock - validatorSubscribers.previousToLastMissedBlock) === 1) {
            let consecutiveCounter = validatorSubscribers.consecutiveCounter + 1;
            let alertLevel = subscriberUtils.getAlertLevel(consecutiveCounter);
            dataUtils.updateOne(dataUtils.subscriberCollection, query, {
                $set: {
                    counter: validatorSubscribers.counter + 1,
                    consecutiveCounter: consecutiveCounter,
                    alertLevel: alertLevel,
                    lastMissedBlock: blockHeight,
                    previousToLastMissedBlock: validatorSubscribers.lastMissedBlock
                }
            })
                .then((result) => {
                    let blocksLevel = subscriberUtils.getBlocksLevel(alertLevel);
                    if (consecutiveCounter % blocksLevel === 0) {
                        checkJailedStatusAndSendMessage(validatorSubscribers.operatorAddress);
                        sendMissedMsgToSubscribers(validatorDetails.description.moniker, validatorSubscribers.subscribers, consecutiveCounter)
                            .catch(err => errors.Log(err))
                    }
                })
                .catch(err => errors.Log(err, 'UPDATING_COUNTER_AND_SENDING_MESSAGE'));
        } else {
            dataUtils.updateOne(dataUtils.subscriberCollection, query, {
                $set: {
                    counter: validatorSubscribers.counter + 1,
                    consecutiveCounter: 0,
                    alertLevel: 1,
                    lastMissedBlock: blockHeight,
                    previousToLastMissedBlock: validatorSubscribers.lastMissedBlock
                }
            })
                .catch(err => errors.Log(err, 'UPDATING_COUNTER_AND_SENDING_MESSAGE'));
        }
    }
}

async function sendMissedMsgToSubscribers(moniker, subscribersList, consecutiveCounter) {
    subscribersList.forEach((subscriber) => {
        botUtils.sendMessage(bot, subscriber.chatID, `Alert: \`${moniker}\` has consecutively missed \`${consecutiveCounter}\` blocks.`, {
            parseMode: 'Markdown',
            notification: true
        });
    });
}

function checkJailedStatusAndSendMessage(operatorAddress) {
    httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/validators/${operatorAddress}`)
        .then(async data => {
            let json = JSON.parse(data);
            if (json.error) {
                errors.Log('Invalid Operator Address', 'UPDATE_COUNTER_QUERY_VALIDATOR')
            } else {
                let validator = json.result;       // with cosmos version upgrade, change here
                if (validator.jailed) {
                    dataUtils.updateOne(dataUtils.validatorCollection, query, {
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

async function sendReports(blockHeight) {
    dataUtils.find(dataUtils.subscriberCollection, {})
        .then((result, err) => {
            if (err) {
                errors.Log(err, 'SEND_REPORT');
                return;
            }
            let oldBlockHeight = blockHeight - 10000;
            result.forEach((validatorSubscribers) => {
                httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/validators/${validatorSubscribers.operatorAddress}?height=${oldBlockHeight}`)
                    .then(data => JSON.parse(data))
                    .then(async (json) => {
                        let oldValidatorDetails = json.result;       // with cosmos version upgrade, change here
                        httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/validators/${validatorSubscribers.operatorAddress}?height=${blockHeight}`)
                            .then(data => JSON.parse(data))
                            .then(json => {
                                let latestValidatorDetails = json.result;       // with cosmos version upgrade, change here
                                httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/pool?height=${oldBlockHeight}`)
                                    .then(data => JSON.parse(data))
                                    .then(json => {
                                        let oldTotalBondedTokens = json.result.bonded_tokens;       // with cosmos version upgrade, change here
                                        httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/pool?height=${blockHeight}`)
                                            .then(data => JSON.parse(data))
                                            .then(json => {
                                                let newTotalBondedTokens = json.result.bonded_tokens;       // with cosmos version upgrade, change here
                                                let message = validatorUtils.getValidatorReport(oldValidatorDetails, latestValidatorDetails, oldTotalBondedTokens, newTotalBondedTokens, blockHeight, validatorSubscribers.counter, validatorSubscribers.initHeight);
                                                validatorSubscribers.subscribers.forEach((subscriber) => {
                                                    botUtils.sendMessage(bot, subscriber.chatID, message, {
                                                        parseMode: 'Markdown',
                                                        notification: true
                                                    });
                                                });
                                            })
                                            .catch(err => {
                                                errors.Log(err, 'SENDING_REPORTS');
                                            });

                                    })
                                    .catch(err => {
                                        errors.Log(err, 'SENDING_REPORTS');
                                    });
                            })
                            .catch(err => {
                                errors.Log(err, 'SENDING_REPORTS');
                            });
                    })
                    .catch(err => {
                        errors.Log(err, 'SENDING_REPORTS');
                    });
            });
        })
        .catch(err => errors.Log(err, 'SENDING_REPORTS'))
}

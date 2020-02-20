const config = require('../config');
const HttpUtils = require('./httpRequest');
const validatorUtils = require('./validator');
const botUtils = require('./bot');
const dataUtils = require('./data');
const errors = require('./errors');
const subscriberUtils = require('./subscriber');
const httpUtils = new HttpUtils();

const queries = {
    sendLastBlock(bot, chatID) {
        httpUtils.httpGet(botUtils.nodeURL, config.node.abciPort, '/status')
            .then(data => JSON.parse(data))
            .then(json => queries.sendBlockInfo(bot, chatID, parseInt(json.result.sync_info.latest_block_height, 10)))
            .catch(e => botUtils.handleErrors(bot, chatID, e, 'SEND_LAST_BLOCK'));
    },
    sendValidatorsCount(bot, chatID) {
        httpUtils.httpGet(botUtils.nodeURL, config.node.abciPort, `/validators`)
            .then(data => {
                let json = JSON.parse(data);
                botUtils.sendMessage(bot, chatID, `There are \`${json.result.validators.length}\` validators in total at Block \`${json.result.block_height}\`.`, {parseMode: 'Markdown'})
            })
            .catch(e => botUtils.handleErrors(bot, chatID, e, 'SEND_VALIDATORS_COUNT'));
    },
    sendValidators(bot, chatID, blockHeight) {
        httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/pool?height=${blockHeight}`)
            .then(data => JSON.parse(data))
            .then(json => {
                const totalBondedToken = parseInt(json.result.bonded_tokens, 10);      // with cosmos version upgrade, change here
                httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/validators?height=${blockHeight}`)
                    .then(data => {
                        let json = JSON.parse(data);
                        let validatorsList = json.result;   // with cosmos version upgrade, change here
                        dataUtils.find(dataUtils.subscriberCollection, {})
                            .then(async (validatorSubscribers, err) => {
                                if (err) {
                                    errors.Log(err, 'SEND_VALIDATORS_LIST');
                                    return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                                }
                                let message = `\`${validatorsList.length}\` validators in total at height ${blockHeight}.`;
                                let i = 1;
                                for (let validator of validatorsList) {
                                    let filteredValidator = validatorSubscribers.filter(validatorSubscriber => (validatorSubscriber.operatorAddress === validator.operator_address));
                                    if (filteredValidator.length !== 0) {
                                        let upTime = subscriberUtils.calculateUptime(filteredValidator[0].blocksHistory);
                                        let valMsg = validatorUtils.getValidatorMessage(validator, totalBondedToken, upTime, filteredValidator[0].blocksHistory.length, blockHeight);
                                        message = message + `\n\n` + `(${i})\n\n` + valMsg;
                                    }
                                    if (i % 10 === 0) {
                                        await bot.sendMessage(chatID, message, {parseMode: 'Markdown'});
                                        message = ``;
                                    }
                                    i++;
                                }
                                if (message !== ``) {
                                    botUtils.sendMessage(bot, chatID, message, {parseMode: 'Markdown'});
                                }
                            })
                            .catch(e => botUtils.handleErrors(bot, chatID, e, 'SEND_VALIDATORS_INFO'));
                    })
                    .catch(e => botUtils.handleErrors(bot, chatID, e, 'SEND_VALIDATORS_LIST'));
            })
            .catch(e => botUtils.handleErrors(bot, chatID, e, 'SEND_VALIDATORS_LIST'));
    },
    sendValidatorInfo(bot, chatID, operatorAddress, blockHeight) {
        httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/pool`)
            .then(data => JSON.parse(data))
            .then(json => {
                const totalBondedToken = parseInt(json.result.bonded_tokens, 10);      // with cosmos version upgrade, change here
                httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/validators/${operatorAddress}?height=${blockHeight}`)
                    .then(data => {
                        let json = JSON.parse(data);
                        if (json.error) {
                            botUtils.sendMessage(bot, chatID, `Invalid operator address!`);
                        } else {
                            let validator = json.result;       // with cosmos version upgrade, change here
                            dataUtils.find(dataUtils.subscriberCollection, {operatorAddress: operatorAddress})
                                .then((validatorSubscriber, err) => {
                                    if (err) {
                                        errors.Log(err, 'VALIDATOR_INFO');
                                        return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                                    }
                                    if (validatorSubscriber.length === 0) {
                                        errors.Log('No validatorSubscriber found.', 'VALIDATOR_INFO');
                                        subscriberUtils.initializeValidatorSubscriber(operatorAddress, blockHeight);
                                        return botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
                                    } else {
                                        let upTime = subscriberUtils.calculateUptime(validatorSubscriber[0].blocksHistory);
                                        let message = validatorUtils.getValidatorMessage(validator, totalBondedToken, upTime, validatorSubscriber[0].blocksHistory.length, blockHeight);
                                        botUtils.sendMessage(bot, chatID, `At Block \`${blockHeight}\`\n\n`
                                            + message, {parseMode: 'Markdown'});
                                    }
                                })
                                .catch(e => botUtils.handleErrors(bot, chatID, e, 'SEND_VALIDATORS_INFO'));
                        }
                    })
                    .catch(e => botUtils.handleErrors(bot, chatID, e, 'SEND_VALIDATORS_INFO'));
            })
            .catch(e => botUtils.handleErrors(bot, chatID, e, 'SEND_VALIDATORS_INFO'));
    },
    sendBalance(bot, chatID, addr, blockHeight) {
        httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/auth/accounts/${addr}?height=${blockHeight}`)
            .then(data => {
                let json = JSON.parse(data);
                if (json.error) {
                    botUtils.sendMessage(bot, chatID, `Invalid address!`);
                } else {
                    let coins = '0';
                    if (!errors.isEmpty(json.result.value.coins)) {
                        coins = '';
                        json.result.value.coins.forEach((coin) => {                // with cosmos version upgrade, change here
                            let cn = Number(parseFloat(coin.amount).toFixed(2)).toLocaleString("en-US");
                            coins = coins + `${cn} ${coin.denom}, `
                        });
                    }
                    botUtils.sendMessage(bot, chatID, `At height \`${blockHeight}\` \n\n Coins: \`${coins}\`\n`,
                        {parseMode: 'Markdown'});
                }
            })
            .catch(e => botUtils.handleErrors(bot, chatID, e, 'SEND_BALANCE'));
    },
    sendDelRewards(bot, chatID, addr, blockHeight) {
        httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/distribution/delegators/${addr}/rewards?height=${blockHeight}`)
            .then(data => {
                let json = JSON.parse(data);
                if (json.error) {
                    botUtils.sendMessage(bot, chatID, `Invalid address!`);
                } else {
                    let total = '0';
                    if (!errors.isEmpty(json.result.total)) {
                        total = '';
                        json.result.total.forEach((rwd) => {        // with cosmos version upgrade, change here
                            let rwdAmt = Number(parseFloat(rwd.amount).toFixed(2)).toLocaleString("en-US");
                            total = total + `${rwdAmt} ${rwd.denom}, `;
                        });
                    }
                    botUtils.sendMessage(bot, chatID, `At height ${blockHeight} \n\n Total Rewards: \`${total}\``,      // with cosmos version upgrade, change here
                        {parseMode: 'Markdown'});
                }
            })
            .catch(e => botUtils.handleErrors(bot, chatID, e, 'SEND_DELEGATOR_REWARDS'));
    },
    sendValRewards(bot, chatID, addr, blockHeight) {
        httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/distribution/validators/${addr}`)
            .then(data => {
                let json = JSON.parse(data);
                if (json.error) {
                    botUtils.sendMessage(bot, chatID, `Invalid validator's operator address.`);
                } else {
                    let selfRewards = '0';
                    if (!errors.isEmpty(json.result.self_bond_rewards)) {
                        selfRewards = '';
                        json.result.self_bond_rewards.forEach((reward) => {        // with cosmos version upgrade, change here
                            let rewardAmt = Number(parseFloat(reward.amount).toFixed(2)).toLocaleString("en-US");
                            selfRewards = selfRewards + `${rewardAmt} ${reward.denom}, `;
                        });
                    }
                    let commission = '0';
                    if (!errors.isEmpty(json.result.val_commission)) {
                        commission = '';
                        json.result.val_commission.forEach((comm) => {        // with cosmos version upgrade, change here
                            let commAmt = Number(parseFloat(comm.amount).toFixed(2)).toLocaleString("en-US");
                            commission = commission + `${commAmt} ${comm.denom}, `;
                        });
                    }
                    botUtils.sendMessage(bot, chatID, `At height \`${blockHeight}\` \n\nSelf Bond Rewards: \`${selfRewards}\`\n\nCommission: \`${commission}\`\n`, {parseMode: 'Markdown'});
                }
            })
            .catch(e => botUtils.handleErrors(bot, chatID, e, 'SEND_VALIDATOR_REWARDS'));
    },
    sendTxByHash(bot, chatID, hash) {
        httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/txs/${hash}`)  // on abciPort query with /tx/0x${hash}
            .then(async data => {
                let json = JSON.parse(data);
                if (json.error) {
                    botUtils.sendMessage(bot, chatID, `Invalid Tx Hash.`);
                } else {
                    let raw_log = JSON.parse(json.raw_log);
                    let numMsgs = raw_log.length;
                    let statuses = [];
                    raw_log.forEach((log) => {
                        statuses.push(log.success);
                    });
                    let status = JSON.stringify(statuses);
                    botUtils.sendMessage(bot, chatID, `Height: \`${json.height}\`\n\n`
                        + `Number of Messages: \`${numMsgs}\`\n\n`
                        + `Messages Status: \`${status}\`\n\n`
                        + `Gas Wanted: \`${json.gas_wanted}\`\n\n`
                        + `Gas Used: \`${json.gas_used}\`\n\u200b\n`, {parseMode: 'Markdown'});
                }
            })
            .catch(e => botUtils.handleErrors(bot, chatID, e, 'SEND_TX_BY_HASH'));
    },
    sendTxByHeight(bot, chatID, height) {
        httpUtils.httpGet(botUtils.nodeURL, config.node.abciPort, `/tx_search?query="tx.height=${height}"`)
            .then(async data => {
                let json = JSON.parse(data);
                if (json.error) {
                    botUtils.sendMessage(bot, chatID, 'Invalid height.');
                } else {
                    if (json.result.txs[0]) {
                        await bot.sendMessage(chatID, `\`${json.result.txs.length}\` transactions at height \`${height}\`.`, {parseMode: 'Markdown'});
                        for (let i = 0; i < json.result.txs.length; i++) {
                            let raw_log = JSON.parse(json.result.txs[i].tx_result.log);
                            let numMsgs = raw_log.length;
                            let statuses = [];
                            raw_log.forEach((log) => {
                                statuses.push(log.success);
                            });
                            let status = JSON.stringify(statuses);
                            await bot.sendMessage(chatID, `(${i + 1})\n\n`
                                + `Tx Hash: \`${json.result.txs[i].hash}\`\n\n`
                                + `Number of Messages: \`${numMsgs}\`\n\n`
                                + `Messages Status: \`${status}\`\n\n`
                                + `Gas Wanted: \`${json.result.txs[i].tx_result.gasWanted}\`\n\n`
                                + `Gas Used: \`${json.result.txs[i].tx_result.gasUsed}\`\n\u200b\n`, {parseMode: 'Markdown'});
                        }
                    } else {
                        botUtils.sendMessage(bot, chatID, `No transactions at height \`${height}\`.`, {parseMode: 'Markdown'});
                    }
                }
            })
            .catch(e => botUtils.handleErrors(bot, chatID, e, 'SEND_TX_BY_HEIGHT'));
    },
    sendBlockInfo(bot, chatID, height) {
        httpUtils.httpGet(botUtils.nodeURL, config.node.abciPort, `/block?height=${height}`)
            .then(data => {
                let json = JSON.parse(data);
                if (json.error) {
                    botUtils.sendMessage(bot, chatID, 'Invalid height or height is greater than the current blockchain height.');
                } else {
                    dataUtils.find(dataUtils.validatorCollection, {hexAddress: json.result.block.header.proposer_address})
                        .then(async (result, err) => {
                            if (err) {
                                return botUtils.handleErrors(bot, chatID, err, 'SEND_BLOCK_INFO');
                            }
                            if (result.length === 0) {
                                validatorUtils.initializeValidatorDB();
                                return botUtils.handleErrors(bot, chatID, {message: 'Cannot find validator in DB'}, 'SEND_BLOCK_INFO');
                            }
                            let validator = result[0];
                            let moniker = validator.description.moniker;
                            await bot.sendMessage(chatID, `Height: \`${height}\`\n\n`
                                + `Proposer: \`${moniker}\`\n\u200b\n`, {parseMode: 'Markdown'});
                            queries.sendTxByHeight(bot, chatID, height);
                        })
                        .catch(e => botUtils.handleErrors(bot, chatID, e, 'SEND_BLOCK_INFO'));

                }
            })
            .catch(e => botUtils.handleErrors(bot, chatID, e, 'SEND_BLOCK_INFO'));
    }
};

module.exports = {queries};
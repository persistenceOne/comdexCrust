const config = require('../config');
const HttpUtils = require('./httpRequest');
const validatorUtils = require('./validator');
const botUtils = require('./bot');
const dataUtils = require('./data');
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
        sendValidators(bot, chatID) {
            httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/pool`)
                .then(data => JSON.parse(data))
                .then(json => {
                    const totalBondedToken = parseInt(json.result.bonded_tokens, 10);      // with cosmos version upgrade, change here
                    httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/validators`)
                        .then(async data => {
                            let json = JSON.parse(data);
                            let validatorsList = json.result;   // with cosmos version upgrade, change here
                            await bot.sendMessage(chatID, `\`${validatorsList.length}\` validators in total at current height.`, {parseMode: 'Markdown'});
                            let i = 1;
                            for (let validator of validatorsList) {
                                let message = validatorUtils.getValidatorMessage(validator, totalBondedToken);
                                await bot.sendMessage(chatID, `(${i})\n\n` + message, {parseMode: 'Markdown'});
                                i++;
                            }
                        })
                        .catch(e => botUtils.handleErrors(bot, chatID, e, 'SEND_VALIDATORS_LIST'));
                })
                .catch(e => botUtils.handleErrors(bot, chatID, e, 'SEND_VALIDATORS_LIST'));
        },
        sendValidatorInfo(bot, chatID, operatorAddress) {
            httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/pool`)
                .then(data => JSON.parse(data))
                .then(json => {
                    const totalBondedToken = parseInt(json.result.bonded_tokens, 10);      // with cosmos version upgrade, change here
                    httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/validators/${operatorAddress}`)
                        .then(data => {
                            let json = JSON.parse(data);
                            if (json.error) {
                                botUtils.sendMessage(bot, chatID, `Invalid operator address!`);
                            } else {
                                let validator = json.result;       // with cosmos version upgrade, change here
                                let message = validatorUtils.getValidatorMessage(validator, totalBondedToken);
                                botUtils.sendMessage(bot, chatID, message, {parseMode: 'Markdown'});
                            }
                        })
                        .catch(e => botUtils.handleErrors(bot, chatID, e, 'SEND_VALIDATORS_LIST'));
                })
                .catch(e => botUtils.handleErrors(bot, chatID, e, 'SEND_VALIDATORS_LIST'));
        },
        sendBalance(bot, chatID, addr) {
            httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/auth/accounts/${addr}`)
                .then(data => {
                    let json = JSON.parse(data);
                    if (json.error) {
                        botUtils.sendMessage(bot, chatID, `Invalid address!`);
                    } else {
                        let coins = '';
                        json.result.value.coins.forEach((coin) => {                // with cosmos version upgrade, change here
                            coins = coins + `${coin.amount} ${coin.denom}, `
                        });
                        botUtils.sendMessage(bot, chatID, `Coins: \`${coins}\`\n`,
                            {parseMode: 'Markdown'});
                    }
                })
                .catch(e => botUtils.handleErrors(bot, chatID, e, 'SEND_BALANCE'));
        },
        sendDelRewards(bot, chatID, addr) {
            httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/distribution/delegators/${addr}/rewards`)
                .then(data => {
                    let json = JSON.parse(data);
                    if (json.error) {
                        botUtils.sendMessage(bot, chatID, `Invalid address!`);
                    } else {
                        botUtils.sendMessage(bot, chatID, `Total Rewards: \`${json.result.total[0].amount} ${json.result.total[0].denom}\``,      // with cosmos version upgrade, change here
                            {parseMode: 'Markdown'});
                    }
                })
                .catch(e => botUtils.handleErrors(bot, chatID, e, 'SEND_DELEGATOR_REWARDS'));
        },
        sendValRewards(bot, chatID, addr) {
            httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/distribution/validators/${addr}`)
                .then(data => {
                    let json = JSON.parse(data);
                    if (json.error) {
                        botUtils.sendMessage(bot, chatID, `Invalid validator's operator address.`);
                    } else {
                        let selfRewards = '';
                        json.result.self_bond_rewards.forEach((reward) => {        // with cosmos version upgrade, change here
                            selfRewards = selfRewards + `${reward.amount} ${reward.denom}, `;
                        });
                        let commission = '';
                        json.result.val_commission.forEach((comm) => {             // with cosmos version upgrade, change here
                            commission = commission + `${comm.amount} ${comm.denom}, `;
                        });
                        botUtils.sendMessage(bot, chatID, `Self Bond Rewards: \`${selfRewards}\`\n\nCommission: \`${commission}\`\n`, {parseMode: 'Markdown'});
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
                        botUtils.sendMessage(bot, chatID, `Height: \`${json.height}\`\n\n`
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
                                await bot.sendMessage(chatID, `(${i + 1})\n\n`
                                    + `Tx Hash: \`${json.result.txs[i].hash}\`\n\n`
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
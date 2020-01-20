const config = require('../config');
const errors = require('./errors');
const HttpUtil = require('./httpRequest');
const validatorUtils = require('./validator');
const dataUtil = require('./data');
const httpUtil = new HttpUtil();

const queries = {
        sendLastBlock(bot, chatID) {
            httpUtil.httpGet(config.node.url, config.node.abciPort, '/status')
                .then(data => JSON.parse(data))
                .then(json => bot.sendMessage(chatID, `\`${json.result.sync_info.latest_block_height}\``, {parseMode: 'Markdown'}))
                .catch(e => handleErrors(bot, chatID, e, 'SEND_LAST_BLOCK'));
        },

        sendNodeInfo(bot, chatID) {
            httpUtil.httpGet(config.node.url, config.node.abciPort, '/status')
                .then(data => {
                    let json = JSON.parse(data);
                    let syncedUP = '';
                    if (!json.result.sync_info.catching_up) {
                        syncedUP = 'Synced Up';
                    } else {
                        syncedUP = 'Not Synced Up';
                    }
                    bot.sendMessage(chatID, `Chain: \`${json.result.node_info.network}\`\n\n`
                        + `ID: \`${json.result.node_info.id}\`\n\n`
                        + `Moniker: \`${json.result.node_info.moniker}\`\n\n`
                        + `Address: \`${json.result.validator_info.address}\`\n\n`
                        + `Voting Power: \`${json.result.validator_info.voting_power}\`\n\n`
                        + `\`${syncedUP}\`\n\u200b\n`
                        , {parseMode: 'Markdown'})
                })
                .catch(e => handleErrors(bot, chatID, e, 'SEND_NODE_INFO'));
        },
        sendPeersCount(bot, chatID) {
            httpUtil.httpGet(config.node.url, config.node.abciPort, '/net_info')
                .then(data => {
                    let json = JSON.parse(data);
                    bot.sendMessage(chatID, `The node has \`${json.result.n_peers}\` peers in total.`, {parseMode: 'Markdown'})
                })
                .catch(e => handleErrors(bot, chatID, e, 'SEND_PEERS_COUNT'));
        },
        sendPeersList(bot, chatID) {
            httpUtil.httpGet(config.node.url, config.node.abciPort, '/net_info')
                .then(async data => {
                    let json = JSON.parse(data);
                    await bot.sendMessage(chatID, `${json.result.n_peers} peers in total.`, {parseMode: 'Markdown'});
                    let peers = json.result.peers;
                    let i = 1;
                    for (let peer of peers) {
                        await bot.sendMessage(chatID, `(${i})\n\nNode ID: \`${peer.node_info.id}\`\n\n`
                            + `Moniker: \`${peer.node_info.moniker}\`\n\u200b\n`,
                            {parseMode: 'Markdown'});
                        i += 1;
                    }
                })
                .catch(e => handleErrors(bot, chatID, e, 'SEND_PEERS_LIST'));
        },
        sendConsensusState(bot, chatID) {
            httpUtil.httpGet(config.node.url, config.node.abciPort, '/dump_consensus_state')
                .then(data => {
                    let json = JSON.parse(data);
                    let roundState = json.result.round_state;
                    bot.sendMessage(chatID, `Round State Height: \`${roundState.height}\`\n\n`
                        + `Round: \`${roundState.round}\`\n\n`
                        + `Step: \`${roundState.step}\`\n\n`
                        + `Proposer: \`${roundState.validators.proposer.address}\`\n\u200b\n`, {parseMode: 'Markdown'})
                })
                .catch(e => handleErrors(bot, chatID, e, 'SEND_CONSENSUS_STATE'));
        },
        sendConsensusParams(bot, chatID) {
            httpUtil.httpGet(config.node.url, config.node.abciPort, `/consensus_params`)
                .then(data => {
                    let json = JSON.parse(data);
                    let blockSize = json.result.consensus_params.block.max_bytes;
                    bot.sendMessage(chatID, `Height: \`${json.result.block_height}\`\n\n`
                        + `Block Max Size: \`${blockSize}\`\n\n`
                        + `Evidence Max Age: \`${json.result.consensus_params.evidence.max_age}\`\n\n`
                        + `Validator Pubkey Type(s): \`${json.result.consensus_params.validator.pub_key_types[0]}\`\n\u200b\n`, {parseMode: 'Markdown'})
                })
                .catch(e => handleErrors(bot, chatID, e, 'SEND_CONSENSUS_PARAMS'));
        },
        sendValidatorsCount(bot, chatID) {
            httpUtil.httpGet(config.node.url, config.node.abciPort, `/validators`)
                .then(data => {
                    let json = JSON.parse(data);
                    bot.sendMessage(chatID, `There are \`${json.result.validators.length}\` validators in total at Block \`${json.result.block_height}\`.`, {parseMode: 'Markdown'})
                })
                .catch(e => handleErrors(bot, chatID, e, 'SEND_VALIDATORS_COUNT'));
        },
        sendValidators(bot, chatID) {
            httpUtil.httpGet(config.node.url, config.node.lcdPort, `/staking/validators`)
                .then(async data => {
                    let json = JSON.parse(data);
                    let validatorsList = json.result;
                    await bot.sendMessage(chatID, `\`${validatorsList.length}\` validators in total at current height \`${json.height}\`.`, {parseMode: 'Markdown'});
                    let i = 1;
                    for (let validator of validatorsList) {
                        let selfDelegationAddress = validatorUtils.getDelegatorAddrFromOperatorAddr(validator.operator_address);
                        let message = `(${i})\n\nOperator Address: \`${validator.operator_address}\`\n\n`
                            + `Self Delegation Address: \`${selfDelegationAddress}\`\n\n`
                            + `Moniker: \`${validator.description.moniker}\`\n\n`
                            + `Details: \`${validator.description.details}\`\n\n`
                            + `Website: ${validator.description.website}\n\u200b\n`;
                        await bot.sendMessage(chatID, message, {parseMode: 'Markdown'});
                        i++;
                    }
                })
                .catch(e => {
                    handleErrors(bot, chatID, e, 'SEND_VALIDATORS_LIST');
                });
        },
        sendValidatorInfo(bot, chatID, operatorAddress) {
            httpUtil.httpGet(config.node.url, config.node.lcdPort, `/staking/validators/${operatorAddress}`)
                .then(data => {
                    let json = JSON.parse(data);
                    if (json.error) {
                        bot.sendMessage(chatID, `Invalid operator address!`);
                    } else {
                        let validator = json.result;
                        let selfDelegationAddress = validatorUtils.getDelegatorAddrFromOperatorAddr(validator.operator_address);
                        bot.sendMessage(chatID, `Operator Address: \`${validator.operator_address}\`\n\n`
                            + `Self Delegation Address: \`${selfDelegationAddress}\`\n\n`
                            + `Moniker: \`${validator.description.moniker}\`\n\n`
                            + `Details: \`${validator.description.details}\`\n\n`
                            + `Jailed: \`${validator.jailed}\`\n\n`
                            + `Website: ${validator.description.website}\n\u200b\n`,
                            {parseMode: 'Markdown'});
                    }
                })
                .catch(e => handleErrors(bot, chatID, e, 'SEND_VALIDATORS_LIST'));
        },
        sendBalance(bot, chatID, addr) {
            httpUtil.httpGet(config.node.url, config.node.lcdPort, `/auth/accounts/${addr}`)
                .then(data => {
                    let json = JSON.parse(data);
                    if (json.error) {
                        bot.sendMessage(chatID, `Invalid address!`);
                    } else {
                        let coins = '';
                        json.result.value.coins.forEach((coin) => {
                            coins = coins + `${coin.amount} ${coin.denom}, `
                        });
                        bot.sendMessage(chatID, `Height: \`${json.height}\`\n\n`
                            + `Coins: \`${coins}\`\n`,
                            {parseMode: 'Markdown'});
                    }
                })
                .catch(e => handleErrors(bot, chatID, e, 'SEND_BALANCE'));
        },
        sendDelRewards(bot, chatID, addr) {
            httpUtil.httpGet(config.node.url, config.node.lcdPort, `/distribution/delegators/${addr}/rewards`)
                .then(data => {
                    let json = JSON.parse(data);
                    if (json.error) {
                        bot.sendMessage(chatID, `Invalid address!`);
                    } else {
                        bot.sendMessage(chatID, `Total Rewards: \`${json.result.total[0].amount} ${json.result.total[0].denom}\``,
                            {parseMode: 'Markdown'});
                    }
                })
                .catch(e => handleErrors(bot, chatID, e, 'SEND_DELEGATOR_REWARDS'));
        },
        sendValRewards(bot, chatID, addr) {
            httpUtil.httpGet(config.node.url, config.node.lcdPort, `/distribution/validators/${addr}`)
                .then(data => {
                    let json = JSON.parse(data);
                    if (json.error) {
                        bot.sendMessage(chatID, `Invalid validator's operator address.`);
                    } else {
                        let selfRewards = '';
                        json.result.self_bond_rewards.forEach((reward) => {
                            selfRewards = selfRewards + `${reward.amount} ${reward.denom}, `;
                        });
                        let commission = '';
                        json.result.val_commission.forEach((comm) => {
                            commission = commission + `${comm.amount} ${comm.denom}, `;
                        });
                        bot.sendMessage(chatID, `Self Bond Rewards: \`${selfRewards}\`\n\nCommission: \`${commission}\`\n`, {parseMode: 'Markdown'});
                    }
                })
                .catch(e => handleErrors(bot, chatID, e, 'SEND_VALIDATOR_REWARDS'));
        },
        sendValSigningInfo(bot, chatID, addr) {
            httpUtil.httpGet(config.node.url, config.node.lcdPort, `/slashing/validators/${addr}/signing_info`)
                .then(data => {
                    let json = JSON.parse(data);
                    if (json.error) {
                        bot.sendMessage(chatID, `Invalid validator's consensus public address.`);
                    } else {
                        let jailed = !(json.result.jailed_until === '1970-01-01T00:00:00Z');
                        bot.sendMessage(chatID, `Start height: \`${json.result.start_height}\`\n\n`
                            + `Jailed: \`${jailed}\`\n\n`
                            + `Tombstoned: \`${json.result.tombstoned}\`\n\n`
                            + `Missed Blocks Counter: \`${json.result.missed_blocks_counter}\`\n\u200b\n`, {parseMode: 'Markdown'});
                    }
                })
                .catch(e => handleErrors(bot, chatID, e, 'SEND_VALIDATOR_SIGNING_INFO'));
        },
        sendSlashingParams(bot, chatID) {
            httpUtil.httpGet(config.node.url, config.node.lcdPort, `/slashing/parameters`)
                .then(data => {
                    let json = JSON.parse(data);
                    bot.sendMessage(chatID, `Max evidence age: \`${parseInt(json.result.max_evidence_age, 10) / (1e9 * 24 * 3600)} days\`\n\n`
                        + `Downtime Jail Duration: \`${parseInt(json.result.downtime_jail_duration, 10) / 1e9} seconds\`\n\n`
                        + `Double Sign Slashing Fraction: \`${parseFloat(json.result.slash_fraction_double_sign, 10) * 100}\`\n\n`
                        + `Downtime Slashing Fraction: \`${parseFloat(json.result.slash_fraction_downtime, 10) * 100}\`\n\n`
                        + `Signed Blocks Window: \`${parseInt(json.result.signed_blocks_window, 10)}\`\n\u200b\n`, {parseMode: 'Markdown'});
                })
                .catch(e => handleErrors(bot, chatID, e, 'SEND_SLASHING_PARAMS'));
        },
        sendMintingParams(bot, chatID) {
            httpUtil.httpGet(config.node.url, config.node.lcdPort, `/minting/parameters`)
                .then(data => {
                    let json = JSON.parse(data);
                    bot.sendMessage(chatID, `Mint Denom: \`${json.result.mint_denom}\`\n\n`
                        + `Inflation Rate Change: \`${parseFloat(json.result.inflation_rate_change, 10) * 100.0}\`\n\n`
                        + `Inflation Max Rate: \`${parseFloat(json.result.inflation_max, 10) * 100.0}\`\n\n`
                        + `Inflation Min Rate: \`${parseFloat(json.result.inflation_min, 10) * 100.0}\`\n\n`
                        + `Goal Bonded: \`${parseFloat(json.result.goal_bonded, 10) * 100.0}\`\n\n`
                        + `Blocks per year: \`${json.result.blocks_per_year}\`\n\u200b\n`, {parseMode: 'Markdown'});
                })
                .catch(e => handleErrors(bot, chatID, e, 'SEND_MINTING_PARAMS'));
        },
        sendMintingInflation(bot, chatID) {
            httpUtil.httpGet(config.node.url, config.node.lcdPort, `/minting/inflation`)
                .then(data => {
                    let json = JSON.parse(data);
                    bot.sendMessage(chatID, `Minting inflation \`${parseFloat(json.result, 10) * 100}\` at height \`${json.height}\`.`,
                        {parseMode: 'Markdown'});
                })
                .catch(e => handleErrors(bot, chatID, e, 'SEND_MINTING_INFLATION'));
        },
        sendStakingPool(bot, chatID) {
            httpUtil.httpGet(config.node.url, config.node.lcdPort, `/staking/pool`)
                .then(data => {
                    let json = JSON.parse(data);
                    bot.sendMessage(chatID, `Height: \`${json.height}\`\n\n`
                        + `Bonded Tokens: \`${parseInt(json.result.bonded_tokens, 10) / 1000000} ${config.token}\`\n\n`
                        + `Not Bonded Tokens: \`${parseInt(json.result.not_bonded_tokens, 10) / 1000000} ${config.token}\`\n\u200b\n`,
                        {parseMode: 'Markdown'});
                })
                .catch(e => handleErrors(bot, chatID, e, 'SEND_STAKING_POOL'));
        },
        sendStakingParams(bot, chatID) {
            httpUtil.httpGet(config.node.url, config.node.lcdPort, `/staking/parameters`)
                .then(data => {
                    let json = JSON.parse(data);
                    bot.sendMessage(chatID, `Height: \`${json.height}\`\n\n`
                        + `Unbonding Time: \`${json.result.unbonding_time / (1e9 * 24 * 3600)} days\`\n\n`
                        + `Max Validators: \`${json.result.max_validators}\`\n\n`
                        + `Max Entries: \`${json.result.max_entries}\`\n\n`
                        + `Bond Denom: \`${json.result.bond_denom}\`\n\u200b\n`, {parseMode: 'Markdown'});
                })
                .catch(e => handleErrors(bot, chatID, e, 'SEND_STAKING_PARAMS'));
        },
        sendTxByHash(bot, chatID, hash) {
            httpUtil.httpGet(config.node.url, config.node.lcdPort, `/txs/${hash}`)  // on abciPort query with /tx/0x${hash}
                .then(async data => {
                    let json = JSON.parse(data);
                    if (json.error) {
                        bot.sendMessage(chatID, `Invalid Tx Hash.`);
                    } else {
                        bot.sendMessage(chatID, `Height: \`${json.height}\`\n\n`
                            + `Gas Wanted: \`${json.gas_wanted}\`\n\n`
                            + `Gas Used: \`${json.gas_used}\`\n\u200b\n`, {parseMode: 'Markdown'});
                    }
                })
                .catch(e => handleErrors(bot, chatID, e, 'SEND_TX_BY_HASH'));
        },
        sendTxByHeight(bot, chatID, height) {
            httpUtil.httpGet(config.node.url, config.node.abciPort, `/tx_search?query="tx.height=${height}"&per_page=30`)
                .then(async data => {
                    let json = JSON.parse(data);
                    if (json.error) {
                        bot.sendMessage(chatID, 'Invalid height.');
                    } else {
                        if (json.result.txs[0]) {
                            await bot.sendMessage(chatID, `\`${json.result.txs.length}\` transactions at height \`${height}\`.`);
                            for (let i = 0; i < json.result.txs.length; i++) {
                                await bot.sendMessage(chatID, `(${i + 1})\n\n`
                                    + `Tx Hash: \`${json.result.txs[i].hash}\`\n\n`
                                    + `Gas Wanted: \`${json.result.txs[i].tx_result.gasWanted}\`\n\n`
                                    + `Gas USed: \`${json.result.txs[i].tx_result.gasUsed}\`\n\u200b\n`, {parseMode: 'Markdown'});
                            }
                        } else {
                            bot.sendMessage(chatID, `No transactions at height \`${height}\`.`, {parseMode: 'Markdown'});
                        }
                    }
                })
                .catch(e => handleErrors(bot, chatID, e, 'SEND_TX_BY_HEIGHT'));
        },
        sendBlockInfo(bot, chatID, height) {
            httpUtil.httpGet(config.node.url, config.node.abciPort, `/block?height=${height}`)
                .then(async data => {
                    let json = JSON.parse(data);
                    if (json.error) {
                        bot.sendMessage(chatID, 'Invalid height.');
                    } else {
                        await bot.sendMessage(chatID, `Block Hash: \`${json.result.block_meta.block_id.hash}\`\n\n`
                            + `Proposer: \`${json.result.block.header.proposer_address}\`\n\n`
                            + `Evidence: \`${json.result.block.evidence.evidence}\`\n\u200b\n`, {parseMode: 'Markdown'});
                        queries.sendTxByHeight(bot, chatID, height);
                    }
                })
                .catch(e => handleErrors(bot, chatID, e, 'SEND_BLOCK_INFO'));
        }
    }
;

function handleErrors(bot, chatID, err, method = '') {
    console.log(JSON.stringify(err));
    errors.Log(err, method);
    if (err.statusCode === 400) {
        bot.sendMessage(chatID, errors.INVALID_REQUEST, {parseMode: 'Markdown'});
    } else {
        bot.sendMessage(chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
    }
}

function updateValidatorDetails(operatorAddress) {
    httpUtil.httpGet(config.node.url, config.node.lcdPort, `/staking/validators/${operatorAddress}`)
        .then(data => JSON.parse(data))
        .then(json => {
            let validator = json.result;
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


module.exports = {queries, updateValidatorDetails, newValidatorObject};
const config = require('../config');

const buttons = {
    chain: {
        command: '/chain'
    },
    hide: {
        command: '/hide_keyboard'
    },
    home: {
        command: '/home'
    },
    back: {
        command: '/back'
    },
    validatorQuery: {
        command: '/validator_queries'
    },
    chainQuery: {
        command: '/chain_queries'
    },
    analyticsQuery: {
        command: '/analytics_queries'
    },
    votingPower: {
        command: '/top_validators_wrt_voting_power'
    },
    commission: {
        command: '/top_validators_wrt_commission'
    },
    uptime: {
        command: '/top_validators_wrt_uptime'
    },
    topValidator: {
        command: '/top_validators_wrt_commission_voting_power'
    },
    subscribe: {
        command: '/subscribe'
    },
    validator: {
        command: '/one_validator'
    },
    allValidators: {
        command: '/all_validator'
    },
    unsubValidator: {
        command: '/unsub_one_validator'
    },
    unsubAllValidators: {
        command: '/unsub_all_validator'
    },
    lastBlock: {
        command: '/last_block'
    },
    validatorsCount: {
        command: '/validators_count'
    },
    validatorsList: {
        command: '/validators_list'
    },
    lastMissedBlock: {
        command: '/last_missed_blocks'
    },
    validatorInfo: {
        command: '/validator_info'
    },
    validatorReport: {
        command: '/validator_report'
    },
    blockLookup: {
        command: '/block_lookup'
    },
    txLookup: {
        command: '/tx_lookup'
    },
    txByHeight: {
        command: '/tx_by_height'
    },
    accountBalance: {
        command: '/account_balance'
    },
    delegatorRewards: {
        command: '/delegator_rewards'
    },
    validatorRewards: {
        command: '/validator_rewards'
    },
};

module.exports = buttons;
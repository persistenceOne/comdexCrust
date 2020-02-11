const config = require('../config');

const buttons = {
    chain: {
        label: config.chainAppName,
        command: '/chain'
    },
    hide: {
        label: '⌨ Hide Keyboard',
        command: '/hide'
    },
    home: {
        label: '⌂ Home',
        command: '/home'
    },
    back: {
        label: '↵ Back',
        command: '/back'
    },
    validatorQuery: {
        label: 'Validator Queries',
        command: '/validator_queries'
    },
    chainQuery: {
        label: 'Chain Queries',
        command: '/chain_queries'
    },
    analyticsQuery: {
        label: 'Analytics',
        command: '/analytics_queries'
    },
    votingPower: {
        label: 'Voting Power',
        command: '/voting_power'
    },
    commission: {
        label: '' +
            'Commission',
        command: '/commission'
    },
    alerts: {
        label: 'Alerts',
        command: '/alerts'
    },
    subscribe: {
        label: 'Subscribe',
        command: '/subscribe'
    },
    unsubscribe: {
        label: 'Unsubscribe',
        command: '/unsubscribe'
    },
    lastBlock: {
        label: 'Last Block',
        command: '/last_block'
    },
    validatorsCount: {
        label: 'Validators Count',
        command: '/validators_count'
    },
    validatorsList: {
        label: 'Validators List',
        command: '/validators_list'
    },
    validatorInfo: {
        label: 'Validator Info',
        command: '/validator_info'
    },
    blockLookup: {
        label: 'Block Look Up',
        command: '/block_lookup'
    },
    txLookup: {
        label: 'Tx Look Up',
        command: '/tx_lookup'
    },
    txByHeight: {
        label: 'Tx By Height',
        command: '/tx_by_height'
    },
    accountBalance: {
        label: 'Account Balance',
        command: '/account_balance'
    },
    delegatorRewards: {
        label: 'Delegator Rewards',
        command: '/delegator_rewards'
    },
    validatorRewards: {
        label: 'Validator Rewards',
        command: '/validator_rewards'
    },
};

module.exports = buttons;
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
    nodeQuery: {
        label: 'Node Queries',
        command: '/node_queries'
    },
    chainQuery: {
        label: 'Chain Queries',
        command: '/chain_queries'
    },
    lcdQuery: {
        label: 'Light Client Queries',
        command: '/lcd_queries'
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
    nodeStatus: {
        label: 'Status',
        command: '/node_status'
    },
    lastBlock: {
        label: 'Last Block',
        command: '/last_block'
    },
    peersCount: {
        label: 'Peers Count',
        command: '/peers_count'
    },
    peersList: {
        label: 'Peers List',
        command: '/peers_list'
    },
    consensusState: {
        label: 'Consensus State',
        command: '/consensus_state'
    },
    consensusParams: {
        label: 'Consensus Parameters',
        command: '/consensus_params'
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
    stakingPool: {
        label: 'Staking Pool',
        command: '/staking_pool'
    },
    stakingParams: {
        label: 'Staking Params',
        command: '/staking_params'
    },
    mintingInflation: {
        label: 'Minting Inflation',
        command: '/minting_inflation'
    },
    slashingParams: {
        label: 'Slashing Params',
        command: '/slashing_params'
    },
    mintingParams: {
        label: 'Minting Params',
        command: '/minting_params'
    },
    validatorSigning: {
        label: 'Validator Signing',
        command: '/validator_signing'
    },
};

module.exports = buttons;
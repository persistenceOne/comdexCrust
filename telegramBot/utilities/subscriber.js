const HttpUtils = require('./httpRequest');
const errors = require('./errors');
const httpUtils = new HttpUtils();
const config = require('../config.json');
const dataUtils = require('./data');
const botUtils = require('./bot');

function initializeSubscriberDB() {
    httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/validators`, 60)
        .then(data => JSON.parse(data))
        .then(json => {
            let validators = json.result;       // with cosmos version upgrade, change here
            let latestBlockHeight = json.height;    // with cosmos version upgrade, change here
            validators.forEach((validator) => {
                initializeValidatorSubscriber(validator.operator_address, latestBlockHeight);
            });
        })
        .catch(err => errors.exitProcess(err, 'INITIALIZING_SUBSCRIBER_DB'));
}

function initializeValidatorSubscriber(operatorAddress, latestBlockHeight) {
    dataUtils.find(dataUtils.subscriberCollection, {operatorAddress: operatorAddress})
        .then((subscribersResult, err) => {
                if (err) {
                    errors.exitProcess('DB_INITIALIZING_VALIDATOR_SUBSCRIBER');
                }
                if (subscribersResult.length === 0) {
                    let validatorSubscriber = newValidatorSubscribers(operatorAddress, latestBlockHeight, []);
                    dataUtils.insertOne(dataUtils.subscriberCollection, validatorSubscriber)
                        .catch(err => errors.exitProcess(err, 'DB_INITIALIZING_VALIDATOR_SUBSCRIBER'));
                } else {
                    dataUtils.updateOne(dataUtils.subscriberCollection, {operatorAddress: operatorAddress}, {
                        $set: {
                            consecutiveCounter: 0,
                            alertLevel: 1,
                            blocksHistory: [],
                            lastMissedBlock: 0
                        }
                    })
                        .catch(err => errors.Log(err, 'RESET_VALIDATOR_SUBSCRIBER'));
                }
            }
        )
        .catch(err => errors.Log(err, 'INITIALIZING_VALIDATOR_SUBSCRIBER'));
}

function newValidatorSubscribers(operatorAddress, latestBlockHeight, subscribers) {
    return {
        operatorAddress: operatorAddress,
        consecutiveCounter: 0,
        alertLevel: 1,
        blocksHistory: [],
        subscribers: subscribers,
        lastMissedBlock: 0
    };
}

function getBlocksLevel(slashingWindow, alertLevel) {
    return Math.ceil(slashingWindow ** (alertLevel/config.alertLevels));
}

function getAlertLevel(slashingWindow, consecutiveCounter) {
    if (consecutiveCounter === 0 || consecutiveCounter === 1){
        return 1;
    } else {
        return Math.ceil(config.alertLevels * Math.log(consecutiveCounter) / Math.log(slashingWindow));
    }
}

function calculateUptime(blocksHistory) {
    let missed = 0;
    for (let i = 0; i < blocksHistory.length; i++) {
        if (!blocksHistory[i].found) {
            missed = missed + 1;
        }
    }
    return ((1.00 - (missed / blocksHistory.length)) * 100.00).toFixed(2);
}

module.exports = {
    initializeSubscriberDB,
    initializeValidatorSubscriber,
    newValidatorSubscribers,
    getBlocksLevel,
    getAlertLevel,
    calculateUptime
};
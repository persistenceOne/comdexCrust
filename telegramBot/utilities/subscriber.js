const HttpUtils = require('./httpRequest');
const errors = require('./errors');
const httpUtils = new HttpUtils();
const config = require('../config.json');
const dataUtils = require('./data');
const botUtils = require('./bot');

function initializeSubscriberDB() {
    httpUtils.httpGet(botUtils.nodeURL, config.node.lcdPort, `/staking/validators`)
        .then(data => JSON.parse(data))
        .then(json => {
            let validators = json.result;       // with cosmos version upgrade, change here
            validators.forEach((validator) => {
                initializeValidatorSubscriber(validator.operator_address, json.height);      // with cosmos version upgrade, change here
            });
        })
        .catch(err => {
            errors.exitProcess(err, 'INITIALIZING_SUBSCRIBER_DB');
        });
}

function initializeValidatorSubscriber(operatorAddress, latestBlockHeight) {
    dataUtils.find(dataUtils.subscriberCollection, {operatorAddress: operatorAddress})
        .then((result, err) => {
            if (err) {
                errors.exitProcess('DB_INITIALIZING_VALIDATOR_SUBSCRIBER');
            }
            if (result.length === 0) {
                let validatorSubscriber = newValidatorSubscribers(operatorAddress, latestBlockHeight, []);
                dataUtils.insertOne(dataUtils.subscriberCollection, validatorSubscriber)
                    .catch(err => errors.exitProcess(err, 'DB_INITIALIZING_VALIDATOR_SUBSCRIBER'));
            } else {
                resetSubscriberCounters(operatorAddress, latestBlockHeight);
            }
        })
        .catch(err => errors.Log(err, 'INITIALIZING_VALIDATOR_SUBSCRIBER'));
}

function resetSubscriberCounters(operatorAddress, latestBlockHeight) {
    dataUtils.find(dataUtils.subscriberCollection, {operatorAddress: operatorAddress})
        .then((result, err) => {
            if (err) {
                return errors.Log('DB_INITIALIZING_VALIDATOR_SUBSCRIBER', 'RESET_VALIDATOR_SUBSCRIBER');
            }
            dataUtils.updateOne(dataUtils.subscriberCollection, {operatorAddress: operatorAddress}, {
                $set: {
                    counter: 0,
                    consecutiveCounter: 0,
                    alertLevel: 1,
                    initHeight: latestBlockHeight,
                    counterHeight: latestBlockHeight,
                }
            })
        })
        .catch(err => errors.Log(err, 'RESET_VALIDATOR_SUBSCRIBER'));
}

function newValidatorSubscribers(operatorAddress, latestBlockHeight, subscribers) {
    return  {
        operatorAddress: operatorAddress,
        counter: 0,
        consecutiveCounter: 0,
        alertLevel: 1,
        initHeight: latestBlockHeight,
        counterHeight: latestBlockHeight,
        subscribers: subscribers,

    };
}

function getBlocksLevel(alertLevel) {
    switch (alertLevel) {
        case 1:
            return 5;
        case 2:
            return 10;
        case 3:
            return 15;
        case 4:
            return 20;
        case 5:
            return 20;
        default:
            return 20;
    }
}

function getAlertLevel(consecutiveCounter) {
    switch (true) {
        case consecutiveCounter <= 5:
            return 1;
        case 5 <= consecutiveCounter && consecutiveCounter <= 10:
            return 2;
        case 10 <= consecutiveCounter && consecutiveCounter <= 15:
            return 3;
        case 15 <= consecutiveCounter && consecutiveCounter <= 20:
            return 4;
        case consecutiveCounter >= 20:
            return 5;
        default:
            return 1;
    }
}

module.exports = {initializeSubscriberDB, initializeValidatorSubscriber, newValidatorSubscribers, getBlocksLevel, getAlertLevel};
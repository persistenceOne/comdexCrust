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
                initializeValidatorSubscriber(validator, json.height);      // with cosmos version upgrade, change here
            });
        })
        .catch(err => {
            errors.exitProcess(err, 'INITIALIZING_SUBSCRIBER_DB');
        });
}

function initializeValidatorSubscriber(validator, latestBlockHeight) {
    let operatorAddress = validator.operator_address;
    dataUtils.find(dataUtils.subscriberCollection, {operatorAddress: operatorAddress})
        .then((result, err) => {
            if (err) {
                errors.exitProcess('DB_INITIALIZING_VALIDATOR_SUBSCRIBER');
            }
            if (result.length === 0) {
                dataUtils.insertOne(dataUtils.subscriberCollection, {
                    operatorAddress: operatorAddress,
                    counter: 0,
                    height: latestBlockHeight,
                    subscribers: []
                })
                    .catch(err => errors.exitProcess(err, 'DB_INITIALIZING_VALIDATOR_SUBSCRIBER'));
            }
        })
        .catch(err => errors.Log(err, 'INITIALIZING_VALIDATOR_SUBSCRIBER'));
}

module.exports = {initializeSubscriberDB, initializeValidatorSubscriber};
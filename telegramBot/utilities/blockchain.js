const HttpUtils = require('./httpRequest');
const errors = require('./errors');
const httpUtils = new HttpUtils();
const config = require('../config.json');
const dataUtils = require('./data');
const botUtils = require('./bot');

function initializeBlockchainDB() {
    dataUtils.deleteMany(dataUtils.blockchainCollection, {})
        .catch(e => errors.exitProcess(e, 'INITIALIZING_BLOCKHCHAIN_DB'));
}

async function updateBlock(blockHeight) {
    httpUtils.httpGet(botUtils.nodeURL, config.node.abciPort, `/block?height=${blockHeight}`)
        .then(data => {
            let json = JSON.parse(data);
            let proposerHexAddress = json.result.block.header.proposer_address;
            let numTxs = 0;
            if (!errors.isEmpty(json.result.block.data.txs)) {
                numTxs = json.result.block.data.txs.length;
            }
            if (json.error) {
                errors.log('Invalid height or height is greater than the current blockchain height.', 'UPDATE_BLOCK_DB');
            } else {
                dataUtils.findSorted(dataUtils.blockchainCollection, {}, {height: 1}, {projection:{ _id: 0 }})
                    .then(async (blockchainDetails, err) => {
                        if (err) {
                            errors.Log(err, 'COMMISSION');
                        } else {
                            let newBlock = {height: blockHeight, proposer: proposerHexAddress, numTxs: numTxs};
                            if (blockchainDetails.length > config.blockchainHistoryLimit) {
                                dataUtils.deleteOne(dataUtils.blockchainCollection, {height: blockchainDetails[0].height})
                                    .then((result, err) => {
                                        dataUtils.insertOne(dataUtils.blockchainCollection, newBlock)
                                            .catch(err => errors.Log(err, 'SUBSCRIBE_INSERT'));
                                    })
                                    .catch(err => errors.Log(err, 'SUBSCRIBE_INSERT'));
                            } else {
                                dataUtils.insertOne(dataUtils.blockchainCollection, newBlock)
                                    .catch(err => errors.Log(err, 'SUBSCRIBE_INSERT'));
                            }
                        }
                    })
            }
        })
        .catch(e => errors.log(e, 'UPDATE_BLOCK_DB'));
}

module.exports = {
    initializeBlockchainDB,
    updateBlock
};

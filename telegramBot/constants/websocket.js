const config = require('../config');

const queryNewBlockString = `tm.event='NewBlock'`;
const queryNewTxString = `tm.event='Tx'`;

const constants = {
    subscribeNewBlockMsg: {
        "jsonrpc": "2.0",
        "method": "subscribe",
        "id": "0",
        "params": {
            "query": `${queryNewBlockString}`,
        },
    },
    subscribeTxMsg: {
        "jsonrpc": "2.0",
        "method": "subscribe",
        "id": "0",
        "params": {
            "query": `${queryNewTxString}`,
        },
    },
    unsubscribeAllMsg: {
        "jsonrpc": "2.0",
        "method": "unsubscribe_all",
        "id": "0",
        "params": {},
    },
    url: `ws://${config.node.url}:${config.node.abciPort}/websocket`,
    backupURL: `ws://${config.node.backupURL}:${config.node.abciPort}/websocket`,
};

module.exports = constants;
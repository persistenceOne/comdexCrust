const config = require('../config');
const bech32 = require('bech32');
const hash = require('tendermint/lib/hash.js');
const tmhash = hash.tmhash;

const addressOperations = {
    pubkeyToBech32(pubkey, prefix) {
        let pubkeyAminoPrefix = Buffer.from('1624DE6420', 'hex');
        let buffer = Buffer.alloc(37);
        pubkeyAminoPrefix.copy(buffer, 0);
        Buffer.from(pubkey, 'base64').copy(buffer, pubkeyAminoPrefix.length);
        return bech32.encode(prefix, bech32.toWords(buffer));
    },
    bech32ToPubkey(pubkey) {
        let pubkeyAminoPrefix = Buffer.from('1624DE6420', 'hex')
        let buffer = Buffer.from(bech32.fromWords(bech32.decode(pubkey).words));
        return buffer.slice(pubkeyAminoPrefix.length).toString('base64');
    },
    getHexAddress(pubkeyValue) {
        let bytes = Buffer.from(pubkeyValue, 'base64');
        return tmhash(bytes).slice(0, 20).toString('hex').toUpperCase();
    },
    toPubKey(address) {
        return bech32.decode(config.prefix, address);
    },
    createAddress(publicKey) {
        const message = CryptoJS.enc.Hex.parse(publicKey.toString(`hex`));
        const hash = CryptoJS.RIPEMD160(CryptoJS.SHA256(message)).toString();
        const addr = Buffer.from(hash, `hex`);
        return bech32ify(addr, config.prefix);
    },
    getDelegatorAddrFromOperatorAddr(operatorAddr) {
        let address = bech32.decode(operatorAddr);
        return bech32.encode(config.prefix, address.words);
    },
    verifyValidatorOperatorAddress(validatorOperatorAddr) {
        const validatorOperatorAddrRegex = new RegExp('\^' + config.prefix + 'valoper' + '\[a-z0-9]{39}$');
        return validatorOperatorAddrRegex.test(validatorOperatorAddr);
    },
};

function bech32ify(address, prefix) {
    const words = bech32.toWords(address);
    return bech32.encode(prefix, words);
}

module.exports = addressOperations;
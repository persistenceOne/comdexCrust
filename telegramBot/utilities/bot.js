const errors = require('./errors');
const config = require('../config.json');

//When using labels, bot directly goes into subEvent when `bot.sendMessage(chatID, message, {ask: subEvent})`
// is called without waiting for the user to reply. This method makes it await somehow.
async function waitForUserReply(bot, chatID, message, subEvent, options) {
    if (!options) {
        await bot.eventList;
        bot.sendMessage(chatID, message, {ask: subEvent})
        .catch((err) => errors.Log(err))
    } else {
        await bot.eventList;
        bot.sendMessage(chatID, message, {ask: subEvent}, options)
        .catch((err) => errors.Log(err))
    }
}

function sendMessage(bot, chatID, message, options) {
    if (!options) {
        bot.sendMessage(chatID, message)
        .catch((err) => errors.Log(err.description, 'SENDING_MESSAGE'))
    } else {
        bot.sendMessage(chatID, message, options)
        .catch((err) => errors.Log(err.description, 'SENDING_MESSAGE'))
    }
}

function handleErrors(bot, chatID, err, method = '') {
    console.log(JSON.stringify(err));
    errors.Log(err, method);
    if (err.statusCode === 400 || err.statusCode === 404) {
        botUtils.sendMessage(bot, chatID, errors.INVALID_REQUEST, {parseMode: 'Markdown'});
    } else {
        botUtils.sendMessage(bot, chatID, errors.INTERNAL_ERROR, {parseMode: 'Markdown'});
    }
}

let nodeURL = config.node.url;

module.exports = {waitForUserReply, sendMessage, handleErrors, nodeURL};
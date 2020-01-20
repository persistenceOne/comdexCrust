const INVALID_ADDRESS = 'Address is invalid!';
const INTERNAL_ERROR = 'Internal Server Error!';
const ERROR_RESPONSE = 'ERROR_RESPONSE';
const INVALID_REQUEST = 'Not a proper query';

const isEmpty = (object) => {
    return !object || Object.keys(object).length === 0;
};

function exitProcess(err = '', code = 1) {
    let log;
    if (isEmpty(err.message)) {
        log = err;
    } else {
        log = err.message;
    }
    console.log('Server Stopped:    ' + log);
    process.exit(code);
}

function Log(err, method = '') {
    let message = '';
    if (method !== '') {
        message = message + ' METHOD: ' + method + ';';
    }
    if (!isEmpty(err.statusCode)) {
        message = message + ' Status Code: ' + err.statusCode + ';';
    }
    if (isEmpty(err.message)) {
        message = message + err;
    } else {
        message = message + ' MESSAGE: ' + err.message;
    }
    console.log(new Date() + ' ERROR' + message);
}

module.exports = {isEmpty, Log, exitProcess, INVALID_ADDRESS, INTERNAL_ERROR, ERROR_RESPONSE, INVALID_REQUEST};
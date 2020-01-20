const https = require('https');
const http = require('http');
const errors = require('./errors');
const config = require('../config');
const REQ_TIMEOUT = config.requestTimeout;

class HttpUtils {

    httpsGet(url, port, path, req_timeout = REQ_TIMEOUT) {
        let options = {
            hostname: url,
            port: port,
            path: path,
            method: 'GET'
        };
        // create new promise
        return new Promise((resolve, reject) => {
            // request timeout
            setTimeout(() => {
                reject(new Error(`${new Date()} - ${req_timeout}s timeout exceeded`));
            }, req_timeout * 1000);

            https.request(options, (res) => {
                // response status check
                if (res.statusCode < 200 || res.statusCode > 299) {
                    reject({statusCode: res.statusCode, message: errors.ERROR_RESPONSE});
                }
                // var to store res body
                let res_body = '';
                // get body (by chunks)
                res.on('data', (data) => {
                    res_body += data;
                });
                // resolve promise(return body as text)
                res.on('end', () => {
                    resolve(res_body);
                });
            }).on('error', function (err) {
                reject({statusCode: err.code, message: err.message})
            }).end();
        });
    }

    httpGet(url, port, path, req_timeout = REQ_TIMEOUT) {
        let options = {
            hostname: url,
            port: port,
            path: path,
            method: 'GET'
        };
        // create new promise
        return new Promise((resolve, reject) => {
            // request timeout
            setTimeout(() => {
                reject(new Error(`${new Date()} - ${req_timeout}s timeout exceeded`));
            }, req_timeout * 1000);

            http.request(options, (res) => {
                // response status check
                if (res.statusCode < 200 || res.statusCode > 299) {
                    console.log('asdasdas');
                    reject({statusCode: res.statusCode, message: errors.ERROR_RESPONSE});
                }
                // var to store res body
                let res_body = '';
                // get body (by chunks)
                res.on('data', (data) => {
                    res_body += data;
                });
                // resolve promise(return body as text)
                res.on('end', () => {
                    resolve(res_body);
                });
            }).on('error', function (err) {
                reject({statusCode: err.code, message: err.message})
            }).end();
        });
    }
}

module.exports = HttpUtils;
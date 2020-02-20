const MongoClient = require('mongodb').MongoClient;
const errors = require('./errors');
const config = require('../config');

const mongoURL = config.dbURL + config.dbName;

const subscriberCollection = 'subscribers';
const subscriberAllCollection = 'subscribeAll';
const validatorCollection = 'validators';
const blockchainCollection = 'blockchain';

let dbo;    //Not to export.

function SetupDB(callback) {
    console.log('Intialzing DB...');
    MongoClient.connect(mongoURL, {useUnifiedTopology: true})
        .then((db, err) => {
            if (err) throw  err;
            dbo = db.db(config.dbName);
            console.log('DB Initialization complete.');
            callback();
        })
        .catch(err => {
            errors.exitProcess(err);
        });
}

function find(collection, query, options = {}) {
    return dbo.collection(collection).find(query, options).toArray();
}

function findSorted(collection, query, sortingOption, options = {}) {
    return dbo.collection(collection).find(query, options).sort(sortingOption).toArray();
}

function insertOne(collection, data) {
    return dbo.collection(collection).insertOne(data);
}

function insertMany(collection, data) {
    return dbo.collection(collection).insertMany(data);
}

function updateOne(collection, query, data) {
    return dbo.collection(collection).updateOne(query, data);
}

function deleteOne(collection, query) {
    return dbo.collection(collection).deleteOne(query);
}

function deleteMany(collection, query) {
    return dbo.collection(collection).deleteMany(query);
}

function upsertOne(collection, query, data) {
    return dbo.collection(collection).updateOne(query, data, {upsert: true});
}

module.exports = {
    subscriberCollection,
    validatorCollection,
    blockchainCollection,
    subscriberAllCollection,
    SetupDB,
    find,
    findSorted,
    insertOne,
    insertMany,
    updateOne,
    deleteOne,
    deleteMany,
    upsertOne
};

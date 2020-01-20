const MongoClient = require('mongodb').MongoClient;
const errors = require('./errors');
const config = require('../config');

const mongoURL = config.dbURL + config.dbName;

const subscriberCollection = 'subscribers';
const validatorCollection = 'validators';

let dbo;    //Not to export.

function InitializeDB() {
    console.log('Intialzing DB...');
    MongoClient.connect(mongoURL, {useUnifiedTopology: true})
        .then((db, err) => {
            if (err) throw  err;
            dbo = db.db(config.dbName);
            console.log('DB Initialization complete.')
        })
        .catch(err => {
            errors.exitProcess(err);
        });
}

function find(collection, query) {
    return dbo.collection(collection).find(query).toArray();
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

function upsertOne(collection, query, data) {
    return dbo.collection(collection).updateOne(query, data, {upsert: true});
}

module.exports = {
    subscriberCollection,
    validatorCollection,
    InitializeDB,
    find,
    insertOne,
    insertMany,
    updateOne,
    upsertOne
};

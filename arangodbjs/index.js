'use strict';

const arangojs = require('arangojs');
const config = require('../config');
const _ = require('lodash');
const handler = require('./handler');


let arangoDbClient = null;
let acquireHostListTimer = null;

async function waitDatabaseReady(db) {
    async function checkConnection(checkDb) {
        try {
            await checkDb.version();
            return true;
        } catch (err) {
            if (_.isNumber(err.code) && err.code < 500) {
                return true;
            }
            return false;
        }
    }

    if (!(await checkConnection(db))) {
        return new Promise((resolve, reject) => {
            const timer = setInterval(() => {
                console.debug('Check Database Connection...');
                checkConnection(db).then((result) => {
                    if (result) {
                        console.debug('Database Ready to served.');
                        clearInterval(timer);
                        resolve(true);
                    }
                });
            }, 1000);
        });
    }
    return true;
}

async function login(db) {
    if (!_.isString(config.ARANGODB_USER) || !_.isString(config.ARANGODB_PASSWORD) || config.ARANGODB_USER.length === 0) {
        return true;
    }

    try {
        const token = await db.login(config.ARANGODB_USER, config.ARANGODB_PASSWORD);
        db.useBearerAuth(token);
    } catch (err) {
        if (err instanceof arangojs.ArangoError) {
            // ARANGO_NO_AUTH enable
            if (err.code === 404) {
                return true;
            }
        }
        throw err;
    }
    return true;
}

module.exports.init = async function () {
    const arangodbConnectionTimeout = setTimeout(() => {
        throw new Error('arangodbConnectionTimeout, check arangodb plz');
    }, 15000);
    try {
        const arangoNodes = config.ARANGODB_URI.split(',');
        console.log(config.ARANGODB_DATABASE)
        arangoDbClient = new arangojs.Database({
            url: arangoNodes,
            agentOptions: {
                maxSockets: 30,
                keepAlive: true,
                keepAliveMsecs: 10000,
            },
            loadBalancingStrategy: 'ROUND_ROBIN',
        });

        console.debug('Wait for ArangoDB to be ready...');
        await waitDatabaseReady(arangoDbClient);

        await login(arangoDbClient);

        if (arangoNodes.length > 1) {
            // Update all coordinators at startup
            await arangoDbClient.acquireHostList();
            // Updates the URL list by requesting a list of all coordinators in the cluster
            // and adding any endpoints not initially specified in the url configuration.
            // Update it per hour
            acquireHostListTimer = setInterval(() => {
                arangoDbClient.acquireHostList();
            }, 1000 * 60 * 60);
        }

        const availDatabases = await arangoDbClient.listDatabases();

        if (!availDatabases.includes(config.ARANGODB_DATABASE)) {
            const user = {
                username: config.ARANGODB_USER,
                passwd: config.ARANGODB_PASSWORD,
            };
            const info = await arangoDbClient.createDatabase(config.ARANGODB_DATABASE, [user]);
            console.debug(`Create arango database: ${config.ARANGODB_DATABASE}`, info);
        }
        // 7.7.0
        // arangoDbClient.useDatabase(config.ARANGODB_DATABASE);
        // 8.0.0
        // arangoDbClient = arangoDbClient.database(config.ARANGODB_DATABASE);

        // // inject AQL helper
        // arangoDbClient.aql = aql;

        clearTimeout(arangodbConnectionTimeout);
        console.debug('ArangoDB initialized.');
    } catch (err) {
        if (err) {
            console.debug('ArangoError error:', err.stack);
        }
        // throw err;
    }
};

module.exports.shutdown = async function () {
    if (acquireHostListTimer) {
        clearInterval(acquireHostListTimer);
        acquireHostListTimer = null;
    }
    arangoDbClient.close();
    return;
};

module.exports.get = function () {
    return {
        db: arangoDbClient.database(config.ARANGODB_DATABASE),

        // inject AQL helper
        aql:  arangojs.aql,
    }
};

// 清除資料表環境
module.exports.clearCollection = async function (collectionName) {
    await handler.clearCollection(collectionName);
};

// 清除資料表
module.exports.clear = async function () {
    await handler.clearPlayer();
    await handler.clearCollection(this.NAMES.CLUBS);
    await handler.clearCollection(this.NAMES.CLUBMEMBERS);
    await handler.clearCollection(this.NAMES.WHITELIST);
    await handler.clearCollection(this.NAMES.FRIENDSHIP);

    // 清除過 theme 需要 從開 gs gm
    // await handler.clearCollection(this.NAMES.THEMES)
    // await handler.inputThemes()
};

module.exports.NAMES = {
    PLAYERS: 'Players',
    ACCOUNTS: 'Accounts',
    CLUBS: 'Clubs',
    CLUBMEMBERS: 'ClubMembers',
    THEMES: 'Themes',
    WHITELIST: 'Whitelist',
    FRIENDSHIP: 'FriendShip',
    COMPETITIONCYCLES: 'CompetitionCycles',
};

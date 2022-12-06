const ArangoContext = require('./index');

module.exports.getSignatureDoc = async (collection, key) => {
    const Arangodb = ArangoContext.get();
    const aql = Arangodb.aql;
    const queryString = aql`
        FOR d IN @@collName
            FILTER d._key == @key
        RETURN d
    `;

    queryString.bindVars['@collName'] = collection;
    queryString.bindVars.key = key;

    return Arangodb.db.query(queryString).then((cur) => cur.next());
}
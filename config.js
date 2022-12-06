const path = require('path');
const Joi = require('joi');

require('dotenv').config({
    path: path.join(__dirname, '..', '.env'),
});

const envSchema = Joi.object()
    .keys({
        ARANGODB_URI: Joi.string().default('http://127.0.0.1:8529'),
        ARANGODB_USER: Joi.string().default('root'),
        ARANGODB_PASSWORD: Joi.string().empty(''),
        ARANGODB_DATABASE: Joi.string().default('Database'),
        ARANGODB_DATA_CLEAR: Joi.boolean().default(false),
    })
    .unknown()
    .required();

const { error, value: config } = envSchema.validate(process.env);
if (error) {
    throw new Error(`Invalid enviroment variable: ${error.message}`);
}

module.exports = config;

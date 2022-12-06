const arango = require('./arangodbjs');
const arangoHandler = require('./arangodbjs/handler');
var ellipticcurve = require("starkbank-ecdsa");
var Ecdsa = ellipticcurve.Ecdsa;
var Signature = ellipticcurve.Signature;
var PublicKey = ellipticcurve.PublicKey;
var File = ellipticcurve.utils.File;

let publicKeyPem = File.read("public.pem");
let publicKey = PublicKey.fromPem(publicKeyPem);
// let signatureDer = File.read("signatureDer.txt", "binary");
let message = File.read("msg.txt");

async function init() {
    await arango.init();
    return true;
}

async function close() {
    await arango.shutdown();
    return true;
}

Promise.resolve()
    .then(() => {
        return init();
    })
    .then(async () => {
        let result = await arangoHandler.getSignatureDoc("erictest", "1670291833718");
        let decodeSig = Buffer.from(result.signature, 'base64');
        let signature = Signature.fromDer(decodeSig.toString('binary'));

        console.log(Ecdsa.verify(message, signature, publicKey));
        await close();
    }).
    catch((err) => {
        close();
        throw err;
    })
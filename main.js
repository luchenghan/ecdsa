var ellipticcurve = require("starkbank-ecdsa");
var Ecdsa = ellipticcurve.Ecdsa;
var Signature = ellipticcurve.Signature;
var PublicKey = ellipticcurve.PublicKey;
var File = ellipticcurve.utils.File;

let publicKeyPem = File.read("public.key");
let signatureDer = File.read("signatureDer.txt", "binary");
let message = File.read("msg.txt");

let publicKey = PublicKey.fromPem(publicKeyPem);
let signature = Signature.fromDer(signatureDer);

console.log(Ecdsa.verify(message, signature, publicKey));
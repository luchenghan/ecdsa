package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"reflect"
)

func encode(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) (string, string) {
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: x509Encoded})

	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return string(pemEncoded), string(pemEncodedPub)
}

func decode(pemEncoded string, pemEncodedPub string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	if err != nil {
		panic(err)
	}
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return privateKey, publicKey
}

// 產公私鑰
func generatePrivateAndPublicKey() {
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	pkStr, pubStr := encode(pk, &pk.PublicKey)
	fmt.Println(pkStr)
	fmt.Println(pubStr)

	f, _ := os.Create("private.key")
	_, writeStringErr := f.WriteString(pkStr)
	if writeStringErr != nil {
		panic(err)
	}

	f2, _ := os.Create("public.key")
	_, writeStringErr = f2.WriteString(pubStr)
	if writeStringErr != nil {
		panic(err)
	}

	priv2, pub2 := decode(pkStr, pubStr)

	if !reflect.DeepEqual(pk, priv2) {
		fmt.Println("Private keys do not match.")
	}
	if !reflect.DeepEqual(&pk.PublicKey, pub2) {
		fmt.Println("Public keys do not match.")
	}
}

func main() {
	generatePrivateAndPublicKey()
	privateBytes, err := os.ReadFile("private.key")
	if err != nil {
		panic(err)
	}

	publicBytes, err := os.ReadFile("public.key")
	if err != nil {
		panic(err)
	}

	privateKey, publicKey := decode(string(privateBytes), string(publicBytes))

	// 數位簽章
	msg := "hello world"
	hash := sha256.Sum256([]byte(msg))
	fmt.Printf("%x\n", hash)
	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	if err != nil {
		panic(err)
	}
	fmt.Printf("signature: %v\n", sig)

	f, _ := os.Create("signatureDer.txt")
	_, writeErr := f.Write(sig)
	if writeErr != nil {
		panic(writeErr)
	}

	// 驗證
	msg2 := "hello world"
	hash2 := sha256.Sum256([]byte(msg2))
	fmt.Println("signature verified:", ecdsa.VerifyASN1(publicKey, hash2[:], sig))
}

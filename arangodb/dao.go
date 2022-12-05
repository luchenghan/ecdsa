package arangodb

type DigtalSignature struct {
	Key       string `json:"_key"`
	Signature []byte `json:"signature"`
}

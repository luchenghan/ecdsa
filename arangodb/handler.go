package arangodb

import (
	"log"
	"os"
	"reflect"

	"github.com/arangodb/go-driver"
	jsoniter "github.com/json-iterator/go"
)

func (h Handler) EnsureCollection(name string, options *driver.CreateCollectionOptions) driver.Collection {
	c, err := h.db.Collection(h.ctx, name)
	if driver.IsNotFoundGeneral(err) {
		c, err = h.db.CreateCollection(h.ctx, name, options)
		if err != nil {
			log.Fatalf("Create collection error: %v", err)
		}
	} else if err != nil {
		log.Fatalf("Open collection error: %v", err)
	}
	return c
}

func (h Handler) TxDigtalSignature(ds DigtalSignature) {
	Error := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	bytes, err := jsoniter.Marshal(ds)
	if err != nil {
		Error.Printf("Marshal error: %v", err)
		return
	}

	const colName = "erictest"
	action := `function (Params) {
		const db = require('@arangodb').db;
		const ds = JSON.parse(Params[0]);
		const erictestCol = db._collection(Params[1]);
		erictestCol.save(ds);
		return 1;
	}`

	txOptions := &driver.TransactionOptions{
		MaxTransactionSize:   100000,
		WriteCollections:     []string{colName},
		ReadCollections:      []string{colName},
		ExclusiveCollections: []string{colName},
		Params:               []interface{}{string(bytes), colName},
	}

	result, err := h.db.Transaction(h.ctx, action, txOptions)
	if !reflect.DeepEqual(1.0, result.(float64)) {
		Error.Printf("Transaction expect: %v, got: %v, error: %v", 1, result, err)
		return
	}
	if err != nil {
		Error.Printf("Transaction error: %v", err)
		return
	}
}

func (h Handler) GetDigtalSignature(name string, key string) []byte {
	c, err := h.db.Collection(h.ctx, name)
	if err != nil {
		log.Fatalf("Collection error: %v", err)
	}
	s := struct {
		Signature []byte `json:"signature"`
	}{}
	_, err = c.ReadDocument(h.ctx, key, &s)
	if err != nil {
		log.Fatalf("Read Document error: %v", err)
	}
	return s.Signature
}

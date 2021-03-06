package service

import (
	"errors"
	"fmt"

	bolt "github.com/coreos/bbolt"
	"github.com/dedis/student_18_omniledger/omniledger/collection"
	"github.com/dedis/student_18_omniledger/omniledger/darc"
	"gopkg.in/dedis/onet.v2/network"
)

func init() {
	network.RegisterMessages(&Transaction{},
		&darc.Signature{})
}

type collectionDB struct {
	db         *bolt.DB
	bucketName []byte
	coll       collection.Collection
}

// newCollectionDB initialises a structure and reads all key/value pairs to store
// it in the collection.
func newCollectionDB(db *bolt.DB, name []byte) *collectionDB {
	c := &collectionDB{
		db:         db,
		bucketName: name,
		coll:       collection.New(collection.Data{}, collection.Data{}),
	}
	c.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(name)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	c.loadAll()
	// TODO: Check the merkle tree root.
	return c
}

func (c *collectionDB) loadAll() {
	c.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(c.bucketName))
		cur := b.Cursor()

		for k, v := cur.First(); k != nil; k, v = cur.Next() {
			sig := b.Get(append(k, []byte("sig")...))
			c.coll.Add(k, v, sig)
		}

		return nil
	})
}

func (c *collectionDB) Store(key, value, sig []byte) error {
	c.coll.Add(key, value, sig)
	err := c.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(c.bucketName))
		if err := bucket.Put(key, value); err != nil {
			return err
		}
		keysig := make([]byte, len(key)+3)
		copy(keysig, key)
		keysig = append(keysig, []byte("sig")...)
		if err := bucket.Put(keysig, sig); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (c *collectionDB) GetValue(key []byte) (value, sig []byte, err error) {
	// err = c.db.Update(func(tx *bolt.Tx) error {
	// 	bucket := tx.Bucket([]byte(c.bucketName))
	// 	value = bucket.Get(key)
	// 	sig = bucket.Get(append(key, []byte("sig")...))
	// 	return nil
	// })
	// return

	// TODO: make it so that the collection only stores the hashes, not the values
	proof, err := c.coll.Get(key).Record()
	if err != nil {
		return
	}
	hashes, err := proof.Values()
	if err != nil {
		return
	}
	if len(hashes) == 0 {
		err = errors.New("nothing stored under that key")
		return
	}
	value, ok := hashes[0].([]byte)
	if !ok {
		err = errors.New("the value is not of type []byte")
		return
	}
	sig, ok = hashes[1].([]byte)
	if !ok {
		err = errors.New("the signature is not of type []byte")
		return
	}
	return
}

// RootHash returns the hash of the root node in the merkle tree.
func (c *collectionDB) RootHash() []byte {
	return c.coll.GetRoot()
}

// Transaction is the struct specifying the modifications to the skipchain.
// Key is the key chosen by the user, Kind is the kind of value to store
// (e.g. a drac...). The key used in the conode's collection will be
// Kind ':' Key, in order to maintain key uniqueness across different kinds
// of values.
// For a Transaction to be valid, there must exist a path from the master-darc
// in the genesis block to the SubjectPK in Signature.
type Transaction struct {
	Key   []byte
	Kind  []byte
	Value []byte
	// type Signature struct {
	//     Signature []byte
	//     Signer SubjectPK
	// }
	// The signature is performed on the concatenation of the []bytes
	Signature darc.Signature
}

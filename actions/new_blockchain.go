package actions

import (
	"github.com/boltdb/bolt"
	"github.com/tclchiam/block_n_go/blockchain"
	"fmt"
)

type NewBlockchainAction struct {
	dbFileName   string
	blocksBucket string
}

func (action *NewBlockchainAction) Execute() (bool, error) {
	db, err := bolt.Open(action.dbFileName, 0600, nil)
	if err != nil {
		return false, fmt.Errorf("opening db: %s", err)
	}

	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		blocksBucketExists := tx.Bucket([]byte(action.blocksBucket)) != nil

		if blocksBucketExists == false {
			genesisBlock := blockchain.NewGenesisBlock()

			bucket, err := tx.CreateBucket([]byte(action.blocksBucket))
			if err != nil {
				return fmt.Errorf("creating block bucket: %s", err)
			}

			blockData, err := genesisBlock.Serialize()
			if err != nil {
				return fmt.Errorf("serializing block: %s", err)
			}

			err = bucket.Put(genesisBlock.Hash, blockData)
			if err != nil {
				return fmt.Errorf("writing block: %s", err)
			}

			err = bucket.Put([]byte("l"), genesisBlock.Hash)
			if err != nil {
				return fmt.Errorf("writing last hash: %s", err)
			}
		}

		return nil
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

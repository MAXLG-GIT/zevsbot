package repo

import (
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"os"

	"time"
	"zevsbot/internal/domain/domain_interfaces"
)

type BadgerRepo struct {
	db *badger.DB
}

func InitBadgerRepo(repoDir string) (domain_interfaces.RepoService, error) {
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		return nil, err
	}
	opts := badger.DefaultOptions(repoDir)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &BadgerRepo{db: db}, nil
}

func (br *BadgerRepo) Save(chatId, token string, expire time.Time) error {
	txn := br.db.NewTransaction(true)

	err := txn.SetEntry(&badger.Entry{
		Key:       []byte(chatId),
		Value:     []byte(token),
		ExpiresAt: uint64(expire.Unix()),
	})
	if err != nil {
		return err
	}
	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (br *BadgerRepo) Get(id string) (string, error) {
	var value string
	return value, br.db.View(
		func(tx *badger.Txn) error {
			item, err := tx.Get([]byte(id))
			if err != nil {
				return fmt.Errorf("getting value: %w", err)
			}
			valCopy, err := item.ValueCopy(nil)
			if err != nil {
				return fmt.Errorf("copying value: %w", err)
			}
			value = string(valCopy)
			return nil
		})
}

func (br *BadgerRepo) Delete(chatId string) error {
	txn := br.db.NewTransaction(true)

	err := txn.Delete([]byte(chatId))
	if err != nil {
		return err
	}
	err = txn.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (br *BadgerRepo) Close() error {
	err := br.db.Close()
	if err != nil {
		return err
	}
	return nil
}

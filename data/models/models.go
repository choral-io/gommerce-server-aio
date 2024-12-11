package models

import (
	"errors"

	"github.com/uptrace/bun"
)

func RegisterModels(bdb bun.IDB) error {
	if _, ok := bdb.(*bun.DB); ok {
		// db.RegisterModel((*ModelType)(nil))
		return nil
	}
	return errors.New("failed to register models, provided db must be *bun.DB")
}

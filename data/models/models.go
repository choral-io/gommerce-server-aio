package models

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/uptrace/bun"
)

// ErrImmutableModel represents an error when attempting to modify an immutable model.
var ErrImmutableModel = errors.New("model is immutable and cannot be updated or deleted")

// RegisterModels registers all models with the provided bun database instance.
func RegisterModels(bdb bun.IDB) error {
	if _, ok := bdb.(*bun.DB); ok {
		// db.RegisterModel((*ModelType)(nil))
		return nil
	}
	return errors.New("failed to register models, provided db must be *bun.DB")
}

// SeedingMode checks if the application is running in seeding mode.
func SeedingMode(_ context.Context) bool {
	return strings.EqualFold(os.Getenv("DATA_SEEDING_MODE"), "true")
}

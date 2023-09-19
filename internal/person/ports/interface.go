package ports

import (
	"context"

	"github.com/alukart32/effective-mobile-test-task/internal/person/model"
)

type personCreator interface {
	CreateFrom(ctx context.Context, fio model.FIO) (string, error)
}

type personFinder interface {
	FindById(ctx context.Context, id string) (model.Person, error)
}

type personCollector interface {
	Collect(ctx context.Context, filter model.PersonFilter, limit, offset int) ([]model.Person, error)
}

type personUpdater interface {
	Update(ctx context.Context, id string, meta model.PersonalMetaData) error
}

type personDeleter interface {
	Delete(ctx context.Context, id string) error
}

type personManager interface {
	personCreator
	personFinder
	personCollector
	personUpdater
	personDeleter
}

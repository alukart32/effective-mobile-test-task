package persondata

import (
	"context"
	"fmt"

	"github.com/alukart32/effective-mobile-test-task/internal/person/model"
)

type metaDataProvider interface {
	AgeByName(context.Context, string) (int, error)
	GenderByName(context.Context, string) (string, error)
	NationByName(context.Context, string) (string, error)
}

type saver interface {
	Save(context.Context, model.Person) error
}

type finder interface {
	FindById(ctx context.Context, id string) (model.Person, error)
}

type collector interface {
	Collect(ctx context.Context, filter model.PersonFilter, limit, offset int) ([]model.Person, error)
}

type updater interface {
	Update(ctx context.Context, id string, meta model.PersonalMetaData) error
}

type deleter interface {
	Delete(ctx context.Context, id string) error
}

type repo interface {
	saver
	finder
	collector
	updater
	deleter
}

type manager struct {
	repo
	metaDataProvider
}

func Manager(r repo, metaData metaDataProvider) (*manager, error) {
	if r == nil {
		return nil, fmt.Errorf("repo is nil")
	}
	if metaData == nil {
		return nil, fmt.Errorf("metaData provider is nil")
	}

	return &manager{
		repo:             r,
		metaDataProvider: metaData,
	}, nil
}

func (m *manager) CreateFrom(ctx context.Context, fio model.FIO) (string, error) {
	age, err := m.metaDataProvider.AgeByName(ctx, fio.Name)
	if err != nil {
		return "", fmt.Errorf("PersonManager.CreateFrom: %w", err)
	}
	gender, err := m.metaDataProvider.GenderByName(ctx, fio.Name)
	if err != nil {
		return "", fmt.Errorf("PersonManager.CreateFrom: %w", err)
	}
	nation, err := m.metaDataProvider.NationByName(ctx, fio.Name)
	if err != nil {
		return "", fmt.Errorf("PersonManager.CreateFrom: %w", err)
	}

	person := model.NewPerson(fio, model.PersonalMetaData{
		Nation: nation,
		Gender: gender,
		Age:    age,
	})

	if err = m.repo.Save(ctx, person); err != nil {
		return "", fmt.Errorf("PersonManager.CreateFrom: %w", err)
	}
	return person.Id, nil
}

func (m *manager) FindById(ctx context.Context, id string) (model.Person, error) {
	if len(id) == 0 {
		return model.Person{}, fmt.Errorf("PersonManager.FindById: empty id")
	}

	person, err := m.repo.FindById(ctx, id)
	if err != nil {
		return model.Person{}, fmt.Errorf("PersonManager.FindById: %w", err)
	}
	if person.IsEmpty() {
		return model.Person{}, fmt.Errorf("PersonManager.FindById: not found")
	}
	return person, nil
}

func (m *manager) Collect(ctx context.Context, filter model.PersonFilter, limit, offset int) ([]model.Person, error) {
	persons, err := m.repo.Collect(ctx, filter, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("PersonManager.Collect: %w", err)
	}
	return persons, nil
}

func (m *manager) Update(ctx context.Context, id string, meta model.PersonalMetaData) error {
	if len(id) == 0 {
		return fmt.Errorf("PersonManager.Update: empty id")
	}
	if meta.IsEmpty() {
		return fmt.Errorf("PersonManager.Update: no data for update")
	}

	err := m.repo.Update(ctx, id, meta)
	if err != nil {
		return fmt.Errorf("PersonManager.Update: %w", err)
	}
	return nil
}

func (m *manager) Delete(ctx context.Context, id string) error {
	if len(id) == 0 {
		return fmt.Errorf("PersonManager.Delete: empty id")
	}

	err := m.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("PersonManager.Delete: %w", err)
	}
	return nil
}

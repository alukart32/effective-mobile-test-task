package persondata

import (
	"context"
	"fmt"
	"testing"

	"github.com/alukart32/effective-mobile-test-task/internal/person/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type metaDataProviderMock struct {
	AgeByNameFn    func(context.Context, string) (int, error)
	GenderByNameFn func(context.Context, string) (string, error)
	NationByNameFn func(context.Context, string) (string, error)
}

func (m *metaDataProviderMock) AgeByName(ctx context.Context, name string) (int, error) {
	if m != nil && m.AgeByNameFn != nil {
		return m.AgeByNameFn(ctx, name)
	}
	return 0, fmt.Errorf("can't get age")
}

func (m *metaDataProviderMock) GenderByName(ctx context.Context, name string) (string, error) {
	if m != nil && m.GenderByNameFn != nil {
		return m.GenderByNameFn(ctx, name)
	}
	return "", fmt.Errorf("can't get gender")
}

func (m *metaDataProviderMock) NationByName(ctx context.Context, name string) (string, error) {
	if m != nil && m.NationByNameFn != nil {
		return m.NationByNameFn(ctx, name)
	}
	return "", fmt.Errorf("can't get nation")
}

type saverMock struct {
	SaveFn func(context.Context, model.Person) error
}

func (m *saverMock) Save(ctx context.Context, person model.Person) error {
	if m != nil && m.SaveFn != nil {
		return m.SaveFn(ctx, person)
	}
	return fmt.Errorf("can't save person")
}

type finderMock struct {
	FindByIdFn func(ctx context.Context, id string) (model.Person, error)
}

func (m *finderMock) FindById(ctx context.Context, id string) (model.Person, error) {
	if m != nil && m.FindByIdFn != nil {
		return m.FindByIdFn(ctx, id)
	}
	return model.Person{}, fmt.Errorf("can't find person")
}

type collectorMock struct {
	CollectFn func(ctx context.Context, filter model.PersonFilter, limit, offset int) ([]model.Person, error)
}

func (m *collectorMock) Collect(ctx context.Context, filter model.PersonFilter, limit, offset int) ([]model.Person, error) {
	if m != nil && m.CollectFn != nil {
		return m.CollectFn(ctx, filter, limit, offset)
	}
	return nil, fmt.Errorf("can't collect persons")
}

type updaterMock struct {
	UpdateFn func(ctx context.Context, id string, meta model.PersonalMetaData) error
}

func (m *updaterMock) Update(ctx context.Context, id string, meta model.PersonalMetaData) error {
	if m != nil && m.UpdateFn != nil {
		return m.UpdateFn(ctx, id, meta)
	}
	return fmt.Errorf("can't update person")
}

type deleterMock struct {
	DeleteFn func(ctx context.Context, id string) error
}

func (m *deleterMock) Delete(ctx context.Context, id string) error {
	if m != nil && m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return fmt.Errorf("can't delete person")
}

type repoMock struct {
	saverMock
	finderMock
	collectorMock
	updaterMock
	deleterMock
}

func (m *repoMock) Save(ctx context.Context, p model.Person) error {
	return m.saverMock.Save(ctx, p)
}

func (m *repoMock) FindById(ctx context.Context, id string) (model.Person, error) {
	return m.finderMock.FindById(ctx, id)
}

func (m *repoMock) Collect(ctx context.Context, filter model.PersonFilter, limit, offset int) ([]model.Person, error) {
	return m.collectorMock.Collect(ctx, filter, limit, offset)
}

func (m *repoMock) Update(ctx context.Context, id string, meta model.PersonalMetaData) error {
	return m.updaterMock.Update(ctx, id, meta)
}

func (m *repoMock) Delete(ctx context.Context, id string) error {
	return m.deleterMock.Delete(ctx, id)
}

func TestManager_CreateFrom(t *testing.T) {
	type services struct {
		metaProvider metaDataProviderMock
		saver        saverMock
	}
	type want struct {
		err bool
	}
	tests := []struct {
		name string
		fio  model.FIO
		want want
		serv services
	}{
		{
			name: "Valid fio, no error",
			fio: model.FIO{
				Name:       "test",
				Surname:    "test",
				Patronymic: "optional",
			},
			want: want{
				err: false,
			},
			serv: services{
				saver: saverMock{
					SaveFn: func(ctx context.Context, p model.Person) error {
						return nil
					},
				},
				metaProvider: metaDataProviderMock{
					AgeByNameFn: func(ctx context.Context, s string) (int, error) {
						return 20, nil
					},
					GenderByNameFn: func(ctx context.Context, s string) (string, error) {
						return "test", nil
					},
					NationByNameFn: func(ctx context.Context, s string) (string, error) {
						return "go", nil
					},
				},
			},
		},
		{
			name: "Can't get age, error",
			fio: model.FIO{
				Name:       "test",
				Surname:    "test",
				Patronymic: "optional",
			},
			want: want{
				err: true,
			},
			serv: services{
				metaProvider: metaDataProviderMock{
					AgeByNameFn: func(ctx context.Context, s string) (int, error) {
						return 0, fmt.Errorf("error")
					},
				},
			},
		},
		{
			name: "Can't get gender, error",
			fio: model.FIO{
				Name:       "test",
				Surname:    "test",
				Patronymic: "optional",
			},
			want: want{
				err: true,
			},
			serv: services{
				metaProvider: metaDataProviderMock{
					AgeByNameFn: func(ctx context.Context, s string) (int, error) {
						return 20, nil
					},
					GenderByNameFn: func(ctx context.Context, s string) (string, error) {
						return "", fmt.Errorf("error")
					},
				},
			},
		},
		{
			name: "Can't get nation, error",
			fio: model.FIO{
				Name:       "test",
				Surname:    "test",
				Patronymic: "optional",
			},
			want: want{
				err: true,
			},
			serv: services{
				metaProvider: metaDataProviderMock{
					AgeByNameFn: func(ctx context.Context, s string) (int, error) {
						return 20, nil
					},
					GenderByNameFn: func(ctx context.Context, s string) (string, error) {
						return "test", nil
					},
					NationByNameFn: func(ctx context.Context, s string) (string, error) {
						return "", fmt.Errorf("error")
					},
				},
			},
		},
		{
			name: "Can't save person, error",
			fio: model.FIO{
				Name:       "test",
				Surname:    "test",
				Patronymic: "optional",
			},
			want: want{
				err: false,
			},
			serv: services{
				saver: saverMock{
					SaveFn: func(ctx context.Context, p model.Person) error {
						return fmt.Errorf("error")
					},
				},
				metaProvider: metaDataProviderMock{
					AgeByNameFn: func(ctx context.Context, s string) (int, error) {
						return 20, nil
					},
					GenderByNameFn: func(ctx context.Context, s string) (string, error) {
						return "test", nil
					},
					NationByNameFn: func(ctx context.Context, s string) (string, error) {
						return "go", nil
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := Manager(&repoMock{saverMock: tt.serv.saver}, &tt.serv.metaProvider)
			require.NoError(t, err)

			_, err = manager.CreateFrom(context.Background(), tt.fio)
			if tt.want.err {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestManager_FindById(t *testing.T) {
	type services struct {
		finder finderMock
	}
	type want struct {
		person model.Person
		err    error
	}
	tests := []struct {
		name string
		id   string
		want want
		serv services
	}{
		{
			name: "Found by id, no error",
			id:   "person_1",
			want: want{
				person: model.Person{
					Id: "person_1",
				},
				err: nil,
			},
			serv: services{
				finder: finderMock{
					FindByIdFn: func(ctx context.Context, id string) (model.Person, error) {
						return model.Person{
							Id: "person_1",
						}, nil
					},
				},
			},
		},
		{
			name: "Empty id, no error",
			id:   "",
			want: want{
				err: fmt.Errorf("PersonManager.FindById: empty id"),
			},
			serv: services{
				finder: finderMock{},
			},
		},
		{
			name: "Not found, error",
			id:   "person_unknown",
			want: want{
				err: fmt.Errorf("PersonManager.FindById: not found"),
			},
			serv: services{
				finder: finderMock{
					FindByIdFn: func(ctx context.Context, id string) (model.Person, error) {
						return model.Person{}, nil
					},
				},
			},
		},
		{
			name: "Finder error",
			id:   "person_1",
			want: want{
				err: fmt.Errorf("PersonManager.FindById: internal error"),
			},
			serv: services{
				finder: finderMock{
					FindByIdFn: func(ctx context.Context, id string) (model.Person, error) {
						return model.Person{}, fmt.Errorf("internal error")
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := Manager(&repoMock{finderMock: tt.serv.finder}, &metaDataProviderMock{})
			require.NoError(t, err)

			_, err = manager.FindById(context.Background(), tt.id)
			if tt.want.err != nil {
				assert.EqualError(t, err, tt.want.err.Error())
			}
		})
	}
}

func TestManager_Collect(t *testing.T) {
	type services struct {
		collector collectorMock
	}
	type args struct {
		limit   int
		offset  int
		filters model.PersonFilter
	}
	type want struct {
		persons []model.Person
		err     error
	}
	tests := []struct {
		name string
		args args
		want want
		serv services
	}{
		{
			name: "Collect, no error",
			args: args{
				limit:  3,
				offset: 20,
				filters: model.PersonFilter{
					OlderThan:   15,
					YoungerThan: 40,
					Nations:     []string{"T1", "T2", "T3"},
				},
			},
			want: want{
				err: nil,
				persons: []model.Person{
					{
						Id: "1",
						FIO: model.FIO{
							Name: "person_1",
						},
						PersonalMetaData: model.PersonalMetaData{
							Nation: "T1",
							Age:    17,
						},
					},
					{
						Id: "2",
						FIO: model.FIO{
							Name: "person_2",
						},
						PersonalMetaData: model.PersonalMetaData{
							Nation: "T2",
							Age:    20,
						},
					},
					{
						Id: "3",
						FIO: model.FIO{
							Name: "person_3",
						},
						PersonalMetaData: model.PersonalMetaData{
							Nation: "T3",
							Age:    21,
						},
					},
				},
			},
			serv: services{
				collector: collectorMock{
					CollectFn: func(ctx context.Context, filter model.PersonFilter, limit, offset int) ([]model.Person, error) {
						return []model.Person{
							{
								Id: "1",
								FIO: model.FIO{
									Name: "person_1",
								},
								PersonalMetaData: model.PersonalMetaData{
									Nation: "T1",
									Age:    17,
								},
							},
							{
								Id: "2",
								FIO: model.FIO{
									Name: "person_2",
								},
								PersonalMetaData: model.PersonalMetaData{
									Nation: "T2",
									Age:    20,
								},
							},
							{
								Id: "3",
								FIO: model.FIO{
									Name: "person_3",
								},
								PersonalMetaData: model.PersonalMetaData{
									Nation: "T3",
									Age:    21,
								},
							},
						}, nil
					},
				},
			},
		},
		{
			name: "Collector error",
			args: args{
				limit:  3,
				offset: 20,
				filters: model.PersonFilter{
					OlderThan:   15,
					YoungerThan: 40,
					Nations:     []string{"T1", "T2", "T3"},
				},
			},
			want: want{
				err: fmt.Errorf("PersonManager.Collect: internal error"),
			},
			serv: services{
				collector: collectorMock{
					CollectFn: func(ctx context.Context, filter model.PersonFilter, limit, offset int) ([]model.Person, error) {
						return []model.Person{}, fmt.Errorf("internal error")
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := Manager(&repoMock{collectorMock: tt.serv.collector}, &metaDataProviderMock{})
			require.NoError(t, err)

			persons, err := manager.Collect(context.Background(), tt.args.filters, tt.args.limit, tt.args.offset)
			if tt.want.err != nil {
				assert.EqualError(t, err, tt.want.err.Error())
				return
			}

			for i, p := range persons {
				assert.EqualValues(t, tt.want.persons[i].Id, p.Id)
				assert.EqualValues(t, tt.want.persons[i].FIO, p.FIO)
				assert.EqualValues(t, tt.want.persons[i].PersonalMetaData, p.PersonalMetaData)
			}
		})
	}
}

func TestManager_Update(t *testing.T) {
	type services struct {
		updater updaterMock
	}
	type args struct {
		id   string
		meta model.PersonalMetaData
	}
	type want struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want want
		serv services
	}{
		{
			name: "Updated, no error",
			args: args{
				id: "person_1",
				meta: model.PersonalMetaData{
					Nation: "T1",
					Age:    25,
				},
			},
			want: want{
				err: nil,
			},
			serv: services{
				updater: updaterMock{
					UpdateFn: func(ctx context.Context, id string, meta model.PersonalMetaData) error {
						return nil
					},
				},
			},
		},
		{
			name: "Empty id, error",
			args: args{
				id: "",
				meta: model.PersonalMetaData{
					Nation: "GB",
					Age:    67,
				},
			},
			want: want{
				err: fmt.Errorf("PersonManager.Update: empty id"),
			},
			serv: services{
				updater: updaterMock{},
			},
		},
		{
			name: "Empty personal meta data, error",
			args: args{
				id:   "person_1",
				meta: model.PersonalMetaData{},
			},
			want: want{
				err: fmt.Errorf("PersonManager.Update: no data for update"),
			},
			serv: services{
				updater: updaterMock{},
			},
		},
		{
			name: "Updater error",
			args: args{
				id: "person_1",
				meta: model.PersonalMetaData{
					Nation: "T1",
					Age:    25,
				},
			},
			want: want{
				err: fmt.Errorf("PersonManager.Update: internal error"),
			},
			serv: services{
				updater: updaterMock{
					UpdateFn: func(ctx context.Context, id string, meta model.PersonalMetaData) error {
						return fmt.Errorf("internal error")
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := Manager(&repoMock{updaterMock: tt.serv.updater}, &metaDataProviderMock{})
			require.NoError(t, err)

			err = manager.Update(context.Background(), tt.args.id, tt.args.meta)
			if tt.want.err != nil {
				assert.EqualError(t, err, tt.want.err.Error())
			}
		})
	}
}

func TestManager_Delete(t *testing.T) {
	type services struct {
		deleter deleterMock
	}
	type want struct {
		err error
	}
	tests := []struct {
		name string
		id   string
		want want
		serv services
	}{
		{
			name: "Deleted, no error",
			id:   "person_1",
			want: want{
				err: nil,
			},
			serv: services{
				deleter: deleterMock{
					DeleteFn: func(ctx context.Context, id string) error {
						return nil
					},
				},
			},
		},
		{
			name: "Empty id, error",
			id:   "",
			want: want{
				err: fmt.Errorf("PersonManager.Delete: empty id"),
			},
			serv: services{
				deleter: deleterMock{},
			},
		},
		{
			name: "Deleter error",
			id:   "person_1",
			want: want{
				err: fmt.Errorf("PersonManager.Delete: internal error"),
			},
			serv: services{
				deleter: deleterMock{
					DeleteFn: func(ctx context.Context, id string) error {
						return fmt.Errorf("internal error")
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := Manager(&repoMock{deleterMock: tt.serv.deleter}, &metaDataProviderMock{})
			require.NoError(t, err)

			err = manager.Delete(context.Background(), tt.id)
			if tt.want.err != nil {
				assert.EqualError(t, err, tt.want.err.Error())
			}
		})
	}
}

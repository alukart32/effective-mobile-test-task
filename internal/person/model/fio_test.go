package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFio(t *testing.T) {
	type args struct {
		name       string
		surname    string
		patronymic string
	}
	type want struct {
		fio FIO
		err error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Valid args, no error",
			args: args{
				name:       "test",
				surname:    "test",
				patronymic: "test",
			},
			want: want{
				fio: FIO{
					Name:       "test",
					Surname:    "test",
					Patronymic: "test",
				},
				err: nil,
			},
		},
		{
			name: "Empty name, error",
			args: args{
				surname:    "test",
				patronymic: "test",
			},
			want: want{
				err: fmt.Errorf("empty required name"),
			},
		},
		{
			name: "Invalid name, error",
			args: args{
				name:       "test@test",
				surname:    "test",
				patronymic: "test",
			},
			want: want{
				err: fmt.Errorf("name contains invalid characters: test@test"),
			},
		},
		{
			name: "Empty surname, error",
			args: args{
				name:       "test",
				patronymic: "test",
			},
			want: want{
				err: fmt.Errorf("empty required surname"),
			},
		},
		{
			name: "Invalid surname, error",
			args: args{
				name:       "test",
				surname:    "test@test",
				patronymic: "test",
			},
			want: want{
				err: fmt.Errorf("surname contains invalid characters: test@test"),
			},
		},
		{
			name: "Invalid patronymic, error",
			args: args{
				name:       "test",
				surname:    "test",
				patronymic: "test@test",
			},
			want: want{
				err: fmt.Errorf("patronymic contains invalid characters: test@test"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fio, err := NewFIO(tt.args.name, tt.args.surname, tt.args.patronymic)
			if tt.want.err != nil {
				assert.EqualError(t, err, tt.want.err.Error())
				return
			}

			assert.EqualValues(t, tt.want.fio.Name, fio.Name)
			assert.EqualValues(t, tt.want.fio.Surname, fio.Surname)
			assert.EqualValues(t, tt.want.fio.Patronymic, fio.Patronymic)
		})
	}
}

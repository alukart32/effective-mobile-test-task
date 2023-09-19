package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPersonFilter(t *testing.T) {
	type want struct {
		filter PersonFilter
		err    error
	}
	tests := []struct {
		name    string
		filters []string
		want    want
	}{
		{
			name: "Valid filters, no error",
			filters: []string{
				"older-than.1",
				"younger-than.99",
				"gender.test",
				"nation.T1",
			},
			want: want{
				filter: PersonFilter{
					OlderThan:   1,
					YoungerThan: 99,
					Gender:      "test",
					Nations:     []string{"T1"},
				},
			},
		},
		{
			name: "Unknown filter, error",
			filters: []string{
				"older-than.1",
				"younger-than.99",
				"gender.test",
				"nation.T1",
				"unknown.test",
			},
			want: want{
				err: fmt.Errorf("unsupported %s filter", "unknown.test"),
			},
		},
		{
			name: "Invalid older-than filter, error",
			filters: []string{
				"older-than.txt",
			},
			want: want{
				err: fmt.Errorf("filter parsing error: %s", "older-than.txt"),
			},
		},
		{
			name: "Invalid older-than filter, error",
			filters: []string{
				"younger-than.txt",
			},
			want: want{
				err: fmt.Errorf("filter parsing error: %s", "younger-than.txt"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := NewPersonFilter(tt.filters)
			if tt.want.err != nil {
				assert.EqualError(t, err, tt.want.err.Error())
				return
			}
			assert.EqualValues(t, tt.want.filter, filter)
		})
	}
}

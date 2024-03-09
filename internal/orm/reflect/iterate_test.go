package reflect

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterateArrayOrSlice(t *testing.T) {
	testCases := []struct{
		name  	string

		entity any

		wantRes []any
		wantErr error
	}{
		{
			name: "[3]int",

			entity: [3]int{1,2,3},

			wantRes: []any{1,2,3},
		},
		{
			name: "[]int",
			entity: []int{1,2,3},
			wantRes: []any{1,2,3},
		},
	}
	
	for _, tc := range  testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := IterateArrayOrSlice(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, val)
		})
	}
}



func TestIterateMap(t *testing.T) {
	testCases := []struct{
		name  string

		entity any

		wantKeys   []any
		wantValues  []any
		wantErr    error
	}{
		{
			name: "map",

			entity: map[string]string{
				"123": "hello",
				"456": "world",
			},

			wantKeys: []any{"123","456"},
			wantValues: []any{"hello", "world"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			keys, values, err := IterateMap(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.EqualValues(t, tc.wantKeys, keys)
			assert.EqualValues(t, tc.wantValues, values)
		})
	}
}
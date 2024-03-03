package reflect

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterateFunc(t *testing.T) {
	testCases := []struct{
		name string
		entity any

		wantRes map[string]FuncInfo
		wantErr  error
	}{
		{
			name: "struct",
		},
		
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := IterateFields(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

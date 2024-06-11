package valuer

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ExerciseCoding/template/internal/orm/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func BenchmarkSetColumns(b *testing.B) {
	fn := func(b *testing.B, creator Creator) {
		mockDB, mock, err := sqlmock.New()
		defer mockDB.Close()
		require.NoError(b, err)
		// 我们需要跑N次，也就是需要准备N行
		mockRows := mock.NewRows([]string{"id", "first_name", "age", "last_name"})

		row := []driver.Value{"1", "Tom", "18", "Jerry"}
		for i := 0; i < b.N; i++ {
			mockRows.AddRow(row...)
		}
		mock.ExpectQuery("SELECT XXX").WillReturnRows(mockRows)

		rows, err := mockDB.Query("SELECT XXX")
		require.NoError(b, err)
		r := model.NewRegistry()
		m, err := r.Get(&TestModel{})
		require.NoError(b, err)

		// 重置计时器
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			rows.Next()
			val := creator(m, &TestModel{})
			_ = val.SetColums(rows)
		}
	}
	b.Run("flect", func(b *testing.B) {
		fn(b, NewReflectValue)
	})

	b.Run("unsafe", func(b *testing.B) {
		fn(b, NewUnsafeValue)
	})
}

package berr_test

import (
	"errors"
	"testing"

	"github.com/tuingking/berr"
)

func Test_Printer(t *testing.T) {
	t.Run("error is wrapped", func(t *testing.T) {
		root := errors.New("root cause")
		err := berr.WrapWithCode(berr.CodeSQLInsert, root, "data not found")
		err = berr.WrapWithCode(berr.CodeInternal, err, "failed to fetch data")
		err = berr.Wrapx(err)

		t.Run("PrintStack", func(t *testing.T) {
			berr.PrintStack(err, berr.PrintWithFile(true), berr.PrintWithLine(true), berr.PrintWithShortFunc(true))
		})
		t.Run("PrintJSON", func(t *testing.T) {
			berr.PrintJSON(err)
		})
	})

	t.Run("error is not wrapped", func(t *testing.T) {
		err := errors.New("root cause")
		t.Run("PrintStack", func(t *testing.T) {
			berr.PrintStack(err)
		})
		t.Run("PrintJSON", func(t *testing.T) {
			berr.PrintJSON(err)
		})
	})
}

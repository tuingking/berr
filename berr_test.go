package berr_test

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/tuingking/berr"

	"gotest.tools/assert"
)

func TestError_Wrap(t *testing.T) {
	t.Run("Wrap", func(t *testing.T) {
		root := errors.New("root cause")
		err := berr.Wrap(root, "message 3")
		err = berr.Wrap(err, "message 2")
		err = berr.Wrap(err, "message 1")
		assert.Error(t, err, "root cause")
		assert.Equal(t, berr.GetCode(err), berr.CodeInternal)
		assert.Equal(t, berr.GetMsg(err), "message 1: message 2: message 3")
		assert.Error(t, berr.GetErrRoot(err), "root cause")
		assert.Equal(t, len(berr.GetStack(err)), 3)
		assert.Equal(t, berr.Is(err, root), true)
		assert.Equal(t, errors.Is(berr.GetErrRoot(err), root), true)
		berr.PrintStack(err, berr.PrintWithShortFunc(false))
	})

	t.Run("Wrap nil", func(t *testing.T) {
		err := berr.Wrap(nil, "message 3")
		assert.NilError(t, err)
	})
}

func TestError_WrapWithCode(t *testing.T) {
	t.Run("WrapWithCode", func(t *testing.T) {
		root := errors.New("root cause")
		err := berr.WrapWithCode(berr.CodeSQLInsert, root, "message 3")
		err = berr.WrapWithCode(berr.CodeDB, err, "message 2")
		err = berr.WrapWithCode(berr.CodeBusiness, err, "message 1")
		assert.Error(t, err, "root cause")
		assert.Equal(t, berr.GetCode(err), berr.CodeSQLInsert)
		assert.Equal(t, berr.GetMsg(err), "message 1: message 2: message 3")
		assert.Error(t, berr.GetErrRoot(err), "root cause")
		assert.Equal(t, len(berr.GetStack(err)), 3)
		assert.Equal(t, berr.Is(err, root), true)
		assert.Equal(t, errors.Is(berr.GetErrRoot(err), root), true)
		berr.PrintStack(err)
	})

	t.Run("WrapWithCode nil", func(t *testing.T) {
		err := berr.WrapWithCode(berr.CodeDB, nil, "message 1")
		assert.NilError(t, err)
	})
}

func TestError_Wrapx(t *testing.T) {
	t.Run("Wrapx", func(t *testing.T) {
		root := errors.New("root cause")
		err := berr.Wrapx(root)
		err = berr.Wrapx(err)
		err = berr.Wrapx(err)
		assert.Error(t, err, "root cause")
		assert.Equal(t, berr.GetCode(err), berr.CodeInternal)
		assert.Equal(t, berr.GetMsg(err), "")
		assert.Error(t, berr.GetErrRoot(err), "root cause")
		assert.Equal(t, len(berr.GetStack(err)), 3)
		assert.Equal(t, berr.Is(err, root), true)
		assert.Equal(t, errors.Is(berr.GetErrRoot(err), root), true)
		berr.PrintStack(err)
	})

	t.Run("Wrapx nil", func(t *testing.T) {
		err := berr.Wrapx(nil)
		assert.NilError(t, err)
	})
}

func TestError_New(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		err := berr.New("root cause")
		assert.Error(t, err, "root cause")
		assert.Equal(t, berr.GetCode(err), berr.CodeInternal)
		assert.Equal(t, berr.GetMsg(err), "root cause")
		assert.Error(t, berr.GetErrRoot(err), "root cause")
		assert.Equal(t, len(berr.GetStack(err)), 1)
		err2 := err
		assert.Equal(t, berr.Is(err, err2), true)
		berr.PrintStack(err)
	})
}

func TestError_Newx(t *testing.T) {
	t.Run("Newx", func(t *testing.T) {
		err := berr.Newx("root cause")
		assert.Error(t, err, "root cause")
		assert.Equal(t, berr.GetCode(err), berr.CodeInternal)
		assert.Equal(t, berr.GetMsg(err), "root cause")
		assert.Error(t, berr.GetErrRoot(err), "root cause")
		assert.Equal(t, len(berr.GetStack(err)), 0)
		err2 := err
		assert.Equal(t, berr.Is(err, err2), true)
		berr.PrintStack(err)
	})
}

func TestError_Combination(t *testing.T) {
	t.Run("combination", func(t *testing.T) {
		root := errors.New("root cause")
		err := fmt.Errorf("foo: %w", root)
		err = berr.Wrapx(root)
		err = berr.Wrap(err, "message 2")
		err = berr.WrapWithCode(berr.CodeBusiness, err, "message 1")
		err = berr.Wrapx(err)
		assert.Error(t, err, "root cause")
		assert.Equal(t, berr.GetCode(err), berr.CodeBusiness)
		assert.Equal(t, berr.GetMsg(err), "message 1: message 2")
		assert.Error(t, berr.GetErrRoot(err), "root cause")
		assert.Equal(t, len(berr.GetStack(err)), 4)
		assert.Equal(t, berr.Is(err, root), true)
		assert.Equal(t, errors.Is(berr.GetErrRoot(err), root), true)
		berr.PrintStack(err)
	})
}

func Test_Utils(t *testing.T) {
	t.Run("error is wrapped", func(t *testing.T) {
		root := sql.ErrNoRows
		err := berr.WrapWithCode(berr.CodeSQLInsert, root, "data not found")
		err = berr.WrapWithCode(berr.CodeInternal, err, "failed to fetch data")
		err = berr.Wrapx(err)

		assert.Equal(t, berr.Is(err, root), true)
		assert.Error(t, err, "sql: no rows in result set")
		assert.Error(t, berr.GetErrRoot(err), "sql: no rows in result set")
		assert.Error(t, berr.GetErrChain(err), "func1|132: func1|131: func1|130: sql: no rows in result set")
		assert.Equal(t, berr.GetCode(err), berr.CodeSQLInsert)
		assert.Equal(t, berr.GetMsg(err), "failed to fetch data: data not found")
		assert.Equal(t, berr.GetMsgRoot(err), "data not found")
		assert.Equal(t, len(berr.GetStack(err)), 3)
		assert.DeepEqual(t, berr.GetStack(err).Compact(berr.PrintWithShortFile(true)), []string{
			"func1: berr_test.go:130",
			"func1: berr_test.go:131",
			"func1: berr_test.go:132",
		})
	})

	t.Run("error is not wrapped", func(t *testing.T) {
		err := errors.New("root cause")
		assert.Error(t, err, "root cause")
		assert.Error(t, berr.GetErrRoot(err), "root cause")
		assert.Equal(t, berr.GetMsg(err), "root cause")
		assert.Equal(t, berr.GetMsgRoot(err), "root cause")
		assert.Error(t, berr.GetErrChain(err), "root cause")
		assert.Equal(t, len(berr.GetStack(err)), 0)
	})

	t.Run("nil error", func(t *testing.T) {
		var err error
		assert.NilError(t, err, nil)
		assert.NilError(t, berr.GetErrRoot(err))
		assert.Equal(t, berr.GetMsg(err), "")
		assert.Equal(t, berr.GetMsgRoot(err), "")
		assert.NilError(t, berr.GetErrChain(err))
		assert.Equal(t, len(berr.GetStack(err)), 0)
	})
}

func TestError_Must(t *testing.T) {
	t.Run("Must", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				assert.Error(t, r.(error), "root cause")
			}
		}()
		callback := func() error {
			return errors.New("root cause")
		}
		berr.Must(callback())
	})
}

package test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// UT is a utility to implement unit tests.
type UT interface {
	Ctx() context.Context
	// Assert continues to the end of a test even though there is an assertion error at somewhere in the test
	Assert() *assert.Assertions
	// Require stops at the first assertion error
	Require() *require.Assertions
	// AssertError is used verify an error in both case nil and not nil
	AssertError(expectedErr, actualErr error)
	// RequireError is used verify an error in both case nil and not nil
	RequireError(expectedErr, actualErr error)
}

//nolint:containedctx
type ut struct {
	t       *testing.T
	ctx     context.Context
	assert  *assert.Assertions
	require *require.Assertions
}

func NewUT(t *testing.T) UT {
	t.Helper()

	return newUT(t)
}

func newUT(t *testing.T) *ut {
	t.Helper()

	// SetMode cause DATA RACE in tests because it update a global variable
	// gin.SetMode(gin.TestMode)
	// nolint: tenv
	os.Setenv("NEW_RELIC_ENABLED", "false")
	return &ut{
		t:       t,
		ctx:     context.Background(),
		assert:  assert.New(t),
		require: require.New(t),
	}
}

func (u *ut) Ctx() context.Context {
	return u.ctx
}

func (u *ut) Assert() *assert.Assertions {
	return u.assert
}

func (u *ut) Require() *require.Assertions {
	return u.require
}

func (u *ut) AssertError(expectedErr, actualErr error) {
	if expectedErr == nil {
		u.assert.NoError(actualErr)
	} else {
		u.assert.ErrorIs(actualErr, expectedErr)
	}
}

func (u *ut) RequireError(expectedErr, actualErr error) {
	if expectedErr == nil {
		u.require.NoError(actualErr)
	} else {
		u.require.ErrorIs(actualErr, expectedErr)
	}
}

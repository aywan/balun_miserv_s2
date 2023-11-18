package db

import (
	"context"
	"testing"
)

const (
	testingTxKey key = "testing_tx"
)

type TestTxManager struct {
	txLevel int
}

func (t TestTxManager) ReadCommitted(ctx context.Context, f TxHandler) error {
	t.txLevel++
	err := f(ctx)
	t.txLevel--

	return err
}

func NewTestTxManager(t *testing.T) *TestTxManager {
	m := &TestTxManager{}

	t.Cleanup(func() {
		if m.txLevel != 0 {
			t.Errorf("not closed transaction")
		}
	})

	return m
}

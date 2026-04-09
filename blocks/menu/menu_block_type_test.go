package menu

import (
	"testing"

	"github.com/dracory/cmsstore/testutils"
	_ "modernc.org/sqlite"
)

func TestMenuBlockType_GetCustomVariables(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	bt := NewMenuBlockType(store, &testLogger{})
	vars := bt.GetCustomVariables()
	if vars != nil {
		t.Errorf("expected nil custom variables for menu block, got %v", vars)
	}
}

// testLogger is a minimal logger for tests
type testLogger struct{}

func (l *testLogger) Error(msg string, args ...interface{}) {}

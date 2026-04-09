package html

import (
	"testing"
)

func TestHTMLBlockType_GetCustomVariables(t *testing.T) {
	bt := NewHTMLBlockType()
	vars := bt.GetCustomVariables()
	if vars != nil {
		t.Errorf("expected nil custom variables for HTML block, got %v", vars)
	}
}

package translation

import (
	"sort"
	"testing"
)

// Check this here to avoid unnecessary import of sort package.
var _ = sort.Interface(make(SortableByID, 0, 0))

func TestNewSingleTranslation(t *testing.T) {
	t.Skipf("not implemented")
}

func TestNewPluralTranslation(t *testing.T) {
	t.Skipf("not implemented")
}

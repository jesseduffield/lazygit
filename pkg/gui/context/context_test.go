package context

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// The config package validates a custom command's context against its own copy of
// these names, being unable to import this package. A name in one list but not the
// other would be either a context that validation rejects although you can bind to
// it, or one it accepts although binding to it exits lazygit.
func TestValidCustomCommandContextsMatchesAllContextKeys(t *testing.T) {
	keys := lo.Map(AllContextKeys, func(key types.ContextKey, _ int) string {
		return string(key)
	})

	assert.Equal(t, keys, config.ValidCustomCommandContexts)
}

package components

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
	"github.com/stretchr/testify/assert"
)

type fakeGuiDriver struct {
	failureMessage string
	pressedKeys    []string
}

var _ integrationTypes.GuiDriver = &fakeGuiDriver{}

type GuiDriver interface {
	PressKey(string)
	Keys() config.KeybindingConfig
	CurrentContext() types.Context
	Model() *types.Model
	Fail(message string)
	// These two log methods are for the sake of debugging while testing. There's no need to actually
	// commit any logging.
	// logs to the normal place that you log to i.e. viewable with `lazygit --logs`
	Log(message string)
	// logs in the actual UI (in the commands panel)
	LogUI(message string)
	CheckedOutRef() *models.Branch
}

func (self *fakeGuiDriver) PressKey(key string) {
	self.pressedKeys = append(self.pressedKeys, key)
}

func (self *fakeGuiDriver) Keys() config.KeybindingConfig {
	return config.KeybindingConfig{}
}

func (self *fakeGuiDriver) CurrentContext() types.Context {
	return nil
}

func (self *fakeGuiDriver) Model() *types.Model {
	return &types.Model{Commits: []*models.Commit{}}
}

func (self *fakeGuiDriver) Fail(message string) {
	self.failureMessage = message
}

func (self *fakeGuiDriver) Log(message string) {
}

func (self *fakeGuiDriver) LogUI(message string) {
}

func (self *fakeGuiDriver) CheckedOutRef() *models.Branch {
	return nil
}

func TestAssertionFailure(t *testing.T) {
	test := NewIntegrationTest(NewIntegrationTestArgs{
		Description: unitTestDescription,
		Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
			input.PressKeys("a")
			input.PressKeys("b")
			assert.CommitCount(2)
		},
	})
	driver := &fakeGuiDriver{}
	test.Run(driver)
	assert.EqualValues(t, []string{"a", "b"}, driver.pressedKeys)
	assert.Equal(t, "Expected 2 commits present, but got 0", driver.failureMessage)
}

func TestManualFailure(t *testing.T) {
	test := NewIntegrationTest(NewIntegrationTestArgs{
		Description: unitTestDescription,
		Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
			assert.Fail("blah")
		},
	})
	driver := &fakeGuiDriver{}
	test.Run(driver)
	assert.Equal(t, "blah", driver.failureMessage)
}

func TestSuccess(t *testing.T) {
	test := NewIntegrationTest(NewIntegrationTestArgs{
		Description: unitTestDescription,
		Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
			input.PressKeys("a")
			input.PressKeys("b")
			assert.CommitCount(0)
		},
	})
	driver := &fakeGuiDriver{}
	test.Run(driver)
	assert.EqualValues(t, []string{"a", "b"}, driver.pressedKeys)
	assert.Equal(t, "", driver.failureMessage)
}

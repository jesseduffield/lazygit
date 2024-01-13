package components

import (
	"testing"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
	"github.com/stretchr/testify/assert"
)

// this file is for testing our test code (meta, I know)

type coordinate struct {
	x, y int
}

type fakeGuiDriver struct {
	failureMessage     string
	pressedKeys        []string
	clickedCoordinates []coordinate
}

var _ integrationTypes.GuiDriver = &fakeGuiDriver{}

func (self *fakeGuiDriver) PressKey(key string) {
	self.pressedKeys = append(self.pressedKeys, key)
}

func (self *fakeGuiDriver) Click(x, y int) {
	self.clickedCoordinates = append(self.clickedCoordinates, coordinate{x: x, y: y})
}

func (self *fakeGuiDriver) Keys() config.KeybindingConfig {
	return config.KeybindingConfig{}
}

func (self *fakeGuiDriver) CurrentContext() types.Context {
	return nil
}

func (self *fakeGuiDriver) ContextForView(viewName string) types.Context {
	return nil
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

func (self *fakeGuiDriver) MainView() *gocui.View {
	return nil
}

func (self *fakeGuiDriver) SecondaryView() *gocui.View {
	return nil
}

func (self *fakeGuiDriver) View(viewName string) *gocui.View {
	return nil
}

func (self *fakeGuiDriver) SetCaption(string) {
}

func (self *fakeGuiDriver) SetCaptionPrefix(string) {
}

func (self *fakeGuiDriver) NextToast() *string {
	return nil
}

func (self *fakeGuiDriver) CheckAllToastsAcknowledged() {}

func TestManualFailure(t *testing.T) {
	test := NewIntegrationTest(NewIntegrationTestArgs{
		Description: unitTestDescription,
		Run: func(t *TestDriver, keys config.KeybindingConfig) {
			t.Fail("blah")
		},
	})
	driver := &fakeGuiDriver{}
	test.Run(driver)
	assert.Equal(t, "blah", driver.failureMessage)
}

func TestSuccess(t *testing.T) {
	test := NewIntegrationTest(NewIntegrationTestArgs{
		Description: unitTestDescription,
		Run: func(t *TestDriver, keys config.KeybindingConfig) {
			t.press("a")
			t.press("b")
			t.click(0, 1)
			t.click(2, 3)
		},
	})
	driver := &fakeGuiDriver{}
	test.Run(driver)
	assert.EqualValues(t, []string{"a", "b"}, driver.pressedKeys)
	assert.EqualValues(t, []coordinate{{0, 1}, {2, 3}}, driver.clickedCoordinates)
	assert.Equal(t, "", driver.failureMessage)
}

func TestGitVersionRestriction(t *testing.T) {
	scenarios := []struct {
		testName          string
		gitVersion        GitVersionRestriction
		expectedShouldRun bool
	}{
		{
			testName:          "AtLeast, current is newer",
			gitVersion:        AtLeast("2.24.9"),
			expectedShouldRun: true,
		},
		{
			testName:          "AtLeast, current is same",
			gitVersion:        AtLeast("2.25.0"),
			expectedShouldRun: true,
		},
		{
			testName:          "AtLeast, current is older",
			gitVersion:        AtLeast("2.26.0"),
			expectedShouldRun: false,
		},
		{
			testName:          "Before, current is older",
			gitVersion:        Before("2.24.9"),
			expectedShouldRun: false,
		},
		{
			testName:          "Before, current is same",
			gitVersion:        Before("2.25.0"),
			expectedShouldRun: false,
		},
		{
			testName:          "Before, current is newer",
			gitVersion:        Before("2.26.0"),
			expectedShouldRun: true,
		},
		{
			testName:          "Includes, current is included",
			gitVersion:        Includes("2.23.0", "2.25.0"),
			expectedShouldRun: true,
		},
		{
			testName:          "Includes, current is not included",
			gitVersion:        Includes("2.23.0", "2.27.0"),
			expectedShouldRun: false,
		},
	}

	currentGitVersion := git_commands.GitVersion{Major: 2, Minor: 25, Patch: 0}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			test := NewIntegrationTest(NewIntegrationTestArgs{
				Description: unitTestDescription,
				GitVersion:  s.gitVersion,
			})
			shouldRun := test.ShouldRunForGitVersion(&currentGitVersion)
			assert.Equal(t, shouldRun, s.expectedShouldRun)
		})
	}
}

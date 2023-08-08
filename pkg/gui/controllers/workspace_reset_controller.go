package controllers

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// this is in its own file given that the workspace controller file is already quite long

func (self *FilesController) createResetMenu() error {
	red := style.FgRed

	nukeStr := "git reset --hard HEAD && git clean -fd"
	if len(self.c.Model().Submodules) > 0 {
		nukeStr = fmt.Sprintf("%s (%s)", nukeStr, self.c.Tr.AndResetSubmodules)
	}

	menuItems := []*types.MenuItem{
		{
			LabelColumns: []string{
				self.c.Tr.DiscardAllChangesToAllFiles,
				red.Sprint(nukeStr),
			},
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.NukeWorkingTree)
				if err := self.c.Git().WorkingTree.ResetAndClean(); err != nil {
					return self.c.Error(err)
				}

				if self.c.UserConfig.Gui.AnimateExplosion {
					self.animateExplosion()
				}

				return self.c.Refresh(
					types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}},
				)
			},
			Key:     'x',
			Tooltip: self.c.Tr.NukeDescription,
		},
		{
			LabelColumns: []string{
				self.c.Tr.DiscardAnyUnstagedChanges,
				red.Sprint("git checkout -- ."),
			},
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.DiscardUnstagedFileChanges)
				if err := self.c.Git().WorkingTree.DiscardAnyUnstagedFileChanges(); err != nil {
					return self.c.Error(err)
				}

				return self.c.Refresh(
					types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}},
				)
			},
			Key: 'u',
		},
		{
			LabelColumns: []string{
				self.c.Tr.DiscardUntrackedFiles,
				red.Sprint("git clean -fd"),
			},
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.RemoveUntrackedFiles)
				if err := self.c.Git().WorkingTree.RemoveUntrackedFiles(); err != nil {
					return self.c.Error(err)
				}

				return self.c.Refresh(
					types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}},
				)
			},
			Key: 'c',
		},
		{
			LabelColumns: []string{
				self.c.Tr.DiscardStagedChanges,
				red.Sprint("stash staged and drop stash"),
			},
			Tooltip: self.c.Tr.DiscardStagedChangesDescription,
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.RemoveStagedFiles)
				if !self.c.Helpers().WorkingTree.IsWorkingTreeDirty() {
					return self.c.ErrorMsg(self.c.Tr.NoTrackedStagedFilesStash)
				}
				if err := self.c.Git().Stash.SaveStagedChanges("[lazygit] tmp stash"); err != nil {
					return self.c.Error(err)
				}
				if err := self.c.Git().Stash.DropNewest(); err != nil {
					return self.c.Error(err)
				}

				return self.c.Refresh(
					types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}},
				)
			},
			Key: 'S',
		},
		{
			LabelColumns: []string{
				self.c.Tr.SoftReset,
				red.Sprint("git reset --soft HEAD"),
			},
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.SoftReset)
				if err := self.c.Git().WorkingTree.ResetSoft("HEAD"); err != nil {
					return self.c.Error(err)
				}

				return self.c.Refresh(
					types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}},
				)
			},
			Key: 's',
		},
		{
			LabelColumns: []string{
				"mixed reset",
				red.Sprint("git reset --mixed HEAD"),
			},
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.MixedReset)
				if err := self.c.Git().WorkingTree.ResetMixed("HEAD"); err != nil {
					return self.c.Error(err)
				}

				return self.c.Refresh(
					types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}},
				)
			},
			Key: 'm',
		},
		{
			LabelColumns: []string{
				self.c.Tr.HardReset,
				red.Sprint("git reset --hard HEAD"),
			},
			OnPress: func() error {
				self.c.LogAction(self.c.Tr.Actions.HardReset)
				if err := self.c.Git().WorkingTree.ResetHard("HEAD"); err != nil {
					return self.c.Error(err)
				}

				return self.c.Refresh(
					types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES}},
				)
			},
			Key: 'h',
		},
	}

	return self.c.Menu(types.CreateMenuOptions{Title: "", Items: menuItems})
}

func (self *FilesController) animateExplosion() {
	self.Explode(self.c.Views().Files, func() {
		err := self.c.PostRefreshUpdate(self.c.Contexts().Files)
		if err != nil {
			self.c.Log.Error(err)
		}
	})
}

// Animates an explosion within the view by drawing a bunch of flamey characters
func (self *FilesController) Explode(v *gocui.View, onDone func()) {
	width := v.InnerWidth()
	height := v.InnerHeight() + 1
	styles := []style.TextStyle{
		style.FgLightWhite.SetBold(),
		style.FgYellow.SetBold(),
		style.FgRed.SetBold(),
		style.FgBlue.SetBold(),
		style.FgBlack.SetBold(),
	}

	self.c.OnWorker(func(_ gocui.Task) {
		max := 25
		for i := 0; i < max; i++ {
			image := getExplodeImage(width, height, i, max)
			style := styles[(i*len(styles)/max)%len(styles)]
			coloredImage := style.Sprint(image)
			self.c.OnUIThread(func() error {
				_ = v.SetOrigin(0, 0)
				v.SetContent(coloredImage)
				return nil
			})
			time.Sleep(time.Millisecond * 20)
		}
		self.c.OnUIThread(func() error {
			v.Clear()
			onDone()
			return nil
		})
	})
}

// Render an explosion in the given bounds.
func getExplodeImage(width int, height int, frame int, max int) string {
	// Predefine the explosion symbols
	explosionChars := []rune{'*', '.', '@', '#', '&', '+', '%'}

	// Initialize a buffer to build our string
	var buf bytes.Buffer

	// Initialize RNG seed
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	// calculate the center of explosion
	centerX, centerY := width/2, height/2

	// calculate the max radius (hypotenuse of the view)
	maxRadius := math.Hypot(float64(centerX), float64(centerY))

	// calculate frame as a proportion of max, apply square root to create the non-linear effect
	progress := math.Sqrt(float64(frame) / float64(max))

	// calculate radius of explosion according to frame and max
	radius := progress * maxRadius * 2

	// introduce a new radius for the inner boundary of the explosion (the shockwave effect)
	var innerRadius float64
	if progress > 0.5 {
		innerRadius = (progress - 0.5) * 2 * maxRadius
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// calculate distance from center, scale x by 2 to compensate for character aspect ratio
			distance := math.Hypot(float64(x-centerX), float64(y-centerY)*2)

			// if distance is less than radius and greater than innerRadius, draw explosion char
			if distance <= radius && distance >= innerRadius {
				// Make placement random and less likely as explosion progresses
				if random.Float64() > progress {
					// Pick a random explosion char
					char := explosionChars[random.Intn(len(explosionChars))]
					buf.WriteRune(char)
				} else {
					buf.WriteRune(' ')
				}
			} else {
				// If not explosion, then it's empty space
				buf.WriteRune(' ')
			}
		}
		// End of line
		if y < height-1 {
			buf.WriteRune('\n')
		}
	}

	return buf.String()
}

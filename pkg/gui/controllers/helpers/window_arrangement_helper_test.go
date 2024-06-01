package helpers

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jesseduffield/lazycore/pkg/boxlayout"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

// The best way to add test cases here is to set your args and then get the
// test to fail and copy+paste the output into the test case's expected string.
// TODO: add more test cases
func TestGetWindowDimensions(t *testing.T) {
	getDefaultArgs := func() WindowArrangementArgs {
		return WindowArrangementArgs{
			Width:               75,
			Height:              30,
			UserConfig:          config.GetDefaultConfig(),
			CurrentWindow:       "files",
			CurrentSideWindow:   "files",
			CurrentStaticWindow: "files",
			SplitMainPanel:      false,
			ScreenMode:          types.SCREEN_NORMAL,
			AppStatus:           "",
			InformationStr:      "information",
			ShowExtrasWindow:    false,
			InDemo:              false,
			IsAnyModeActive:     false,
			InSearchPrompt:      false,
			SearchPrefix:        "",
		}
	}

	type Test struct {
		name       string
		mutateArgs func(*WindowArrangementArgs)
		expected   string
	}

	tests := []Test{
		{
			name:       "default",
			mutateArgs: func(args *WindowArrangementArgs) {},
			expected: `
			╭status─────────────────╮╭main────────────────────────────────────────────╮
			│                       ││                                                │
			╰───────────────────────╯│                                                │
			╭files──────────────────╮│                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			╰───────────────────────╯│                                                │
			╭branches───────────────╮│                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			╰───────────────────────╯│                                                │
			╭commits────────────────╮│                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			╰───────────────────────╯│                                                │
			╭stash──────────────────╮│                                                │
			│                       ││                                                │
			╰───────────────────────╯╰────────────────────────────────────────────────╯
			<options──────────────────────────────────────────────────────>A<B────────>
			A: statusSpacer1
			B: information
			`,
		},
		{
			name: "stash focused",
			mutateArgs: func(args *WindowArrangementArgs) {
				args.CurrentSideWindow = "stash"
			},
			expected: `
			╭status─────────────────╮╭main────────────────────────────────────────────╮
			│                       ││                                                │
			╰───────────────────────╯│                                                │
			╭files──────────────────╮│                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			╰───────────────────────╯│                                                │
			╭branches───────────────╮│                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			╰───────────────────────╯│                                                │
			╭commits────────────────╮│                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			╰───────────────────────╯│                                                │
			╭stash──────────────────╮│                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			╰───────────────────────╯╰────────────────────────────────────────────────╯
			<options──────────────────────────────────────────────────────>A<B────────>
			A: statusSpacer1
			B: information
			`,
		},
		{
			name: "expandFocusedSidePanel",
			mutateArgs: func(args *WindowArrangementArgs) {
				args.UserConfig.Gui.ExpandFocusedSidePanel = true
			},
			expected: `
			╭status─────────────────╮╭main────────────────────────────────────────────╮
			│                       ││                                                │
			╰───────────────────────╯│                                                │
			╭files──────────────────╮│                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			╰───────────────────────╯│                                                │
			╭branches───────────────╮│                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			╰───────────────────────╯│                                                │
			╭commits────────────────╮│                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			╰───────────────────────╯│                                                │
			╭stash──────────────────╮│                                                │
			│                       ││                                                │
			╰───────────────────────╯╰────────────────────────────────────────────────╯
			<options──────────────────────────────────────────────────────>A<B────────>
			A: statusSpacer1
			B: information
			`,
		},
		{
			name: "expandSidePanelWeight",
			mutateArgs: func(args *WindowArrangementArgs) {
				args.UserConfig.Gui.ExpandFocusedSidePanel = true
				args.UserConfig.Gui.ExpandedSidePanelWeight = 4
			},
			expected: `
			╭status─────────────────╮╭main────────────────────────────────────────────╮
			│                       ││                                                │
			╰───────────────────────╯│                                                │
			╭files──────────────────╮│                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			│                       ││                                                │
			╰───────────────────────╯│                                                │
			╭branches───────────────╮│                                                │
			│                       ││                                                │
			│                       ││                                                │
			╰───────────────────────╯│                                                │
			╭commits────────────────╮│                                                │
			│                       ││                                                │
			│                       ││                                                │
			╰───────────────────────╯│                                                │
			╭stash──────────────────╮│                                                │
			│                       ││                                                │
			╰───────────────────────╯╰────────────────────────────────────────────────╯
			<options──────────────────────────────────────────────────────>A<B────────>
			A: statusSpacer1
			B: information
			`,
		},
		{
			name: "half screen mode, enlargedSideViewLocation left",
			mutateArgs: func(args *WindowArrangementArgs) {
				args.Height = 20 // smaller height because we don't more here
				args.ScreenMode = types.SCREEN_HALF
				args.UserConfig.Gui.EnlargedSideViewLocation = "left"
			},
			expected: `
			╭status──────────────────────────────╮╭main───────────────────────────────╮
			│                                    ││                                   │
			│                                    ││                                   │
			│                                    ││                                   │
			│                                    ││                                   │
			│                                    ││                                   │
			│                                    ││                                   │
			│                                    ││                                   │
			│                                    ││                                   │
			│                                    ││                                   │
			│                                    ││                                   │
			│                                    ││                                   │
			│                                    ││                                   │
			│                                    ││                                   │
			│                                    ││                                   │
			│                                    ││                                   │
			│                                    ││                                   │
			│                                    ││                                   │
			╰────────────────────────────────────╯╰───────────────────────────────────╯
			<options──────────────────────────────────────────────────────>A<B────────>
			A: statusSpacer1
			B: information
			`,
		},
		{
			name: "half screen mode, enlargedSideViewLocation top",
			mutateArgs: func(args *WindowArrangementArgs) {
				args.Height = 20 // smaller height because we don't more here
				args.ScreenMode = types.SCREEN_HALF
				args.UserConfig.Gui.EnlargedSideViewLocation = "top"
			},
			expected: `
			╭status───────────────────────────────────────────────────────────────────╮
			│                                                                         │
			│                                                                         │
			│                                                                         │
			│                                                                         │
			│                                                                         │
			╰─────────────────────────────────────────────────────────────────────────╯
			╭main─────────────────────────────────────────────────────────────────────╮
			│                                                                         │
			│                                                                         │
			│                                                                         │
			│                                                                         │
			│                                                                         │
			│                                                                         │
			│                                                                         │
			│                                                                         │
			│                                                                         │
			│                                                                         │
			╰─────────────────────────────────────────────────────────────────────────╯
			<options──────────────────────────────────────────────────────>A<B────────>
			A: statusSpacer1
			B: information
			`,
		},
		{
			name: "search mode",
			mutateArgs: func(args *WindowArrangementArgs) {
				args.InSearchPrompt = true
				args.SearchPrefix = "Search: "
				args.Height = 6 // small height cos we only care about the bottom line
			},
			expected: `
			<status─────────────────>╭main────────────────────────────────────────────╮
			<files──────────────────>│                                                │
			<branches───────────────>│                                                │
			<commits────────────────>│                                                │
			<stash──────────────────>╰────────────────────────────────────────────────╯
			<A─────><search───────────────────────────────────────────────────────────>
			A: searchPrefix
			`,
		},
		{
			name: "app status present",
			mutateArgs: func(args *WindowArrangementArgs) {
				args.AppStatus = "Rebasing /"
				args.Height = 6 // small height cos we only care about the bottom line
			},
			// We expect single-character spacers between the windows of the bottom line
			expected: `
			<status─────────────────>╭main────────────────────────────────────────────╮
			<files──────────────────>│                                                │
			<branches───────────────>│                                                │
			<commits────────────────>│                                                │
			<stash──────────────────>╰────────────────────────────────────────────────╯
			<A───────>B<options───────────────────────────────────────────>C<D────────>
			A: appStatus
			B: statusSpacer2
			C: statusSpacer1
			D: information
			`,
		},
		{
			name: "information present without options",
			mutateArgs: func(args *WindowArrangementArgs) {
				args.Height = 6                            // small height cos we only care about the bottom line
				args.UserConfig.Gui.ShowBottomLine = false // this hides the options window
				args.IsAnyModeActive = true                // this means we show the bottom line despite the user config
			},
			// We expect a spacer on the left of the bottom line so that the information
			// window is right-aligned
			expected: `
			<status─────────────────>╭main────────────────────────────────────────────╮
			<files──────────────────>│                                                │
			<branches───────────────>│                                                │
			<commits────────────────>│                                                │
			<stash──────────────────>╰────────────────────────────────────────────────╯
			<statusSpacer1────────────────────────────────────────────────>A<B────────>
			A: statusSpacer2
			B: information
			`,
		},
		{
			name: "app status present without information or options",
			mutateArgs: func(args *WindowArrangementArgs) {
				args.Height = 6                            // small height cos we only care about the bottom line
				args.UserConfig.Gui.ShowBottomLine = false // this hides the options window
				args.IsAnyModeActive = false
				args.AppStatus = "Rebasing /"
			},
			// We expect the app status window to take up all the available space
			expected: `
			<status─────────────────>╭main────────────────────────────────────────────╮
			<files──────────────────>│                                                │
			<branches───────────────>│                                                │
			<commits────────────────>│                                                │
			<stash──────────────────>╰────────────────────────────────────────────────╯
			<appStatus────────────────────────────────────────────────────────────────>
			`,
		},
		{
			name: "app status present with information but without options",
			mutateArgs: func(args *WindowArrangementArgs) {
				args.Height = 6                            // small height cos we only care about the bottom line
				args.UserConfig.Gui.ShowBottomLine = false // this hides the options window
				args.IsAnyModeActive = true
				args.AppStatus = "Rebasing /"
			},
			expected: `
			<status─────────────────>╭main────────────────────────────────────────────╮
			<files──────────────────>│                                                │
			<branches───────────────>│                                                │
			<commits────────────────>│                                                │
			<stash──────────────────>╰────────────────────────────────────────────────╯
			<A───────><statusSpacer1──────────────────────────────────────>B<C────────>
			A: appStatus
			B: statusSpacer2
			C: information
			`,
		},
		{
			name: "app status present with very long information but without options",
			mutateArgs: func(args *WindowArrangementArgs) {
				args.Height = 6                            // small height cos we only care about the bottom line
				args.Width = 55                            // smaller width so that not all bottom line views fit
				args.UserConfig.Gui.ShowBottomLine = false // this hides the options window
				args.IsAnyModeActive = true
				args.AppStatus = "Rebasing /"
				args.InformationStr = "Showing output for: git diff deadbeef fa1afe1 -- (Reset)"
			},
			expected: `
			<status───────────>╭main──────────────────────────────╮
			<files────────────>│                                  │
			<branches─────────>│                                  │
			<commits──────────>│                                  │
			<stash────────────>╰──────────────────────────────────╯
			<A───────>B<information──────────────────────────────────────────>
			A: appStatus
			B: statusSpacer2
			`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := getDefaultArgs()
			test.mutateArgs(&args)
			windows := GetWindowDimensions(args)
			output := renderLayout(windows)
			// removing tabs so that it's easier to paste the expected output
			expected := strings.ReplaceAll(test.expected, "\t", "")
			expected = strings.TrimSpace(expected)
			if output != expected {
				fmt.Println(output)
				t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, output)
			}
		})
	}
}

func renderLayout(windows map[string]boxlayout.Dimensions) string {
	// Each window will be represented by a letter.
	windowMarkers := map[string]string{}
	shortLabels := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	currentShortLabelIdx := 0
	windowNames := lo.Keys(windows)
	// Sort first by name, then by position. This means our short labels will
	// increment in the order that the windows appear on the screen.
	slices.Sort(windowNames)
	slices.SortStableFunc(windowNames, func(a, b string) bool {
		dimensionsA := windows[a]
		dimensionsB := windows[b]
		if dimensionsA.Y0 < dimensionsB.Y0 {
			return true
		}
		if dimensionsA.Y0 > dimensionsB.Y0 {
			return false
		}
		return dimensionsA.X0 < dimensionsB.X0
	})

	// Uniquefy windows by dimensions (so perfectly overlapping windows are de-duped). This prevents getting 'fileshes' as a label where the files and branches windows overlap.
	// branches windows overlap.
	windowNames = lo.UniqBy(windowNames, func(windowName string) boxlayout.Dimensions {
		return windows[windowName]
	})

	// excluding the limit window because it overlaps with everything. In future
	// we should have a concept of layers and then our test can assert against
	// each layer.
	windowNames = lo.Without(windowNames, "limit")

	// get width/height by getting the max values of the dimensions
	width := 0
	height := 0
	for _, dimensions := range windows {
		if dimensions.X1+1 > width {
			width = dimensions.X1 + 1
		}
		if dimensions.Y1+1 > height {
			height = dimensions.Y1 + 1
		}
	}

	screen := make([][]string, height)
	for i := range screen {
		screen[i] = make([]string, width)
	}

	// Draw each window
	for _, windowName := range windowNames {
		dimensions := windows[windowName]

		zeroWidth := dimensions.X0 == dimensions.X1+1
		if zeroWidth {
			continue
		}

		singleRow := dimensions.Y0 == dimensions.Y1
		oneOrTwoColumns := dimensions.X0 == dimensions.X1 || dimensions.X0+1 == dimensions.X1

		assignShortLabel := func(windowName string) string {
			windowMarkers[windowName] = shortLabels[currentShortLabelIdx]
			currentShortLabelIdx++
			return windowMarkers[windowName]
		}

		if singleRow {
			y := dimensions.Y0
			// If our window only occupies one (or two) columns we'll just use the short
			// label once (or twice) i.e. 'A' or 'AA'.
			if oneOrTwoColumns {
				shortLabel := assignShortLabel(windowName)

				for x := dimensions.X0; x <= dimensions.X1; x++ {
					screen[y][x] = shortLabel
				}
			} else {
				screen[y][dimensions.X0] = "<"
				screen[y][dimensions.X1] = ">"
				for x := dimensions.X0 + 1; x < dimensions.X1; x++ {
					screen[y][x] = "─"
				}

				// Now add the label
				label := windowName
				// If we can't fit the label we'll use a one-character short label
				if len(label) > dimensions.X1-dimensions.X0-1 {
					label = assignShortLabel(windowName)
				}
				for i, char := range label {
					screen[y][dimensions.X0+1+i] = string(char)
				}
			}
		} else {
			// Draw box border
			for y := dimensions.Y0; y <= dimensions.Y1; y++ {
				for x := dimensions.X0; x <= dimensions.X1; x++ {
					if x == dimensions.X0 && y == dimensions.Y0 {
						screen[y][x] = "╭"
					} else if x == dimensions.X1 && y == dimensions.Y0 {
						screen[y][x] = "╮"
					} else if x == dimensions.X0 && y == dimensions.Y1 {
						screen[y][x] = "╰"
					} else if x == dimensions.X1 && y == dimensions.Y1 {
						screen[y][x] = "╯"
					} else if y == dimensions.Y0 || y == dimensions.Y1 {
						screen[y][x] = "─"
					} else if x == dimensions.X0 || x == dimensions.X1 {
						screen[y][x] = "│"
					} else {
						screen[y][x] = " "
					}
				}
			}

			// Add the label
			label := windowName
			// If we can't fit the label we'll use a one-character short label
			if len(label) > dimensions.X1-dimensions.X0-1 {
				label = assignShortLabel(windowName)
			}
			for i, char := range label {
				screen[dimensions.Y0][dimensions.X0+1+i] = string(char)
			}
		}
	}

	// Draw the screen
	output := ""
	for _, row := range screen {
		for _, marker := range row {
			output += marker
		}
		output += "\n"
	}

	// Add a legend
	for _, windowName := range windowNames {
		if !lo.Contains(lo.Keys(windowMarkers), windowName) {
			continue
		}
		marker := windowMarkers[windowName]
		output += fmt.Sprintf("%s: %s\n", marker, windowName)
	}

	output = strings.TrimSpace(output)

	return output
}

// author_colors_repo creates a git repository for checking that the colors
// lazygit gives to authors are readable.
//
// If gui.authorColors names no color for an author, lazygit derives one from a
// hash of their name. The authors of an ordinary repository rarely land near the
// edges of the range that this color is picked from. Every commit in the
// repository created here is by an author at one of those edges: the lowest or
// highest lightness, combined with the lowest or highest saturation, at twelve
// hues around the color wheel. The commit subject says which edge it is.
//
// Usage:
//
//	go run ./cmd/author_colors_repo <path>
//
// Then open the repository with lazygit in each terminal theme you want to check.
package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"time"

	"github.com/jesseduffield/lazygit/pkg/gui/presentation/authors"
)

type extreme int

const (
	lowest extreme = iota
	highest
)

func (self extreme) String() string {
	if self == lowest {
		return "min"
	}
	return "max"
}

func (self extreme) matches(fraction float64) bool {
	if self == lowest {
		return fraction < 0.01
	}
	return fraction >= 0.99
}

const (
	numHues = 12

	// Small enough that the windows of neighbouring hues don't overlap
	hueTolerance = 0.02
)

type commit struct {
	author  string
	subject string
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <path>", os.Args[0])
	}
	path := os.Args[1]
	if _, err := os.Stat(path); err == nil {
		log.Fatalf("%s already exists", path)
	}

	extremes := []extreme{lowest, highest}
	commits := make([]commit, 0, len(extremes)*len(extremes)*numHues)
	for _, lightness := range extremes {
		for _, saturation := range extremes {
			for i := range numHues {
				hue := float64(i) / numHues
				author, actualHue := findAuthor(lightness, saturation, hue)
				subject := fmt.Sprintf("lightness %s, saturation %s, hue %.0f°", lightness, saturation, actualHue*360)
				commits = append(commits, commit{author: author, subject: subject})
			}
		}
	}

	if err := runGit("", nil, "init", "-q", path); err != nil {
		log.Fatal(err)
	}

	// Commit in reverse, so that the commits panel lists them in the order above
	startTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := range commits {
		c := commits[len(commits)-1-i]
		date := fmt.Sprintf("%d +0000", startTime.Add(time.Duration(i)*time.Hour).Unix())
		env := []string{
			"GIT_AUTHOR_NAME=" + c.author,
			"GIT_AUTHOR_EMAIL=author@example.com",
			"GIT_AUTHOR_DATE=" + date,
			"GIT_COMMITTER_NAME=" + c.author,
			"GIT_COMMITTER_EMAIL=author@example.com",
			"GIT_COMMITTER_DATE=" + date,
		}
		if err := runGit(path, env, "commit", "-q", "--allow-empty", "-m", c.subject); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("Created %s with %d commits\n", path, len(commits))
}

// findAuthor returns the first name of the form "L<lightness> S<saturation> <n>"
// whose color lies at the given extremes of lightness and saturation, and close
// to the given hue. It also returns the hue that the name lands on.
//
// Every name of this form has the same initials, so the color is the only
// thing that differs between authors in the commits panel.
func findAuthor(lightness extreme, saturation extreme, hue float64) (string, float64) {
	for n := 1; ; n++ {
		name := fmt.Sprintf("L%s S%s %d", lightness, saturation, n)
		h, s, l := authors.ColorPosition(name)
		if lightness.matches(l) && saturation.matches(s) && hueDistance(h, hue) <= hueTolerance {
			return name, h
		}
	}
}

// hueDistance is the distance between two hues on the color wheel, where each
// hue is a fraction of a full turn.
func hueDistance(a float64, b float64) float64 {
	d := math.Abs(a - b)
	return math.Min(d, 1-d)
}

func runGit(dir string, env []string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	// Keep the user's git config out of it, so that hooks, commit signing and
	// the like don't apply, and every run creates the same commits
	cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL="+os.DevNull, "GIT_CONFIG_NOSYSTEM=1")
	cmd.Env = append(cmd.Env, env...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

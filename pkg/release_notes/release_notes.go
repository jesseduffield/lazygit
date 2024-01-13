package release_notes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Represents the JSON structure of a GitHub release
type Release struct {
	Name        string `json:"name"`
	PublishedAt string `json:"published_at"`
	Body        string `json:"body"`
}

// Fetches the last 5 releases from the lazygit repository and returns the information as a string
func GetLazyGitReleases() (string, error) {
	url := "https://api.github.com/repos/jesseduffield/lazygit/releases?per_page=5"

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var releases []Release
	err = json.Unmarshal(body, &releases)
	if err != nil {
		return "", err
	}

	var sb strings.Builder

	sb.WriteString("_These are the latest release notes for lazygit. They are pulled straight from GitHub so if you're not on the latest version, you can see what you're missing out on!_\n\n")

	for _, release := range releases {
		sb.WriteString(fmt.Sprintf("# Version %s\n\n", release.Name))
		sb.WriteString(fmt.Sprintf("_Released on %s_\n", formatDate(release.PublishedAt)))
		sb.WriteString(fmt.Sprintf("%s\n", release.Body))
		sb.WriteString("--------------------------------------------------\n")
	}

	return sb.String(), nil
}

func formatDate(dateStr string) string {
	parsedDate, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return "Unknown date"
	}

	day := fmt.Sprintf("%d", parsedDate.Day())
	switch day {
	case "1", "21", "31":
		day += "st"
	case "2", "22":
		day += "nd"
	case "3", "23":
		day += "rd"
	default:
		day += "th"
	}

	return parsedDate.Format("January " + day + " 2006")
}

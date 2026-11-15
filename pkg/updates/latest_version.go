package updates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// fetchLatestVersion asks GitHub for the version of the latest lazygit release.
// We use the release page instead of GitHub's REST API, because the API limits
// the number of unauthenticated requests per IP address.
func fetchLatestVersion(ctx context.Context, client *http.Client, repoURL string) (*types.VersionNumber, error) {
	var data struct {
		TagName string `json:"tag_name"`
	}
	if err := fetchJSON(ctx, client, repoURL+"/releases/latest", &data); err != nil {
		return nil, err
	}

	return types.ParseVersionNumber(data.TagName)
}

func fetchJSON(ctx context.Context, client *http.Client, url string, data any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s returned %s", url, resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(data)
}

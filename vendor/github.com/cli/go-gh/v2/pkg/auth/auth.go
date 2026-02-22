// Package auth is a set of functions for retrieving authentication tokens
// and authenticated hosts.
package auth

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cli/go-gh/v2/internal/set"
	"github.com/cli/go-gh/v2/pkg/config"
	"github.com/cli/safeexec"
)

const (
	codespaces            = "CODESPACES"
	defaultSource         = "default"
	ghEnterpriseToken     = "GH_ENTERPRISE_TOKEN"
	ghHost                = "GH_HOST"
	ghToken               = "GH_TOKEN"
	github                = "github.com"
	githubEnterpriseToken = "GITHUB_ENTERPRISE_TOKEN"
	githubToken           = "GITHUB_TOKEN"
	hostsKey              = "hosts"
	localhost             = "github.localhost"
	oauthToken            = "oauth_token"
)

// TokenForHost retrieves an authentication token and the source of that token for the specified
// host. The source can be either an environment variable, configuration file, or the system
// keyring. In the latter case, this shells out to "gh auth token" to obtain the token.
//
// Returns "", "default" if no applicable token is found.
func TokenForHost(host string) (string, string) {
	if token, source := TokenFromEnvOrConfig(host); token != "" {
		return token, source
	}

	ghExe := os.Getenv("GH_PATH")
	if ghExe == "" {
		ghExe, _ = safeexec.LookPath("gh")
	}

	if ghExe != "" {
		if token, source := tokenFromGh(ghExe, host); token != "" {
			return token, source
		}
	}

	return "", defaultSource
}

// TokenFromEnvOrConfig retrieves an authentication token from environment variables or the config
// file as fallback, but does not support reading the token from system keyring. Most consumers
// should use TokenForHost.
func TokenFromEnvOrConfig(host string) (string, string) {
	cfg, _ := config.Read(nil)
	return tokenForHost(cfg, host)
}

func tokenForHost(cfg *config.Config, host string) (string, string) {
	host = normalizeHostname(host)
	if IsEnterprise(host) {
		if token := os.Getenv(ghEnterpriseToken); token != "" {
			return token, ghEnterpriseToken
		}
		if token := os.Getenv(githubEnterpriseToken); token != "" {
			return token, githubEnterpriseToken
		}
		if isCodespaces, _ := strconv.ParseBool(os.Getenv(codespaces)); isCodespaces {
			if token := os.Getenv(githubToken); token != "" {
				return token, githubToken
			}
		}
		if cfg != nil {
			token, _ := cfg.Get([]string{hostsKey, host, oauthToken})
			return token, oauthToken
		}
	}
	if token := os.Getenv(ghToken); token != "" {
		return token, ghToken
	}
	if token := os.Getenv(githubToken); token != "" {
		return token, githubToken
	}
	if cfg != nil {
		token, _ := cfg.Get([]string{hostsKey, host, oauthToken})
		return token, oauthToken
	}
	return "", defaultSource
}

func tokenFromGh(path string, host string) (string, string) {
	cmd := exec.Command(path, "auth", "token", "--secure-storage", "--hostname", host)
	result, err := cmd.Output()
	if err != nil {
		return "", "gh"
	}
	return strings.TrimSpace(string(result)), "gh"
}

// KnownHosts retrieves a list of hosts that have corresponding
// authentication tokens, either from environment variables
// or from the configuration file.
// Returns an empty string slice if no hosts are found.
func KnownHosts() []string {
	cfg, _ := config.Read(nil)
	return knownHosts(cfg)
}

func knownHosts(cfg *config.Config) []string {
	hosts := set.NewStringSet()
	if host := os.Getenv(ghHost); host != "" {
		hosts.Add(host)
	}
	if token, _ := tokenForHost(cfg, github); token != "" {
		hosts.Add(github)
	}
	if cfg != nil {
		keys, err := cfg.Keys([]string{hostsKey})
		if err == nil {
			hosts.AddValues(keys)
		}
	}
	return hosts.ToSlice()
}

// DefaultHost retrieves an authenticated host and the source of host.
// The source can be either an environment variable or from the
// configuration file.
// Returns "github.com", "default" if no viable host is found.
func DefaultHost() (string, string) {
	cfg, _ := config.Read(nil)
	return defaultHost(cfg)
}

func defaultHost(cfg *config.Config) (string, string) {
	if host := os.Getenv(ghHost); host != "" {
		return host, ghHost
	}
	if cfg != nil {
		keys, err := cfg.Keys([]string{hostsKey})
		if err == nil && len(keys) == 1 {
			return keys[0], hostsKey
		}
	}
	return github, defaultSource
}

// TenancyHost is the domain name of a tenancy GitHub instance.
const tenancyHost = "ghe.com"

// IsEnterprise determines if a provided host is a GitHub Enterprise Server instance,
// rather than GitHub.com or a tenancy GitHub instance.
func IsEnterprise(host string) bool {
	normalizedHost := normalizeHostname(host)
	return normalizedHost != github && normalizedHost != localhost && !IsTenancy(normalizedHost)
}

// IsTenancy determines if a provided host is a tenancy GitHub instance,
// rather than GitHub.com or a GitHub Enterprise Server instance.
func IsTenancy(host string) bool {
	normalizedHost := normalizeHostname(host)
	return strings.HasSuffix(normalizedHost, "."+tenancyHost)
}

func normalizeHostname(host string) string {
	hostname := strings.ToLower(host)
	if strings.HasSuffix(hostname, "."+github) {
		return github
	}
	if strings.HasSuffix(hostname, "."+localhost) {
		return localhost
	}
	// This has been copied over from the cli/cli NormalizeHostname function
	// to ensure compatible behaviour but we don't fully understand when or
	// why it would be useful here. We can't see what harm will come of
	// duplicating the logic.
	if before, found := cutSuffix(hostname, "."+tenancyHost); found {
		idx := strings.LastIndex(before, ".")
		return fmt.Sprintf("%s.%s", before[idx+1:], tenancyHost)
	}
	return hostname
}

// Backport strings.CutSuffix from Go 1.20.
func cutSuffix(s, suffix string) (string, bool) {
	if !strings.HasSuffix(s, suffix) {
		return s, false
	}
	return s[:len(s)-len(suffix)], true
}

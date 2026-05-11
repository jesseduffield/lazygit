package components

var RandomCommitMessages = []string{
	`Refactor HTTP client for better error handling`,
	`Integrate pagination in user listings`,
	`Fix incorrect type in updateUser function`,
	`Create initial setup for postgres database`,
	`Add unit tests for authentication service`,
	`Improve efficiency of sorting algorithm in util package`,
	`Resolve intermittent test failure in CartTest`,
	`Introduce cache layer for product images`,
	`Revamp User Interface of the settings page`,
	`Remove deprecated uses of api endpoints`,
	`Ensure proper escaping of SQL queries`,
	`Implement feature flag for dark mode`,
	`Add functionality for users to reset password`,
	`Optimize performance of image loading on home screen`,
	`Correct argument type in the sendEmail function`,
	`Merge feature branch 'add-payment-gateway'`,
	`Add validation to signup form fields`,
	`Refactor User model to include middle name`,
	`Update README with new setup instructions`,
	`Extend session expiry time to 24 hours`,
	`Implement rate limiting on login attempts`,
	`Add sorting feature to product listing page`,
	`Refactor logic in Lazygit Diff view`,
	`Optimize Lazygit startup time`,
	`Fix typos in documentation`,
	`Move global variables to environment config`,
	`Upgrade Rails version to 6.1.4`,
	`Refactor user notifications system`,
	`Implement user blocking functionality`,
	`Improve Dockerfile for more efficient builds`,
	`Introduce Redis for session management`,
	`Ensure CSRF protection for all forms`,
	`Implement bulk delete feature in admin panel`,
	`Harden security of user password storage`,
	`Resolve race condition in transaction handling`,
	`Migrate legacy codebase to Typescript`,
	`Update UX of password reset feature`,
	`Add internationalization support for German`,
	`Enhance logging in production environment`,
	`Remove hardcoded values from payment module`,
	`Introduce retry mechanism in network calls`,
	`Handle edge case for zero quantity in cart`,
	`Revamp error handling in user registration`,
	`Replace deprecated lifecycle methods in React components`,
	`Update styles according to new design guidelines`,
	`Handle database connection failures gracefully`,
	`Ensure atomicity of transactions in payment system`,
	`Refactor session management using JWT`,
	`Enhance user search with fuzzy matching`,
	`Move constants to a separate config file`,
	`Add TypeScript types to User module`,
	`Implement automated backups for database`,
	`Fix broken links on the help page`,
	`Add end-to-end tests for checkout flow`,
	`Add loading indicators to improve UX`,
	`Improve accessibility of site navigation`,
	`Refactor error messages for better clarity`,
	`Enable gzip compression for faster page loads`,
	`Set up CI/CD pipeline using GitHub actions`,
	`Add a user-friendly 404 page`,
	`Implement OAuth login with Google`,
	`Resolve dependency conflicts in package.json`,
	`Add proper alt text to all images for SEO`,
	`Implement comment moderation feature`,
	`Fix double encoding issue in URL parameters`,
	`Resolve flickering issue in animation`,
	`Update dependencies to latest stable versions`,
	`Set proper cache headers for static assets`,
	`Add structured data for better SEO`,
	`Refactor to remove circular dependencies`,
	`Add feature to report inappropriate content`,
	`Implement mobile-friendly navigation menu`,
	`Update privacy policy to comply with GDPR`,
	`Fix memory leak issue in event listeners`,
	`Improve form validation feedback for user`,
	`Implement API versioning`,
	`Improve resilience of system by adding circuit breaker`,
	`Add sitemap.xml for better search engine indexing`,
	`Set up performance monitoring with New Relic`,
	`Introduce service worker for offline support`,
	`Enhance email notifications with HTML templates`,
	`Ensure all pages are responsive across devices`,
	`Create helper functions to reduce code duplication`,
	`Add 'remember me' feature to login`,
	`Increase test coverage for User model`,
	`Refactor error messages into a separate module`,
	`Optimize images for faster loading`,
	`Ensure correct HTTP status codes for all responses`,
	`Implement auto-save feature in post editor`,
	`Update user guide with new screenshots`,
	`Implement load testing using Gatling`,
	`Add keyboard shortcuts for commonly used actions`,
	`Set up staging environment similar to production`,
	`Ensure all forms use POST method for data submission`,
	`Implement soft delete for user accounts`,
	`Add Webpack for asset bundling`,
	`Handle session timeout gracefully`,
	`Remove unused code and libraries`,
	`Integrate support for markdown in user posts`,
	`Fix bug in timezone conversion.`,
}

type RandomFile struct {
	Name    string
	Content string
}

var RandomFiles = []RandomFile{
	{Name: `http_client.go`, Content: `package httpclient`},
	{Name: `user_listings.go`, Content: `package listings`},
	{Name: `user_service.go`, Content: `package service`},
	{Name: `database_setup.sql`, Content: `CREATE TABLE`},
	{Name: `authentication_test.go`, Content: `package auth_test`},
	{Name: `utils/sorting.go`, Content: `package utils`},
	{Name: `tests/cart_test.go`, Content: `package tests`},
	{Name: `cache/product_images.go`, Content: `package cache`},
	{Name: `ui/settings_page.jsx`, Content: `import React`},
	{Name: `api/deprecated_endpoints.go`, Content: `package api`},
	{Name: `db/sql_queries.go`, Content: `package db`},
	{Name: `features/dark_mode.go`, Content: `package features`},
	{Name: `user/password_reset.go`, Content: `package user`},
	{Name: `performance/image_loading.go`, Content: `package performance`},
	{Name: `email/send_email.go`, Content: `package email`},
	{Name: `merge/payment_gateway.go`, Content: `package merge`},
	{Name: `forms/signup_validation.go`, Content: `package forms`},
	{Name: `models/user.go`, Content: `package models`},
	{Name: `README.md`, Content: `# Project`},
	{Name: `config/session.go`, Content: `package config`},
	{Name: `security/rate_limit.go`, Content: `package security`},
	{Name: `product/sort_list.go`, Content: `package product`},
	{Name: `lazygit/diff_view.go`, Content: `package lazygit`},
	{Name: `performance/lazygit.go`, Content: `package performance`},
	{Name: `docs/documentation.go`, Content: `package docs`},
	{Name: `config/global_variables.go`, Content: `package config`},
	{Name: `Gemfile`, Content: `source 'https://rubygems.org'`},
	{Name: `notification/user_notification.go`, Content: `package notification`},
	{Name: `user/blocking.go`, Content: `package user`},
	{Name: `Dockerfile`, Content: `FROM ubuntu:18.04`},
	{Name: `redis/session_manager.go`, Content: `package redis`},
	{Name: `security/csrf_protection.go`, Content: `package security`},
	{Name: `admin/bulk_delete.go`, Content: `package admin`},
	{Name: `security/password_storage.go`, Content: `package security`},
	{Name: `transactions/transaction_handling.go`, Content: `package transactions`},
	{Name: `migrations/typescript_migration.go`, Content: `package migrations`},
	{Name: `ui/password_reset.jsx`, Content: `import React`},
	{Name: `i18n/german.go`, Content: `package i18n`},
	{Name: `logging/production_logging.go`, Content: `package logging`},
	{Name: `payment/hardcoded_values.go`, Content: `package payment`},
	{Name: `network/retry.go`, Content: `package network`},
	{Name: `cart/zero_quantity.go`, Content: `package cart`},
	{Name: `registration/error_handling.go`, Content: `package registration`},
	{Name: `components/deprecated_methods.jsx`, Content: `import React`},
	{Name: `styles/new_guidelines.css`, Content: `.class {}`},
	{Name: `db/connection_failure.go`, Content: `package db`},
	{Name: `payment/transaction_atomicity.go`, Content: `package payment`},
	{Name: `session/jwt_management.go`, Content: `package session`},
	{Name: `search/fuzzy_matching.go`, Content: `package search`},
	{Name: `config/constants.go`, Content: `package config`},
	{Name: `models/user_types.go`, Content: `package models`},
	{Name: `backup/database_backup.go`, Content: `package backup`},
	{Name: `help_page/links.go`, Content: `package help_page`},
	{Name: `tests/checkout_test.sql`, Content: `DELETE ALL TABLES;`},
	{Name: `ui/loading_indicator.jsx`, Content: `import React`},
	{Name: `navigation/site_navigation.go`, Content: `package navigation`},
	{Name: `error/error_messages.go`, Content: `package error`},
	{Name: `performance/gzip_compression.go`, Content: `package performance`},
	{Name: `.github/workflows/ci.yml`, Content: `name: CI`},
	{Name: `pages/404.html`, Content: `<html></html>`},
	{Name: `oauth/google_login.go`, Content: `package oauth`},
	{Name: `package.json`, Content: `{}`},
	{Name: `seo/alt_text.go`, Content: `package seo`},
	{Name: `moderation/comment_moderation.go`, Content: `package moderation`},
	{Name: `url/double_encoding.go`, Content: `package url`},
	{Name: `animation/flickering.go`, Content: `package animation`},
	{Name: `upgrade_dependencies.sh`, Content: `#!/bin/sh`},
	{Name: `security/csrf_protection2.go`, Content: `package security`},
	{Name: `admin/bulk_delete2.go`, Content: `package admin`},
	{Name: `security/password_storage2.go`, Content: `package security`},
	{Name: `transactions/transaction_handling2.go`, Content: `package transactions`},
	{Name: `migrations/typescript_migration2.go`, Content: `package migrations`},
	{Name: `ui/password_reset2.jsx`, Content: `import React`},
	{Name: `i18n/german2.go`, Content: `package i18n`},
	{Name: `logging/production_logging2.go`, Content: `package logging`},
	{Name: `payment/hardcoded_values2.go`, Content: `package payment`},
	{Name: `network/retry2.go`, Content: `package network`},
	{Name: `cart/zero_quantity2.go`, Content: `package cart`},
	{Name: `registration/error_handling2.go`, Content: `package registration`},
	{Name: `components/deprecated_methods2.jsx`, Content: `import React`},
	{Name: `styles/new_guidelines2.css`, Content: `.class {}`},
	{Name: `db/connection_failure2.go`, Content: `package db`},
	{Name: `payment/transaction_atomicity2.go`, Content: `package payment`},
	{Name: `session/jwt_management2.go`, Content: `package session`},
	{Name: `search/fuzzy_matching2.go`, Content: `package search`},
	{Name: `config/constants2.go`, Content: `package config`},
	{Name: `models/user_types2.go`, Content: `package models`},
	{Name: `backup/database_backup2.go`, Content: `package backup`},
	{Name: `help_page/links2.go`, Content: `package help_page`},
	{Name: `tests/checkout_test2.go`, Content: `package tests`},
	{Name: `ui/loading_indicator2.jsx`, Content: `import React`},
	{Name: `navigation/site_navigation2.go`, Content: `package navigation`},
	{Name: `error/error_messages2.go`, Content: `package error`},
	{Name: `performance/gzip_compression2.go`, Content: `package performance`},
	{Name: `.github/workflows/ci2.yml`, Content: `name: CI`},
	{Name: `pages/4042.html`, Content: `<html></html>`},
	{Name: `oauth/google_login2.go`, Content: `package oauth`},
	{Name: `package2.json`, Content: `{}`},
	{Name: `seo/alt_text2.go`, Content: `package seo`},
	{Name: `moderation/comment_moderation2.go`, Content: `package moderation`},
}

var RandomFileContents = []string{
	`package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/fs"
	"os"
	"strings"

	"github.com/samber/lo"
)

func main() {
	code := generateCode()

	formattedCode, err := format.Source(code)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("test_list.go", formattedCode, 0o644); err != nil {
		panic(err)
	}
}
`,
	`
package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazycore/pkg/utils"
	"github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/samber/lo"
)

func GetTests() []*components.IntegrationTest {
	// first we ensure that each test in this directory has actually been added to the above list.
	testCount := 0

	testNamesSet := set.NewFromSlice(lo.Map(
		tests,
		func(test *components.IntegrationTest, _ int) string {
			return test.Name()
		},
	))
}
`,
	`
package components

import (
	"os"
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/config"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

// IntegrationTest describes an integration test that will be run against the lazygit gui.

// our unit tests will use this description to avoid a panic caused by attempting
// to get the test's name via it's file's path.
const unitTestDescription = "test test"

const (
	defaultWidth  = 100
	defaultHeight = 100
)
`,
	`package components

import (
	"fmt"
	"time"

	"github.com/atotto/clipboard"
	"github.com/jesseduffield/lazygit/pkg/config"
	integrationTypes "github.com/jesseduffield/lazygit/pkg/integration/types"
)

type TestDriver struct {
	gui        integrationTypes.GuiDriver
	keys       config.KeybindingConfig
	inputDelay int
	*assertionHelper
	shell *Shell
}

func NewTestDriver(gui integrationTypes.GuiDriver, shell *Shell, keys config.KeybindingConfig, inputDelay int) *TestDriver {
	return &TestDriver{
		gui:             gui,
		keys:            keys,
		inputDelay:      inputDelay,
		assertionHelper: &assertionHelper{gui: gui},
		shell:           shell,
	}
}

// key is something like 'w' or '<space>'. It's best not to pass a direct value,
// but instead to go through the default user config to get a more meaningful key name
func (self *TestDriver) press(keyStr string) {
	self.SetCaption(fmt.Sprintf("Pressing %s", keyStr))
	self.gui.PressKey(keyStr)
	self.Wait(self.inputDelay)
}
`,
	`package updates

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/go-errors/errors"

	"github.com/kardianos/osext"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// Updater checks for updates and does updates
type Updater struct {
	*common.Common
	Config    config.AppConfigurer
	OSCommand *oscommands.OSCommand
}

// Updaterer implements the check and update methods
type Updaterer interface {
	CheckForNewUpdate()
	Update()
}
`,
	`
package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// IsValidEmail checks if an email address is valid
func IsValidEmail(email string) bool {
	// Using a regex pattern to validate email addresses
	// This is a simple example and might not cover all edge cases
	emailPattern := ` + "`" + `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$` + "`" + `
	match, _ := regexp.MatchString(emailPattern, email)
	return match
}
`,
	`
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jesseduffield/lazygit/pkg/utils"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, the current time is: %s", time.Now().Format(time.RFC3339))
	})

	port := 8080
	utils.PrintMessage(fmt.Sprintf("Server is listening on port %d", port))
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
`,
	`
package logging

import (
	"fmt"
	"os"
	"time"
)

// LogMessage represents a log message with its timestamp
type LogMessage struct {
	Timestamp time.Time
	Message   string
}

// Log writes a message to the log file along with a timestamp
func Log(message string) {
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()

	logEntry := LogMessage{
		Timestamp: time.Now(),
		Message:   message,
	}

	logLine := fmt.Sprintf("[%s] %s\n", logEntry.Timestamp.Format("2006-01-02 15:04:05"), logEntry.Message)
	_, err = logFile.WriteString(logLine)
	if err != nil {
		fmt.Println("Error writing to log file:", err)
	}
}
`,
	`
package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// Encrypt encrypts a plaintext using AES-GCM encryption
func Encrypt(key []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
	return append(nonce, ciphertext...), nil
}
`,
}

var RandomBranchNames = []string{
	"hotfix/fix-bug",
	"r-u-fkn-srs",
	"iserlohn-build",
	"hotfix/fezzan-corridor",
	"terra-investigation",
	"quash-rebellion",
	"feature/attack-on-odin",
	"feature/peace-time",
	"feature/repair-brunhild",
	"feature/iserlohn-backdoor",
	"bugfix/resolve-crash",
	"enhancement/improve-performance",
	"experimental/new-feature",
	"release/v1.0.0",
	"release/v2.0.0",
	"chore/update-dependencies",
	"docs/add-readme",
	"refactor/cleanup-code",
	"style/update-css",
	"test/add-unit-tests",
}

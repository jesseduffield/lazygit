package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
)

//counterfeiter:generate . ISyncMgr
type ISyncMgr interface {
	Push(opts PushOpts) (bool, error)
	Fetch(opts FetchOptions) error
	FetchInBackground(opts FetchOptions) error
	FastForward(branchName string, remoteName string, remoteBranchName string) error
	FetchRemote(remoteName string) error
	PushRef(remoteName string, refName string) error
	DeleteRemoteRef(remoteName string, ref string) error
	SetCredentialHandlers(promptUserForCredential func(CredentialKind) string, handleCredentialError func(error))
}

type SyncMgr struct {
	ICommander

	config IGitConfigMgr
	os     oscommands.IOS

	// callbacks to be provided from the gui package
	promptUserForCredential func(CredentialKind) string
	handleCredentialError   func(error)
}

func NewSyncMgr(
	commander ICommander,
	config IGitConfigMgr,
	os oscommands.IOS,
) *SyncMgr {
	return &SyncMgr{
		ICommander: commander,
		config:     config,
		os:         os,
	}
}

func (c *SyncMgr) SetCredentialHandlers(promptUserForCredential func(CredentialKind) string, handleCredentialError func(error)) {
	c.promptUserForCredential = promptUserForCredential
	c.handleCredentialError = handleCredentialError
}

type PushOpts struct {
	Force             bool
	SetUpstream       bool
	DestinationRemote string
	DestinationBranch string
}

func (c *SyncMgr) Push(opts PushOpts) (bool, error) {
	cmdObj := BuildGitCmdObj("push", []string{opts.DestinationRemote, opts.DestinationBranch},
		map[string]bool{
			"--follow-tags":      c.config.GetConfigValue("push.followTags") != "false",
			"--force-with-lease": opts.Force,
			"--set-upstream":     opts.SetUpstream,
		})

	err := c.runCommandWithCredentialsPrompt(cmdObj)

	if isRejectionErr(err) {
		return true, nil
	}

	c.handleCredentialError(err)

	return false, nil
}

func isRejectionErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Updates were rejected")
}

type FetchOptions struct {
	RemoteName string
	BranchName string
}

// Fetch fetch git repo
func (c *SyncMgr) Fetch(opts FetchOptions) error {
	cmdObj := GetFetchCommandObj(opts)

	return c.runCommandWithCredentialsHandling(cmdObj)
}

// FetchInBackground fails if credentials are requested
func (c *SyncMgr) FetchInBackground(opts FetchOptions) error {
	cmdObj := GetFetchCommandObj(opts)

	cmdObj = c.failOnCredentialsRequest(cmdObj)
	return c.Run(cmdObj)
}

func GetFetchCommandObj(opts FetchOptions) ICmdObj {
	return BuildGitCmdObj("fetch", []string{opts.RemoteName, opts.BranchName}, nil)
}

func (c *SyncMgr) FastForward(branchName string, remoteName string, remoteBranchName string) error {
	cmdObj := BuildGitCmdObj("fetch", []string{remoteName, remoteBranchName + ":" + branchName}, nil)
	return c.runCommandWithCredentialsHandling(cmdObj)
}

func (c *SyncMgr) FetchRemote(remoteName string) error {
	cmdObj := BuildGitCmdObj("fetch", []string{remoteName}, nil)
	return c.runCommandWithCredentialsHandling(cmdObj)
}

func (c *SyncMgr) DeleteRemoteRef(remoteName string, ref string) error {
	return c.runCommandWithCredentialsHandling(
		BuildGitCmdObjFromStr(fmt.Sprintf("push %s --delete %s", remoteName, ref)),
	)
}

func (c *SyncMgr) PushRef(remoteName string, ref string) error {
	return c.runCommandWithCredentialsHandling(
		BuildGitCmdObjFromStr(fmt.Sprintf("push %s %s", remoteName, ref)),
	)
}

// runCommandWithCredentialsPrompt detect a username / password / passphrase question in a command
// promptUserForCredential is a function that gets executed when this function detect you need to fillin a password or passphrase
// The promptUserForCredential argument will be "username", "password" or "passphrase" and expects the user's password/passphrase or username back
func (c *SyncMgr) runCommandWithCredentialsPrompt(cmdObj ICmdObj) error {
	ttyText := ""
	err := c.os.RunAndParseWords(cmdObj, func(word string) string {
		ttyText = ttyText + " " + word

		prompts := map[string]CredentialKind{
			`.+'s password:`:                         PASSWORD,
			`Password\s*for\s*'.+':`:                 PASSWORD,
			`Username\s*for\s*'.+':`:                 USERNAME,
			`Enter\s*passphrase\s*for\s*key\s*'.+':`: PASSPHRASE,
		}

		for pattern, askFor := range prompts {
			if match, _ := regexp.MatchString(pattern, ttyText); match {
				ttyText = ""
				return c.promptUserForCredential(askFor)
			}
		}

		return ""
	})

	return err
}

// this goes one step beyond runCommandWithCredentialsPrompt and handles a credential error
func (c *SyncMgr) runCommandWithCredentialsHandling(cmdObj ICmdObj) error {
	err := c.runCommandWithCredentialsPrompt(cmdObj)
	c.handleCredentialError(err)
	return nil
}

func (c *SyncMgr) failOnCredentialsRequest(cmdObj ICmdObj) ICmdObj {
	lazyGitPath := c.os.GetLazygitPath()

	cmdObj.AddEnvVars(
		"LAZYGIT_CLIENT_COMMAND=EXIT_IMMEDIATELY",
		// prevents git from prompting us for input which would freeze the program. Only works for git v2.3+
		"GIT_TERMINAL_PROMPT=0",
		"GIT_ASKPASS="+lazyGitPath,
	)

	return cmdObj
}

# Update Checks

lazygit checks for new releases and tells the user about them. If it was
installed from one of our release archives, it can also install them. The code
is in `pkg/updates` and in `UpdateHelper` (`pkg/gui/controllers/helpers`).

Whether all of this works for users only shows in the release after the one
that contains a change, because only then is there a newer version for the
released code to find. So test changes to this feature by hand before
releasing them, following the test plan below.

## How it works

### Install kinds

`DetectInstallation` decides how lazygit was installed. It looks at the build
source (`-X main.buildSource=...`), the version, and the path of the
executable:

| Kind | How we detect it | What we offer |
|---|---|---|
| Release | build source `binaryRelease`, not in a package manager's folder | install in place, if we can write to the folder |
| Homebrew | build source `Homebrew`, or the path contains `/Cellar/lazygit/` | show the command |
| go install | build source `goInstall` (see `mergeBuildInfo`) | show the command |
| Scoop, winget, Chocolatey | path of the executable | show the command |
| Other | any other build with a release version | tell the user to check their package manager |
| Dev build | the version isn't `X.Y.Z` | nothing; `u` is disabled |

Scoop, winget, Chocolatey, mise and others install our release archives
unchanged, so the build source of their binaries is `binaryRelease`. We must
not update those in place, because the package manager would then report the
wrong version.

### Checks

- At startup, if the last successful check is at least `update.days` ago
  (default 1). Every hour, lazygit tests again whether a check is due.
- With `u` in the status panel.

All install kinds compare with our latest release on GitHub. Package managers
get a new release some time later; the menu tells the user so.

The result goes into the `updates` section of `state.yml`, which all running
instances share.

### Telling the user

- A note in the bottom-right corner ("Update available", or "Restart to use
  vX" after an install). Clicking it opens the update menu.
- The update menu, when the user hasn't used lazygit for a minute and no popup
  is open. It has "Remind me tomorrow" and "Skip this version", and a link to
  the release notes in its first line. Whatever the user chooses, including
  closing the menu, counts as a response, and the menu doesn't appear again for
  a day. Clicking the link doesn't count.
- A message after lazygit exits, for users who never leave lazygit idle. It
  appears at most once a day, and not on a day when the user responded to the
  menu. It doesn't count as a response.
- The menu at startup, if the update has been available for a week and the
  user never responded to the menu for it.

Every idle instance shows the menu. Otherwise an instance that the user forgot
in some terminal tab would get the reminder every day, and the others never.
When the user responds in one instance, the others close their menus within
ten seconds.

With `update.method: background`, lazygit installs a new release without
asking if it can. With `never`, it still checks and shows the note, but doesn't
remind the user.

## Automated tests

The unit tests in `pkg/updates` cover the install kinds, the asset names (against
the names of a real release), the build source in `.goreleaser.yml`, fetching
the latest version, downloading and installing against a test server,
replacing the executable on Unix and Windows, and when checks and reminders are
due. CI runs them on Windows too.

The integration test `status/check_for_update_in_dev_build` checks that `u` is
disabled for dev builds.

## Manual test plan

### Setup

Build a binary that pretends to be an older release. It finds the real latest
release on GitHub and updates to it. The code of the release that it installs
doesn't matter, so this tests exactly what the next release will do.

```sh
mkdir -p /tmp/lg/bin /tmp/lg/config
go build -ldflags='-X main.buildSource=binaryRelease -X main.version=0.60.0' -o /tmp/lg/bin/lazygit
printf 'disableStartupPopups: true\n' > /tmp/lg/config/config.yml
```

Run it with `CONFIG_DIR=/tmp/lg/config /tmp/lg/bin/lazygit -debug` so that it
uses its own `config.yml` and `state.yml`, and watch the log with
`lazygit --logs` in another terminal. Edit the `updates` section of
`/tmp/lg/config/state.yml` to set up a scenario. Setting `lastcheck` to 0 makes
lazygit check at the next start; setting `lastresponse` and `lastexitmessage`
to 0 allows reminders right away.

For the other install kinds, use these build sources: `Homebrew`, `goInstall`,
or `unknown` (Other). Set `os.openLink` in the config to something like
`echo {{link}} >> /tmp/lg/links` to see which links lazygit opens without
opening a browser. On Windows, put the binary into a folder like
`...\scoop\apps\lazygit\0.60.0\` to test the package manager detection.

To test a real `go install` build with the code of your branch, serve the
branch to `go install` as a made-up release version through a local module
proxy. Pick a version that doesn't exist and is older than the latest release.
The separate module cache keeps the made-up version out of your real one.

```sh
V=v0.60.99
P=/tmp/lg/proxy/github.com/jesseduffield/lazygit/@v
mkdir -p $P
git archive --format=zip --prefix=github.com/jesseduffield/lazygit@$V/ -o $P/$V.zip HEAD
git show HEAD:go.mod > $P/$V.mod
echo "{\"Version\":\"$V\"}" > $P/$V.info
echo $V > $P/list
GOMODCACHE=/tmp/lg/modcache GOBIN=/tmp/lg/bin GOFLAGS=-modcacherw \
  GOPROXY="file:///tmp/lg/proxy,https://proxy.golang.org" \
  GONOSUMDB=github.com/jesseduffield/lazygit \
  go install github.com/jesseduffield/lazygit@$V
```

### Release archive in a writable folder

- [ ] At startup, the note says "Update available". The log shows one check.
- [ ] Starting again on the same day shows the note without a check.
- [ ] `u` shows the menu with "Install now". Clicking "Release notes" in its
      first line opens the releases page, and the menu stays open.
- [ ] Installing shows a waiting
      status, then a toast. The note says "Restart to use vX". `lazygit
      --version` shows the new version, and no temporary files are left in the
      folder.
- [ ] `u` after installing shows the toast again instead of the menu.
- [ ] Clicking the note opens the menu.
- [ ] Without input, the menu appears after a minute, not earlier. It doesn't
      appear while another popup is open, or while lazygit runs a
      subprocess (e.g. the editor).
- [ ] Esc and "Remind me tomorrow" close the menu; the menu doesn't appear
      again that day, and neither does the message after exiting. "Skip this
      version" stops reminders for this version; the note stays.
- [ ] Quitting prints the message to the terminal. Quitting again on the same
      day doesn't. An idle instance still shows the menu afterwards.
- [ ] With `latestversionfoundat` a week ago and `respondedversion` empty,
      the menu appears at startup.
- [ ] Two instances in two terminals, both idle: both show the menu.
      Responding in one closes the menu in the other within ten seconds.
- [ ] After installing in one instance, the other one closes its menu and
      doesn't remind any more. Its note says "Restart to use vX" after its
      next hourly check.
- [ ] `update.method: background` installs at startup without asking.
- [ ] `update.method: never`: the check happens and the note appears, but
      there's no menu and no message after exiting; `u` works.
- [ ] Installed through a symlink (e.g. `~/bin/lazygit` pointing to the
      binary): the update replaces the binary, and the symlink stays.

### Release archive in a folder that isn't writable

- [ ] Make the folder read-only, or use `/usr/local/bin` owned by root. The
      menu explains that lazygit can't update itself and offers the
      installation instructions instead of "Install now". The exit message
      says the same.

### Other kinds

- [ ] `Homebrew`: the menu and the exit message show `brew upgrade lazygit`,
      the menu says that package managers can take a while, and "Copy
      command" copies the command.
- [ ] `goInstall`, or a real `go install` build (see above): `lazygit
      --version` shows `build source=goInstall`, and the menu and the exit
      message show
      `go install github.com/jesseduffield/lazygit@latest`, and "Copy command"
      copies it.
- [ ] `unknown` (Other): the menu tells the user to check their package
      manager and offers the installation instructions.
- [ ] A binary in a mise folder (`.../mise/installs/lazygit/0.60.0/`) with
      build source `binaryRelease` counts as Other.
- [ ] Dev build (`just build`, or `go install` from a branch): the
      bottom-right corner shows the release that the build is based on and
      its commit, e.g. `v0.66.3+d0d00783`, with `-dirty` if there are
      uncommitted changes. `u` shows a popup that says why it can't check. No
      note, and the log shows no check.

### Windows

- [ ] Installing replaces the running `lazygit.exe`. The old one remains as
      `.lazygit-old-<number>` until the next start of lazygit, which deletes
      it.
- [ ] A binary in `C:\Program Files\lazygit` (not writable without admin
      rights) gets the explanation instead of "Install now".
- [ ] Scoop, winget and Chocolatey folders get their commands.

### Errors

- [ ] Without network, the startup check only logs an error; `u` shows it.
- [ ] Installing with a broken network shows the error, and the old binary
      still works.

### After a release

- [ ] The release assets have the names that `assetName` produces
      (`TestAssetName` lists the names of a real release).
- [ ] In the release after the next one, the released version finds and
      installs the new one.

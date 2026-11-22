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

- At startup, if the last successful check was at least `update.days`
  calendar days ago (default 1). With the default, this is the first start of
  each day. Every hour, lazygit tests again whether a check is due.
- With `u` in the status panel.

All install kinds compare with our latest release on GitHub. Package managers
get a new release some time later; the menu tells the user so.

The result goes into the `updates` section of `state.yml`, which all running
instances share.

### Telling the user

- A note in the bottom-right corner. It says "Update available", or "Restart
  to use vX" after an update was installed in place of the running lazygit.
  This lazygit or another one running from the same executable can have
  installed it; lazygit tells by the modification time of the executable.
  The lazygit that installs an update records its version in `state.yml`,
  so vX is the installed version even if a newer release has come out since.
  Clicking the note opens the update menu.
- The update menu at startup, if an earlier check found an update. It appears
  after the user closes the intro or breaking-changes popup, and it replaces
  the recent repos menu. It has "Remind me tomorrow" and "Skip this version",
  and a link to the release notes in its first line. After any choice other
  than skipping, including closing the menu, the menu doesn't appear at startup
  again until the next day. Skipping keeps lazygit quiet until a newer version
  comes out. If the user responds in several instances, the last response
  counts.

A check that finds an update also shows the menu, but only if the user hasn't
pressed a key, clicked, scrolled or pasted since lazygit started. With a fast
network, this is usually the case for the check at startup. Otherwise the menu
would appear while the user is working, and a key that they meant for something
else could choose one of its items; the note tells them about the update until
the next start. Like the menu at startup, it waits for the intro or
breaking-changes popup, and replaces the recent repos menu. If the intro popup
is open, the check at startup has usually finished by the time the user closes
it.

With `update.method: background`, lazygit installs a new release without
asking if it can, and doesn't show the menu for it. With `never`, it still
checks and shows the note, but doesn't show the menu.

If two instances install an update at the same time, only one of them replaces
the executable. They do this while holding the lock of `state.yml`, so the
second one finds the executable replaced already.

## Automated tests

The unit tests in `pkg/updates` cover the install kinds, the asset names (against
the names of a real release), the build source in `.goreleaser.yml`, fetching
the latest version, downloading and installing against a test server,
replacing the executable on Unix and Windows, two instances installing at the
same time, which version an installed update is, when checks and reminders are
due (also with a clock that was wrong), and how responses to the menu combine.
CI runs them on Windows too.

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
lazygit check at the next start; setting `remindafter` to 0 and removing
`skippedversion` brings the menu back at the next start.

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

- [ ] At the first start, the note says "Update available", the log shows one
      check, and the menu appears when the check has finished.
- [ ] At the first start, press a key or click before the check finishes: the
      note appears, but the menu doesn't. To have enough time, slow down the
      network, e.g. with the Network Link Conditioner on macOS. Moving the
      mouse without clicking doesn't keep the menu away.
- [ ] Starting again on the same day shows the note without a check, and the
      menu appears right away.
- [ ] With `lastcheck` set to yesterday evening, less than 24 hours ago,
      lazygit checks at startup.
- [ ] An instance that nobody touched since it started shows the menu after
      its hourly check, if a reminder is due by then, e.g. because another
      instance found the update.
- [ ] `u` shows the menu with "Install now". Clicking "Release notes" in its
      first line opens the releases page, and the menu stays open.
- [ ] Installing shows a waiting status, then a toast. Switching repos works
      while it runs. The note says "Restart to use vX". `lazygit --version`
      shows the new version, and no temporary files are left in the folder.
- [ ] `u` after installing shows the toast again instead of the menu.
- [ ] Copy the binary from the setup over the installed one again, as a
      manual downgrade, and start it: the note says "Update available"
      again, and `u` offers "Install now".
- [ ] Clicking the note opens the menu.
- [ ] Esc and "Remind me tomorrow" close the menu, and the menu doesn't
      appear at startup again until the next day; `remindafter` is the
      coming midnight. "Skip this version" stops the menu for this version;
      the note stays.
- [ ] With the intro popup (remove `disableStartupPopups` from the config and
      set `startuppopupversion` to 0), the menu appears after closing the
      intro. The same with the breaking-changes popup (set `lastversion` to
      0.40.0).
- [ ] Started outside a repo with `notARepository: skip` in the config, the
      menu replaces the recent repos menu.
- [ ] Two instances started at about the same time both show the menu. Choose
      "Skip this version" in one and "Remind me tomorrow" in the other; the
      menu comes back the next day. In the opposite order, the version stays
      skipped.
- [ ] "Install now" in two instances at about the same time: both show the
      toast, the binary is the new version, and no temporary files are left.
      After installing in one instance, "Install now" in the other shows the
      toast without downloading anything. Its note says "Restart to use vX"
      after its next hourly check.
- [ ] `update.method: background` installs at startup without asking, and
      the menu doesn't appear.
- [ ] `update.method: never`: the check happens and the note appears, but
      the menu doesn't; `u` works.
- [ ] Installed through a symlink (e.g. `~/bin/lazygit` pointing to the
      binary): the update replaces the binary, and the symlink stays.

### Release archive in a folder that isn't writable

- [ ] Make the folder read-only, or use `/usr/local/bin` owned by root. The
      menu explains that lazygit can't update itself and offers the
      installation instructions instead of "Install now".

### Other kinds

- [ ] `Homebrew`: the menu shows `brew upgrade lazygit` and says that
      package managers can take a while, and "Copy command" copies the
      command.
- [ ] `goInstall`, or a real `go install` build (see above): `lazygit
      --version` shows `build source=goInstall`, the menu shows
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
- [ ] "Install now" in two instances at about the same time: neither shows an
      error, and the binary is the new version.
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

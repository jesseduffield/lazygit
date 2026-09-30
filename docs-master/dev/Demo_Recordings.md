# Demo Recordings

We want our demo recordings to be consistent and easy to update if we make changes to Lazygit's UI. Luckily for us, we have an existing recording system for the sake of our integration tests, so we can piggyback on that.

You'll want to familiarise yourself with how integration tests are written: see [here](../../pkg/integration/README.md).

## Prerequisites

Ideally we'd run this whole thing through docker but we haven't got that working. So you will need:
```
# for recording; vhs drives ttyd and ffmpeg under the hood
brew install ttyd ffmpeg

# vhs 0.12.0 runs the tape, reports success and writes no video at all
# (https://github.com/charmbracelet/vhs/issues/787), so pin the release
# before it
go install github.com/charmbracelet/vhs@v0.11.0

# font with icons
wget https://github.com/ryanoasis/nerd-fonts/releases/download/v3.0.2/SourceCodePro.tar.xz && \
  tar -xf SourceCodePro.tar.xz -C ~/Library/Fonts && \
  rm SourceCodePro.tar.xz
```

## Creating a demo

Demos are found in `pkg/integration/tests/demo/`. They are like regular integration tests but have `IsDemo: true` which has a few effects:
* The bottom row of the UI is quieter so that we can render captions
* Fetch/Push/Pull have artificial latency to mimic a network request
* The loader at the bottom-right does not appear

In demos, we don't need to be as strict in our assertions as we are in tests. But it's still good to have some basic assertions so that if we automate the process of updating demos we'll know if one of them has broken.

You can use the same flow as we use with integration tests when you're writing a demo:
* Setup the repo
* Run the demo in sandbox mode to get a feel of what needs to happen
* Come back and write the code to make it happen

### Adding captions

It's good to add captions explaining what task if being performed. Use the existing demos as a guide.

### Setting up the assets worktree

We store assets (which includes demo recordings) in the `assets` branch, which is a branch that shares no history with the main branch and exists purely for storing assets. Storing them separately means we don't clog up the code branches with large binaries.

The scripts and demo definitions live in the code branches but the output lives in the assets branch so to be able to create a video from a demo you'll need to create a linked worktree for the assets branch which you can do with:

```sh
git worktree add .worktrees/assets assets
```

The mp4 of the recording will be stored in `.worktrees/assets/demo/`.

### Recording the demo

Once you're happy with your demo you can record it using:
```sh
scripts/record_demo.sh <path>
# e.g.
scripts/record_demo.sh pkg/integration/tests/demo/interactive_rebase.go
```

The terminal size, font and colours live in `demo/settings.tape`, which the
script sources into the tape it generates for the demo.

While you are still working on how a demo looks, pass `--no-upload`. That
leaves the video in `demo/output` (which is git-ignored) and stops there, so
you can watch it without uploading anything or touching the assets worktree:

```sh
scripts/record_demo.sh --no-upload pkg/integration/tests/demo/interactive_rebase.go
```

### Including demos in README/docs

Recording a demo does three things with the mp4: it writes it to your assets
worktree, it uploads a copy to GitHub's attachment store, and it posts that
copy as a comment on the issue named by `PUBLISH_ISSUE` in the script. Then it
prints the tag to embed:

```html
<video src="https://github.com/user-attachments/assets/<uuid>" controls></video>
```

GitHub plays a video in a README only when it is served from its own attachment
store. If you commit a video to the assets branch and link it the way we link
the images, GitHub drops the whole `<video>` element when it renders the page.
So the README reads the uploaded copy rather than the one in the assets
worktree. Keep that one anyway, so that we still have the file if we ever need
to upload it again. Stage it and raise a PR against the assets branch as you
would for any other asset.

Attachment URLs are opaque and have no path we can predict, so a new recording
of an existing demo means a new URL and an edit to the page that embeds it.

That comment on `PUBLISH_ISSUE` is not bookkeeping; it is what makes the video
watchable. An uploaded attachment is readable only by people signed in to
GitHub until some posted comment in the repository refers to it, and a README
on a branch does not count. Skip that step and the video plays for you and
404s for everyone else, which is easy to miss because you are signed in. The
script waits until the video can be fetched without a token before it prints
the tag. Referring to an attachment once is enough and cannot be undone, so
the comments could be deleted later, but leaving them gives us a dated list of
every recording.

Uploading needs push access to the lazygit repository. If you don't have it,
record the demo, then ask a maintainer to upload the mp4 for you.

### Re-recording every demo on a page

Write the name of the demo above each video, so that we can find our way from
a page back to the demo that produced it:

```html
<!-- demo: commit_and_push -->
<video src="https://github.com/user-attachments/assets/<uuid>" controls></video>
```

GitHub drops the comment when it renders the page. With it in place, a change
to `demo/settings.tape` or to lazygit's own appearance can be rolled out across
every recording at once:

```sh
scripts/rerecord_demos.sh
# or, to look before you upload anything:
scripts/rerecord_demos.sh --no-upload
```

That re-records every demo the README embeds and rewrites each URL in place.
Name other pages as arguments to do the same for them.

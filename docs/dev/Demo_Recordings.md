# Demo Recordings

We want our demo recordings to be consistent and easy to update if we make changes to Lazygit's UI. Luckily for us, we have an existing recording system for the sake of our integration tests, so we can piggyback on that.

You'll want to familiarise yourself with how integration tests are written: see [here](../../pkg/integration/README.md).

## Prerequisites

Ideally we'd run this whole thing through docker but we haven't got that working. So you will need:
```
# for recording
npm i -g terminalizer
# for gif compression
npm i -g gifsicle
# for mp4 conversion
brew install ffmpeg

# font with icons
wget https://github.com/ryanoasis/nerd-fonts/releases/download/v3.0.2/DejaVuSansMono.tar.xz && \
  tar -xf DejaVuSansMono.tar.xz -C /usr/local/share/fonts && \
  rm DejaVuSansMono.tar.xz
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

Outputs will be stored in `.worktrees/assets/demos/`. We'll store three separate things:
* the yaml of the recording
* the original gif
* either the compressed gif or the mp4 depending on the output you chose (see below)

### Recording the demo

Once you're happy with your demo you can record it using:
```sh
scripts/record_demo.sh [gif|mp4] <path>
# e.g.
scripts/record_demo.sh gif pkg/integration/tests/demo/interactive_rebase.go
```

~~The gif format is for use in the first video of the readme (it has a larger size but has auto-play and looping)~~
~~The mp4 format is for everything else (no looping, requires clicking, but smaller size).~~

Turns out that you can't store mp4s in a repo and link them from a README so we're gonna just use gifs across the board for now.

### Including demos in README/docs

If you've followed the above steps you'll end up with your output in your assets worktree.

Within that worktree, stage all three output files and raise a PR against the assets branch.

Then back in the code branch, in the doc, you can embed the recording like so:
```md
![Nuke working tree](../assets/demo/interactive_rebase-compressed.gif)
```

This means we can update assets without needing to update the docs that embed them.

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

### Recording the demo

Once you're happy with your demo you can record it using:
```sh
scripts/record_demo.sh <path>
# e.g.
scripts/record_demo.sh pkg/integration/tests/demo/interactive_rebase.go
```

### Storing demos

This part is subject to change. I'm thinking of storing all gifs in the `assets` branch. But yet to finalize on that.
For now, feel free to upload `demo/demo-compressed.gif` to GitHub by dragging and dropping it in a file in the browser (e.g. the README).

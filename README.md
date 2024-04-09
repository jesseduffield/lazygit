# lazytask

A simple terminal UI for tasks

## Table of contents

- [Table of contents](#table-of-contents)
- [Features](#features)
- [Usage](#usage)
  - [Keybindings](#keybindings)
  - [Changing Directory On Exit](#changing-directory-on-exit)
  - [Undo/Redo](#undoredo)
- [Configuration](#configuration)
  - [Custom Pagers](#custom-pagers)
- [Contributing](#contributing)
  - [Debugging Locally](#debugging-locally)
- [FAQ](#faq)
  - [What do the task colors represent?](#what-do-the-task-colors-represent)

## Features

## Usage

Call `lazytask` in your terminal inside a tome.

```sh
lazytask
```

### Alias

If you want, you can also add an alias for this with 

```sh
echo "alias lt='lazytask'" >> ~/.zshrc` 
```
(or whichever rc file you're using).

TODO: test `profile.ps1`

### Keybindings

You can check out the list of keybindings [here](/docs/keybindings).

### Changing Directory On Exit

If you change tomes in `lazytask` and want your shell to change directory into that tome on exiting `lazytask`, add this to your `~/.zshrc` (or other rc file):

```sh
lt()
{
  export LAZYTASK_NEW_DIR_FILE=~/.lazytask/newdir

  lazytask "$@"

  if [ -f $LAZYTASK_NEW_DIR_FILE ]; then
    cd "$(cat $LAZYTASK_NEW_DIR_FILE)"
    rm -f $LAZYTASK_NEW_DIR_FILE > /dev/null
  fi
}
```

TODO: test `profile.ps1`

Then `source ~/.zshrc` and from now on when you call `lt` and exit you'll switch directories to whereever you were in inside `lazytask`. To override this behaviour you can exit using `shift+Q` rather than just `q`.

### Undo/Redo

See the [docs](/docs/Undoing.md)

## Configuration

Check out the [configuration docs](docs/Config.md).

### Custom Pagers

See the [docs](docs/Custom_Pagers.md)


## Contributing

We love your input! Please check out the [contributing guide](CONTRIBUTING.md).

### Debugging Locally

Run `lazytask --debug` in one terminal tab and `lazytask --logs` in another to view the program and its log output side by side.


## FAQ

### What do the task colours represent?

- **Green:** the task is done
- **Yellow:** the task is in progress
- **Red:** the task isn't started


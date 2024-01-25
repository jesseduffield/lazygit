# Lazygit Codebase Guide

## Packages

* `pkg/app`: Contains startup code, inititalises a bunch of stuff like logging, the user config, etc, before starting the gui. Catches and handles some errors that the gui raises.
* `pkg/app/daemon`: Contains code relating to the lazygit daemon. This could be better named: it's is not a daemon in the sense that it's a long-running background process; rather it's a short-lived background process that we pass to git for certain tasks, like GIT_EDITOR for when we want to set the TODO file for an interactive rebase.
* `pkg/cheatsheet`: Generates the keybinding cheatsheets in `docs/keybindings`.
* `pkg/commands/git_commands`: All communication to the git binary happens here. So for example there's a `Checkout` method which calls `git checkout`.
* `pkg/commands/oscommands`: Contains code for talking to the OS, and for invoking commands in general
* `pkg/commands/git_config`: Reading of the git config all happens here.
* `pkg/commands/hosting_service`: Contains code that is specific to git hosting services (aka forges).
* `pkg/commands/models`: Contains model structs that represent commits, branches, files, etc.
* `pkg/commands/patch`: Contains code for parsing and working with git patches
* `pkg/common`: Contains the `Common` struct which holds common dependencies like the logger, i18n, and the user config. Most structs in the code will have a field named `c` which holds a common struct (or a derivative of the common struct).
* `pkg/config`: Contains code relating to the Lazygit user config. Specifically `pkg/config/user_config/go` defines the user config struct and its default values.
* `pkg/constants`: Contains some constant strings (e.g. links to docs)
* `pkg/env`: Contains code relating to setting/getting environment variables
* `pkg/i18n`: Contains internationalised strings
* `pkg/integration`: Contains end-to-end tests
* `pkg/jsonschema`: Contains generator for user config JSON schema.
* `pkg/logs`: Contains code for instantiating the logger and for tailing the logs via `lazygit --logs`
* `pkg/tasks`: Contains code for running asynchronous tasks: mostly related to efficiently rendering command output to the main window.
* `pkg/theme`: Contains code related to colour themes.
* `pkg/updates`: Contains code related to Lazygit updates (checking for update, download and installing the update)
* `pkg/utils`: Contains lots of low-level helper functions
* `pkg/gui`: Contains code related to the gui. We've still got a God Struct in the form of our Gui struct, but over time code has been moved out into contexts, controllers, and helpers, and we intend to continue moving more code out over time.
* `pkg/gui/context`: Contains code relating to contexts. There is a context for each view e.g. a branches context, a tags context, etc. Contexts manage state related to the view and receive keypresses.
* `pkg/gui/controllers`: Contains code relating to controllers. Controllers define a list of keybindings and their associated handlers. One controller can be assigned to multiple contexts, and one context can contain multiple controllers.
* `pkg/gui/controllers/helpers`: Contains code that is shared between multiple controllers.
* `pkg/gui/filetree`: Contains code relating to the representation of filetrees.
* `pkg/gui/keybindings`: Contains code for mapping between keybindings and their labels
* `pkg/gui/mergeconflicts`: Contains code relating to the handling of merge conflicts
* `pkg/gui/modes`: Contains code relating to the state of different modes e.g. cherry picking mode, rebase mode.
* `pkg/gui/patch_exploring`: Contains code relating to the state of patch-oriented views like the staging view.
* `pkg/gui/popup`: Contains code that lets you easily raise popups
* `pkg/gui/presentation`: Contains presentation code i.e. code concerned with rendering content inside views
* `pkg/gui/services/custom_commands`: Contains code related to user-defined custom commands.
* `pkg/gui/status`: Contains code for invoking loaders and toasts
* `pkg/gui/style`: Contains code for specifying text styles (colour, bold, etc)
* `pkg/gui/types`: Contains various gui-specific types and interfaces. Lots of code lives here to avoid circular dependencies
* `vendor/github.com/jesseduffield/gocui`: Gocui is the underlying library used for handling the gui event loop, handling keypresses, and rendering the UI. It defines the View struct which our own context structs build upon.

## Important files

* `pkg/config/user_config.go`: defines the user config and default values
* `pkg/gui/keybindings.go`: defines keybindings which have not yet been moved into a controller (originally all keybindings were defined here)
* `pkg/gui/controllers.go`: links up controllers with contexts
* `pkg/gui/controllers/helpers/helpers.go`: defines all the different helper structs
* `pkg/commands/git.go`: defines all the different git command structs
* `pkg/gui/gui.go`: defines the top-level gui state and gui initialisation/run code
* `pkg/gui/layout.go`: defines what happens on each render
* `pkg/gui/controllers/helpers/window_arrangement_helper.go`: defines the layout of the UI and the size/position of each window
* `pkg/gui/context/context.go`: defines the different contexts
* `pkg/gui/context/setup.go`: defines initialisation code for all contexts
* `pkg/gui/context.go`: manages the lifecycle of contexts, the context stack, and focus changes.
* `pkg/gui/types/views.go`: defines views
* `pkg/gui/views.go`: defines the ordering of views (front to back) and their initialisation code
* `pkg/gui/gui_common.go`: defines gui-specific methods that all controllers and helpers have access to
* `pkg/i18n/english.go`: defines the set of i18n strings and their English values
* `pkg/gui/controllers/helpers/refresh_helper.go`: manages refreshing of models. The refresh helper is typically invoked at the end of an action to re-load affected models from git (e.g. re-load branches after doing a git pull)
* `pkg/gui/controllers/quit_actions.go`: contains code that runs when you hit 'escape' on a view (assuming the view doesn't define its own escape handler)
* `vendor/github.com/jesseduffield/gocui/gui.go`: defines the gocui gui struct
* `vendor/github.com/jesseduffield/gocui/view.go`: defines the gocui view struct

## Concepts

* **View**: Views are defined in the gocui package, and they maintain an internal buffer of content which is rendered each time the screen is drawn.
* **Context**: A context is tied to a view and contains some additional state and logic specific to that view e.g. the branches context has code relating specifically to branches, and writes the list of branches to the branches view. Views and contexts share some responsibilities for historical reasons.
* **Controller**: A controller defined keybindings with associated handlers. One controller can be assigned to multiple contexts and one context can have multiple controllers. For example the list controller handles keybindings relating to navigating a list, and is assigned to all list contexts (e.g. the branches context).
* **Helper**: A helper defines shared code used by controllers, or used by some other parts of the application. Often a controller will have a method that ends up needing to be used by another controller, so in that case we move the method out into a helper so that both controllers can use it. We need to do this because controllers cannot refer to other controllers' methods.

In terms of dependencies, controllers sit at the highest level, so they can refer to helpers, contexts, and views (although it's preferable for view-specific code to live in contexts). Helpers can refer to contexts and views, and contexts can only refer to views. Views can't refer to contexts, controllers, or helpers.

* **Window**: A window is a section of the screen which will render a view. Windows are named after the default view that appears there, so for example there is a 'stash' window that is so named because by default the stash view appears there. But if you press enter on a stash entry, the stash entry's files will be shown in a different view, but in the same window.
* **Panel**: The term 'panel' is still used in a few places to refer to either a view or a window, and it's a term that is now deprecated in favour of 'view' and 'window'.
* **Tab**: Each tab in a window (e.g. Files, Worktrees, Submodules) actually has a corresponding view which we bring to the front upon changing tabs.
* **Model**: Representation of a git object e.g. commits, branches, files.
* **ViewModel**: Used by a context to maintain state related to the view.
* **Keybinding**: A keybinding associates a _key_ with an _action_. For example if you press the 'down' arrow, the action performed will be your cursor moving down a list by one.
* **Action**: An action is the thing that happens when you press a key. Often an action will invoke a git command, but not always: for example, navigation actions don't involve git.
* **Common structs**: Most structs have a field named `c` which contains a 'common' struct: a struct containing a bag of dependencies that most structs of the same layer require. For example if you want to access a helper from a controller you can do so with `self.c.Helpers.MyHelper`.

## Event loop and threads

The event loop is managed in the `MainLoop` function of `vendor/github.com/jesseduffield/gocui/gui.go`. Any time there is an event like a key press or a window resize, the event will be processed and then the screen will be redrawn. This involves calling the `layout` function defined in `pkg/gui/layout.go`, which lays out the windows and invokes some on-render hooks.

Often, as part of handling a keypress, we'll want to run some code asynchronously so that it doesn't block the UI thread. For this we'll typically run `self.c.OnWorker(myFunc)`. If the worker wants to then do something on the UI thread again it can call `self.c.OnUIThread(myOtherFunc)`.

## Legacy code structure

Before we had controllers and contexts, all the code lived directly in the gui package under a gui God Struct. This was fairly bloated and so we split things out to have a better separation of concerns. Nonetheless, it's a big effort to migrate all the code so we still have some logic in the gui struct that ought to live somewhere else. Likewise, we have some keybindings defined in `pkg/gui/keybindings.go` that ought to live on a controller (all keybindings used to be defined in that one file).

The new structure has its own problems: we don't have a clear guide on whether code should live in a controller or helper. The current approach is to put code in a controller until it's needed by another controller, and to then extract it out into a helper. We may be better off just putting code in helpers to start with and leaving controllers super-thin, with the responsibility of just pairing keys with corresponding helper functions. But it's not clear to me if that would be better than the current approach.


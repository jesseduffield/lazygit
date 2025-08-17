# Vision and Design Principles

## Vision

Lazygit's vision is to be the most enjoyable UI for git.

## Design Principles

There are seven (sometimes contradictory) design principles we follow:

- Dicoverability
- Simplicity
- Safety
- Power
- Speed
- Conformity with git
- Think of the codebase

### Discoverability

TUI's are notoriously hard to learn, thanks to limited screen real-estate to provide contextual help and a general lack of effort on the part of developers to make things obvious. We want Lazygit to buck the trend and be easy for a new user to grok.

Examples:

- Clearly document all the features/configuration options
  - e.g. gifs in the README
- Document how to solve various git problems with Lazygit
  - This is something we don't have yet but should: a section in the docs explaining how Lazygit can help you in various scenarios
- Use tooltips to explain what actions will do
- Make it easy for users to ask questions and get answers from the community
- Make it easy to find entities and actions from within Lazygit
- Use visual elements to make things obvious
  - e.g. '<-- YOU ARE HERE' label when rebasing
- Don't require the user to memorise keybindings
  - e.g. when the user is mid-rebase, we prominently show that the keybinding for viewing rebase options is 'm'
- When the user performs an action in Lazygit, make the impact obvious
  - If the affected entity isn't visible, show a toast notification
- If a keybinding is disabled, give a reason why

### Simplicity

The git CLI is very complex but most git use cases are simple. Lazygit needs to ensure that simple use cases are easy to satisfy.

- Make the most common use cases dead-simple (staging files, committing, pulling/pushing)
- Don't overwhelm the user with options
- Use sensible defaults
- We already have too many configuration options: think hard before adding any new ones

### Safety

It's easy to screw things up in git so Lazygit should try to protect the user from screwing things up.

- Prompt for a confirmation before doing anything that's hard to reverse
- Make it easy to correct mistakes
  - e.g. undo action
  - the escape key should get you out of most transient situations (rebasing, diffing, etc)

## Power

Users shouldn't have to drop down the CLI _too_ often. Lazygit should be able to handle some complex use cases.

- Make complex (but common) CLI flows simple
  - e.g. interactive rebasing
- Use the custom commands system to handle the really rare complex edge-cases

### Speed

Pro users should be able to move at lightning speed with Lazygit.

- Always think about the number of keypresses involved in a given UX flow
- Make lazygit performant and responsive
- Think about the individual commands being run and how fast they are
- Startup should be FAST. If you want to run something at startup that is slow, make it non-blocking.
- Support muscle-memory
  - Prefer disabling menu items instead of hiding them so that muscle memory can be used to select the desired menu item
  - Try to make keybinding intuitions to transfer across contexts (e.g. 'd' for destroy)
  - When changing keybindings in a new release, always consider what will happen if a user does not read the release notes and relies on muscle memory.

### Conformity with git

Satisfying the use-cases of git users is more important than perfectly conforming to git's API, but even obscure parts of git's API were motivated by real use-cases.

- Users should only have to drop down to the git CLI in rare circumstances
- Honour the git config
  - Don't override anything set in the git config without the user's permission
- Work with git, not against it.
  - Too much magic will get us into trouble
- Avoid storing Lazygit-specific session state that could instead be stored in git
- Ensure that Lazygit can represent the state of any repo
- Sometimes git's default behaviour is just silly and we'll make the call to override but it should be a well-considered decision.

### Think of the codebase

Will somebody PLEASE think of the codebase!

Some features are not worth the added complexity in the codebase. The more this codebase grows, the harder it will be to make the changes that everybody wants.

## Resolving conflicts

Many of the above objectives are directly antithetical to one another. If you add an extra confirmation prompt for the sake of _safety_, you're sacrificing _speed_. If you support toggling various git flags in the name of _power_, you're sacrificing _simplicity_. There are a few things to say here.

When there are conflicts, we need to make a judgement call. In general we should err on the side of safety and simplicity as the default, with the ability for users to make things faster / more powerful either through configuration or separate keybindings.

This does not mean for example that force pushes should be impossible without being manually enabled: force pushes are table stakes for anybody who rebases. But it does mean that a confirmation popup should appear when force pushing.

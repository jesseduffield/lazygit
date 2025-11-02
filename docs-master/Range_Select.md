# Range Select

Some actions can be performed on a range of contiguous items. For example:
* staging multiple files at once
* squashing multiple commits at once
* copying (for cherry-pick) multiple commits at once

There are two ways to select a range of items:
1. Sticky range select: Press 'v' to toggle range select, then expand the selection using the up/down arrow key. To reset the selection, press 'v' again.
2. Non-sticky range select: Press shift+up or shift+down to expand the selection. To reset the selection, press up/down without shift.

The sticky option will be more familiar to vim users, and the second option will feel more natural to users who aren't used to doing things in a modal way.

In order to perform an action on a range of items, simply press the normal key for that action. If the action only works on individual items, it will raise an error. This is a new feature and the plan is to incrementally support range select for more and more actions. If there is an action you would like to support range select which currently does not, please raise an issue in the repo.

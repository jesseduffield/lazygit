# Knowing when Lazygit is busy/idle

## The use-case

This topic deserves its own doc because there there are a few touch points for it. We have a use-case for knowing when Lazygit is idle or busy because integration tests follow the following process:
1) press a key
2) wait until Lazygit is idle
3) run assertion / press another key
4) repeat

In the past the process was:
1) press a key
2) run assertion
3) if assertion fails, wait a bit and retry
4) repeat

The old process was problematic because an assertion may give a false positive due to the contents of some view not yet having changed since the last key was pressed.

## The solution

First, it's important to distinguish three different types of goroutines:
* The UI goroutine, of which there is only one, which infinitely processes a queue of events
* Worker goroutines, which do some work and then typically enqueue an event in the UI goroutine to display the results
* Background goroutines, which periodically spawn worker goroutines (e.g. doing a git fetch every minute)

The point of distinguishing worker goroutines from background goroutines is that when any worker goroutine is running, we consider Lazygit to be 'busy', whereas this is not the case with background goroutines. It would be pointless to have background goroutines be considered 'busy' because then Lazygit would be considered busy for the entire duration of the program!

In gocui, the underlying package we use for managing the UI and events, we keep track of how many busy goroutines there are using the `Task` type. A task represents some work being done by lazygit. The gocui Gui struct holds a map of tasks and allows creating a new task (which adds it to the map), pausing/continuing a task, and marking a task as done (which removes it from the map). Lazygit is considered to be busy so long as there is at least one busy task in the map; otherwise it's considered idle. When Lazygit goes from busy to idle, it notifies the integration test.

It's important that we play by the rules below to ensure that after the user does anything, all the processing that follows happens in a contiguous block of busy-ness with no gaps.

### Spawning a worker goroutine

Here's the basic implementation of `OnWorker` (using the same flow as `WaitGroup`s):

```go
func (g *Gui) OnWorker(f func(*Task)) {
	task := g.NewTask()
	go func() {
		f(task)
		task.Done()
	}()
}
```

The crucial thing here is that we create the task _before_ spawning the goroutine, because it means that we'll have at least one busy task in the map until the completion of the goroutine. If we created the task within the goroutine, the current function could exit and Lazygit would be considered idle before the goroutine starts, leading to our integration test prematurely progressing.

You typically invoke this with `self.c.OnWorker(f)`. Note that the callback function receives the task. This allows the callback to pause/continue the task (see below).

### Spawning a background goroutine

Spawning a background goroutine is as simple as:

```go
go utils.Safe(f)
```

Where `utils.Safe` is a helper function that ensures we clean up the gui if the goroutine panics.

### Programmatically enqueing a UI event

This is invoked with `self.c.OnUIThread(f)`. Internally, it creates a task before enqueuing the function as an event (including the task in the event struct) and once that event is processed by the event queue (and any other pending events are processed) the task is removed from the map by calling `task.Done()`.

### Pressing a key

If the user presses a key, an event will be enqueued automatically and a task will be created before (and `Done`'d after) the event is processed.

## Special cases

There are a couple of special cases where we manually pause/continue the task directly in the client code. These are subject to change but for the sake of completeness:

### Writing to the main view(s)

If the user focuses a file in the files panel, we run a `git diff` command for that file and write the output to the main view. But we only read enough of the command's output to fill the view's viewport: further loading only happens if the user scrolls. Given that we have a background goroutine for running the command and writing more output upon scrolling, we create our own task and call `Done` on it as soon as the viewport is filled.

### Requesting credentials from a git command

Some git commands (e.g. git push) may request credentials. This is the same deal as above; we use a worker goroutine and manually pause continue its task as we go from waiting on the git command to waiting on user input. This requires passing the task through to the `Push` method so that it can be paused/continued.

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

Lazygit is considered to be 'busy' so long as the counter remains greater than zero, and as soon as it hits zero, Lazygit is considered 'idle' and the integration test is notified. So it's important that we play by the rules below to ensure that after the user does anything, all the processing that follows happens in a contiguous block of busy-ness with no gaps.

In gocui, the underlying package we use for managing the UI and events, we keep track of how many busy goroutines there are with a `busyCount` counter.

### Spawning a worker goroutine

Here's the basic implementation of `OnWorker`:

```go
func (g *Gui) OnWorker(f func()) {
	g.IncrementBusyCount()
	go func() {
		f()
		g.DecrementBusyCount()
	}()
}
```

The crucial thing here is that we increment the busy count _before_ spawning the goroutine, because it means that our counter never goes to zero while there's still work being done. If we incremented the busy count within the goroutine, the current function could exit and decrement the counter to zero before the goroutine starts.

You typically invoke this with `self.c.OnWorker(f)`

### Spawning a background goroutine

Spawning a background goroutine is as simple as:

```go
go utils.Safe(f)
```

Where `utils.Safe` is a helper function that ensures we clean up the gui if the goroutine panics.

### Programmatically enqueing a UI event

This is invoked with `self.c.OnUIThread(f)`. Internally, it increments the counter before enqueuing the function as an event and once that event is processed by the event queue (and any other pending events are processed) the counter is decremented again.

### Pressing a key

If the user presses a key, an event will be enqueued automatically and the counter will be incremented before (and decremented after) the event is processed.

## Special cases

There are a couple of special cases where we manually increment/decrement the counter in the code. These are subject to change but for the sake of completeness:

### Writing to the main view(s)

If the user focuses a file in the files panel, we run a `git diff` command for that file and write the output to the main view. But we only read enough of the command's output to fill the view's viewport: further loading only happens if the user scrolls. Given that we have a background goroutine for running the command and writing more output upon scrolling, we manually increment the busy count within that goroutine and then decrement it once the viewport is filled.

### Requesting credentials from a git command

Some git commands (e.g. git push) may request credentials. This is the same deal as above; we use a background goroutine and manually increment/decrement the counter as we go from waiting on the git command to waiting on user input.

## Future improvements

### Better API

The current approach is fairly simple in terms of the API which, except for the special cases above, encapsulates the incrementing/decrementing of the busy counter. But the counter is a form of global state and in future we may switch to an API where we have objects representing a task in progress, and those objects have `Start()`, `Finish()`, and `Pause()` methods. This would better defend against bugs caused by a random goroutine accidentally decrementing twice, for example.

### More applications

We could use the concept of idle/busy to show a loader whenever Lazygit is busy. But our current situation is pretty good: we have the `WithWaitingStatus()` method for running a function on a worker goroutine along with a message to show within the loader e.g. 'Refreshing branches'. If we find a situation where we're a function is taking a while and a loader isn't appearing, that's because we're running the code on the UI goroutine and we should just wrap the code in `WithWaitingStatus()`.

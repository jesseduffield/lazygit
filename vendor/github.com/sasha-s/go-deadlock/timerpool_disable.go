//go:build deadlock_disable

package deadlock

// shouldDisableTimerPool always returns true when deadlock detection is disabled
// since there's no timer pool or deadlock detection happening anyway
func shouldDisableTimerPool() bool {
	return true
}

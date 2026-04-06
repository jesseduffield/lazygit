//go:build (goexperiment.synctest || deadlock_synctest) && !deadlock_disable

package deadlock

// shouldDisableTimerPool determines if timer pooling should be disabled
// In synctest builds, timer pooling is disabled by default to avoid cross-bubble issues
func shouldDisableTimerPool() bool {
	switch Opts.TimerPool {
	case TimerPoolDefault:
		return true // Default: disable timer pooling for synctest compatibility
	case TimerPoolEnabled:
		return false
	case TimerPoolDisabled:
		return true
	default:
		return true
	}
}

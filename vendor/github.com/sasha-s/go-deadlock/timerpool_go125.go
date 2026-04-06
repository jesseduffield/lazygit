//go:build !goexperiment.synctest && !deadlock_synctest && !deadlock_disable && go1.25

package deadlock

// shouldDisableTimerPool determines if timer/entry pooling should be disabled.
// In Go 1.25, pooling is enabled by default for performance.
func shouldDisableTimerPool() bool {
	switch Opts.TimerPool {
	case TimerPoolDefault:
		return false // Default: enable timer pooling for performance
	case TimerPoolEnabled:
		return false
	case TimerPoolDisabled:
		return true
	default:
		return false
	}
}

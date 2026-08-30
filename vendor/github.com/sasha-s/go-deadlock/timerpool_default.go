//go:build !goexperiment.synctest && !deadlock_synctest && !deadlock_disable && !go1.25

package deadlock

// shouldDisableTimerPool determines if timer pooling should be disabled
// In normal builds, timer pooling is enabled by default for performance
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

package presentation

import (
	"time"

	"github.com/jesseduffield/lazygit/pkg/config"
)

// Loader dumps a string to be displayed as a loader
func Loader(now time.Time, config config.SpinnerConfig) string {
	milliseconds := now.UnixMilli()
	index := milliseconds / int64(config.Rate) % int64(len(config.Frames))
	return config.Frames[index]
}

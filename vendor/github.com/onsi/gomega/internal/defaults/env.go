package defaults

import (
	"fmt"
	"time"
)

func SetDurationFromEnv(getDurationFromEnv func(string) string, varSetter func(time.Duration), name string) {
	durationFromEnv := getDurationFromEnv(name)

	if len(durationFromEnv) == 0 {
		return
	}

	duration, err := time.ParseDuration(durationFromEnv)

	if err != nil {
		panic(fmt.Sprintf("Expected a duration when using %s!  Parse error %v", name, err))
	}

	varSetter(duration)
}

package justest

import (
	"os"
	"strconv"
	"time"
)

func transformDurationIfNecessary(t T, d time.Duration) time.Duration {
	if v, found := os.LookupEnv(SlowFactorEnvVarName); found {
		if factor, err := strconv.ParseInt(v, 0, 0); err != nil {
			t.Logf("Ignoring value of '%s' environment variable: %+v", SlowFactorEnvVarName, err)
			return d
		} else {
			oldSeconds := int64(d.Seconds())
			newSeconds := oldSeconds * factor
			return time.Duration(newSeconds) * time.Second
		}
	}
	return d
}

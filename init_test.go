package justest_test

import "os"

func init() {
	_ = os.Setenv("JUSTEST_DISABLE_SOURCE_HIGHLIGHT", "false")
}

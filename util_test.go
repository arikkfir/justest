package justest

import (
	"testing"
	"time"
)

func TestTransformDurationIfNecessary(t *testing.T) {
	With(t).Verify(transformDurationIfNecessary(t, 5*time.Second)).Will(EqualTo(5 * time.Second)).OrFail()
	t.Setenv(SlowFactorEnvVarName, "2")
	With(t).Verify(transformDurationIfNecessary(t, 5*time.Second)).Will(EqualTo(10 * time.Second)).OrFail()
	t.Setenv(SlowFactorEnvVarName, "3")
	With(t).Verify(transformDurationIfNecessary(t, 5*time.Second)).Will(EqualTo(15 * time.Second)).OrFail()
}

package justest

import (
	"testing"
	"time"
)

func TestTransformDurationIfNecessary(t *testing.T) {
	With(t).VerifyThat(transformDurationIfNecessary(t, 5*time.Second)).Will(EqualTo(5 * time.Second)).Now()
	t.Setenv(SlowFactorEnvVarName, "2")
	With(t).VerifyThat(transformDurationIfNecessary(t, 5*time.Second)).Will(EqualTo(10 * time.Second)).Now()
	t.Setenv(SlowFactorEnvVarName, "3")
	With(t).VerifyThat(transformDurationIfNecessary(t, 5*time.Second)).Will(EqualTo(15 * time.Second)).Now()
}

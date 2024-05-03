package justest_test

import (
	. "github.com/arikkfir/justest"
	"testing"
)

func TestBeNil(t *testing.T) {
	t.Parallel()
	t.Run("Nil", func(t *testing.T) {
		t.Parallel()
		defer VerifyTestOutcome(t, ExpectSuccess, "")
		With(NewMockT(t)).Verify(nil).Will(BeNil()).OrFail()
	})
	t.Run("Not nil", func(t *testing.T) {
		t.Parallel()
		defer VerifyTestOutcome(t, ExpectFailure, `Expected actual to be nil, but it is not: abc`)
		With(NewMockT(t)).Verify("abc").Will(BeNil()).OrFail()
	})
}

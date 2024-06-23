package justest_test

import (
	"testing"

	. "github.com/arikkfir/justest"
)

func TestBeNil(t *testing.T) {
	t.Parallel()
	t.Run("Nil", func(t *testing.T) {
		t.Parallel()
		mt := NewMockT(t)
		defer mt.Verify(SuccessVerifier())
		With(mt).VerifyThat(nil).Will(BeNil()).Now()
	})
	t.Run("Not nil", func(t *testing.T) {
		t.Parallel()
		mt := NewMockT(t)
		defer mt.Verify(FailureVerifier(`Expected actual to be nil, but it is not: abc`))
		With(mt).VerifyThat("abc").Will(BeNil()).Now()
	})
}

package backoff

import (
	"testing"
	"time"
)

// func TestNextBackOffMillis(t *testing.T) {
// 	subtestNextBackOff(t, 0, new(ZeroBackOff))
// 	// subtestNextBackOff(t, Stop, new(StopBackOff))
// }

// func subtestNextBackOff(t *testing.T, expectedValue time.Duration, backOffPolicy BackOff) {
// 	for i := 0; i < 10; i++ {
// 		next := backOffPolicy.NextBackOff()
// 		t.Log(next)
// 		if next != expectedValue {
// 			t.Errorf("got: %d expected: %d", next, expectedValue)
// 		}
// 	}
// }

// func TestConstantBackOff(t *testing.T) {
// 	backoff := NewConstantBackOff(time.Second)
// 	for i := 0; i < 10; i++ {
// 		t.Log(backoff.NextBackOff())
// 	}
// 	if backoff.NextBackOff() != time.Second {
// 		t.Error("invalid interval")
// 	}
// }
func TestConstantBackOff2(t *testing.T) {
	backoff := NewExponentialBackOff()
	backoff.MaxElapsedTime = 10 * time.Second
	for i := 0; i < 20; i++ {
		v := backoff.NextBackOff()
		t.Log(i, v, v == Stop)
	}
	// if backoff.NextBackOff() != time.Second {
	// 	t.Error("invalid interval")
	// }
}

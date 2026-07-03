package cache

import (
	"errors"
	"testing"
	"time"
)

func TestGetOrFetch_FreshValueServedWithoutRefetch(t *testing.T) {
	c := NewCache[int](time.Hour)
	calls := 0
	fetch := func() (int, error) {
		calls++
		return 42, nil
	}

	for i := 0; i < 3; i++ {
		v, err := c.GetOrFetch(fetch)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v != 42 {
			t.Fatalf("expected 42, got %d", v)
		}
	}

	if calls != 1 {
		t.Fatalf("expected fetch to be called once, got %d", calls)
	}
}

func TestGetOrFetch_StaleTriggersRefetch(t *testing.T) {
	c := NewCache[int](10 * time.Millisecond)
	calls := 0
	fetch := func() (int, error) {
		calls++
		return calls, nil
	}

	first, _ := c.GetOrFetch(fetch)
	if first != 1 {
		t.Fatalf("expected first value 1, got %d", first)
	}

	time.Sleep(20 * time.Millisecond)

	second, _ := c.GetOrFetch(fetch)
	if second != 2 {
		t.Fatalf("expected refetch to produce 2, got %d", second)
	}
	if calls != 2 {
		t.Fatalf("expected 2 fetch calls, got %d", calls)
	}
}

func TestGetOrFetch_ErrorWithPriorValueReturnsStale(t *testing.T) {
	c := NewCache[int](10 * time.Millisecond)

	// Prime the cache with a good value.
	_, err := c.GetOrFetch(func() (int, error) { return 99, nil })
	if err != nil {
		t.Fatalf("unexpected error priming cache: %v", err)
	}

	time.Sleep(20 * time.Millisecond)

	// Now the upstream fails; the stale value should still be returned.
	v, err := c.GetOrFetch(func() (int, error) { return 0, errors.New("upstream down") })
	if err != nil {
		t.Fatalf("expected no error when a stale value exists, got %v", err)
	}
	if v != 99 {
		t.Fatalf("expected stale value 99, got %d", v)
	}
}

func TestGetOrFetch_ErrorOnColdCacheReturnsError(t *testing.T) {
	c := NewCache[int](time.Hour)

	_, err := c.GetOrFetch(func() (int, error) { return 0, errors.New("upstream down") })
	if err == nil {
		t.Fatal("expected an error on a cold cache with a failing fetch")
	}
}

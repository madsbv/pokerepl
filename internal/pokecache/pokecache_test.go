package pokecache

import (
	"testing"
	"time"
)

func TestCacheGet(t *testing.T) {
	c := New(1 * time.Second)
	want := []byte{1, 2}
	c.Add("key", want)
	val, exists := c.Get("key")
	if len(val) != len(want) || val[0] != want[0] || val[1] != want[1] || !exists {
		t.Fatalf(`Get("key") = %v, %v, want %v, true.`, val, exists, want)
	}
}

func TestCacheReap(t *testing.T) {
	c := New(10 * time.Millisecond)
	temp := []byte{1}
	c.Add("key", temp)
	time.Sleep(20 * time.Millisecond)
	val, exists := c.Get("key")
	if val != nil || exists {
		t.Fatalf(`Get("key") = %v, %v, but this entry should have been evicted from the cache by now.`, val, exists)
	}
}

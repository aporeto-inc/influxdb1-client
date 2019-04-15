package models

import (
	"hash/fnv"
	"testing"
	"testing/quick"
)

func TestInlineFNV64aEquivalenceFuzz(t *testing.T) {
	f := func(data []byte) bool {
		stdlibFNV := fnv.New64a()
		stdlibFNV.Write(data) // nolint
		want := stdlibFNV.Sum64()

		inlineFNV := NewInlineFNV64a()
		inlineFNV.Write(data) // nolint
		got := inlineFNV.Sum64()

		return want == got
	}
	cfg := &quick.Config{
		MaxCount: 10000,
	}
	if err := quick.Check(f, cfg); err != nil {
		t.Fatal(err)
	}
}

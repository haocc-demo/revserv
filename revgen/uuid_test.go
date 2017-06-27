// Copyright 2017 <CompanyName>, Inc. All Rights Reserved.

package revgen

import (
	"strings"
	"testing"
	"unicode"
)

// Run 'go test -bench=".*"' to get an idea of how many UUID are generated per second.
func BenchmarkUuidGen(b *testing.B) {

	for i := 0; i < b.N; i++ {
		_ = getUuid()
	}
}

func TestUuidGen(t *testing.T) {

	sizehint := 1000
	var counters map[string]int = make(map[string]int, sizehint)

	for i := 0; i < sizehint; i++ {
		uuid := getUuid()
		counters[uuid]++
	}

	for k, v := range counters {
		if v > 1 {
			t.Errorf("Duplicate value %q", k)
		}
	}
}

// The generated UUID is used as a key, it must not have trailing space.
func TestUuidLength(t *testing.T) {

	// Result must be a string
	var uuid string = getUuid()

	// Result must have no spaces
	len1 := len(uuid)
	var uuidPrime = strings.TrimSpace(uuid)
	if len1 != len(uuidPrime) {
		t.Error("Result UUID has spaces")
	}

	if len(uuid) != 36 {
		t.Errorf("Expected length 36, got %d", len(uuid))
	}

	// Resulting bytes must all be printable
	for _, c := range uuid {

		if ok := unicode.IsPrint(rune(c)); !ok {
			t.Error("Unprintable character")
		}
	}
}

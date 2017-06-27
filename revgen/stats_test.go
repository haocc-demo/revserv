// Copyright 2017 <CompanyName>, Inc. All Rights Reserved.

package revgen

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestAccumulate(t *testing.T) {

	var s Stats
	for i := 0; i < 20; i++ {
		start := time.Now()
		time.Sleep(100 * time.Millisecond)
		s.Accumulate(start)
	}
	// Because accumulate is in the background, let some time elapse
	// before looking at results. Poor practice, use concurrency mechanisms instead.
	time.Sleep(300 * time.Millisecond)
	s.RLock()
	count := s.requestCount
	elapsed := s.elapsed.Seconds()
	fmt.Printf("count: %d\n", s.requestCount)
	fmt.Printf("elapsed: %f\n", elapsed)
	s.RUnlock()

	if count != 20 {
		t.Errorf("Expected count = %d, actual = %d", 20, count)
	}
	delta := math.Abs(elapsed - float64(2.0))
	fmt.Printf("delta: %f\n", delta)
	if delta > float64(0.2) {
		t.Errorf("Expected delta less than 0.2, actual = %f", delta)
	}
}

func TestGetJson(t *testing.T) {

	var s Stats
	for i := 0; i < 20; i++ {
		start := time.Now()
		time.Sleep(100 * time.Millisecond)
		s.Accumulate(start)
	}
	// Because accumulate is in the background, let some time elapse
	// before looking at results. Poor practice, use concurrency mechanisms instead.
	time.Sleep(100 * time.Millisecond)

	b := s.GetJson()
	fmt.Println(string(b) + "\n")
}

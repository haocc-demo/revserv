// Copyright 2017 <CompanyName>, Inc. All Rights Reserved.

package revgen

import (
	"encoding/json"
	"log"
	"math"
	"sync"
	"time"
)

type StatsMessage struct {

	// A request count
	Total int `json:"total,omitempty"`

	// Average duration for handling counted requests
	Average int64 `json:"average,omitempty"`

	// Units for average duration
	Units string `json:"units,omitempty"`
}

// A thread-safe container for tracking request duration
type Stats struct {

	// Same as adding the methods of sync.RWMutex
	sync.RWMutex

	// Total number of requests, for example to GET or POST /reverse
	requestCount int

	// Total elapsed time for handling counted requests
	elapsed time.Duration
}

func (s *Stats) Accumulate(start time.Time) {

	go func() {
		// Yields a duration
		elapsed := time.Since(start)
		s.Lock()
		defer s.Unlock()
		// For a long running process, protects against overflow
		if s.requestCount == math.MaxInt64 {
			s.requestCount = 1
			s.elapsed = elapsed
			return
		}
		s.requestCount++
		// TODO: best way to protect against overflow here?
		s.elapsed += elapsed
		return
	}()
}

// Marshals stats into JSON message response
func (s *Stats) GetJson() []byte {

	if s == nil {
		return []byte("{}")
	}

	var avg int64 = 0
	s.RLock()
	count := s.requestCount
	t := s.elapsed.Nanoseconds()
	s.RUnlock()
	// Protects against divide by zero
	if count > 0 {
		avg = t / int64(count)
	}
	m := StatsMessage{s.requestCount, avg, "ns"}
	b, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
	}
	return b
}

// Copyright 2017 <CompanyName>, Inc. All Rights Reserved.

// Provides server and utilities for generating reverse strings.
// In a package for reuse, since one can't import a main package.
package revgen

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"time"
)

const sizeHint = 8192

// Server bits
var (
	host      string
	port      string
	revServer *http.Server
	stats     Stats
	// Evicts after sizeHint entries, but all records persisted to disk.
	// The purpose is to put a limit on server memory usage for recording
	// job output. Of course, the persistent backing file needs monitoring
	// to alert when it approaches partition size.
	cache *LruCache = NewCache(sizeHint, "backup.json")
)

// For monitoring shutdown request
var (
	// Exported so that calling process can terminate gracefully
	Exit = make(chan struct{})

	// Used to internally coordinate graceful shutdown
	exit = make(chan struct{})

	// The shutdown request is Ctrl-c
	quit = make(chan os.Signal)

	// Reference count for jobs
	waiter sync.WaitGroup
)

// State toggles on shutdown request, so that additional jobs are blocked from starting
var shutdownState struct {
	sync.RWMutex        // same as adding the methods of sync.RWMutex
	isShutdownRequested bool
}

// A duration used for simulating work, units unspecified
const simDuration = 5

func HandleGetResult(w http.ResponseWriter, req *http.Request) {

	// TODO: Might be more elegant to wrap this handler in a timer handler
	start := time.Now()
	defer stats.Accumulate(start)

	// Validate parameters (last element of path)
	uuid := path.Base(req.URL.Path)
	if len(uuid) != 36 {
		http.Error(w, "Error: illegal Id", http.StatusBadRequest)
		return
	}

	userRecord, ok := cache.Get(uuid)
	if !ok {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	io.WriteString(w, string(userRecord.RevBytes)+"\n")
}

func HandleReverseRequest(w http.ResponseWriter, req *http.Request) {

	// TODO: Might be more elegant to wrap this handler in a timer handler
	start := time.Now()
	defer stats.Accumulate(start)

	// We can easily support a timeout by creating a context here.
	// See https://blog.golang.org/context

	// Validate parameters
	req.ParseForm()
	inputPhrase := req.FormValue("phrase")
	if len(inputPhrase) == 0 {
		http.Error(w, "Error: missing phrase", http.StatusBadRequest)
		return
	}

	// While holding read lock on shutdown not requested, add to wait group
	shutdownState.RLock()
	shutdownPending := shutdownState.isShutdownRequested
	if shutdownPending == false {
		waiter.Add(1)
	}
	shutdownState.RUnlock()
	if shutdownPending {
		http.Error(w, "Error: service shutdown pending", http.StatusServiceUnavailable)
		return
	}

	// Computes unique job identifier
	uuid := getUuid()

	// Launch a goroutine to compute the result
	go func() {
		// Requirement: Must sleep for n seconds.
		time.Sleep(simDuration * time.Second)
		// Decrements the waiter when goroutine completes
		defer waiter.Done()
		var r Reversable = []byte(inputPhrase)
		revBytes := r.Reverse()
		cache.Add(uuid, revBytes)
	}()

	// Returns the job identifier
	io.WriteString(w, uuid+"\n")
}

// Returns StatsMessage as Json
func HandleStatsRequest(w http.ResponseWriter, req *http.Request) {

	io.WriteString(w, string(stats.GetJson())+"\n")
}

func StartServer(host string, port string) *http.Server {

	s := []string{host, port}
	address := strings.Join(s, ":")
	revServer = &http.Server{Addr: address}

	// Register all endpoints
	// Measure duration for GET and POST handlers
	http.HandleFunc("/reverse", HandleReverseRequest)
	http.HandleFunc("/reverse/", HandleGetResult)
	http.HandleFunc("/stats", HandleStatsRequest)

	// Subscribe to SIGINT
	signal.Notify(quit, os.Interrupt)

	// The server supports a shutdown request. Work in progress completes while shutdown
	// pending. Therefore, run server in goroutine.
	go func() {

		// In the background, wait for a quit (Ctrl-c) request
		go func() {
			// If this function completes with a panic, run recovery.
			defer func() {
				if err := recover(); err != nil {
					log.Printf("Recovered: %+v\n", err)
				}
			}()
			// Blocks routine pending a quit (Ctrl-c) request
			<-quit
			// All jobs must complete prior to shutdown,
			// with the output persisted to disk.
			log.Println("Shutdown requested...")
			// Take a write lock for mutating state
			shutdownState.Lock()
			shutdownState.isShutdownRequested = true
			shutdownState.Unlock()
			// Waits for all jobs to complete
			waiter.Wait()
			log.Println("Done waiting")
			log.Println("Shutting down server...")
			// This is a closure, we can use revServer  from calling function
			if err := revServer.Shutdown(nil); err != nil {
				log.Fatalf("Failed to shutdown: %v", err)
			}
			// Notify thread monitoring the 'exit' channel
			close(exit)
		}()

		log.Printf("Starting server at %q\n", address)
		if err := revServer.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatalf("Httpserver: ListenAndServe() error: %s\n", err)
			}
			// Blocks routine pending graceful shutdown
			<-exit
			log.Println("Server shutdown complete. Bye!")
			// Notify main() monitoring the 'Exit' channel
			close(Exit)
		}
	}()

	return revServer
}

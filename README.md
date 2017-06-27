# revserv - String Reversing Service

## Summary

Provides an executable (revserv) that runs the string reversing server.

## Dependencies

Linux is supported.
* Depends on Go standard library
* Depends on package `haocc-demo/revserv/revgen` for server and utilties for reuse.
* OS commands `uuidgen` and `grep`

## Install

Clone `haocc-demo/revserv` to a standard Go workspace.

Example build and install:
```
$ cd $GOPATH/src/github.com/haocc-demo/revserv
$ go build
$ go install
```

### Starting the `revserv` server

```
$ $GOPATH/bin/revserv --host <hostname> --port <port number>
```

For a help message, use `revserv -h`.

### Endpoints

* POST /reverse to submit a job

Example:
```
$ curl --data "phrase=reverse me" http://localhost:8080/reverse
48765c95-9ea8-43bd-bc0b-5ecb6897f2b9
```
Submits a job.
Returns a unique job id that can be used to fetch the results for the given phrase.

* GET /reverse/{id}

Returns result of given job if available.
Example:
```
$ curl http://localhost:8080/reverse/48765c95-9ea8-43bd-bc0b-5ecb6897f2b9
em esrever
```

* GET /stats

Returns statistics for all GET and POST reverse requests for current server instance.

Example:
```
$ curl http://localhost:8080/stats
{"total":2001,"average":2543630,"units":"ns"}
```

### Stopping the `revserv` server

The server supports a graceful shutdown. Simply issue a SIGINT to the process.
All running jobs will complete before shutdown.
No additional job requests can start while shutdown is pending.
Example:
```
2017/06/01 12:58:39 Starting server at "localhost:8080"
2017/06/01 12:58:47 getting record with uuid: 48765c95-9ea8-43bd-bc0b-5ecb6897f2b9
^C2017/06/01 12:58:58 Shutdown requested...
2017/06/01 12:58:58 Done waiting
2017/06/01 12:58:58 Shutting down server...
2017/06/01 12:58:58 Server shutdown complete. Bye!
```

# Development

## Testing

The `revgen` package has tests.
```
$ cd revgen
$ go test -v
```

### Code Coverage

Sample coverage check:

```
$ cd revgen
$ go test -coverprofile=cover.out
...
$ go tool cover -func=cover.out
github.com/haocc-demo/revserv/revgen/dao.go:24:		New			100.0%
github.com/haocc-demo/revserv/revgen/dao.go:29:		Append			71.4%
github.com/haocc-demo/revserv/revgen/dao.go:65:		Get			100.0%
github.com/haocc-demo/revserv/revgen/lru.go:40:		NewCache		100.0%
github.com/haocc-demo/revserv/revgen/lru.go:56:		Add			84.6%
github.com/haocc-demo/revserv/revgen/lru.go:81:		Evict			88.9%
github.com/haocc-demo/revserv/revgen/lru.go:95:		Get			94.1%
github.com/haocc-demo/revserv/revgen/reversable.go:8:	Reverse			100.0%
github.com/haocc-demo/revserv/revgen/server.go:58:	HandleGetResult		100.0%
github.com/haocc-demo/revserv/revgen/server.go:82:	HandleReverseRequest	91.3%
github.com/haocc-demo/revserv/revgen/server.go:133:	HandleStatsRequest	100.0%
github.com/haocc-demo/revserv/revgen/server.go:138:	StartServer		90.3%
github.com/haocc-demo/revserv/revgen/stats.go:38:	Accumulate		72.7%
github.com/haocc-demo/revserv/revgen/stats.go:59:	GetJson			85.7%
github.com/haocc-demo/revserv/revgen/uuid.go:11:	getUuid			75.0%
total:							(statements)		88.0%
```

### Benchmarking

Some of the tests support benchmarks.

```
$ cd revgen
$ go test -bench=".*"
```

## Implementation Choices

### Job Id
The use of uuid instead of monotonically increasing id for a job identifier allows server
restart without having to persist the next id. Furthermore, no synchronization is required
on next id.

A drawback is that, without a database back end, ordering of uuid is less efficient for
persistent storage. A consequence is that records that are evicted from the cache will
take longer to find.

### LRU cache
The purpose for using the cache is to put an upper bound on memory usage for the server.

The design intent is that n most recently used jobs are held in cache, all are passed thru
to persistent storage. On a cache miss, job data could be fetched from persistent storage.
The default backing file is called 'backup.json'. One would prefer to use a performant data
store instead. The Go standard library has sql API, but no database drivers are included.

## Outstanding Development Items

* Raise code coverage
* Consider API versioning scheme
* Add command line option for backing file name
* Add guard against running the executable twice on same machine?
* Add safety for when backing file size approaches partition size
* Consider adding an option to reverse delimited words

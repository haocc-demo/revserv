// Copyright 2017 <Company Name>, Inc. All Rights Reserved.

package main

import (
	"flag"
	"fmt"
	"github.com/haocc-demo/revserv/revgen"
	"os"
)

// Filled by command line flags
var (
	host string
	port string
)

func init() {
	const (
		defaultHost = "localhost"
		defaultPort = "8080"
	)
	flag.StringVar(&host, "host", defaultHost, "Server Host")
	flag.StringVar(&port, "port", defaultPort, "Listen Port")
	// Custom usage message
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Starts string reversing server\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {

	flag.Parse()
	revgen.StartServer(host, port)

	// Block (waiting for exit signal)
	<-revgen.Exit
}

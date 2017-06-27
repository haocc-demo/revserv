// Copyright 2017 <CompanyName>, Inc. All Rights Reserved.

package revgen

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"sync"
)

// No database drivers are included in the Go standard library.
// Data access used in lieu of a database.
// Appends job results to file.

type UserAccountFile struct {
	sync.RWMutex

	filename string
}

func New(filename string) *UserAccountFile {

	return &UserAccountFile{sync.RWMutex{}, filename}
}

func (f *UserAccountFile) Append(value interface{}) {

	if f == nil {
		return
	}
	// Background file operation
	go func() {
		b, err := json.Marshal(value)
		if err != nil {
			log.Fatal(err)
		}
		f.Lock()
		defer f.Unlock()
		ofile, err := os.OpenFile(f.filename, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			if os.IsNotExist(err) {
				_, err = os.Create(f.filename)
				if err != nil {
					log.Fatalf("Failed to create/open %s", f.filename)
				}
			}
			ofile, err = os.OpenFile(f.filename, os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				log.Fatal(err)
			}
		}
		if _, err := ofile.WriteString(string(b) + "\n"); err != nil {
			log.Fatal(err)
		}
		// Don't defer writes due to possible OS caching
		if err := ofile.Close(); err != nil {
			log.Fatal(err)
		}
	}()
}

func (f *UserAccountFile) Get(uuid string) (record []byte, ok bool) {

	// Without an indexed database, finding any given uuid is a runtime concern.
	// We could use an elaborate sorted data file plus binary search, but this
	// is a toy application. Try the simplest thing that can work: grep!

	if _, err := os.Stat(f.filename); os.IsNotExist(err) {
		// No file, so not found
		return []byte(""), false
	}

	cmd := exec.Command("grep", uuid, f.filename)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return []byte(""), false
	}
	result := out.String()
	log.Println(result)
	return []byte(result), true
}

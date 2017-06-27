// Copyright 2017 <CompanyName>, Inc. All Rights Reserved.

package revgen

import (
	"log"
	"os"
	"testing"
	"time"
)

var userRecords = []UserRecord{

	{"id1", []byte("result1")},
	{"id2", []byte("result2")},
	{"id3", []byte("result3")},
	{"id4", []byte("result4")},
}

func TestAppendFileExists(t *testing.T) {

	var fh *UserAccountFile = New("test1.json")
	fh.Append(userRecords[0])
	fh.Append(userRecords[1])
}

func TestAppendNoFileExists(t *testing.T) {

	if _, err := os.Stat("test2.json"); !os.IsNotExist(err) {
		// File exists, delete it
		err = os.Remove("test2.json")
		if err != nil {
			t.Fatal(err)
		}
	}

	var fh *UserAccountFile = New("test2.json")
	fh.Append(userRecords[2])
}

func TestGet(t *testing.T) {

	var fh *UserAccountFile = New("test3.json")
	fh.Append(userRecords[0])
	fh.Append(userRecords[3])
	fh.Append(userRecords[2])
	// Poor practice, use other concurrency techniques instead
	time.Sleep(20 * time.Millisecond)
	result, ok := fh.Get(userRecords[3].Uuid)
	if !ok {
		t.Errorf("Record not found for uuid %s", userRecords[3].Uuid)
	}
	log.Println(result)

	if _, ok := fh.Get(userRecords[1].Uuid); ok {
		t.Errorf("Unexpected record found for uuid %s", userRecords[1].Uuid)
	}
}

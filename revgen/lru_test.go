// Copyright 2017 <CompanyName>, Inc. All Rights Reserved.

package revgen

import (
	"fmt"
	"os"
	"testing"
)

var testCases = []struct {
	desc     string
	addUuid  string
	getUuid  string
	expectOk bool
}{
	{"Add then Get", "abcd", "abcd", true},
	{"No Add then Get", "abcd", "xyz", false},
	{"Add 2", "a8c37fd2-ade0-4b7b-b62c-528016d73e1b", "a8c37fd2-ade0-4b7b-b62c-528016d73e1b", true},
}

func TestAddAndGet(t *testing.T) {

	if _, err := os.Stat("cachetest1.json"); !os.IsNotExist(err) {
		os.Remove("cachetest1.json")
	}
	cache := NewCache(0, "cachetest1.json")
	if cache.Capacity != 1000 {
		t.Fatalf("Expected capacity 1000, got %d", cache.Capacity)
	}

	for _, testcase := range testCases {
		fmt.Printf("testcase: %s\n", testcase.desc)
		cache.Add(testcase.addUuid, []byte("12345"))
		value, ok := cache.Get(testcase.getUuid)
		if ok != testcase.expectOk {
			t.Fatalf("%s expected Get to return %t, got %t\n", testcase.desc, testcase.expectOk, ok)
		}
		if ok && value == nil {
			t.Fatalf("%s failed\n", testcase.desc)
		} else if value != nil {
			fmt.Printf("type value is %T", value)
			fmt.Printf("uuid: %s revbytes: %s\n", value.Uuid, string(value.RevBytes))
		}
	}
}

func TestEvict(t *testing.T) {

	var userRecords = []UserRecord{
		UserRecord{"id1", []byte("result1")},
		UserRecord{"id2", []byte("result2")},
		UserRecord{"id3", []byte("result3")},
		UserRecord{"id4", []byte("result4")},
		UserRecord{"id5", []byte("result5")},
		UserRecord{"id6", []byte("result6")},
		UserRecord{"id7", []byte("result7")},
	}
	if _, err := os.Stat("cachetest2.json"); !os.IsNotExist(err) {
		os.Remove("cachetest2.json")
	}
	// Creates a cache with a small capacity
	cache := NewCache(3, "cachetest2.json")
	if cache.Capacity != 3 {
		t.Fatalf("Expected capacity 3, got %d", cache.Capacity)
	}
	for _, record := range userRecords {
		cache.Add(record.Uuid, record.RevBytes)
	}
	for _, record := range userRecords {
		if r, ok := cache.Get(record.Uuid); !ok {
			t.Errorf("Record %s not found", record.Uuid)
		} else {
			fmt.Println(r)
		}
	}
}

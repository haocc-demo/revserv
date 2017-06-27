// Copyright 2017 <CompanyName>, Inc. All Rights Reserved.

package revgen

import (
    "fmt"
//    "encoding/base64"
    "reflect"
    "testing"
)

func TestReversable(t *testing.T) {

    var testCases = []struct {
        input string
        expectedResult string
    }{
        {"word1 word2", "2drow 1drow"},
        {"12345 678", "876 54321"},
        {"ABCdef$#", "#$fedCBA"},
    }

    for _, testcase := range testCases {

        var byteBuffer Reversable = []byte(testcase.input)
        result := byteBuffer.Reverse()
        //encodedResult := base64.URLEncoding.EncodeToString(result)
        fmt.Printf("phrase: %q expected result: %q\n", testcase.input, result)

        expectedResult := []byte(testcase.expectedResult)
        if !reflect.DeepEqual(result, expectedResult) {
            t.Errorf("Expected %q, got %q", expectedResult, result)
        }
    }
}


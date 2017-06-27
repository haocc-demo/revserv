// Copyright 2017 <CompanyName>, Inc. All Rights Reserved.

package revgen

import (
	"log"
	"os/exec"
	"strings"
)

func getUuid() string {

	result, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	// UUID to be used as a key, must have no trailing newline.
	return strings.TrimSpace(string(result))
}

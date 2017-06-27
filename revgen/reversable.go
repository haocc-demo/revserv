// Copyright 2017 <Company Name>. All Rights Reserved.

package revgen

// A reversable byte array.
type Reversable []byte

func (r Reversable) Reverse() []byte {

	mutable := []byte(r)

	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		mutable[i], mutable[j] = mutable[j], mutable[i]
	}
	return mutable
}

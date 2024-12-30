package hash

import (
	"testing"
)

func TestMakeSHA256(t *testing.T) {
	hash, _ := MakeSHA256([]byte(`Hello world!`), `key`)
	if hash != `d6979d591536f5aee0003a13c9b6deedaca21b8ee2da1c7a7cdec4d51d1ae67d` {
		t.Error("hash error")
	}
	t.Log(hash)
}

//852d2fec4bda6add8f12c5c1dff8420510ac5b85ef432140c7097aaee3c270ca
